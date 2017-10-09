# Azure Service Broker

[Azure Service Broker](https://github.com/deis/azure-service-broker) is the
open source, [Open Service Broker](https://www.openservicebrokerapi.org/)
compatible API server that provisions managed services in the Microsoft
Azure public cloud.

This chart bootstraps Azure Service Broker in your Kubernetes cluster.

## Prerequisites

- [Kubernetes](https://kubernetes.io/) 1.7+ with RBAC enabled
- The
  [Kubernetes Service Catalog](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install-1.7.md)
  software has been installed
- An [Azure subscription](https://azure.microsoft.com/en-us/free/)
- The [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- A _service principal_ with the `Contributor` role on your Azure subscription

## Obtain Your Subscription ID

```console
$ export AZURE_SUBSCRIPTION_ID=$(az account show | grep '"id":' | awk '{print $2}' | awk '{gsub(/\"|,/,"")}1')
```

## Creating a Service Principal
  
```console
$ az ad sp create-for-rbac
```

The new service principal will be assigned, by default, to the `Contributor`
role. The output of the command above will be similar to the following:

```console
{
  "appId": "redacted",
  "displayName": "azure-cli-xxxxxx",
  "name": "http://azure-cli-xxxxxx",
  "password": "redacted",
  "tenant": "redacted"
}
```

For convenience in subsequent steps, we will export several of the fields above
as environment variables:

```console
$ export AZURE_TENANT_ID=<tenant>
$ export AZURE_CLIENT_ID=<appId>
$ export AZURE_CLIENT_SECRET=<password>
```

## Installing the Chart

Because this chart is not yet hosted in its own Helm repository at this time,
installation currently requires cloning this repository.

We assume your system is configured for Go development and that the environment
variable `GOPATH` is therefore defined. If this is not the case, start by
exporting this environment variable. Use your discretion in choosing a path,
but the path used below should generally be adequate:

```console
export GOPATH=~/Code/go
```

Then proceed with cloning this repository:

```console
$ mkdir -p $GOPATH/src/github.com/Azure
$ git clone git@github.com:deis/azure-service-broker.git \
    $GOPATH/src/github.com/Azure/azure-service-broker
```

To install the chart in the `asb` namespace with the release name `asb`:

```console
$ cd $GOPATH/src/github.com/Azure/azure-service-broker/contrib/k8s/charts/azure-service-broker
$ helm install . --name asb --namespace asb \
    --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
    --set azure.tenantId=$AZURE_TENANT_ID \
    --set azure.clientId=$AZURE_CLIENT_ID \
    --set azure.clientSecret=$AZURE_CLIENT_SECRET
```

This command deploys the Azure Service Broker on your Kubernetes cluster in the
default configuration. The [configuration](#configuration) section lists the
parameters that can optionally be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `asb` deployment:

```console
$ helm delete asb --purge
```

The command removes all the Kubernetes components associated with the chart and
deletes the release.

## Configuration

The following tables lists the configurable parameters of the Azure Service
Broker chart and their default values.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `image.repository` | Docker image location, _without_ the tag. | `"quay.io/deisci/azure-service-broker"` |
| `image.tag` | Tag / version of the Docker image. | `"canary"`; This references an image built from the very latest, possibly unreleased source code. This value is temporary and will change once ASB and its chart both stabilize. At that time, each revision of the chart will reference a specific ASB version using an _immutable_ tag, such as a semantic version. |
| `image.pullPolicy` | `"IfNotPresent"`, `"Always"`, or `"Never"`; When launching a pod, this option indicates when to pull the ASB Docker image. | `"Always"`; This policy complements the use of the mutable `canary` tag. This value is temporary and will change once ASB and its chart both stabalize and begin to reference images using _immutable_ tags, such as semantic versions. |
| `registerBroker` | Whether to register this broker with the Kubernetes Service Catalog. If true, the Kubernetes Service Catalog must already be installed on the cluster. Marking this option false is useful for scenarios wherein one wishes to host the broker in a separate cluster than the Service Catalog (or other client) that will access it. | `true` |
| `service.type` | Type of service; valid values are `"ClusterIP"`, `"LoadBalancer"`, and `"NodePort"`. `"ClusterIP"` is sufficient in the average case where ASB only receives traffic from within the cluster-- e.g. from Kubernetes Service Catalog. | `"ClusterIP"` |
| `service.nodePort.port` | _If and only if_ `service.type` is set to `"NodePort"`, `service.nodePort.port` indicates the port this service should bind to on each Kubernetes node. | `30080` |
| `azure.environment` | Indicates which Azure public cloud to use. Valid values are `"AzureCloud"` and `"AzureChinaCloud"`. | `"AzureCloud"` |
| `azure.subscriptionId` | Identifies the Azure subscription into which ASB will provision services. | none |
| `azure.tenantId` | Identifies the Azure Active Directory to which the _service principal_ used by ASB to access the Azure subscription belongs. | none |
| `azure.clientId` | Identifies the _service principal_ used by ASB to access the Azure subscription. | none |
| `azure.clientSecret` | Key/password for the _service principal_ used by ASB to access the Azure subscription. | none |
| `basicAuth.username` | Specifies the basic auth username that clients (e.g. the Kubernetes Service Catalog) must use when connecting to ASB. | `"username"`; __Do not use this default value in production!__ |
| `basicAuth.password` | Specifies the basic auth password that clients (e.g. the Kubernetes Service Catalog) must use when connecting to ASB. | `"password"`; __Do not use this default value in production!__ |
| `encryptionKey` | Specifies the key used by ASB for applying AES-256 encryption to sensitive (or potentially sensitive) data. | `"This is a key that is 256 bits!!"`; __Do not use this default value in production!__ |
| `modules.minStability` | Specifies the minimum level of stability an ASB module must meet for the services and plans it provides to be included in ASB's catalog of offerings. Valid values are `"ALPHA"`, `"BETA"`, and `"STABLE"`. | `"ALPHA"`; __Only use `"STABLE"` modules in production!__ |
| `redis.embedded` | ASB uses Redis for data persistence and as a message queue. This option indicates whether an on-cluster Redis deployment should be included when installing this chart. If set to `false`, connection details for a remote Redis cache must be provided. | `true`; __Do not use the embedded Redis cache in production!__ |
| `redis.host` | _If and only if_ `redis.embedded` is `false`, this option specifies the location of the remote Redis cache. | none |
| `redis.port` | Specifies the Redis port. | `6379` |
| `redis.redisPassword` | Specifies the Redis password. If `redis.embedded` is `true`, this option sets the password on the Redis cache itself _and_ provides it to ASB. If `redis.embedded` is `false`, this option only provides the password to ASB. In that case, the value of this option must be the correct password for the remote Redis cache. | `"password"`; __Do not use this default value in production!__ |

Specify a value for each option using the `--set <key>=<value>` switch on the
`helm install` command. That switch can be invoked multiple times to set
multiple options.

Alternatively, copy the charts default values to a file, edit the file to your
liking, and reference that file in your `helm install` command:

```console
$ helm inspect values azure-service-broker > values.yaml
$ vim my-values.yaml
$ helm install . --name asb --namespace asb --values my-values.yaml
```

