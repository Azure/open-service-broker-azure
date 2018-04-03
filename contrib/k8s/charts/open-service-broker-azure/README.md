# Open Service Broker for Azure

[Open Service Broker for Azure](https://github.com/Azure/open-service-broker-azure) is the
open source, [Open Service Broker](https://www.openservicebrokerapi.org/)
compatible API server that provisions managed services in the Microsoft
Azure public cloud.

This chart bootstraps Open Service Broker for Azure in your Kubernetes cluster.

## Prerequisites

- [Kubernetes](https://kubernetes.io/) 1.7+ with RBAC enabled
- The
  [Kubernetes Service Catalog](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md)
  software has been installed
- An [Azure subscription](https://azure.microsoft.com/en-us/free/)
- The Azure CLI: You can
[install it locally](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)
or use it in the
[Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/cloud-shell/overview?view=azure-cli-latest)
- A _service principal_ with the `Contributor` role on your Azure subscription

## Obtain Your Subscription ID

```console
$ export AZURE_SUBSCRIPTION_ID=$(az account show --query id | sed s/\"//g)
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

Installation of this chart is simple. First, ensure that you've [added the
`azure` repository](../README.md#installing-charts). Then, install from the
`azure` repo:

```console
$ helm install azure/open-service-broker-azure --name osba --namespace osba \
  --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
  --set azure.tenantId=$AZURE_TENANT_ID \
  --set azure.clientId=$AZURE_CLIENT_ID \
  --set azure.clientSecret=$AZURE_CLIENT_SECRET
```

If you'd like to customize the installation, please see the 
[configuration](#configuration) section to see options that can be
configured during installation.

To verify the service broker has been deployed and show installed service classes and plans:

```console
$ kubectl get clusterservicebroker -o yaml

$ kubectl get clusterserviceclasses -o=custom-columns=NAME:.metadata.name,EXTERNAL\ NAME:.spec.externalName

$ kubectl get clusterserviceplans -o=custom-columns=NAME:.metadata.name,EXTERNAL\ NAME:.spec.externalName,SERVICE\ CLASS:.spec.clusterServiceClassRef.name --sort-by=.spec.clusterServiceClassRef.name
```

## Uninstalling the Chart

To uninstall/delete the `osba` deployment:

```console
$ helm delete osba --purge
```

The command removes all the Kubernetes components associated with the chart and
deletes the release.

## Configuration

The following tables lists the configurable parameters of the Azure Service
Broker chart and their default values.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `image.repository` | Docker image location, _without_ the tag. | `"microsoft/azure-service-broker"` |
| `image.tag` | Tag / version of the Docker image. | `"v0.4.0-alpha"` |
| `image.pullPolicy` | `"IfNotPresent"`, `"Always"`, or `"Never"`; When launching a pod, this option indicates when to pull the OSBA Docker image. | `"IfNotPresent"` |
| `registerBroker` | Whether to register this broker with the Kubernetes Service Catalog. If true, the Kubernetes Service Catalog must already be installed on the cluster. Marking this option false is useful for scenarios wherein one wishes to host the broker in a separate cluster than the Service Catalog (or other client) that will access it. | `true` |
| `service.type` | Type of service; valid values are `"ClusterIP"`, `"LoadBalancer"`, and `"NodePort"`. `"ClusterIP"` is sufficient in the average case where OSBA only receives traffic from within the cluster-- e.g. from Kubernetes Service Catalog. | `"ClusterIP"` |
| `service.nodePort.port` | _If and only if_ `service.type` is set to `"NodePort"`, `service.nodePort.port` indicates the port this service should bind to on each Kubernetes node. | `30080` |
| `azure.environment` | Indicates which Azure public cloud to use. Valid values are `"AzureCloud"` and `"AzureChinaCloud"`. | `"AzureCloud"` |
| `azure.subscriptionId` | Identifies the Azure subscription into which OSBA will provision services. | none |
| `azure.tenantId` | Identifies the Azure Active Directory to which the _service principal_ used by OSBA to access the Azure subscription belongs. | none |
| `azure.clientId` | Identifies the _service principal_ used by OSBA to access the Azure subscription. | none |
| `azure.clientSecret` | Key/password for the _service principal_ used by OSBA to access the Azure subscription. | none |
| `basicAuth.username` | Specifies the basic auth username that clients (e.g. the Kubernetes Service Catalog) must use when connecting to OSBA. | `"username"`; __Do not use this default value in production!__ |
| `basicAuth.password` | Specifies the basic auth password that clients (e.g. the Kubernetes Service Catalog) must use when connecting to OSBA. | `"password"`; __Do not use this default value in production!__ |
| `encryptionKey` | Specifies the key used by OSBA for applying AES-256 encryption to sensitive (or potentially sensitive) data. | `"This is a key that is 256 bits!!"`; __Do not use this default value in production!__ |
| `modules.minStability` | Specifies the minimum level of stability an OSBA module must meet for the services and plans it provides to be included in OSBA's catalog of offerings. Valid values are `"EXPERIMENTAL"`, `"PREVIEW"`, and `"STABLE"`. | `"PREVIEW"`; __Only use `"STABLE"` modules in production!__ |
| `redis.embedded` | OSBA uses Redis for data persistence and as a message queue. This option indicates whether an on-cluster Redis deployment should be included when installing this chart. If set to `false`, connection details for a remote Redis cache must be provided. | `true`; __Do not use the embedded Redis cache in production!__ |
| `redis.host` | _If and only if_ `redis.embedded` is `false`, this option specifies the location of the remote Redis cache. | none |
| `redis.port` | _If and only if_ `redis.embedded` is `false`, this option specifies the port to connect to on the remote Redis host. | `6380` |
| `redis.redisPassword` | Specifies the Redis password. If `redis.embedded` is `true`, this option sets the password on the Redis cache itself _and_ provides it to OSBA. If `redis.embedded` is `false`, this option only provides the password to OSBA. In that case, the value of this option must be the correct password for the remote Redis cache. | `"password"`; __Do not use this default value in production!__ |
| `redis.enableTls` | _If and only if_ `redis.embedded` is `false`, this option specifies whether to use a secure connection to the remote Redis host. | `true` |

Specify a value for each option using the `--set <key>=<value>` switch on the
`helm install` command. That switch can be invoked multiple times to set
multiple options.

Alternatively, copy the charts default values to a file, edit the file to your
liking, and reference that file in your `helm install` command:

```console
$ helm inspect values azure/open-service-broker-azure > values.yaml
$ vim my-values.yaml
$ helm install azure/open-service-broker-azure --name osba --namespace osba --values my-values.yaml
```

