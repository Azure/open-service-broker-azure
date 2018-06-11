# Install and Operate Open Service Broker for Azure on an Azure Container Service managed cluster

Open Service Broker for Azure allows you to provision Azure services from your Kubernetes cluster. OSBA is integrated with Kubernetes using [Service Catalog](https://github.com/kubernetes-incubator/service-catalog). Both Service Catalog and OSBA have data persistence needs. When using Service Catalog and OSBA for development, it is sufficient to use the embedded storage options available with each application. For production use cases, however, we recommend more robust solutions. This guide provides details on setting up OSBA and associated software including Service Catalog, etcd and Redis for production scenarios. If you are new to OSBA, you may find the [AKS](quickstart-aks.md) or [Minikube](quickstart-minikube.md) Quickstart guides useful.

* [Prerequisites](#prerequisites)
  * [Existing Clusters](#existing-cluster)
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
* [Next Steps](#next-steps)

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

* As AKS is currently in preview, you will need to enable it in your subscription
* AKS currently does _not_ support RBAC, so we will need to explicity disable that when we install service catalog.

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
    az aks create --resource-group aks-group --name aks-cluster --generate-ssh-keys --kubernetes-version 1.9.6
    ```

    Note: Service Catalog may not work with Kubernetes versions less than 1.9.0. If you are attempting to use an older AKS cluster, you will need to upgrade. The earliest 1.9.x release available from AKS is 1.9.1, so you will need to upgrade to at least that version.

1. Configure kubectl to use the new cluster
    ```console
    az aks get-credentials --resource-group aks-group --name aks-cluster
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
helm init
```

### Install etcd

#### Create Storage Account

In order to configure etcd operator to use Azure Blob Storage for backup and recovery purposes, you will need to first create a storage account:

```console
az storage account create -n etcdoperator -g aks-group
```

Once the account has been created, retrieve the keys.

```console
az storage account keys list -n etcdoperator -g aks-group -o table
```

Next, you will need to create a container for backup storage.

```console
az storage container create --name etcd-backups --account-name etcdoperator --account-key <STORAGE_KEY>
```

#### Install etcd Operator

You will use Helm to install etcd Operator. We have included sample values.yaml file to setup etcd operator along with the backup and recovery operator. The default etcd Operator version installed by the Helm chart does not support Azure based backup and recovery. This values file uses a newer version of etcd Operator that provides Azure storage support.

```console
helm install --name etcd-operator stable/etcd-operator --values=contrib/k8s/etcd-operator/etcd-operator-values.yaml
```

This will create three deployments: etcd-operator, etcd-backup-operator, and restore-operator. This will also create three new Custom Resource Definitions: `etcdclusters.etcd.database.coreos.com`, `etcdbackups.etcd.database.coreos.com` and `etcdrestores.etcd.database.coreos.com`.

#### Create etcd Cluster

Once etcd Operator has been installed, you can create a cluster. For production scenarios, we recommend a three node cluster. For convienence purposes, we have included a sample cluster in contrib/k8s/etcd-operator.

```console
$ cat contrib/k8s/etcd-operator/svc-cat-cluster.yaml
apiVersion: "etcd.database.coreos.com/v1beta2"
kind: "EtcdCluster"
metadata:
  name: "svc-cat-etcd-cluster"
spec:
  size: 3
```

You can use kubectl to create a new etcd cluster using this file. 

```console
kubectl create -f contrib/k8s/etcd-operator/svc-cat-cluster.yaml
```

Once completed, you should see several etcd pods running:

```console
$ kubectl get pods

NAME                                                              READY     STATUS    RESTARTS   AGE
etcd-operator-etcd-operator-etcd-backup-operator-6b697d96c95fgv   1/1       Running   0          11m
etcd-operator-etcd-operator-etcd-operator-676764c476-n4ftv        1/1       Running   0          11m
etcd-operator-etcd-operator-etcd-restore-operator-7c8d6879rgkjv   1/1       Running   0          11m
svc-cat-etcd-cluster-5xtj4vlhx8                                   1/1       Running   0          1m
svc-cat-etcd-cluster-chfwgmjdph                                   1/1       Running   0          47s
svc-cat-etcd-cluster-jj87b2hmwg                                   1/1       Running   0          31s
```

You should also have an etcd service:

```console
$ kubectl get service
NAME                          TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)             AGE
etcd-restore-operator         ClusterIP   10.0.231.143   <none>        19999/TCP           14m
kubernetes                    ClusterIP   10.0.0.1       <none>        443/TCP             56m
svc-cat-etcd-cluster          ClusterIP   None           <none>        2379/TCP,2380/TCP   3m
svc-cat-etcd-cluster-client   ClusterIP   10.0.2.237     <none>        2379/TCP            3m
```

The `cluster-client` service is what you will use to configure Service Catalog.

### Install Service Catalog

Once you have created an etcd cluster, it is time to install Service Catalog. You will use Helm to install Service Catalog. There are a few values you will need to override for your Service Catalog installation:

* RBAC must be disabled (pending AKS support)
* Embedded etcd must be disabled
* You must point the installation at an external etcd.

You can provide these values via `--set` operations when using the Helm CLI, but we recommend creating a `values.yaml` file for your Service Catalog installation.

```yaml
rbacEnable: false
apiserver:
  storage:
    etcd:
      useEmbedded: false
      servers: http://svc-cat-etcd-cluster-client.default.svc.cluster.local:2379
```

The `servers` attribute should point to the service endpoint for your etcd cluster. If you ran the commands above without providing a Namespace, the file above should be sufficient. This file can be found in contrib/k8s/etcd-operator.

```console
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm install svc-cat/catalog --name catalog --namespace catalog \
   --values contrib/k8s/etcd-operator/svc-cat-values.yaml
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
