# Development Guide

This document supplements the contribution guidelines with the technical details
that contributors will require to successfully amend, build, and test the Azure
Service Broker.

## Cloning the Repository

We assume your system is configured for Go development and that the environment
variable `GOPATH` is therefore defined. If this is not the case, start by
exporting this environment variable. Use your discretion in choosing a path,
but the path used below should generally be adequate:

```console
export GOPATH=~/Code/go
```

Then, create the proper directory and clone this repository to it:

```console
$ mkdir -p $GOPATH/src/github.com/Azure
$ git clone git@github.com:Azure/open-service-broker-azure.git \
    $GOPATH/src/github.com/Azure/open-service-broker-azure
$ cd $GOPATH/src/github.com/Azure/open-service-broker-azure
```

Note also that all of the above is for the benefit of one's IDE (integrated
development environment), which may use your system's natively installed Go
tools to introspect, compile, and test code. Because Open Service Broker for Azure
depends upon a consistent, containerized build / test environment, the placement
of code on one's system does not impact one's ability to build and test code
using the commands documented in the following sections.

## The Containerized Development Environment

Open Service Broker for Azure utilizes a consistent, containerized build / test
environment to ensure all contributors build and test their patches using the
_exact_ same environment employed by CI. This approach eliminates the "it worked
on my machine" factor and minimizes the set of tools that contributors must have
installed natively on their system to ensure success. e.g. It is not even
necessary to have Go (or any particular version thereof) installed.

### Prerequisites

The prerequisites for successfully building and testing Open Service Broker for Azure
are:

- make
- Docker (version 17.04 or greater)
- Docker Compose (version 1.16.1 or greater)

### Building, Testing, and Running

Building and testing Open Service Broker for Azure code is facilitated through the use
of a few easy-to-use make targets that mostly execute tasks within the
containerized development environment.

#### Linting

Open Service Broker for Azure project observes quite a large number of style
conventions that are common within the Go community. Adherence to these
conventions is enforced via the
[Go Meta Linter](https://github.com/alecthomas/gometalinter) tool.

All code changes should be linted before a PR is opened, as any defiance of
conventions will cause CI to fail.

To run the linter:

```console
$ make lint
```

#### Running Unit Tests

To execute unit tests:

```console
$ make test-unit
```

#### Running "Lifecycle" Tests

Open Service Broker for Azure is used to facilitate provisioning and binding to various
managed services provided by Microsoft Azure. To assert integration with these
many services works as expected, the project contains a "lifecycle" test suite.
These tests are integration tests that exercise the provision / bind / unbind /
deprovision lifecycle against a live Azure subscription for each service
integration (module).

Executing lifecycle tests requires an Azure subscription and some further
setup.

First, obtain your Azure subscription ID and export it as an environment
variable:

```console
$ export AZURE_SUBSCRIPTION_ID=$(az account show --query id -o tsv)
```

Next create a service principal (service account) in your Azure Active Directory
tenant. This is the identity that Open Service Broker for Azure will use when
authenticating to Azure endpoints.
  
```console
$ az ad sp create-for-rbac
```

The new service principal will be assigned, by default, to the `Contributor`
role, which gives it adequate access to provision and deprovision _any_
resources in your Azure subscription. Guard these credentials carefully.

The output of the command above will be similar to the following:

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

You may wish to export these environment variables in your shell environment's
profile to avoid the need to repeat these steps in the future.

Finally, execute the lifecycle tests.

__WARNING: Executing these tests will provision lots of real services within
your Azure subscription. This will cost you money! The tests do a good job of
cleaning up after themselves, so generally, the services stick around only for
as long as it takes the tests to complete (generally around 20 minutes).
HOWERVER, if the tests are interrupted, tests may not successfully clean up
after themselves.__

```console
$ make test-service-lifecycles
```

If you'd like to run a subset of the service lifecycle tests, you can select
a subset of the services by a comma delimited list of module names to the 
TEST_MODULES environment variable. For example, to run only the `postgresql` and 
`mssql` tests:

```console
$ TEST_MODULES=postgresql,mssql make test-service-lifecycles
```

Regardless of success or failure, after tests have completed, you can verify
that they have cleaned up after themselves by searching for resource groups
named using the scheme "test=*":

```console
$ az group list | grep name | grep "test-"
```

If tests ever appear not to have cleaned up after themselves, the
following command can be used to manually clean up.

__WARNING: This command will INDISCRIMINATELY delete all resource groups from
your Azure subscription(s) if they match the naming scheme "test-*". Understand
the implications of that well before executing this command.__

```console
$ for g in $( \
    az group list \
    | grep name \
    | grep "test-" \
    | awk '{print $2}' \
    | awk '{gsub(/\"|\,/,"")}1'
    ); \
  do \
    echo "Deleting ${g}..."; \
    az group delete --name $g --yes --no-wait; \
    echo "Done deleting ${g}."; \
  done
```

#### Running OSB API Compliance Tests

The Azure Service Broker implements the OSB API. To assert compliance with
the API specification, the project includes a set of compliance tests.
These tests verify API response codes for various operations such as 
provisionining, binding and deprovision. These tests use a mock service module.

To run the compliance tests:

```console
$ make test-api-compliance
```

#### Linting and Running All Tests

To execute unit tests, lifecycle tests, _and_ compliance tests together:

```console
$ make test
```

To add lint checks to the above:

```console
$ make lint test
```

It is also advisable to ensure that there is no difficulty bundling the built
software into a Docker image:

```console
$ make lint test docker-build
```

Any changes passing these three tests locally should pass the same tests in CI.

#### Updating Dependencies

Updating project dependnecies is not a matter of course for most contributions,
but the need for it arises from time to time, especially in the case that a new
service module is being added to the broker.

Open Service Broker for Azure employs the [dep](https://github.com/golang/dep) tool
to manage dependencies. Dep tracks developer intent (the dependencies you
_want_) in a file called `Gopkg.toml`. How these intentions are resolved by the
tool is tracked in a manifest called `Gopkg.lock`.

If you update the desired dependencies in `Gopkg.toml`, be sure to run the
following afterwards:

```console
$ make dep
```

This will update both vendored code in the `vendor/` directory _and_ the
manifest in `Gopkg.lock`.

__Do not edit `Gopkg.lock` directly.__

If the `make dep` step is accidentally omitted after updates to `Gopkg.toml`,
the CI process will catch the mistake and fail the build.

#### Running Open Service Broker for Azure Locally (development mode)

Running Open Service Broker for Azure requires the use of a live Azure subscription.
Refer to the [Running "Lifecycle" Tests](#running-lifecycle-tests) section for
further details on required setup.

To build and launch Open Service Broker for Azure in a container:

```console
$ make run
```

Through its use of Docker Compose, this make target will not only launch the
Open Service Broker for Azure, but will also launch a containerized Redis database that
the broker will use for persistence and reliable queueing of aynchronous tasks.

Note this method of running is only advisable during development. To deploy
Open Service Broker for Azure to Kubernetes or Pivotal Cloud Foundry, see the
relevant documentation for each:

- Deploying on Kubernetes
- Deploying on Pivotal Cloud Foundry

#### Cleaning Up

If at any time, the state of _anything_ is in doubt, _everything_ can be reset:

```console
$ make clean
```

This will delete binaries, remove running containers, and delete images of both
the containerized development environment and the broker.

Subsequent invocation of any make targets will start from a clean slate.

#### Stepping into the Development Environment

If one wishes to "poke around" the containerized development environment, it
is achievable like so:

```console
$ make dev
```

The above make target will launch the containerized development environment
interactively, leaving the user with a TTY and a bash prompt.

## Interacting with the Broker

While CI relies on the tests documented in previous sections, its a natual
human tendency to want to interact with software we are contributing to. Because
the OSB specification implemented by Open Service Broker for Azure is complex, the
`curl` commands one might use to interact with the broker are also complex.

Rather than burden contributors with the need to craft "artisinal" `curl`
commands for the sake of executing simple actions, Open Service Broker for Azure
project comes with a "bespoke" CLI that is used _only_ to facilitate human
interaction with the broker. (The true clients of the broker are the Kubernetes
Service Catalog and Pivotal Cloud Foundry.)

### Running the CLI

Unless hacking directly on the CLI, it's best to build the CLI and simply
resuse the binary rather than build it each time it's used.

This can be accomplished for your chosen OS using the appropriate command from
the following list:

- Mac OS: `GOOS=darwin make build-broker-cli`
- Linux: `GOOS=linux make build-broker-cli`
- Windows: `GOOS=windows make build-broker-cli`

In all cases, the binary is cross-compiled in the previously-discussed
containerized development environment for an AMD64 architecture. (If you have
a need to cross-compile for another OS/arch combination, that is left as an
exercise for the contributor.)

After building the CLI, it can be invoked using the following command:

```console
$ contrib/bin/broker-cli
```

Invoking the CLI with no arguments as shown above will display help.

Contextually appropriate help can also be shown by appending the `--help` or
`-h` flag to any command. For instance:

```console
$ contrib/bin/broker-cli catalog -h
```

If iteratively hacking on the CLI itself, you may find it more productive to
skip pre-compiling the CLI each time you've made changes and simply use
`go run` instead. We advise that this be done within the containerized
development environment.

First launch the containerized development environment interactively:

```console
$ make dev
```

Then use the following command:

```console
$ go run contrib/cmd/cli/*.go
```

All remaining sections will document CLI sub-commands with the assumption you
have pre-compiled the CLI.

#### Getting the Catalog

To list all services and plans (tiers or variants) thereof, use the following
command:

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    catalog
```

Note the service and plan IDs.

#### Provisioning a Service

To provision a service, use the `provision` sub-command and use the
`--service-id` (or `-sid`) and `--plan-id` (or `-pid`) flags to specify a 
service ID and plan ID.

Note that many services also require a `location` parameter to be set using
the `--parameter` (or `--param`) flag.

The following example provisions a PostgreSQL server (and database on that
server) in Azure's `eastus` region.

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    provision \
    --service-id b43b4bba-5741-4d98-a10b-17dc5cee0175 \
    --plan-id b2ed210f-6a10-4593-a6c4-964e6b6fad62 \
    --parameter location=eastus
```

This will kick of asynchronous provisioning and display input similar to the
following:

```console

Provisioning service instance 1cb9fc31-f2f7-498d-b273-eba8981261de

```

To check on the status of the asynchronous provisioning operation, use the
`poll` sub-command with the `--instance-id` (or `-iid`) and `--operation`
flags set appropriately. For instance:

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    poll \
    --instance-id 1cb9fc31-f2f7-498d-b273-eba8981261de \
    --operation provisioning
```

This will produce output similar to the following:

```console

Instance 1cb9fc31-f2f7-498d-b273-eba8981261de provisioning state: in progress

```

Conveniently, provisioning and polling can be combined into a single command by
making use of the `--poll` flag on the `provision` sub-command like so:

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    provision \
    --service-id b43b4bba-5741-4d98-a10b-17dc5cee0175 \
    --plan-id b2ed210f-6a10-4593-a6c4-964e6b6fad62 \
    --parameter location=eastus \
    --poll
```

The above will poll the status of the asynchronous provisioning operation until
it either succeeds or fails.

#### Binding

To bind to a service, use the `bind` sub-command and specify the service using
the `--instance-id` (or `-iid`) flag:

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    bind \
    --instance-id 1cb9fc31-f2f7-498d-b273-eba8981261de
```

This will produce output similar to the following:

```console

Binding a6df0e0a-5924-4693-9f8b-8a8b7cfdb01b created for service instance 1cb9fc31-f2f7-498d-b273-eba8981261de
Credentials:
   host:                f9d08944-c5d6-4fda-a1d1-a02bc9a6a111.postgres.database.azure.com
   port:                5432
   database:            ce8ctydqu2
   username:            y2sdfvi8v7@f9d08944-c5d6-4fda-a1d1-a02bc9a6a111
   password:            eOImZyoEsF8rYWNr

```

#### Unbinding

To unbind, use the `unbind` sub-command and specify both the instance ID and
binding ID using the `--instance-id` (or `-iid`) and `--binding-id` (or `-bid`)
flags, respectively. For example:

```console:
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    unbind \
    --instance-id 1cb9fc31-f2f7-498d-b273-eba8981261de \
    --binding-id a6df0e0a-5924-4693-9f8b-8a8b7cfdb01b
```

This will produce output similar to the following:

```console

Unbound binding a6df0e0a-5924-4693-9f8b-8a8b7cfdb01b to service instance 1cb9fc31-f2f7-498d-b273-eba8981261de

```

#### Deprovisioning

To deprovision, use the `deprovision` sub-command and specify the instance ID
using the `--instance-id` (or `-iid`) flag:

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    deprovision \
    --instance-id 1cb9fc31-f2f7-498d-b273-eba8981261de
```

This will begin an asynchronous deprovisioning process and produce output
similar to the following:

```console

deprovisioning service instance 1cb9fc31-f2f7-498d-b273-eba8981261de

```

To check on the status of the asynchronous deprovisioning operation, use the
`poll` sub-command with the `--instance-id` (or `-iid`) and `--operation`
flags set appropriately. For instance:

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    poll \
    --instance-id 1cb9fc31-f2f7-498d-b273-eba8981261de \
    --operation deprovisioning
```

This will produce output similar to the following:

```console

Instance 1cb9fc31-f2f7-498d-b273-eba8981261de deprovisioning state: gone

```

Conveniently, deprovisioning and polling can be combined into a single command
by making use of the `--poll` flag on the `deprovision` sub-command like so:

```console
$ contrib/bin/broker-cli \
    --username username \
    --password password \
    deprovision \
    --instance-id 1cb9fc31-f2f7-498d-b273-eba8981261de \
    --poll
```

The above will poll the status of the asynchronous deprovisioning operation
until it either succeeds or fails.
