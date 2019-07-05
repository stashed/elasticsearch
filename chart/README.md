# Elassricseach-stash

[elasticsearch-stash by AppsCode](https://github.com/stashed/elasticsearch-stash) - Elassricseach database backup/restore plugin for [Stash](https://github.com/stashed/).

## TL;DR;

```console
helm repo add appscode https://charts.appscode.com/stable/
helm repo update
helm install appscode/elasticsearch-stash --name=elasticsearch-stash-6.3 --version=6.3
```

## Introduction

This chart installs necessary `Function` and `Task` definition to backup or restore Elassricseach database 6.3 using Stash.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart

- Add AppsCode chart repository to your helm repository list,

```console
helm repo add appscode https://charts.appscode.com/stable/
```

- Update helm repositories to fetch latest charts from the remove repository,

```console
helm repo update
```

- Install the chart with the release name `elasticsearch-stash-6.3` run the following command,

```console
helm install appscode/elasticsearch-stash --name=elasticsearch-stash-6.3 --version=6.3
```

The above commands installs `Functions` and `Task` crds that are necessary to backup Elassricseach database 6.3 using Stash.

## Uninstalling the Chart

To uninstall/delete the `elasticsearch-stash-6.3` run the following command,

```console
helm delete elasticsearch-stash-6.3
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the `postgre-stash` chart and their default values.

| Parameter                   | Description                                                  | Default               |
| --------------------------- | :----------------------------------------------------------- | --------------------- |
| `global.registry`           | Docker registry used to pull respective images               | `appscode`            |
| `global.image`              | Docker image used to backup/restore PosegreSQL database      | `elasticsearch-stash` |
| `global.tag`                | Tag of the image that is used to backup/restore Elassricseach database. This is usually same as the database version it can backup. | `6.3`                 |
| `global.backup.esArgs`      | Optional arguments to pass to `esdump` command  while bakcup |                       |
| `global.restore.esArgs`     | Optional arguments to pass to `psql` command while restore   |                       |
| `global.metrics.enabled`    | Specifies whether to send Prometheus metrics                 | `true`                |
| `global.metrics.labels`     | Optional comma separated labels to add to the Prometheus metrics |                       |
| `persistence.enabled`       | Enable persistence using PVC. If `false`, a `empty directory` volume will be used  for elastic backup directory. | `false`               |
| `persistence.existingClaim` | Provide an existing `PersistentVolumeClaim`, the value is evaluated as a template. | `nil`                 |
| `persistence.namespace`     | The namespace where `pvc` is created. This namespace needs to be same as `backupsession` or  `restoresession` namespace. | `default`             |
| `persistence.storageClass`  | PVC Storage Class for Elasticsearch volume                   | `nil`                 |
| `persistence.accessModes`   | PVC Access Mode for Elasticsearch volume                     | `[ReadWriteOnce]`     |
| `persistence.size`          | PVC Storage Request for Elasticsearch volume                 | `8Gi`                 |
| `persistence.annotations`   | Annotations for the PVC                                      | `{}`                  |

> We have declared all the configurable parameters as global parameter so that the parent chart can overwrite them.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

For example:

```console
helm install --name elasticsearch-stash-6.3 --set global.metrics.enabled=false appscode/elasticsearch-stash
```

**Tips:** Use escape character (`\`) while providing multiple comma-separated labels for `global.metrics.labels`.

```console
 helm install chart/elasticsearch-stash --set global.metrics.labels="k1=v1\,k2=v2"
```
