# Roadmap

This is the official roadmap for the Open Service Broker for Azure (OSBA) project. OSBA is currently available as an alpha release and this roadmap represents our plan to deliver (OSBA) and the services it provides in a stable form.

OSBA releases follow [semantic versioning](https://semver.org/) strictly, and the current stability is pre-`v1.0.0`, which means
that backward incompatible changes are possible.

This roadmap addresses stability in terms of two dimensions:

* Services/Plans
* Database Schemas

These aspects of the broker can mature out of step and need to be considered independently. For example, the broker currently uses Redis to store all data and orchestrate asynchronous queues. As the broker matures, we may make changes that require additional information to be persisted for provisioned services or change tablec or collection names. Before the broker has a `v1.0.0` release, these changes may
introduce incompatibilities with previously populated Redis instances.

Once we declare services and our broker have reached preview status, we will begin to ensure backward compatibility for future releases.

See below for more information concerning stability plans for each.

# Services and Plans

As services are developed, the plans available and features of the service may change and are not guaranteed to exist in the final stable release. When a service is promoted to stable, we will ensure backward compatibility. Services/Plans have three stability tiers:

- `experimental`
- `preview`
- `stable`

Please see [stability documentation](stability.md) for information on each of these tiers.

# July 2018

Currently, the following services are in 'preview':

- [Azure SQL Database](https://azure.microsoft.com/services/sql-database/)
- [Azure Database for MySQL](https://azure.microsoft.com/services/mysql/)
- [Azure Database for PostgreSQL](https://azure.microsoft.com/services/postgresql/)

We expect to declare these services 'stable' by July 2018.

# December 2018

By December 2018, we expect to release the following services as 'stable':

- [Azure CosmosDB](https://azure.microsoft.com/services/cosmos-db/)
- [Azure Redis Cache](https://azure.microsoft.com/services/cache/)
- [Azure Service Bus](https://azure.microsoft.com/services/service-bus/)
- [Azure Event Hubs](https://azure.microsoft.com/services/event-hubs/)

We may release other services listed as 'experimental' or 'preview' during this time.

# Database Schemas

There are two tiers of database schemas to consider:

1. General outline: the broker currently uses Redis to store all data and orchestrate asynchronous queues. 'General outline' means the names and arrangement of the "tables" and layout of asynchronous queues
1. Individual services: the broker stores provisioning and binding data for each service/plan in a "table" in Redis. Individual services means the layout of each provisioning & binding record. When a service is marked `stable`, its individual service schema will be backward compatible.

All releases before `v1.0.0` may introduce backward-incompatible changes to either database schema.

## General Outline

The general outline schema will remain backward compatible within a major release >= `v1.0.0`.
For example, the general outline schema will remain backward compatible between `v1.0.0` and `v2.0.0`.
This follows semver.

## Individual Services

The individual service schemas for `stable` services will remain backward compatible within
a major release version >= `v1.0.0`. For example, the individual service schema for
any `stable` service will remain backward compatible between `v1.0.0` and `v2.0.0`.

The schemas for all `preview` and `experimental` services may change at any time,
possibly in backward incompatible ways.
