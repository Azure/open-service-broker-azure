# Azure Service Broker

[![CircleCI](https://circleci.com/gh/Azure/azure-service-broker.svg?style=svg&circle-token=aa5b73cd7dbb09923f96d9c250b85df671693260)](https://circleci.com/gh/Azure/azure-service-broker)

[Azure Service Broker](https://github.com/Azure/azure-service-broker) is the
open source, [Open Service Broker](https://www.openservicebrokerapi.org/)
compatible API server that provisions managed services in the Microsoft
Azure public cloud.

## Getting Started on Kubernetes

### Installing

To install the Azure Service Broker on a Kubernetes cluster, refer to the
[documentation in the Azure Service Broker's Helm chart](contrib/k8s/charts/azure-service-broker/README.md).

### Examples

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
$ kubectl get serviceinstance my-postgresql-instance -o yaml
```

The folowing is excerpted from the output and shows that asynchronous
provisioning is ongoing:

```console
status:
  asyncOpInProgress: true
  conditions:
  - lastTransitionTime: 2017-10-16T18:28:13Z
    message: The instance is being provisioned asynchronously
    reason: Provisioning
    status: "False"
    type: Ready
  currentOperation: Provision
  inProgressProperties:
    externalClusterServicePlanName: basic50
    parameterChecksum: bf5f464a1a09117e100c9a7bb10409e57570f281d756aab4cd00b428c7be16ac
    parameters:
      location: eastus
      resourceGroup: demo
  lastOperation: provisioning
  operationStartTime: 2017-10-16T18:28:13Z
```

Eventually, status will reflect a success or failure state:

```console
status:
  asyncOpInProgress: false
  conditions:
  - lastTransitionTime: 2017-10-16T18:36:43Z
    message: The instance was provisioned successfully
    reason: ProvisionedSuccessfully
    status: "True"
    type: Ready
  externalProperties:
    externalClusterServicePlanName: basic50
    parameterChecksum: bf5f464a1a09117e100c9a7bb10409e57570f281d756aab4cd00b428c7be16ac
    parameters:
      location: eastus
      resourceGroup: demo
  orphanMitigationInProgress: false
  reconciledGeneration: 1
```

#### Binding

Upon success, bind to the instance:

```console
$ kubectl create -f contrib/k8s/examples/postgresql-binding.yaml
```

To check the status of the binding:

```console
$ kubectl get servicebinding my-postgresql-binding -o yaml
```

The following is excerpted from the output:

```console
spec:
  externalID: 25638746-bd86-44a1-a60e-06069734f2cd
  instanceRef:
    name: my-postgresql-instance
  secretName: my-postgresql-secret
status:
  conditions:
  - lastTransitionTime: 2017-10-16T18:44:38Z
    message: Injected bind result
    reason: InjectedBindResult
    status: "True"
    type: Ready
  externalProperties: {}
  orphanMitigationInProgress: false
  reconciledGeneration: 1
```

The status shows the binding was successful and the `spec.secretName` field
indicates that connection details and credentials have been written into a
secret named `my-postgresql-secret`. You can observe that this secret exists
and has been populated:

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
$ kubectl delete serviceinstance my-postgresql-binding
```

You can observe the status to see that asynchronous deprovisioning is ongoing:

```console
$ kubectl get serviceinstance my-postgresql-instance -o yaml
```

The following is excerpted from the output:

```console
status:
  asyncOpInProgress: true
  conditions:
  - lastTransitionTime: 2017-10-16T19:02:19Z
    message: The instance is being deprovisioned asynchronously
    reason: Deprovisioning
    status: "False"
    type: Ready
  currentOperation: Deprovision
  externalProperties:
    externalClusterServicePlanName: basic50
    parameterChecksum: bf5f464a1a09117e100c9a7bb10409e57570f281d756aab4cd00b428c7be16ac
    parameters:
      location: eastus
      resourceGroup: demo
  lastOperation: deprovisioning
  operationStartTime: 2017-10-16T19:02:19Z
  orphanMitigationInProgress: false
  reconciledGeneration: 1
```

When the asynchronous deprovisioning procress completes, the deletion of the
resource will also be complete:

```console
$ kubectl get serviceinstance my-postgresql-instance
Error from server (NotFound): serviceinstances.servicecatalog.k8s.io "my-postgresql-instance" not found
```

## Getting Started on Cloud Foundry

### Installation

To deploy the Azure Service Broker to Cloud Foundry, [refer to the documentation](contrib/cf/README.md) showing how to deploy and configure it as a Cloud Foundry app.

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
