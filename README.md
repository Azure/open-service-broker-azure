# Azure Service Broker

[![CircleCI](https://circleci.com/gh/Azure/azure-service-broker.svg?style=svg&circle-token=aa5b73cd7dbb09923f96d9c250b85df671693260)](https://circleci.com/gh/Azure/azure-service-broker)

[Azure Service Broker](https://github.com/Azure/azure-service-broker) is the
open source, [Open Service Broker](https://www.openservicebrokerapi.org/)
compatible API server that provisions managed services in the Microsoft
Azure public cloud.

## Getting Started on Kubernetes

### Installing

You'll need a few pre-requisites in order to run these examples on Kubernetes.

Instructions on how to install each prerequisite are linked below:

- [A compatible Kubernetes cluster](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-1-create-a-compatible-kubernetes-cluster)
- [A working Helm installation](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-2-initialize-helm-on-the-cluster)
- [Service Catalog](https://github.com/Azure/helm-charts/blob/master/docs/prerequisities/README.md#step-3-install-service-catalog)

Finally, you'll need [Service Catalog CLI](https://github.com/Azure/service-catalog-cli)
installed to introspect the Kubernetes cluster.

Please refer to the [CLI installation instructions](https://github.com/Azure/service-catalog-cli#install)
for details on how to install it onto your machine.

#### Provisioning

With the Kubernetes Service Catalog software and the Azure Service Broker both
installed on your Kubernetes cluster, try creating a `ServiceInstance` resource
to see service provisioning in action.

The following will provision PostgreSQL on Azure:

```console
$ kubectl create -f contrib/k8s/examples/postgresql-instance.yaml
```

After the `ServiceInstance` resource is submitted, you can view its status:

```console
$ svc-cat get instance my-postgresql-instance
```

You'll see output that includes a status to indicate that asynchronous 
provisioning is ongoing. Eventually that status will change to indicate
that asynchronous provisioning is complete.

#### Binding

Upon success, bind to the instance:

```console
$ kubectl create -f contrib/k8s/examples/postgresql-binding.yaml
```

To check the status of the binding:

```console
$ svc-cat get binding my-postgresql-binding
```

You'll see some output to indicate that the binding was successful. Once it is,
a secret named `my-postgresql-secret` will be written that contains the database
connection details in it.
You can observe that this secret exists and has been populated:

```console
$ kubectl get secret my-postgresql-secret -o yaml
```

This secret can be used just as any other.

#### Unbinding

To unbind:

```console
$ kubectl delete servicebinding my-postgresql-binding
```

Observe that the secret named `my-postgresql-secret` is also deleted:

```console
$ kubectl get secret my-postgresql-secret
Error from server (NotFound): secrets "my-postgresql-secret" not found
```

#### Deprovisioning

To deprovision:

```console
$ kubectl delete serviceinstance my-postgresql-instance
```

You can observe the status to see that asynchronous deprovisioning is ongoing:

```console
$ svc-cat get instance my-postgresql-instance
```

You'll see in the output that asynchronous deprovision is in progress. When
it's complete, the deletion of the resource will also be complete and
the `svc-cat get` command will indicate that the instance no longer exists.

## Getting Started on Cloud Foundry

### Installation

To deploy the Azure Service Broker to Cloud Foundry, please refer to the 
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
