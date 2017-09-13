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

### Install Broker via Helm chart
Set up your REDIS cache in Azure and then make sure to export your REDIS_HOST and REDIS_PASSWORD in addition to the AZURE specific variables above
> Note: We will be providing a local option as well soon

```console

$ helm install azure-service-broker --name asb --namespace asb \
     --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
     --set azure.tenantId=$AZURE_TENANT_ID \
     --set azure.clientId=$AZURE_CLIENT_ID \
     --set azure.clientSecret=$AZURE_CLIENT_SECRET

$ kubectl get pods -n asb -w
# wait for the broker pod to enter a healthy state
```

### Troubleshooting

If you run into issues with deleting the broker, you will need to remove the catalog to clean up and install again
