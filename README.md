# Open Service Broker&trade; for Azure

[![CircleCI](https://circleci.com/gh/Azure/open-service-broker-azure.svg?style=svg&circle-token=aa5b73cd7dbb09923f96d9c250b85df671693260)](https://circleci.com/gh/Azure/open-service-broker-azure)
[![Go Report Card](https://goreportcard.com/badge/github.com/Azure/open-service-broker-azure)](https://goreportcard.com/report/github.com/Azure/open-service-broker-azure)

**Open Service Broker for Azure** is the open source,
[Open Service Broker](https://www.openservicebrokerapi.org/)-compatible API
server that provisions managed services in the Microsoft Azure public cloud.

_Note: This software is alpha quality. Releases use [semantic versioning](https://semver.org) labels, but any release prior to `v1.0.0` may introduce breaking changes if the MINOR version number was incremented. Please see the offical [roadmap](docs/roadmap.md) for more information._

![Open Service Broker for Azure GIF](docs/images/demovideo.gif)

*CLOUD FOUNDRY and OPEN SERVICE BROKER are trademarks of the CloudFoundry.org Foundation in the United States and other countries.*

## Supported Services

* [Azure Container Instances](docs/modules/aci.md)
* [Azure CosmosDB](docs/modules/cosmosdb.md)
* [Azure Database for MySQL](docs/modules/mysqldb.md)
* [Azure Database for PostgreSQL](docs/modules/postgresqldb.md)
* [Azure Event Hubs](docs/modules/eventhubs.md)
* [Azure Key Vault](docs/modules/keyvault.md)
* [Azure Redis Cache](docs/modules/rediscache.md)
* [Azure SQL Database](docs/modules/mssqldb.md)
* [Azure Search](docs/modules/search.md)
* [Azure Service Bus](docs/modules/servicebus.md)
* [Azure Storage](docs/modules/storage.md)

## Quickstarts

Go from "_I have an Azure account that I have never used_" to "_I just deployed WordPress and know what OSBA means!_"

* The [Minikube Quickstart](docs/quickstart-minikube.md) walks through using the
  Open Service Broker for Azure to deploy WordPress on a local Minikube cluster.
* The [AKS Quickstart](docs/quickstart-aks.md) walks through using the
  Open Service Broker for Azure to deploy WordPress on an Azure Managed Kubernetes Cluster (AKS).

Got questions? Ran into trouble? Check out our [Frequently Asked Questions](docs/faq.md).

## Getting Started on Kubernetes

### Installing

#### Prerequisites

You'll need a few prerequisites before you run these examples on Kubernetes.
Instructions on how to install each prerequisite are linked below:

- [A compatible Kubernetes cluster](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-1-create-a-compatible-kubernetes-cluster)
- [A working Helm installation](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-2-initialize-helm-on-the-cluster)
- [Service Catalog](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-3-install-service-catalog)
- [Helm](https://github.com/kubernetes/helm)

#### Service Catalog CLI

Once you've installed the prerequisites, you'll need the Service Catalog CLI, svcat,
installed to introspect the Kubernetes cluster. Please refer to the
[CLI installation instructions](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md#installing-the-service-catalog-cli)
for details on how to install it onto your machine.

#### Helm Chart

Use [Helm](https://helm.sh) to install Open Service Broker for Azure onto your Kubernetes
cluster. Refer to the OSBA [Helm chart](https://github.com/Azure/helm-charts/tree/master/open-service-broker-azure)
for details on how to complete the installation.

### Provisioning

With the Kubernetes Service Catalog software and Open Service Broker for Azure both
installed on your Kubernetes cluster, try creating a `ServiceInstance` resource
to see service provisioning in action.

The following will provision PostgreSQL on Azure:

```console
$ kubectl create -f contrib/k8s/examples/postgresql/postgresql-instance.yaml
```

After the `ServiceInstance` resource is submitted, you can view its status:

```console
$ svcat get instance example-postgresql-all-in-one-instance
```

You'll see output that includes a status indicating that asynchronous
provisioning is ongoing. Eventually, that status will change to indicate
that asynchronous provisioning is complete.

### Binding

Upon provision success, bind to the instance:

```console
$ kubectl create -f contrib/k8s/examples/postgresql/postgresql-binding.yaml
```

To check the status of the binding:

```console
$ svcat get binding example-postgresql-all-in-one-binding
```

You'll see some output indicating that the binding was successful. Once it is,
a secret named `my-postgresql-secret` will be written that contains the database
connection details in it. You can observe that this secret exists and has been
populated:

```console
$ kubectl get secret example-postgresql-all-in-one-secret -o yaml
```

This secret can be used just as any other.

### Unbinding

To unbind:

```console
$ kubectl delete servicebinding my-postgresqldb-binding
```

Observe that the secret named `my-postgresqldb-secret` is also deleted:

```console
$ kubectl get secret my-postgresqldb-secret
Error from server (NotFound): secrets "my-postgresqldb-secret" not found
```

### Deprovisioning

To deprovision:

```console
$ kubectl delete serviceinstance my-postgresqldb-instance
```

You can observe the status to see that asynchronous deprovisioning is ongoing:

```console
$ svcat get instance my-postgresqldb-instance
```

## Getting Started on Cloud Foundry

### Installing

To deploy Open Service Broker for Azure to Cloud Foundry, please refer to the
[CloudFoundry installation documentation](contrib/cf/README.md) for instructions.

### Provisioning

The following will create a Postgres service:

```console
cf create-service azure-postgresqldb basic50 mypostgresdb -c '{"location": "westus2"}'
```

You can check the status of the service instance using the `cf service` command,
which will show output similar to the following:

```console
Service instance: mypostgresdb
Service: azure-postgresqldb
Bound apps:
Tags:
Plan: basic50
Description: Azure Database for PostgreSQL Service
Documentation url:
Dashboard:

Last Operation
Status: create in progress
Message: Creating server uf666164eb31.
Started: 2017-10-17T23:30:07Z
Updated: 2017-10-17T23:30:12Z
```

### Binding

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```console
cf bind-service myapp mypostgresdb
```

Once bound, the connection details for the service (such as its endpoint and
authentication credentials) are available from the `VCAP_SERVICES` environment
variable within the application. You can view the environment variables for a
given application using the `cf env` command:

```console
cf env myapp
```

### Unbinding

To unbind a service from an application, use the cf unbind-service command:

```console
cf unbind-service myapp mypostgresdb
```

### Deprovisioning

To deprovision the service, use the `cf delete-service` command.

```console
cf delete-service mypostgresdb
```

## Contributing

For details on how to contribute to this project, please see
[contributing.md](./docs/contributing.md).

This project welcomes contributions and suggestions. All contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
