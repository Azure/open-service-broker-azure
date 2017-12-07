# Open Service Broker&trade; for Azure

[![CircleCI](https://circleci.com/gh/Azure/open-service-broker-azure.svg?style=svg&circle-token=aa5b73cd7dbb09923f96d9c250b85df671693260)](https://circleci.com/gh/Azure/open-service-broker-azure)

**Open Service Broker for Azure** is the open source,
[Open Service Broker](https://www.openservicebrokerapi.org/)
compatible API server that provisions managed services in the Microsoft Azure
public cloud.

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/150px-Warning.svg.png) | This software is under heavy development. Releases observe [semantic versioning](https://semver.org), but the project is in an ALPHA status and no assurances are made regarding backwards compatibility or stability. All releases prior to v1.0.0 will remain subject to the possibility of breaking changes when the MINOR version number has been incremented. |
|---|---|

![Open Service Broker for Azure GIF](docs/images/demovideo.gif)

*CLOUD FOUNDRY and OPEN SERVICE BROKER are trademarks of the CloudFoundry.org Foundation in the United States and other countries.*

## Supported Services

* [Azure Container Instances](docs/modules/aci.md)
* [Azure Container Registry](docs/modules/acr.md)
* [Azure CosmosDB](docs/modules/cosmosdb.md)
* [Azure Event Hubs](docs/modules/eventhubs.md)
* [Azure Key Vault](docs/modules/keyvault.md)
* [Azure SQL Database](docs/modules/sqldb.md)
* [Azure Database for MySQL](docs/modules/mysqldb.md)
* [Azure Database for PostgreSQL](docs/modules/postgresqldb.md)
* [Azure Redis Cache](docs/modules/rediscache.md)
* [Azure Search](docs/modules/search.md)
* [Azure Service Bus](docs/modules/servicebus.md)
* [Azure Storage](docs/modules/storage.md)

# Quickstart
The [Minikube Quickstart](docs/quickstart-minikube.md) walks-through using the
Open Service Broker for Azure to deploy WordPress on a local Minikube cluster.
Go from "_I have an Azure account that I have never used_" to "_I just deployed WordPress and know what OSBA means!_" �

## Getting Started on Kubernetes

### Installing


You'll need a few pre-requisites before you run these examples on Kubernetes.
Instructions on how to install each prerequisite are linked below:

- [A compatible Kubernetes cluster](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-1-create-a-compatible-kubernetes-cluster)
- [A working Helm installation](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-2-initialize-helm-on-the-cluster)
- [Service Catalog](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-3-install-service-catalog)
- [Helm](https://kubernetes/helm)

Finally, you'll need [Service Catalog CLI](https://github.com/Azure/service-catalog-cli)
installed to introspect the Kubernetes cluster. Please refer to the 
[CLI installation instructions](https://github.com/Azure/service-catalog-cli#install)
for details on how to install it onto your machine.

[Helm](https://helm.sh) is used to install Open Service Broker for Azure onto Kubernetes
clusters. Please refer to the
[Helm chart](https://github.com/Azure/helm-charts/tree/master/open-service-broker-azure)
for details on how to complete the installation.

#### Provisioning

With the Kubernetes Service Catalog software and Open Service Broker for Azure both
installed on your Kubernetes cluster, try creating a `ServiceInstance` resource
to see service provisioning in action.

The following will provision PostgreSQL on Azure:

```console
$ kubectl create -f contrib/k8s/examples/postgresqldb-instance.yaml
```

After the `ServiceInstance` resource is submitted, you can view its status:

```console
$ svcat get instance my-postgresql-instance
```

You'll see output that includes a status to indicate that asynchronous 
provisioning is ongoing. Eventually that status will change to indicate
that asynchronous provisioning is complete.

#### Binding

Upon success, bind to the instance:

```console
$ kubectl create -f contrib/k8s/examples/postgresqldb-binding.yaml
```

To check the status of the binding:

```console
$ svcat get binding my-postgresql-binding
```

You'll see some output to indicate that the binding was successful. Once it is,
a secret named `my-postgresql-secret` will be written that contains the database
connection details in it.
You can observe that this secret exists and has been populated:

```console
$ kubectl get secret my-postgresqldb-secret -o yaml
```

This secret can be used just as any other.

#### Unbinding

To unbind:

```console
$ kubectl delete servicebinding my-postgresqldb-binding
```

Observe that the secret named `my-postgresqldb-secret` is also deleted:

```console
$ kubectl get secret my-postgresqldb-secret
Error from server (NotFound): secrets "my-postgresqldb-secret" not found
```

#### Deprovisioning

To deprovision:

```console
$ kubectl delete serviceinstance my-postgresqldb-instance
```

You can observe the status to see that asynchronous deprovisioning is ongoing:

```console
$ svcat get instance my-postgresql-instance
```

## Getting Started on Cloud Foundry

### Installation

To deploy Open Service Broker for Azure to Cloud Foundry, please refer to the 
[CloudFoundry installation documentation](contrib/cf/README.md) for instructions.

### Usage

#### Provisioning

The following will create a Postgres service:

```console
cf create-service azure-postgresqldb basic50 mypostgresdb -c '{"location": "westus2"}'
```

You can check the status of the service instance using the `cf service` command, which will show output similar to the following:

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

Once the service has been successfully provisioned, you can bind to it, either using `cf bind-service` or by including it in a Cloud Foundry manifest.

```console
cf bind-service myapp mypostgresdb
```

Once bound, the connection details for the service (such as its endpoint and authentication credentaials) are available from the `VCAP_SERVICES` environment variable within the application. You can view the environment variables for a given application using the `cf env` command:

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

# Contributing

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
