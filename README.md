[![Go Report Card](https://goreportcard.com/badge/stash.appscode.dev/elasticsearch)](https://goreportcard.com/report/stash.appscode.dev/elasticsearch)
![CI](https://github.com/stashed/elasticsearch/workflows/CI/badge.svg)
[![Docker Pulls](https://img.shields.io/docker/pulls/stashed/stash-elasticsearch.svg)](https://hub.docker.com/r/stashed/stash-elasticsearch/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/kubestash.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=KubeStash)

# Elasticsearch

Elasticsearch backup and restore plugin for [Stash by AppsCode](https://stash.run).

## Install

Install Elasticsearch 6.8.0 backup or restore plugin for Stash as below.

```console
helm repo add appscode https://charts.appscode.com/stable/
helm repo update
helm install appscode/stash-elasticsearch --name=stash-elasticsearch-6.8.0 --version=6.8.0
```

To install catalog for all supported Elasticsearch versions, please visit [here](https://github.com/stashed/catalog).

## Uninstall

Uninstall Elasticsearch 6.8.0 backup or restore plugin for Stash as below.

```console
helm delete stash-elasticsearch-6.8.0
```

## Support

To speak with us, please leave a message on [our website](https://appscode.com/contact/).

To join public discussions with the KubeDB community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/C8NCX6N23/details/) channel `#stash`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

To receive product annoucements, follow us on [Twitter](https://twitter.com/KubeStash).

If you have found a bug with Stash or want to request new features, please [file an issue](https://github.com/stashed/project/issues/new).
