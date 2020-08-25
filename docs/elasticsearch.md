---
title: Elasticsearch | Stash
description: Backup and restore Elasticsearch database using Stash
menu:
  docs_{{ .version }}:
    identifier: standalone-elasticsearch-{{ .subproject_version }}
    name: Backup & Restore Elasticsearch
    parent: stash-elasticsearch-guides-{{ .subproject_version }}
    weight: 10
product_name: stash
menu_name: docs_{{ .version }}
section_menu_id: stash-addons
---

# Backup and Restore Elasticsearch database using Stash

Stash 0.9.0+ supports backup and restoration of Elasticsearch clusters. This guide will show you how you can backup and restore your Elasticsearch database with Stash.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the `kubectl` command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using Minikube.
- Install Stash in your cluster following the steps [here](/docs/setup/README.md).
- Install Elasticsearch addon for Stash following the steps [here](/docs/addons/elasticsearch/setup/install.md).
- Install [KubeDB](https://kubedb.com) in your cluster following the steps [here](https://kubedb.com/docs/latest/setup/install/). This step is optional. You can deploy your database using any method you want. We are using KubeDB because KubeDB simplifies many of the difficult or tedious management tasks of running a production grade databases on private and public clouds.
- If you are not familiar with how Stash backup and restore Elasticsearch databases, please check the following guide [here](/docs/addons/elasticsearch/overview.md).

You have to be familiar with following custom resources:

- [AppBinding](/docs/concepts/crds/appbinding.md)
- [Function](/docs/concepts/crds/function.md)
- [Task](/docs/concepts/crds/task.md)
- [BackupConfiguration](/docs/concepts/crds/backupconfiguration.md)
- [RestoreSession](/docs/concepts/crds/restoresession.md)

To keep things isolated, we are going to use a separate namespace called `demo` throughout this tutorial. Create `demo` namespace if you haven't created yet.

```console
$ kubectl create ns demo
namespace/demo created
```

>Note: YAML files used in this tutorial are stored [here](https://github.com/stashed/elasticsearch/tree/{{< param "info.subproject_version" >}}/docs/examples).

## Backup Elasticsearch

This section will demonstrate how to backup an Elasticsearch database. Here, we are going to deploy an Elasticsearch database using KubeDB. Then, we are going to backup this database into a GCS bucket. Finally, we are going to restore the backed up data into another Elasticsearch database.

### Deploy Sample Elasticsearch Database

Let's deploy a sample Elasticsearch database and insert some data into it.

**Create Elasticsearch CRD:**

Below is the YAML of a sample Elasticsearch crd that we are going to create for this tutorial:

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: sample-elasticsearch
  namespace: demo
spec:
  version: "6.5.3"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  terminationPolicy: Delete
```

Create the above `Elasticsearch` crd,

```console
$ kubectl apply -f https://github.com/stashed/elasticsearch/raw/{{< param "info.subproject_version" >}}/docs/examples/backup/elasticsearch.yaml
elasticsearch.kubedb.com/sample-elasticsearch created
```

KubeDB will deploy an Elasticsearch database according to the above specification. It will also create the necessary secrets and services to access the database.

Let's check if the database is ready to use,

```console
$ kubectl get es -n demo sample-elasticsearch
NAME                   VERSION       STATUS    AGE
sample-elasticsearch   6.5.3         Running   3m35s
```

The database is `Running`. Verify that KubeDB has created a Secret and a Service for this database using the following commands,

```console
$ kubectl get secret -n demo -l=kubedb.com/name=sample-elasticsearch
NAME                        TYPE     DATA   AGE
sample-elasticsearch-auth   Opaque   9      5m14s
sample-elasticsearch-cert   Opaque   4      5m14s

$ kubectl get service -n demo -l=kubedb.com/name=sample-elasticsearch
NAME                          TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
sample-elasticsearch          ClusterIP   10.102.109.209   <none>        9200/TCP   5m36s
sample-elasticsearch-master   ClusterIP   10.107.148.74    <none>        9300/TCP   5m36s
```

Here, we have to use service `sample-elasticsearch` and secret `sample-elasticsearch-auth` to connect with the database. KubeDB creates an [AppBinding](/docs/concepts/crds/appbinding.md) crd that holds the necessary information to connect with the database.

**Verify AppBinding:**

Verify that the `AppBinding` has been created successfully using the following command,

```console
$ kubectl get appbindings -n demo
NAME                   AGE
sample-elasticsearch   4m16s
```

Let's check the YAML of the above `AppBinding`,

```console
$ kubectl get appbindings -n demo sample-elasticsearch -o yaml
```

```yaml
apiVersion: appcatalog.appscode.com/v1alpha1
kind: AppBinding
metadata:
  labels:
    app.kubernetes.io/component: database
    app.kubernetes.io/instance: sample-elasticsearch
    app.kubernetes.io/managed-by: kubedb.com
    app.kubernetes.io/name: elasticsearch
    app.kubernetes.io/version: 6.5.3
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: sample-elasticsearch
  name: sample-elasticsearch
  namespace: demo
spec:
  clientConfig:
    service:
      name: sample-elasticsearch
      port: 9200
      scheme: http
  secret:
    name: sample-elasticsearch-auth
  secretTransforms:
  - renameKey:
      from: ADMIN_USERNAME
      to: username
  - renameKey:
      from: ADMIN_PASSWORD
      to: password
  type: kubedb.com/elasticsearch
  version: "6.5.3"
```

Stash uses the `AppBinding` crd to connect with the target database. It requires the following two fields to set in AppBinding's `Spec` section.

- `spec.clientConfig.service.name` specifies the name of the service that connects to the database.
- `spec.secret` specifies the name of the secret that holds necessary credentials to access the database.
- `spec.type` specifies the types of the app that this AppBinding is pointing to. KubeDB generated AppBinding follows the following format: `<app group>/<app resource type>`.

**Creating AppBinding Manually:**

If you deploy Elasticsearch database without KubeDB, you have to create the AppBinding crd manually in the same namespace as the service and secret of the database.

The following YAML shows a minimal AppBinding specification that you have to create if you deploy Elasticsearch database without KubeDB.

```yaml
apiVersion: appcatalog.appscode.com/v1alpha1
kind: AppBinding
metadata:
  name: my-custom-appbinding
  namespace: my-database-namespace
spec:
  clientConfig:
    service:
      name: my-database-service
      port: 9200
      scheme: http
  secret:
    name: my-database-credentials-secret
  # type field is optional. you can keep it empty.
  # if you keep it emtpty then the value of TARGET_APP_RESOURCE variable
  # will be set to "appbinding" during auto-backup.
  type: elasticsearch
```

**Connection information:**

- **Address:** localhost:9200
- **Username:** Run following command to get username

```bash
$ kubectl get secrets -n demo sample-elasticsearch-auth -o jsonpath='{.data.\ADMIN_USERNAME}' | base64 -d
admin
```

- **Password:** Run the following command to get the password

```bash
$ kubectl get secrets -n demo sample-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
mti4sbyy
```

**Insert Sample Data:**

Now, we are going to exec into the database pod and create some sample data. At first, find out the database pod using the following command,

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=sample-elasticsearch"
NAME                     READY   STATUS    RESTARTS   AGE
sample-elasticsearch-0   1/1     Running   0          7m33s
```

Now, let's exec into the pod and create a table,

```console
$ kubectl exec -it -n demo sample-elasticsearch-0 bash
~ curl -XPUT --user "admin:mti4sbyy" "localhost:9200/test/snapshot/1?pretty" -H 'Content-Type: application/json' -d'
{
    "title": "Snapshot",
    "text":  "Testing instant backup",
    "date":  "2018/02/13"
}
'

~ curl -XGET --user "admin:mti4sbyy" "localhost:9200/test/snapshot/1?pretty"
{
  "_index" : "test",
  "_type" : "snapshot",
  "_id" : "1",
  "_version" : 1,
  "found" : true,
  "_source" : {
    "title" : "Snapshot",
    "text" : "Testing instant backup",
    "date" : "2018/02/13"
  }
}
```

Now, we are ready to backup this sample database.

### Prepare Backend

We are going to store our backed up data into a GCS bucket. At first, we need to create a secret with GCS credentials then we need to create a `Repository` crd. If you want to use a different backend, please read the respective backend configuration doc from [here](/docs/guides/latest/backends/overview.md).

**Create Storage Secret:**

Let's create a secret called `gcs-secret` with access credentials to our desired GCS bucket,

```console
$ echo -n 'changeit' > RESTIC_PASSWORD
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ cat downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic -n demo gcs-secret \
    --from-file=./RESTIC_PASSWORD \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret/gcs-secret created
```

**Create Repository:**

Now, crete a `Respository` using this secret. Below is the YAML of Repository crd we are going to create,

```yaml
apiVersion: stash.appscode.com/v1alpha1
kind: Repository
metadata:
  name: gcs-repo
  namespace: demo
spec:
  backend:
    gcs:
      bucket: appscode-qa
      prefix: /demo/elasticsearch/sample-elasticsearch
    storageSecretName: gcs-secret
```

Let's create the `Repository` we have shown above,

```console
$ kubectl apply -f https://github.com/stashed/elasticsearch/raw/{{< param "info.subproject_version" >}}/docs/examples/backup/repository.yaml
repository.stash.appscode.com/gcs-repo created
```

Now, we are ready to backup our database to our desired backend.

### Backup

We have to create a `BackupConfiguration` targeting respective AppBinding crd of our desired database. Then Stash will create a CronJob to periodically backup the database.

**Create BackupConfiguration:**

Below is the YAML for `BackupConfiguration` crd to backup the `sample-elasticsearch` database we have deployed earlier.,

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: BackupConfiguration
metadata:
  name: sample-elasticsearch-backup
  namespace: demo
spec:
  schedule: "*/5 * * * *"
  task:
    name: elasticsearch-backup-6.5.3
  repository:
    name: gcs-repo
  target:
    ref:
      apiVersion: appcatalog.appscode.com/v1alpha1
      kind: AppBinding
      name: sample-elasticsearch
  interimVolumeTemplate:
    metadata:
      name: stash-tmp-storage
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: "standard"
      resources:
        requests:
          storage: 1Gi
  retentionPolicy:
    name: keep-last-5
    keepLast: 5
    prune: true
```

Here,

- `spec.schedule` specifies that we want to backup the database at 5 minutes interval.
- `spec.task.name` specifies the name of the task crd that specifies the necessary Function and their execution order to backup an Elasticsearch database.
- `spec.target.ref` refers to the `AppBinding` crd that was created for `sample-elasticsearch` database.
- `spec.interimVolumeTemplate` specifies a PVC template where the dumped data will be stored temporarily before uploading to the backend.

Let's create the `BackupConfiguration` crd we have shown above,

```console
$ kubectl apply -f https://github.com/stashed/elasticsearch/raw/{{< param "info.subproject_version" >}}/docs/examples/backup/backupconfiguration.yaml
backupconfiguration.stash.appscode.com/sample-elasticsearch-backup created
```

**Verify CronJob:**

If everything goes well, Stash will create a CronJob with the schedule specified in `spec.schedule` field of `BackupConfiguration` crd.

Verify that the CronJob has been created using the following command,

```console
$ kubectl get cronjob -n demo
NAME                          SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
sample-elasticsearch-backup   */5 * * * *   False     0        <none>          10s
```

**Wait for BackupSession:**

The `sample-elasticsearch-backup` CronJob will trigger a backup on each scheduled slot by creating a `BackupSession` crd.

Wait for the next schedule. Run the following command to watch `BackupSession` crd,

```console
$ kubectl get backupsession -n demo -w
NAME                                     INVOKER-TYPE          INVOKER-NAME                  PHASE       AGE
sample-elasticsearch-backup-1570098367   BackupConfiguration   sample-elasticsearch-backup   Running     9s
sample-elasticsearch-backup-1570098367   BackupConfiguration   sample-elasticsearch-backup   Running     65s
sample-elasticsearch-backup-1570098367   BackupConfiguration   sample-elasticsearch-backup   Succeeded   79s
```

We can see above that the backup session has succeeded. Now, we are going to verify that the backed up data has been stored in the backend.

**Verify Backup:**

Once a backup is complete, Stash will update the respective `Repository` crd to reflect the backup. Check that the repository `gcs-repo` has been updated by the following command,

```console
$ kubectl get repository -n demo gcs-repo
NAME       INTEGRITY   SIZE        SNAPSHOT-COUNT   LAST-SUCCESSFUL-BACKUP   AGE
gcs-repo   true        1.030 KiB   1                39s                      3m20s
```

Now, if we navigate to the GCS bucket, we are going to see backed up data has been stored in `demo/elasticsearch/sample-elasticsearch` directory as specified by `spec.backend.gcs.prefix` field of Repository crd.

>Note: Stash keeps all the backed up data encrypted. So, data in the backend will not make any sense until they are decrypted.

## Restore Elasticsearch

In this section, we are going to restore the database from the backup we have taken in the previous section. We are going to deploy a new database and initialize it from the backup.

**Stop Taking Backup of the Old Database:**

At first, let's stop taking any further backup of the old database so that no backup is taken during restore process. We are going to pause the `BackupConfiguration` crd that we had created to backup the `sample-elasticsearch` database. Then, Stash will stop taking any further backup for this database.

Let's pause the `sample-elasticsearch-backup` BackupConfiguration,

```console
$ kubectl patch backupconfiguration -n demo sample-elasticsearch-backup --type="merge" --patch='{"spec": {"paused": true}}'
backupconfiguration.stash.appscode.com/sample-elasticsearch-backup patched
```

Now, wait for a moment. Stash will pause the BackupConfiguration. Verify that the BackupConfiguration  has been paused,

```console
$  kubectl get backupconfiguration -n demo sample-elasticsearch-backup
NAME                          TASK                         SCHEDULE      PAUSED   AGE
sample-elasticsearch-backup   elasticsearch-backup-6.5.3   */5 * * * *   true     3m8s
```

Notice the `PAUSED` column. Value `true` for this field means that the BackupConfiguration has been paused.

**Deploy Restored Database:**

Now, we have to deploy the restored database similarly as we have deployed the original `sample-psotgres` database. However, this time there will be the following differences:

- We have to use the same secret that was used in the original database. We are going to specify it using `spec.databaseSecret` field.
- We have to specify `spec.init` section to tell KubeDB that we are going to use Stash to initialize this database from backup. KubeDB will keep the database phase to `Initializing` until Stash finishes its initialization.

Below is the YAML for `Elasticsearch` crd we are going deploy to initialize from backup,

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: restored-elasticsearch
  namespace: demo
spec:
  version: "6.5.3"
  storageType: Durable
  databaseSecret:
    secretName: sample-elasticsearch-auth # use same secret as original the database
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  init:
    stashRestoreSession:
      name: sample-elasticsearch-restore
  terminationPolicy: Delete
```

Here,

- `spec.init.stashRestoreSession.name` specifies the `RestoreSession` crd name that we are going to use to restore this database.

Let's create the above database,

```console
$ kubectl apply -f https://github.com/stashed/elasticsearch/raw/{{< param "info.subproject_version" >}}/docs/examples/restore/restored-elasticsearch.yaml
elasticsearch.kubedb.com/restored-elasticsearch created
```

If you check the database status, you will see it is stuck in `Initializing` state.

```console
$ kubectl get es -n demo restored-elasticsearch
NAME                     VERSION       STATUS         AGE
restored-elasticsearch   6.5.3         Initializing   38s
```

**Create RestoreSession:**

Now, we need to create a `RestoreSession` crd pointing to the AppBinding for this restored database.

Check AppBinding has been created for the `restored-elasticsearch` database using the following command,

```console
$ kubectl get appbindings -n demo restored-elasticsearch
NAME                     AGE
restored-elasticsearch   29s
```

>If you are not using KubeDB to deploy database, create the AppBinding manually.

Below is the YAML for the `RestoreSession` crd that we are going to create to restore backed up data into `restored-elasticsearch` database.

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: RestoreSession
metadata:
  name: sample-elasticsearch-restore
  namespace: demo
  labels:
    kubedb.com/kind: Elasticsearch # this label is mandatory if you are using KubeDB to deploy the database. Otherwise, Elasticsearch crd will be stuck in `Initializing` phase.
spec:
  task:
    name: elasticsearch-restore-6.5.3
  repository:
    name: gcs-repo
  target:
    ref:
      apiVersion: appcatalog.appscode.com/v1alpha1
      kind: AppBinding
      name: restored-elasticsearch
  interimVolumeTemplate:
    metadata:
      name: stash-tmp-storage
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: "standard"
      resources:
        requests:
          storage: 1Gi
  rules:
  - snapshots: [latest]
```

Here,

- `metadata.labels` specifies a `kubedb.com/kind: Elasticsearch` label that is used by KubeDB to watch this `RestoreSession`.
- `spec.task.name` specifies the name of the `Task` crd that specifies the Functions and their execution order to restore an Elasticsearch database.
- `spec.repository.name` specifies the `Repository` crd that holds the backend information where our backed up data has been stored.
- `spec.target.ref` refers to the AppBinding crd for the `restored-elasticsearch` database.
- `spec.interimVolumeTemplate` specifies a PVC template to store the restored data temporarily before inserting into the targeted Elasticsearch database.
- `spec.rules` specifies that we are restoring from the latest backup snapshot of the database.

> **Warning:** Label `kubedb.com/kind: Elasticsearch` is mandatory if you are using KubeDB to deploy the database. Otherwise, the database will be stuck in `Initializing` state.

Let's create the `RestoreSession` crd we have shown above,

```console
$ kubectl apply -f https://github.com/stashed/elasticsearch/raw/{{< param "info.subproject_version" >}}/docs/examples/restore/restoresession.yaml
restoresession.stash.appscode.com/sample-elasticsearch-restore created
```

Once, you have created the `RestoreSession` crd, Stash will create a job to restore. We can watch the `RestoreSession` phase to check if the restore process is succeeded or not.

Run the following command to watch `RestoreSession` phase,

```console
$ kubectl get restoresession -n demo sample-elasticsearch-restore -w
NAME                           REPOSITORY-NAME   PHASE       AGE
sample-elasticsearch-restore   gcs-repo          Running     5s
sample-elasticsearch-restore   gcs-repo          Succeeded   43s
```

So, we can see from the output of the above command that the restore process succeeded.

**Verify Restored Data:**

In this section, we are going to verify that the desired data has been restored successfully. We are going to connect to the database and check whether the table we had created in the original database is restored or not.

At first, check if the database has gone into `Running` state by the following command,

```console
$ kubectl get es -n demo restored-elasticsearch
NAME                     VERSION       STATUS    AGE
restored-elasticsearch   6.5.3         Running   2m16s
```

Now, find out the database pod by the following command,

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=restored-elasticsearch"
NAME                       READY   STATUS    RESTARTS   AGE
restored-elasticsearch-0   1/1     Running   0          48m
```

Now, exec into the database pod and list available tables,

```console
$ kubectl exec -it -n demo restored-elasticsearch-0 bash

~ curl -XGET --user "admin:mti4sbyy" "localhost:9200/test/snapshot/1?pretty"
{
  "_index" : "test",
  "_type" : "snapshot",
  "_id" : "1",
  "_version" : 1,
  "found" : true,
  "_source" : {
    "title" : "Snapshot",
    "text" : "Testing instand backup",
    "date" : "2018/02/13"
  }
}
```

So, from the above output, we can see the document `test` that we had created in the original database `sample-elasticsearch` is restored in the restored database `restored-elasticsearch`.

## Cleanup

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl delete -n demo restoresession sample-elasticsearch-restore
kubectl delete -n demo backupconfiguration sample-elasticsearch-backup
kubectl delete -n demo es sample-elasticsearch restored-elasticsearch
kubectl delete -n demo repository gcs-repo
```
