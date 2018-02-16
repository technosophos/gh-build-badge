# GitHub Build Badges

This is a simple server for displaying build badges on GitHub projects.

It is designed to run within Kubernetes.

## Installation

This assumes you have a Kubernetes cluster and the Helm package manager.

1. Clone this repository and `cd` into it
2. Install with Helm

```
$ helm install charts/gh-build-badge
```

You probably want to either make the service a LoadBalancer or enable an ingress.
To learn more, run `helm inspect values charts/gh-build-badge`

## Building

To build, run `make build`

To build Docker images, run `make docker-build`
