[![Go Report Card](https://goreportcard.com/badge/stash.appscode.dev/elasticsearch)](https://goreportcard.com/report/stash.appscode.dev/elasticsearch)
[![Build Status](https://travis-ci.org/stashed/elasticsearch.svg?branch=master)](https://travis-ci.org/stashed/elasticsearch)
[![Docker Pulls](https://img.shields.io/docker/pulls/stashed/stash-elasticsearch.svg)](https://hub.docker.com/r/stashed/stash-elasticsearch/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)

# Elasticsearch

Elasticsearch backup and restore plugin for [Stash by AppsCode](https://appscode.com/products/stash).

## Install

Install Elasticsearch 7.2.0 backup or restore plugin for Stash as below.

```console
helm repo add appscode https://charts.appscode.com/stable/
helm repo update
helm install appscode/stash-elasticsearch --name=stash-elasticsearch-7.2.0 --version=7.2.0
```

To install catalog for all supported Elasticsearch versions, please visit [here](https://github.com/stashed/catalog).

## Uninstall

Uninstall Elasticsearch 7.2.0 backup or restore plugin for Stash as below.

```console
helm delete stash-elasticsearch-7.2.0
```

## Support

We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/C8NCX6N23/details/) channel `#stash`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

If you have found a bug with Stash or want to request for new features, please [file an issue](https://github.com/stashed/project/issues/new).
