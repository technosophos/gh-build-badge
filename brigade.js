const { events, Job } = require("brigadier");

const ACTION_MOVE = "action_move_card_from_list_to_list"
const chartName = "/src/charts/gh-build-badge"
const relName = "friday-demo"

events.on("exec", (e, p) => {
  console.log("this space intentionally left blank")
});

events.on("webhook", (e, p) => {
  console.log(e.provider)
  var slack = new Job("slack-notify", "technosophos/slack-notify:latest", ["/slack-notify"])
  var m = "Installing Helm chart"
  slack.storage.enabled = false
  slack.env = {
    SLACK_WEBHOOK: p.secrets.SLACK_WEBHOOK,
    SLACK_USERNAME: "FitBit",
    SLACK_TITLE: "Watch Me!",
    SLACK_MESSAGE: m
    //SLACK_ICON: "https://a.trellocdn.com/images/ios/0307bc39ec6c9ff499c80e18c767b8b1/apple-touch-icon-152x152-precomposed.png"
  }
  console.log(`installing ${chartName} into ${relName}`)

  var helm = new Job("helm", "lachlanevenson/k8s-helm:v2.6.1");
  helm.tasks = [
    "helm init --client-only",
    `helm install ${chartName} -n ${relName}`
  ]

  slack.run().then(() => {return helm.run()})
});

events.on("ReplicaSet:SuccessfulCreate", (e, p) => {
  newSlack("Created replica set", p, e).run()
})

events.on("ReplicaSet:SuccessfulDelete", (e, p) => {
  newSlack("Deleted replica set", p, e).run()
});


function newSlack(msg, p, e) {
  var slack = new Job("slack-notify", "technosophos/slack-notify:latest", ["/slack-notify"])
  var m = `${msg} <${hook.model.shortUrl}> <@U0RMKK605>`
  slack.env = {
    SLACK_WEBHOOK: p.secrets.SLACK_WEBHOOK,
    SLACK_USERNAME: "Trello",
    SLACK_TITLE: `Handled Event ${e.type}`,
    SLACK_MESSAGE: m,
    SLACK_ICON: "https://a.trellocdn.com/images/ios/0307bc39ec6c9ff499c80e18c767b8b1/apple-touch-icon-152x152-precomposed.png"
  }
  return slack;
}

