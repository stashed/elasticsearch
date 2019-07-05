[![Go Report Card](https://goreportcard.com/badge/stash.appscode.dev/elasticsearch)](https://goreportcard.com/report/stash.appscode.dev/elasticsearch)
[![Build Status](https://travis-ci.org/stashed/elasticsearch.svg?branch=master)](https://travis-ci.org/stashed/elasticsearch)
[![Docker Pulls](https://img.shields.io/docker/pulls/appscode/elasticsearch-stash.svg)](https://hub.docker.com/r/appscode/elasticsearch-stash/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)

# Elasticsearch

Elasticsearch backup and restore plugin for [Stash by AppsCode](https://appscode.com/products/stash).

## Install

Install Elasticsearch 6.3 backup or restore plugin for Stash as below.

**Chart:**

```console
helm repo add appscode https://charts.appscode.com/stable/
helm repo update
helm install appscode/elasticsearch-stash --name=elasticsearch-stash-6.3 --version=6.3
```

**Script:**

```console
curl -fsSL https://github.com/stashed/elasticsearch/raw/6.3/hack/setup.sh | bash
```

## Uninstall

Uninstall Elasticsearch 6.3 backup or restore plugin for Stash as below.

**Chart:**

```console
helm delete elasticsearch-stash-6.3
```

**Script:**

```console
curl -fsSL https://github.com/stashed/elasticsearch/raw/6.3/hack/setup.sh | bash -s -- --uninstall
```

## More Options 

[Read setup guide](/chart/README.md) to learn about available installation configurations.

[Read quickstart guide](/docs/elasticsearch.md) to learn how to use stash-elasticsearch to take backup from and to restore backup to a elasticsearch.

## Support

We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/C8NCX6N23/details/) channel `#stash`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

If you have found a bug with Stash or want to request for new features, please [file an issue](https://github.com/stashed/stash/issues/new).
