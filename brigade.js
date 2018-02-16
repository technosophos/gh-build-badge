const { events, Job } = require("brigadier");

const ACTION_MOVE = "action_move_card_from_list_to_list"
const chartName = "/src/charts/gh-build-badge"
const relName = "friday-demo"

events.on("exec", (e, p) => {
  console.log(`installing ${chartName} into ${relName}`)

  var helm = new Job("helm", "lachlanevenson/k8s-helm:v2.6.1");
  helm.tasks = [
    "helm init --client-only",
    `helm install ${chartName} -n ${relName}`
  ]

  helm.run()
});

events.on("webhook", (e, p) => {
  console.log(e.provider)
  var slack = new Job("slack-notify", "technosophos/slack-notify:latest", ["/slack-notify"])
  var m = "This is a message from Butcher's watch"
  slack.storage.enabled = false
  slack.env = {
    SLACK_WEBHOOK: p.secrets.SLACK_WEBHOOK,
    SLACK_USERNAME: "FitBit",
    SLACK_TITLE: "Watch Me!",
    SLACK_MESSAGE: m
    //SLACK_ICON: "https://a.trellocdn.com/images/ios/0307bc39ec6c9ff499c80e18c767b8b1/apple-touch-icon-152x152-precomposed.png"
  }
  slack.run()
});

events.on("trello", (e, p) => {
  console.log(e.payload);
  const hook = JSON.parse(e.payload)
  const d = hook.action.display
  if (d.translationKey != ACTION_MOVE) {
    return
  }

  var s = newSlack(p, hook)
  s.run()

  if (d.entities.listAfter.name == "Up") {
    var helm = new Job("helm", "lachlanevenson/k8s-helm:v2.6.1");
    helm.tasks = [
      `helm upgrade ${relName} ${chartName} --set replicaCount=2`
    ];
    helm.run()
  } else if (d.entities.listAfter == "Down") {
    var helm = new Job("helm", "lachlanevenson/k8s-helm:v2.6.1");
    helm.tasks = [
      `helm upgrade ${relName} ${chartName} --set replicaCount=1`
    ];
    helm.run()
  }
});

events.on("ReplicaSet:SuccessfulCreate", (e, p) => {
  var trello = new Job("trello", "technosophos/trello-cli:latest")
  trello.env = {
    APIKEY: p.secrets.trelloKey,
    TOKEN: p.secrets.trelloToken
  }
  trello.tasks = [
    "env2creds",
    "trello refresh",
    `trello add-card ${relName} ${relName} -b Ionic -l New`
  ]

  trello.run()
})

events.on("ReplicaSet:SuccessfulDelete", (e, p) => {
  var trello = new Job("trello", "technosophos/trello-cli:latest")
  trello.env = {
    APIKEY: p.secrets.trelloKey,
    TOKEN: p.secrets.trelloToken
  }
  trello.tasks = [
    "env2creds",
    "trello refresh",
    `trello delete-card ${relName} -b Ionic`
  ]

  trello.run()
});


function newSlack(p, hook) {
  const d = hook.action.display
  var slack = new Job("slack-notify", "technosophos/slack-notify:latest", ["/slack-notify"])
  var m = `From "${d.entities.listBefore.text}" to "${d.entities.listAfter.text}" <${hook.model.shortUrl}> <@U0RMKK605>`  
  slack.storage.enabled = false
  slack.env = {
    SLACK_WEBHOOK: p.secrets.SLACK_WEBHOOK,
    SLACK_USERNAME: "Trello",
    SLACK_TITLE: `Moved "${d.entities.card.text}"`,
    SLACK_MESSAGE: m,
    SLACK_ICON: "https://a.trellocdn.com/images/ios/0307bc39ec6c9ff499c80e18c767b8b1/apple-touch-icon-152x152-precomposed.png"
  }
  return slack;
}
