# Azure Service Broker

[![Build Status](https://travis-ci.com/deis/azure-service-broker.svg?token=KPqT8rJc1x6zpm6Zq2Sw&branch=master)](https://travis-ci.com/deis/azure-service-broker)

## Development Guide

### Docker Development Environment

To get started, clone this repo into your GOPATH:

```
$ git clone git@github.com:deis/azure-service-broker.git $GOPATH/src/github.com/Azure/azure-service-broker

$ cd $GOPATH/src/github.com/Azure/azure-service-broker
```

Then use a Docker-based environment to build and run.

> Note: You will need at least Docker version 17.04 and Docker Compose.

If you already have Docker installed then you can get started with:

```
$ make dev-bootstrap
```

### Create Azure Service Principal

If you do not already have an Azure Service Principal, you can create one using the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) or you can follow [this](https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal) to create Service Principals through the Azure portal.

```
$ az ad sp create-for-rbac

{
  "appId": "redacted",
  "displayName": "azure-cli-xxxxxx",
  "name": "http://azure-cli-xxxxxx",
  "password": "redacted",
  "tenant": "redacted"
}
```

### Run API Server and Redis Locally

Using the output from the last step, set the following environment variables and they will get shimmed into the running container.

```
AZURE_SUBSCRIPTION_ID=<YOUR SUBSCRIPTION ID>
AZURE_TENANT_ID=<YOUR TENANT ID>
AZURE_CLIENT_ID=<appId from last step>
AZURE_CLIENT_SECRET=<password from last step>
```

The following command will start running two Docker containers locally. 
- Runs the api server at `http://0.0.0.0:8080`
- Starts a local Redis server at `localhost:6379`

```
$ make run
```

### Deploy for Development

While the API Server and the Redis Server containers are running, run the following from your GOPATH to get the `service id` and `plan id` of the service you want to deploy.


```
# Get service id and plan ids for the service you want to deploy
$ go run contrib/cmd/cli/*.go --host 127.0.0.1 --username username --password password catalog

# You should see something like the following, e.g. azure-postgresqldb

service: azure-postgresqldb   id: 9569986f-c4d2-47e1-ba65-6763f08c3124
   plan: basic50              id: bb6ddfd0-4d8f-4496-aa33-d64ad9562c1f
   plan: basic100             id: ffc1e3c8-0e24-471d-8683-1b42e100bb14
```

Use the `service id` and `plan id` from the previous step, run the following command to kickoff the Azure ARM deployment of the service:

```
$ go run contrib/cmd/cli/*.go -H localhost -u username -P password provision -sid <service id> -pid <plan id> --param location=eastus --poll

Provisioning service instance 256787f5-8ecb-4ec0-9881-58c1bc4cf62b

...................................................................................

Service instance 256787f5-8ecb-4ec0-9881-58c1bc4cf62b has been successfully provisioned
```

### Bind for Development

Now that we have an instance provisioned on Azure, let's bind it for development.

```
$ go run contrib/cmd/cli/*.go -H localhost -u username -P password bind -iid 256787f5-8ecb-4ec0-9881-58c1bc4cf62b

Binding a109d773-4a20-439f-90c4-215df25a4f97 created for service instance 256787f5-8ecb-4ec0-9881-58c1bc4cf62b
Credentials:
   host:                <redacted>.postgres.database.azure.com
   port:                5432
   database:            <redacted>
   username:            <redacted>
   password:            <redacted>
```

### Test Credentials

With the instance binded, now we can test the credentials from the previous step.

```
$ psql -h <host> -U <username>@<host> -d "db=<database> sslmode=require"
```

### Unbind

To unbind the provisioned instance, run the following command with the Instance id and the Binding id from the previous step. Note: Once you unbind the instance, you can no longer access the provisioned instance with the previously created credentials.

```
$ go run contrib/cmd/cli/*.go -H localhost -u username -P password unbind -iid 256787f5-8ecb-4ec0-9881-58c1bc4cf62b -bid a109d773-4a20-439f-90c4-215df25a4f97

Unbound binding a109d773-4a20-439f-90c4-215df25a4f97 to service instance 256787f5-8ecb-4ec0-9881-58c1bc4cf62b
```

### Deprovision

To delete the provisioned service, run the following command:

```
$ go run contrib/cmd/cli/*.go -H localhost -u username -P password deprovision -iid 256787f5-8ecb-4ec0-9881-58c1bc4cf62b
```


## Before Submitting a PR

To ensure CI passes, run the following command to run all the tests and lint check:
> Note: this will trigger lifecycle tests as well. Make sure you have set all the Azure related environment variables to a live Azure subscription.

```
$ make lint test docker-build
```

### Run Tests

Run all tests:
> Note: this will trigger lifecycle tests as well. Make sure you have set all the Azure related environment variables to a live Azure subscription.

```
$ make test
``` 

### Lint Checks

Run all the lint checks:

```
$ make lint
```
If there are lint errors, run the following command:

```
$ gofmt -s -w <filename>
```

## Install on Kubernetes via Helm

To install the Azure Service Broker on a Kubernetes 1.7+ cluster, first ensure
that the latest release of the Kubernetes Service Catalog software has been
deployed to that cluster using
[these instructions](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install-1.7.md).

With the Kubernetes Service Catalog software installed and running, proceed with
broker installation by creating a [service principal]() that can be used by the
broker to interact with your Azure subscription.

Ensure the Azure CLI (command line interface) is installed on your system:

```console
$ which az
```

If the Azure CLI is not found, it can be installed using
[these instructions](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest).

With the Azure CLI installed, log in and follow the prompts:

```console
az login
```

When the login process has been completed successfully, note the value of the
`id` field of the `az login` command's JSON output. We'll export
this as the value of an environment variables for our own convenience.

```console
$ export AZURE_SUBSCRIPTION_ID=<id>
```

Now use the CLI to create a new service principal with the
`Contributor` role. This will allow the service principal to provision all
Azure services on the broker's behalf.

```console
$ az ad sp create-for-rbac \
    --role="Contributor" \
    --scopes="/subscriptions/${AZURE_SUBSCRIPTION_ID}"
```

Note the values of the `tenant`, `appId`, and `password` fields in the command's
JSON output. We'll export these values as environment variables for our
convenience:

```console
$ export AZURE_TENANT_ID=<tenant>
$ export AZURE_CLIENT_ID=<appId>
$ export AZURE_CLIENT_SECRET=<password>
$ export AZURE_SUBSCRIPTION_ID=<subscriptionId>
```

Now use [Helm](https://helm.sh/) to install the broker using defaults, which includes the used of an embedded Redis database. From the `contrib/k8s/charts` directory, execute the following:

```console
$ helm install azure-service-broker --name asb --namespace asb \
    --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
    --set azure.tenantId=$AZURE_TENANT_ID \
    --set azure.clientId=$AZURE_CLIENT_ID \
    --set azure.clientSecret=$AZURE_CLIENT_SECRET
```

If you have a Redis database outside of the cluster you would like to use for the broker instead of the embedded Redis, execute the following:

```console
$ export REDIS_HOST=<redishost>
$ export REDIS_PASSWORD=<redispassword>

$ helm install azure-service-broker --name asb --namespace asb \
    --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
    --set azure.tenantId=$AZURE_TENANT_ID \
    --set azure.clientId=$AZURE_CLIENT_ID \
    --set azure.clientSecret=$AZURE_CLIENT_SECRET \
    --set redis.host=$REDIS_HOST \
    --set redis.password=$REDIS_PASSWORD \
    --set redis.embedded=false
```

__Advanced: To achieve a secure and stable deployment in a production
environment, please supply a custom `values.yml` file during installation to
override default passwords, keys, and database location.__

After installing, you may wish to monitor the status of the broker pod until it
enters a healthy state. This can be accomplished like so:

```console
$ kubectl get pods -n asb -w
```

## Uninstalling from Kubernetes via Helm

If you followed the installation instructions in the previous section and wish
to uninstall the broker, begin by deleting the Helm release:

```console
$ helm delete asb --purge
```

If you wish to also uninstall the Kubernetes Service Catalog software:

```console
$ helm delete catalog --purge
```

__Note that uninstalling either or both of these packages will NOT effect
deprovisioning of any services that were provisioned while they were running.__
