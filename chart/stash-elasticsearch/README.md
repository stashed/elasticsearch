# Stash-elasticserach

[stash-elasticsearch](https://github.com/stashed/stash-elasticsearch) - Elasticsearch database backup/restore plugin for [Stash by AppsCode](https://appscode.com/products/stash/).

## TL;DR;

```console
helm repo add appscode https://charts.appscode.com/stable/
helm repo update
helm install appscode/stash-elasticsearch --name=stash-elasticsearch-6.3 --version=6.3
```

## Introduction

This chart installs necessary `Function` and `Task` definition to backup or restore Elasticsearch database 6.3 using Stash.

## Prerequisites

- Kubernetes 1.11+

## Installing the Chart

- Add AppsCode chart repository to your helm repository list,

```console
helm repo add appscode https://charts.appscode.com/stable/
```

- Update helm repositories to fetch latest charts from the remove repository,

```console
helm repo update
```

- Install the chart with the release name `stash-elasticsearch-6.3` run the following command,

```console
helm install appscode/stash-elasticsearch --name=stash-elasticsearch-6.3 --version=6.3
```

The above commands installs `Functions` and `Task` crds that are necessary to backup Elasticsearch database 6.3 using Stash.

## Uninstalling the Chart

To uninstall/delete the `stash-elasticsearch-6.3` run the following command,

```console
helm delete stash-elasticsearch-6.3
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the `postgre-stash` chart and their default values.

|          Parameter          |                                                             Description                                                             |        Default        |
| --------------------------- | :---------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| `docker.registry`           | Docker registry used to pull respective images                                                                                      | `stashed`             |
| `docker.image`              | Docker image used to backup/restore PosegreSQL database                                                                             | `stash-elasticsearch` |
| `docker.tag`                | Tag of the image that is used to backup/restore Elasticsearch database. This is usually same as the database version it can backup. | `6.3`                 |
| `backup.esArgs`             | Optional arguments to pass to `multielasticdump` command  while bakcup                                                              |                       |
| `restore.esArgs`            | Optional arguments to pass to `multielasticdump` command while restore                                                              |                       |
| `metrics.enabled`           | Specifies whether to send Prometheus metrics                                                                                        | `true`                |
| `metrics.labels`            | Optional comma separated labels to add to the Prometheus metrics                                                                    |                       |
| `persistence.enabled`       | Enable persistence using PVC. If `false`, a `empty directory` volume will be used  for elastic backup directory.                    | `false`               |
| `persistence.existingClaim` | Provide an existing `PersistentVolumeClaim`, the value is evaluated as a template.                                                  | `nil`                 |
| `persistence.namespace`     | The namespace where `pvc` is created. This namespace needs to be same as `backupsession` or  `restoresession` namespace.            | `default`             |
| `persistence.storageClass`  | PVC Storage Class for Elasticsearch volume                                                                                          | `nil`                 |
| `persistence.accessModes`   | PVC Access Mode for Elasticsearch volume                                                                                            | `[ReadWriteOnce]`     |
| `persistence.size`          | PVC Storage Request for Elasticsearch volume                                                                                        | `8Gi`                 |
| `persistence.annotations`   | Annotations for the PVC                                                                                                             | `{}`                  |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

For example:

```console
helm install --name stash-elasticsearch-6.3 --set metrics.enabled=false appscode/stash-elasticsearch
```

**Tips:** Use escape character (`\`) while providing multiple comma-separated labels for `metrics.labels`.

```console
 helm install chart/stash-elasticsearch --set metrics.labels="k1=v1\,k2=v2"
```
