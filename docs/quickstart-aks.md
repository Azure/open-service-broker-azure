# Quickstart: Open Service Broker for Azure on an Azure Container Service managed cluster

This quickstart walks-through using the Open Service Broker for Azure (OSBA) to
deploy WordPress on an [Azure Container Service (AKS)](https://azure.microsoft.com/en-us/services/container-service/) managed cluster.

WordPress requires a backend MySQL database. Without OSBA, we would create a database
in the Azure portal, and then manually configure the connection information. Now
with OSBA our Kubernetes manifests can provision an Azure Database for MySQL on our behalf,
save the connection information in Kubernetes secrets, and then bind them to our WordPress instance.

* [Prerequisites](#prerequisites)
* [Cluster Setup](#cluster-setup)
  * [Configure your Azure account](#configure-your-azure-account)
  * [Create a resource group](#create-a-resource-group)
  * [Create a service principal](#create-a-service-principal)
  * [Create a Kubernetes cluster using AKS](#create-an-aks-cluster)
  * [Configure the cluster with Open Service Broker for Azure](#configure-the-cluster-with-open-service-broker-for-azure)
* [Deploy WordPress](#deploy-wordpress)
* [Next Steps](#next-steps)

---

## Prerequisites

* A [Microsoft Azure account](https://azure.microsoft.com/en-us/free/).
* Install the [Azure CLI](#install-the-azure-cli).
* Install the [Kubernetes CLI](#install-the-kubernetes-cli).
* Install the [Helm CLI](#install-the-helm-cli).

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
az acs kubernetes install-cli
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

### Create a Kubernetes cluster using AKS

Next we will create a managed Kubernetes cluster using AKS. AKS will create a managed Kubernetes cluster for you. Once the cluster is created, geting started with OSBA is very similar to doing so on [Minikube](quickstart-minikube.md), with a few exceptions: 

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

1. Create a Resource Group for AKS
    ```console
    az group create --name aks-group --location eastus
    ```

1. Create the AKS cluster!
    ```console
    az aks create --resource-group aks-group --name osba-quickstart-cluster --generate-ssh-keys
    ```

1. Configure kubectl to use the new cluster
    ```console
    az aks get-credentials --resource-group aks-group --name osba-quickstart-cluster
    ```

1. Verify your cluster is up and running
    ```console
    kubectl get nodes
    ```

### Configure the cluster with Open Service Broker for Azure

1. Before we can use Helm to install applications such as Service Catalog and
    WordPress on the cluster, we first need to prepare the cluster to work with Helm:
    ```console
    helm init
    ```
1. Deploy Service Catalog on the cluster:
    ```console
    helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
    helm install svc-cat/catalog --name catalog --namespace catalog --set rbacEnable=false
    ```

    Note: the AKS preview does not _currently_ support RBAC, so you must disable RBAC as shown above.

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

1. Check on the status of everything that we have installed by running the
    following command and checking that every pod is in the `Running` state.
    You may need to wait a few minutes, rerunning the command until all of the
    resources are ready.
    ```console
    $ kubectl get pods --namespace catalog
    NAME                                                     READY     STATUS    RESTARTS   AGE
    po/catalog-catalog-apiserver-5999465555-9hgwm            2/2       Running   4          9d
    po/catalog-catalog-controller-manager-554c758786-f8qvc   1/1       Running   11         9d

    $ kubectl get pods --namespace osba
    NAME                                           READY     STATUS    RESTARTS   AGE
    po/osba-azure-service-broker-8495bff484-7ggj6   1/1       Running   0          9d
    po/osba-redis-5b44fc9779-hgnck                  1/1       Running   0          9d
    ```

---

## Deploy WordPress

Now that we have a cluster with Open Service Broker for Azure, we can deploy
WordPress to Kubernetes and OSBA will handle provisioning an Azure Database for MySQL
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

Note:  While provisioning Wordpress and Azure Database for MySQL using Helm, all of the required resources are created in Kubernetes at the same time. As a result of these requests, Service Catalog will create a secret containing the the binding credentials for the database. This secret will not be created until after the Azure Database for MySQL is created, however. The wordpress container will depend on this secret being created before the container will fully start. Kubernetes and Service Catalog both employ a retry backoff, so you may need to wait several minutes for everything to be fully proviisoned.

## Login to WordPress

1. Run the following command to open WordPress in your browser:
    ```console
    export SERVICE_IP=$(kubectl get svc --namespace osba-quickstart osba-quickstart-wordpress -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    open http://$SERVICE_IP/admin
    ```

1. To retrieve the password, run this command:
    ```console
    kubectl get secret osba-quickstart-wordpress -n osba-quickstart -o jsonpath="{.data.wordpress-password}" | base64 --decode
    ```

1. Login using the username `user` and the password you just retrieved.

## Uninstall WordPress

Using Helm to uninstall the `osba-quickstart` release will delete all resources
associated with the release, including the Azure Database for MySQL instance.

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

At this point, the Azure Database of MySQL instance should have been fully deprovisioned.
In the unlikely event that anything has gone wrong, to ensure that you are not
billed for idle resources, you can delete the Azure resource group that
contained the database. In the case of the WordPress chart, Azure Database for MySQL was
provisioned in a resource group whose name matches the Kubernetes namespace into
which WordPress was deployed.

```console
az group delete --name osba-quickstart --yes --no-wait
```

To remove the service principal:

```console
az ad sp delete --id http://osba-quickstart`
```

To tear down the AKS cluster:

```console
az aks delete -resource-group aks-group --name osba-quickstart-cluster --no-wait
```

## Next Steps

Our AKS managed Kubernetes cluster communicated with Azure via OSBA, provisioned an Azure Database for
MySQL instance, and bound our WordPress installation to that new database.

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
