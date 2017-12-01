# Quickstart: Open Service Broker for Azure on a Minikube cluster
This quickstart walks-through using the Open Service Broker for Azure (OSBA) to
deploy WordPress on a local Minikube cluster.

WordPress requires a backend MySQL database. Without OSBA, we would create a database
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

# Prerequisites
* A [Microsoft Azure account](https://azure.microsoft.com/en-us/free/).
* Install [Minikube](#install-minikube).
* Install the [Azure CLI](#install-the-azure-cli).
* Install the [Kubernetes CLI](#install-the-kubernetes-cli).
* Install the [Helm CLI](#install-the-helm-cli).

## Install Minikube
[Minikube](https://github.com/kubernetes/minikube) is a tool that makes it easy to run Kubernetes locally. Minikube runs a single-node Kubernetes cluster inside a VM on your computer.

**MacOS**
```
brew cask install minikube
```

**Windows**
1. Download the [minikube-windows-amd64.exe](https://storage.googleapis.com/minikube/releases/latest/minikube-windows-amd64.exe) file.
1. Rename it to **minikube.exe**.
1. Add it to a directory on your PATH.

**Linux**
```
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
```

## Install the Azure CLI
Install `az` by following the instructions for your operating system.
See the [full installation instructions](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) if yours isn't listed below.

**MacOS**

```
brew install azure-cli
```

**Windows**

To install the CLI on Windows and use it in the Windows command-line, download and run the [Azure CLI Installer (MSI)](https://aka.ms/InstallAzureCliWindows).

**Ubuntu 64-bit**

1. Add the azure-cli repo to your sources:
    ```
    echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ wheezy main" | \
         sudo tee /etc/apt/sources.list.d/azure-cli.list
    ```
1. Run the following commands to install the Azure CLI and its dependencies:
    ```
    sudo apt-key adv --keyserver packages.microsoft.com --recv-keys 52E16F86FEE04B979B07E28DB02C46DF417A0893
    sudo apt-get install apt-transport-https
    sudo apt-get update && sudo apt-get install azure-cli
    ```

## Install the Kubernetes CLI
Install `kubectl` by running the following command:

```
az acs kubernetes install-cli
```

## Install the Helm CLI
[Helm](https://github.com/kubernetes/helm) is a tool for installing pre-configured applications on Kubernetes.
Install `helm` by running the following command:

**MacOS**
```
brew install kubernetes-helm
```

**Windows**
1. Download the latest [Helm release](https://storage.googleapis.com/kubernetes-helm/helm-v2.7.2-windows-amd64.tar.gz).
1. Decompress the tar file.
1. Copy **helm.exe** to a directory on your PATH.

**Linux**
```
curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash
```

---

# Cluster Setup
Now that we have all the tools, we need a Kubernetes cluster with Open Service Broker for Azure configured.

## Configure your Azure account
First let's identify your Azure subscription and save it for use later on in the quickstart.

1. Run `az login` and follow the instructions in the command output to authorize `az` to use your account
1. List your Azure subscriptions:
    ```
    az account list -o table
    ```
1. Copy your subscription ID and save it in an environment variable:

    **Bash**
    ```
    export AZURE_SUBSCRIPTION_ID="<SubscriptionId>"
    ```

    **PowerShell**
    ```
    $env:AZURE_SUBSCRIPTION_ID = "<SubscriptionId>"
    ```

## Create a resource group
We are using a resource group to isolate all the resources created in this quickstart
for easy cleanup later.

1. List the available Azure regions and select a region, for example `centralus`:
    ```
    az account list-locations -o table
    ```
1. Create a resource group for the quickstart:
    ```
    az group create --name osba-quickstart --location <RegionName>
    ```

## Create a service principal
This creates an identity for Open Service Broker for Azure to use when provisioning
resources on your account on behalf of Kubernetes.

1. Create a service principal with RBAC enabled for the quickstart:
    ```
    az ad sp create-for-rbac --name osba-quickstart -o table
    ```
1. Save the values from the command output in environment variables:

    **Bash**
    ```
    export AZURE_TENANT_ID=<DisplayName>
    export AZURE_CLIENT_ID=<AppId>
    export AZURE_CLIENT_SECRET=<Password>
    ```

    **PowerShell**
    ```
    $env:AZURE_TENANT_ID = "<DisplayName>"
    $env:AZURE_CLIENT_ID = "<AppId>"
    $env:AZURE_CLIENT_SECRET = "<Password>"
    ```

## Create a Kubernetes cluster using Minikube
Next we will create a local cluster using Minikube. _Support for AKS is coming soon!_

1. Create an RBAC enabled cluster:
    ```
    minikube start --extra-config=apiserver.Authorization.Mode=RBAC
    ```
1. Grant the `cluster-admin` role to the default system account:
    ```
    kubectl create clusterrolebinding cluster-admin:kube-system \
       --clusterrole=cluster-admin \
       --serviceaccount=kube-system:default
    ```

## Configure the cluster with Open Service Broker for Azure

1. Before we can use Helm to install applications such as Service Catalog and
    WordPress on the cluster, we first need to prepare the cluster to work with Helm:
    ```
    kubectl create -f https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/helm-rbac-config.yaml
    helm init --service-account tiller
    ```
1. Deploy Service Catalog on the cluster:
    ```
    helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
    helm install svc-cat/catalog --name catalog --namespace catalog
    ```
1. Deploy Open Service Broker for Azure on the cluster:
    ```
    helm install azure/azure-service-broker --name osba --namespace osba \
      --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
      --set azure.tenantId=$AZURE_TENANT_ID \
      --set azure.clientId=$AZURE_CLIENT_ID \
      --set azure.clientSecret=$AZURE_CLIENT_SECRET
    ```
1. Check on the status of everything that we have installed by running the
    following command and checking that everything is in the `Running` state.
    You may need to wait a few minutes, rerunning the command until all of the
    resources are ready.
    ```console
    $ kubectl get all --namespace catalog
    NAME                                        DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    deploy/catalog-catalog-apiserver            1         1         1            1           9d
    deploy/catalog-catalog-controller-manager   1         1         1            1           9d

    NAME                                               DESIRED   CURRENT   READY     AGE
    rs/catalog-catalog-apiserver-5999465555            1         1         1         9d
    rs/catalog-catalog-controller-manager-554c758786   1         1         1         9d

    NAME                                        DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    deploy/catalog-catalog-apiserver            1         1         1            1           9d
    deploy/catalog-catalog-controller-manager   1         1         1            1           9d

    NAME                                               DESIRED   CURRENT   READY     AGE
    rs/catalog-catalog-apiserver-5999465555            1         1         1         9d
    rs/catalog-catalog-controller-manager-554c758786   1         1         1         9d

    NAME                                                     READY     STATUS    RESTARTS   AGE
    po/catalog-catalog-apiserver-5999465555-9hgwm            2/2       Running   4          9d
    po/catalog-catalog-controller-manager-554c758786-f8qvc   1/1       Running   11         9d

    NAME                            TYPE       CLUSTER-IP   EXTERNAL-IP   PORT(S)         AGE
    svc/catalog-catalog-apiserver   NodePort   10.0.0.117   <none>        443:30443/TCP   9d

    $ kubectl get all --namespace osba
    NAME                              DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    deploy/osba-azure-service-broker   1         1         1            1           9d
    deploy/osba-redis                  1         1         1            1           9d

    NAME                                     DESIRED   CURRENT   READY     AGE
    rs/osba-azure-service-broker-8495bff484   1         1         1         9d
    rs/osba-redis-5b44fc9779                  1         1         1         9d

    NAME                              DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    deploy/osba-azure-service-broker   1         1         1            1           9d
    deploy/osba-redis                  1         1         1            1           9d

    NAME                                     DESIRED   CURRENT   READY     AGE
    rs/osba-azure-service-broker-8495bff484   1         1         1         9d
    rs/osba-redis-5b44fc9779                  1         1         1         9d

    NAME                                           READY     STATUS    RESTARTS   AGE
    po/osba-azure-service-broker-8495bff484-7ggj6   1/1       Running   0          9d
    po/osba-redis-5b44fc9779-hgnck                  1/1       Running   0          9d

    NAME                           TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
    svc/osba-azure-service-broker   ClusterIP   10.0.0.8     <none>        80/TCP     9d
    svc/osba-redis                  ClusterIP   10.0.0.28    <none>        6379/TCP   9d
    ```

---

# Deploy WordPress
Now that we have a cluster with Open Service Broker for Azure, we can deploy
WordPress to Kubernetes and OSBA will handle provisioning an Azure MySQL database
and binding it to our WordPress installation.

```
helm install azure/wordpress --name quickstart
```

Use the following command to tell when WordPress is ready:

```console
$ kubectl get deploy/quickstart-wordpress -w

NAME                DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
quickstart-wordpress   1         1         1            0           1m
quickstart-wordpress   1         1         1            1           2m
```

## Login to WordPress

1. Run the following command to open WordPress in your browser:
    ```
    open http://$(minikube ip):31215/admin
    ```

    **Note**: We are using the `minikube ip` to get the WordPress URL, instead of
    the command from the WordPress deployment output because we are using Minikube
    and do not have a public IP address.
1. Login using the following credentials:
    ```
    echo Username: user
    echo Password: $(kubectl get secret quickstart-wordpress -o jsonpath="{.data.wordpress-password}" | base64 --decode)
    ```

## Optional: Cleanup
Here's how to remove resources created by this quickstart:

1. `az group delete --name osba-quickstart`
1. `az ad sp delete --name osba-quickstart`
1. `minikubte delete`

# Next Steps
Minikube may seem like an odd choice for an Azure quickstart, but it demonstrates
that Open Service Broker for Azure isn't limited to clusters running on Azure!
Our local Kubernetes cluster communicated via OSBA with Azure, provisioned a cloud
database, and bound our local WordPress installation to that new database.

With OSBA _any_ cluster can rely on Azure to provide all those pesky "as a service"
goodies that make life easier.

Now that you have a cluster with OSBA, adding more services is quick. Try out another
service to see for yourself:

* [Concourse CI](https://github.com/Azure/helm-charts/blob/master/concourse)
* [pbpBB](https://github.com/Azure/helm-charts/blob/master/phpbb)

All of our OSBA-enabled helm charts are available in the [Azure/helm-charts](https://github.com/Azure/helm-charts)
repository.

## Contributing
Do you have an application in mind that you'd like to use with OSBA? We'd love to
have it! Learn how to [contribute a new chart](https://github.com/Azure/helm-charts#creating-a-new-chart)
to our helm repository.
