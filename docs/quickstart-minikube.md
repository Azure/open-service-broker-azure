# Quickstart: Open Service Broker for Azure on a Minikube cluster

This quickstart walks through using the Open Service Broker for Azure (OSBA) to
deploy WordPress on a local Minikube cluster.

WordPress requires a back-end MySQL database. Without OSBA, we would create a database
in the Azure portal, and then manually configure the connection information. Now
with OSBA our Kubernetes manifests can provision an Azure MySQL database on our behalf,
save the connection information in Kubernetes secrets, and then bind them to our WordPress instance.

* [Prerequisites](#prerequisites)
* [Cluster Setup](#cluster-setup)
  * [Configure your Azure account](#configure-your-azure-account)
  * [Create a resource group](#create-a-resource-group)
  * [Create a service principal](#create-a-service-principal)
  * [Create a Kubernetes cluster using Minikube](#create-a-kubernetes-cluster-using-minikube)
  * [Configure the cluster with Open Service Broker for Azure](#configure-the-cluster-with-open-service-broker-for-azure)
* [Deploy WordPress](#deploy-wordpress)
* [Next Steps](#next-steps)

---

## Prerequisites

* A [Microsoft Azure account](https://azure.microsoft.com/en-us/free/).
* Install [Minikube](#install-minikube).
* Install the [Azure CLI](#install-the-azure-cli).
* Install the [Kubernetes CLI](#install-the-kubernetes-cli).
* Install the [Helm CLI](#install-the-helm-cli).

### Install Minikube

[Minikube](https://github.com/kubernetes/minikube) is a tool that makes it easy to run Kubernetes locally. Minikube runs a single-node Kubernetes cluster inside a VM on your computer. For this quickstart guide, you'll want to install Minikube v0.25.

#### MacOS

```console
brew cask install minikube
```

#### Windows

1. Download the [minikube-windows-amd64.exe](https://storage.googleapis.com/minikube/releases/latest/minikube-windows-amd64.exe) file.
1. Rename it to **minikube.exe**.
1. Add it to a directory on your PATH.

#### Linux

```console
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
```

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

---

## Cluster Setup

Now that we have all the tools, we need a Kubernetes cluster with Open Service Broker for Azure configured.

### Configure your Azure account

First let's identify your Azure subscription and save it for use later on in the quickstart.

1. Run `az login` and follow the instructions in the command output to authorize `az` to use your account
1. List your Azure subscriptions:
    ```console
    az account list -o table
    ```
1. Copy your subscription ID and save it in an environment variable:

    **Bash**
    ```console
    export AZURE_SUBSCRIPTION_ID="<SubscriptionId>"
    ```

    **PowerShell**
    ```console
    $env:AZURE_SUBSCRIPTION_ID = "<SubscriptionId>"
    ```

### Create a Resource Group

Create a resource group to contain the resources you'll be creating with the quickstart.

```console
az group create --name osba-quickstart --location eastus
```

### Create a service principal

This creates an identity for Open Service Broker for Azure to use when provisioning
resources on your account on behalf of Kubernetes.

1. Create a service principal with RBAC enabled for the quickstart:
    ```console
    az ad sp create-for-rbac --name osba-quickstart -o table
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

### Create a Kubernetes cluster using Minikube

Next we will create a local cluster using Minikube. You can also [try OSBA on the Azure Container Service (AKS)](quickstart-aks.md).

1. Create a Minikube Cluster:
    ```console
    minikube start --bootstrapper=kubeadm
    ```

### Configure the cluster with Open Service Broker for Azure

1. Before we can use Helm to install applications such as Service Catalog and
    WordPress on the cluster, we first need to prepare the cluster to work with Helm:
    ```console
    kubectl create -f https://raw.githubusercontent.com/Azure/helm-charts/master/docs/prerequisities/helm-rbac-config.yaml
    helm init --service-account tiller
    ```
1. Deploy Service Catalog on the cluster:
    ```console
    helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
    helm install svc-cat/catalog --name catalog --namespace catalog
    ```
1. Check the status of Service Catalog:
    Run the following command and checking that every pod is in the `Running` state.
    You may need to wait a few minutes, rerunning the command until all of the
    resources are ready. 
    ```console
    $ kubectl get pods --namespace catalog
    NAME                                                     READY     STATUS    RESTARTS   AGE
    po/catalog-catalog-apiserver-5999465555-9hgwm            2/2       Running   4          9d
    po/catalog-catalog-controller-manager-554c758786-f8qvc   1/1       Running   11         9d
    ```
1. Deploy Open Service Broker for Azure on the cluster:

    **Bash**
    ```console
    helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
    helm install azure/open-service-broker-azure --name osba --namespace osba \
      --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
      --set azure.tenantId=$AZURE_TENANT_ID \
      --set azure.clientId=$AZURE_CLIENT_ID \
      --set azure.clientSecret=$AZURE_CLIENT_SECRET
    ```

    **PowerShell**
    ```console
    helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
    helm install azure/open-service-broker-azure --name osba --namespace osba `
      --set azure.subscriptionId=$env:AZURE_SUBSCRIPTION_ID `
      --set azure.tenantId=$env:AZURE_TENANT_ID `
      --set azure.clientId=$env:AZURE_CLIENT_ID `
      --set azure.clientSecret=$env:AZURE_CLIENT_SECRET
    ```

1. Check on the status of Open Service Broker for Azure by running the
    following command and checking that every pod is in the `Running` state.
    You may need to wait a few minutes, rerunning the command until all of the
    resources are ready.
    ```console
    $ kubectl get pods --namespace osba
    NAME                                           READY     STATUS    RESTARTS   AGE
    po/osba-azure-service-broker-8495bff484-7ggj6   1/1       Running   0          9d
    po/osba-redis-5b44fc9779-hgnck                  1/1       Running   0          9d
    ```

---

## Deploy WordPress

Now that we have a cluster with Open Service Broker for Azure, we can deploy
WordPress to Kubernetes and OSBA will handle provisioning an Azure MySQL database
and binding it to our WordPress installation.

```console
helm install azure/wordpress --name osba-quickstart --namespace osba-quickstart
```

Use the following command to tell when WordPress is ready:

```console
$ kubectl get deploy osba-quickstart-wordpress -n osba-quickstart -w

NAME                        DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
osba-quickstart-wordpress   1         1         1            0           1m
...
osba-quickstart-wordpress   1         1         1            1           2m
```

## Login to WordPress

1. Run the following command to open WordPress in your browser:
    ```console
    open http://$(minikube ip):$(kubectl get service osba-quickstart-wordpress -n osba-quickstart -o jsonpath={.spec.ports[?\(@.name==\"http\"\)].nodePort})/admin
    ```

    **Note**: We are using the `minikube ip` to get the WordPress URL, instead of
    the command from the WordPress deployment output because with Minikube the
    WordPress service won't have a public IP address assigned.

1. To retrieve the password, run this command:
    ```console
    kubectl get secret osba-quickstart-wordpress -n osba-quickstart -o jsonpath="{.data.wordpress-password}" | base64 --decode
    ```

1. Login using the username `user` and the password you just retrieved.

## Uninstall WordPress

Using Helm to uninstall the `osba-quickstart` release will delete all resources
associated with the release, including the Azure MySQL database.

```console
helm delete osba-quickstart --purge
```

Since deprovisioning occurs asynchronously, the corresponding `serviceinstance`
resource will not be fully deleted until that process is complete. When the
following command returns no resources, deprovisioning is complete:

```console
$ kubectl get serviceinstances -n osba-quickstart
No resources found.
```

## Optional: Further Cleanup

At this point, the Azure MySQL database should have been fully deprovisioned.
In the unlikely event that anything has gone wrong, to ensure that you are not
billed for idle resources, you can delete the Azure resource group that
contained the database. In the case of the WordPress chart, Azure MySQL was
provisioned in a resource group whose name matches the Kubernetes namespace into
which WordPress was deployed.

```console
az group delete --name osba-quickstart --yes --no-wait
```

To remove the service principal:

```console
az ad sp delete --id http://osba-quickstart
```

To tear down minikube:

```console
minikube delete
```

## Next Steps

Minikube may seem like an odd choice for an Azure quickstart, but it demonstrates
that Open Service Broker for Azure isn't limited to clusters running on Azure!
Our local Kubernetes cluster communicated with Azure via OSBA, provisioned an Azure
MySQL database, and bound our local WordPress installation to that new database.

With OSBA _any_ cluster can rely on Azure to provide all those pesky "as a service"
goodies that make life easier.

Now that you have a cluster with OSBA, adding more applications is quick. Try out another to see for yourself:

* [Concourse CI](https://github.com/Azure/helm-charts/blob/master/concourse)
* [phpBB](https://github.com/Azure/helm-charts/blob/master/phpbb)

All of our OSBA-enabled helm charts are available in the [Azure/helm-charts](https://github.com/Azure/helm-charts)
repository.

## Contributing

Do you have an application in mind that you'd like to use with OSBA? We'd love to
have it! Learn how to [contribute a new chart](https://github.com/Azure/helm-charts#creating-a-new-chart)
to our helm repository.
