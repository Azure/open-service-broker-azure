# Install and Operate Open Service Broker for Azure on an Azure Container Service managed cluster

Open Service Broker for Azure allows you to provision Azure services from your Kubernetes cluster. OSBA is integrated with Kubernetes using [Service Catalog](https://github.com/kubernetes-incubator/service-catalog).

Both Service Catalog and OSBA have data persistence needs. When using Service Catalog and OSBA for development, it is sufficient to use the embedded storage options available with each application. For production use cases, however, we recommend more robust solutions. This guide provides details on setting up OSBA and associated software including Service Catalog, etcd and Redis for production scenarios.

If you are new to OSBA, you may find the [AKS](quickstart-aks.md) or [Minikube](quickstart-minikube.md) Quickstart guides useful.

* [Prerequisites](#prerequisites)
  * [Existing Clusters](#existing-clusters)
  * [Create an AKS Cluster](#new-cluster)
* [Cluster Configuration](#cluster-configuration)
  * [Install Helm](#install-helm)
  * [Install etcd](#install-etcd)
    * [Create Storage Account](#create-storage-account)
    * [Install etcd Operator](#install-etcd-operator)
    * [Create etcd Cluster](#create-etcd-cluster)
  * [Install Service Catalog](#install-service-catalog)
* [Create an Azure Redis Cache](#create-azure-redis-cache)
* [Create a service principal](#create-a-service-principal)
* [Configure the cluster with Open Service Broker for Azure](#configure-the-cluster-with-open-service-broker-for-azure)
* [Backup and Recovery](#backup-and-recovery)

## Prerequisites

* A [Microsoft Azure account](https://azure.microsoft.com/en-us/free/).
* The [Azure CLI](#install-the-azure-cli) installed.
* The [Kubernetes CLI](#install-the-kubernetes-cli) installed.
* The [Helm CLI](#install-the-helm-cli) installed.

### Install the Azure CLI

Install `az` by following the instructions for your operating system.
See the [full installation instructions](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) if yours isn't listed below.

#### MacOS

```console
brew install azure-cli
```

#### Windows

Download and run the [Azure CLI Installer (MSI)](https://aka.ms/InstallAzureCliWindows).

#### Ubuntu 64-bit

1. Add the azure-cli repo to your sources:
    ```console
    echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ wheezy main" | \
         sudo tee /etc/apt/sources.list.d/azure-cli.list
    ```
1. Run the following commands to install the Azure CLI and its dependencies:
    ```console
    sudo apt-key adv --keyserver packages.microsoft.com --recv-keys 52E16F86FEE04B979B07E28DB02C46DF417A0893
    sudo apt-get install apt-transport-https
    sudo apt-get update && sudo apt-get install azure-cli
    ```

### Install the Kubernetes CLI

Install `kubectl` by running the following command:

```console
az aks install-cli
```

### Install the Helm CLI

[Helm](https://github.com/kubernetes/helm) is a tool for installing pre-configured applications on Kubernetes.
Install `helm` by running the following command:

#### MacOS

```console
brew install kubernetes-helm
```

#### Windows

1. Download the latest [Helm release](https://storage.googleapis.com/kubernetes-helm/helm-v2.7.2-windows-amd64.tar.gz).
1. Decompress the tar file.
1. Copy **helm.exe** to a directory on your PATH.

#### Linux

```console
curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash
```

### Existing Clusters

Service Catalog currently requires Kubernetes version 1.9.0 or later, so you will therefore need an AKS cluster using Kubernetes 1.9.0 or later. You can check the version of your cluster with the `az aks list` or `az aks show -g <resource-group> -n <name>` CLI commands. If the version of your AKS managed cluster is less than 1.9.0, you will need to upgrade the cluster in order to install Service Catalog. You can use the `az aks upgrade` upgrade command to upgrade your cluster.

You'll also want to ensure that you have at least three nodes in the cluster. You can use the `az aks scale` command to scale the number of worker nodes if needed.

### New Cluster

The following steps will walk you through creation of a new AKS cluster. For more information on cluster creation, please consult [Quickstart: Deploy an Azure Kubernetes Service Cluster](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough) and [Create an Azure Kubernetes Service cluster](https://docs.microsoft.com/en-us/azure/aks/create-cluster).

#### Create a Resource Group for AKS

When you create an AKS cluster, you must provide a resource group. Create one with the az cli using the following command.

```console
az group create --name aks-group --location eastus
```

#### Create a Kubernetes cluster using AKS

Azure requires that you enable the AKS resource provider in your subscription. If you have not done so, the following command will enable it. 

1. Enable AKS in your subscription, use the following command with the az cli:
    ```console
    az provider register -n Microsoft.ContainerService
    ```

You should also ensure that the `Microsoft.Compute` and `Microsoft.Network` providers are registered in your subscription. If you need to enable them:
    ```console
    az provider register -n Microsoft.Compute
    az provider register -n Microsoft.Network
    ```

1. Create the AKS cluster!
    ```console
    az aks create --resource-group aks-group --name aks-cluster --generate-ssh-keys --kubernetes-version 1.9.6 --enable-rbac
    ```

    Note: Service Catalog may not work with Kubernetes versions less than 1.9.0. If you are attempting to use an older AKS cluster, you will need to upgrade. The earliest 1.9.x release available from AKS is 1.9.1, so you will need to upgrade to at least that version.

1. Configure kubectl to use the admin user in the new cluster
    ```console
    az aks get-credentials --resource-group aks-group --name aks-cluster --admin
    ```

1. Verify your cluster is up and running
    ```console
    kubectl get nodes
    ```

---

## Cluster Configuration

Before you can install OSBA onto your cluster, you will first need to install Service Catalog. Service Catalog uses etcd for data persistence and it is *strongly* recommended that you create an etcd cluster with backup and recovery capabilities. We recommend using [etcd operator](https://github.com/coreos/etcd-operator) configured with backup and recovery via Azure Blob Storage for this purpose. The following section provides guidance on how to configure etcd operator on your cluster and then how to install Service Catalog.

### Install Helm

We will use Helm to install etcd Operator, Service Catalog and OSBA. In order to install Helm in your cluster, use the `helm` CLI:

```console
kubectl create -f https://raw.githubusercontent.com/Azure/helm-charts/master/docs/prerequisities/helm-rbac-config.yaml
helm init --service-account tiller
```

### Install etcd

#### Create Storage Account

In order to configure etcd operator to use Azure Blob Storage for backup and recovery purposes, you will need to first create a storage account:

First, save your desired account name to an environment variable. It will be reused a few times.

 **Bash**
 ```console
export AZURE_STORAGE_ACCOUNT=<STORAGE ACCOUNT>
 ```

**PowerShell**
```console
$env:AZURE_STORAGE_ACCOUNT = "<STORAGE ACCOUNT>"
```

Now use the az cli to create a new storage account.

**Bash**
```console
az storage account create -n $AZURE_STORAGE_ACCOUNT -g aks-group
```

**PowerShell**
```console
az storage account create -n $AZURE_STORAGE_ACCOUNT -g aks-group
```

Once the account has been created, retrieve the keys.

**Bash**
```console
az storage account keys list -n $AZURE_STORAGE_ACCOUNT -g aks-group -o table
```

**Powershell**
```console
az storage account keys list -n $env:AZURE_STORAGE_ACCOUNT -g aks-group -o table
```

Save the storage account key to an environment variable for later use.
 **Bash**
 ```console
export AZURE_STORAGE_KEY=<STORAGE_KEY>
 ```

**PowerShell**
```console
$env:AZURE_STORAGE_KEY = "<STORAGE_KEY>"
```

Next, you will need to create a container for backup storage.

**Bash**
```console
az storage container create --name etcd-backups --account-name $AZURE_STORAGE_ACCOUNT --account-key $AZURE_STORAGE_KEY
```

**PowerShell**
```console
az storage container create --name etcd-backups --account-name $env:AZURE_STORAGE_ACCOUNT --account-key $env:AZURE_STORAGE_KEY
```

#### Install etcd Operator

Service Catalog requires etcd to persist state for everything it manages. We recommend using etcd Operator to install and manage an etcd cluster for Service Catalog. We have included sample values.yaml file to setup etcd operator with a version that supports Azure backup. 

```console
helm install --name etcd-operator stable/etcd-operator --values=contrib/k8s/etcd-operator/etcd-operator-values.yaml
```

This will create three deployments: etcd-operator, etcd-backup-operator, and restore-operator. This will also create three new Custom Resource Definitions: `etcdclusters.etcd.database.coreos.com`, `etcdbackups.etcd.database.coreos.com` and `etcdrestores.etcd.database.coreos.com`. These can be used to create an etcd cluster, along with backup and restore operations. 

#### Create etcd Cluster

Once etcd Operator has been installed, you can create a cluster. For production scenarios, we recommend a three node cluster. The etcd operator currently creates ephemeral etcd members. The consequence of this is if an etcd member crashes, it's data will be lost. A three-node cluster will therefore provide a higher level of operational stability. If a single etcd pod crashes, it can be rescheduled to rejoin the cluster. For more severe failures, a recovery operation must be initiated from backup. We have provided a Helm chart that will create a three-node cluster and configure a Kubernetes CronJob to enable automatic backups of the cluster. 

**Bash**
```console
helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
helm install azure/svc-cat-etcd --set azure.storage.account=$AZURE_STORAGE_ACCOUNT --set azure.storage.key=$AZURE_STORAGE_KEY
```

**PowerShell**
```console
helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
helm install azure/svc-cat-etcd --set azure.storage.account=$env:AZURE_STORAGE_ACCOUNT --set azure.storage.key=$env:AZURE_STORAGE_KEY
```

Once completed, you should see several etcd pods running:

```console
$ kubectl get pods

NAME                                                              READY     STATUS    RESTARTS   AGE
etcd-operator-etcd-operator-etcd-backup-operator-6b697d96c95fgv   1/1       Running   0          11m
etcd-operator-etcd-operator-etcd-operator-676764c476-n4ftv        1/1       Running   0          11m
etcd-operator-etcd-operator-etcd-restore-operator-7c8d6879rgkjv   1/1       Running   0          11m
svc-cat-etcd-5xtj4vlhx8                                   1/1       Running   0          1m
svc-cat-etcd-chfwgmjdph                                   1/1       Running   0          47s
svc-cat-etcd-jj87b2hmwg                                   1/1       Running   0          31s
```

You should also have an etcd service:

```console
$ kubectl get service

NAME                          TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)             AGE
etcd-restore-operator         ClusterIP   10.0.231.143   <none>        19999/TCP           14m
kubernetes                    ClusterIP   10.0.0.1       <none>        443/TCP             56m
svc-cat-etcd          ClusterIP   None           <none>        2379/TCP,2380/TCP   3m
svc-cat-etcd-client   ClusterIP   10.0.2.237     <none>        2379/TCP            3m
```

The `svc-cat-etcd-client` service is what you will use to configure Service Catalog.

### Install Service Catalog

Once you have created an etcd cluster, it is time to install Service Catalog. You will use Helm to install Service Catalog. There are a few values you will need to override for your Service Catalog installation:

* Embedded etcd must be disabled
* You must point the installation at an etcd cluster.

```console
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm install svc-cat/catalog --name catalog --namespace catalog \
   --set apiserver.storage.etcd.useEmbedded=false \
   --set apiserver.storage.etcd.servers=http://svc-cat-etcd-client.default.svc.cluster.local:2379
```

## Create Azure Redis Cache

Open Service Broker for Azure uses Redis as a backing store for its state. We recommend using a managed Redis service, such as Azure Redis Cache. By default, Azure Redis Cache only keeps data in memory. For best results, you will want to use the Premium tier in order to configure backups of the Redis data. Please see [How to configure data persistence for a Premium Azure Redis Cache](https://docs.microsoft.com/en-us/azure/redis-cache/cache-how-to-premium-persistence) for instructions on how to create an Azure Redis Cache. Once created, you can obtain the hostname and keys from the Portal or via the CLI.
Save the access key and host to an environment variable for later use:

**Bash**

```console
export REDIS_HOSTNAME=<Redis Host>
export REDIS_PASSWORD=<Redis PrimaryKey>
```

**PowerShell**

```console
$env:REDIS_HOSTNAME = "<Redis Host>"
$env:REDIS_PASSWORD = "<Redis PrimaryKey>"
```

## Create a service principal

This creates an identity for Open Service Broker for Azure to use when provisioning
resources on your account on behalf of Kubernetes.

1. Create a service principal with RBAC enabled:
    ```console
    az ad sp create-for-rbac --name osba -o table
    ```
1. Save the values from the command output in environment variables:

    **Bash**
    ```console
    export AZURE_TENANT_ID=<Tenant>
    export AZURE_CLIENT_ID=<AppId>
    export AZURE_CLIENT_SECRET=<Password>
    ```

    **PowerShell**
    ```console
    $env:AZURE_TENANT_ID = "<Tenant>"
    $env:AZURE_CLIENT_ID = "<AppId>"
    $env:AZURE_CLIENT_SECRET = "<Password>"
    ```

## Configure the cluster with Open Service Broker for Azure

You can now deploy Open Service Broker for Azure on the cluster. Using Helm:

**Bash**

```console
helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
helm install azure/open-service-broker-azure --name osba --namespace osba \
  --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
  --set azure.tenantId=$AZURE_TENANT_ID \
  --set azure.clientId=$AZURE_CLIENT_ID \
  --set azure.clientSecret=$AZURE_CLIENT_SECRET \
  --set redis.embedded=false \
  --set redis.host=$REDIS_HOSTNAME \
  --set redis.redisPassword=$REDIS_PASSWORD
```

**PowerShell**

```console
helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
helm install azure/open-service-broker-azure --name osba --namespace osba `
  --set azure.subscriptionId=$env:AZURE_SUBSCRIPTION_ID `
  --set azure.tenantId=$env:AZURE_TENANT_ID `
  --set azure.clientId=$env:AZURE_CLIENT_ID `
  --set azure.clientSecret=$env:AZURE_CLIENT_SECRET `
  --set redis.embedded=false `
  --set redis.host=$env:REDIS_HOSTNAME `
  --set redis.redisPassword=$env:REDIS_PASSWORD
```

## Backup and Recovery

### etcd

The etcd cluster created above should have automatic backup enabled to Azure Blob Storage. The `etcd-backup` CronJob created by the Helm chart will periodically create new `etcdbackup` custom resources. These will not be removed if you remove the Helm chart.

To initiate a restore of the etcd cluster, you need to create an instance of the `EtcdRestore` Custom Resource. A template for this can be found in contrib/k8s/etcd-operator/restore-operation.yaml

```yaml
apiVersion: etcd.database.coreos.com/v1beta2
kind: EtcdRestore
metadata:
  name: svc-cat-etcd-restore
spec:
  etcdCluster:
    name: svc-cat-etcd
  backupStorageType: ABS
  abs:
    path: <abs continer>/<BACKUP-FILE> 
    absSecret: <abs-credential-secret-name> 
```

You will need to replace the `spec.abs.path` value with the backup you'd like to restore from. You will need the storage container as well as the file name. For example, if your absSecret was called `dandy-clownfish-svc-cat-etcd` and you used etcd-backups as the storage container as directed above and the file name was etcd.backup.2018-06-12_19:31:05, your restore yaml would look like:

```yaml
apiVersion: etcd.database.coreos.com/v1beta2
kind: EtcdRestore
metadata:
  name: svc-cat-etcd-restore
spec:
  etcdCluster:
    name: svc-cat-etcd
  backupStorageType: ABS
  abs:
    path: etcd-backups/etcd.backup.2018-06-12_19:31:05
    absSecret: dandy-clownfish-svc-cat-etcd
```

If this file was saved in the current directory as restore-request.yaml, you would initiate the restore by using kubectl:

```console
kubectl create -f restore-request.yaml
```

This will result in the current etcd pods being terminated and restarted with the specified backup file.