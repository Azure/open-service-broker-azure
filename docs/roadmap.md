# Roadmap

This is the official roadmap for the Open Service Broker for Azure (OSBA) project. OSBA is currently available as an alpha release and this roadmap represents our plan to deliver (OSBA) and the services it provides in a stable form.

This roadmap addresses stability in terms of two dimensions:

* Services/Plans
* Database Schemas

These aspects of the broker can mature out of step and need to be considered independently. For example, the broker currently uses Redis to store all data and orchestrate asynchronous queues. As the broker matures, we may make changes that require additional information to be persisted for provisioned services or change table or collection names. These will result in incompatibilities with previously populated Redis instances. As we make these changes to mature the broker and work toward a stable release, we are not able to guarantee backward compatibility or migration paths for previous broker installations. Once we declare services and our broker have reached preview status, we will begin to ensure backward compatibility for future releases.

See below for more information concerning stability plans for each.

# Services and Plans

As services are developed, the plans available and features of the service may change and are not guaranteed to exist in the final stable release. When a service is promoted to stable, we will ensure backward compatibility. Services/Plans have three stability tiers:

- `experimental`
- `preview`
- `stable`

Please see [stability documentation](stability.md) for information on each of these tiers.

Currently, the following services are in 'preview':

- [Azure SQL Database](https://azure.microsoft.com/services/sql-database/)
- [Azure Database for MySQL](https://azure.microsoft.com/services/mysql/)
- [Azure Database for PostgreSQL](https://azure.microsoft.com/services/postgresql/)

We expect to declare these services 'stable' by July 2018.

In the 2nd half of 2018, we intend to move the following services into 'preview' and ultimately to a 'stable' state:

- [Azure CosmosDB](https://azure.microsoft.com/services/cosmos-db/)
- [Azure Redis Cache](https://azure.microsoft.com/services/cache/)
- [Azure Service Bus](https://azure.microsoft.com/services/service-bus/)
- [Azure Event Hubs](https://azure.microsoft.com/services/event-hubs/)

Other services, including the remainder currently listed as 'experimental' will follow these.

# Database Schemas

There are two tiers of database schemas to consider:

1. General outline: the broker currently uses Redis to store all data and orchestrate asynchronous queues. General outline means the names and arrangement of the "tables" and layout of asynchronous queues
1. Individual services: the broker stores provisioning and binding data for each service/plan in a "table" in Redis. Individual services means the layout of each provisioning & binding record. When a service is marked `stable`, its individual service schema will be backward compatible.

In both areas, maturation of the broker may result in changes to the schema that will be incompatible with previous broker and Redis deployments. Our goal is to have a few services into `preview` before we declare the General outline (#1 above) as backward compatible. Since we're declaring three of them (see above) in `preview` at the end of February 2018, *we are aiming at Mid-march 2018 for a stable general database schema (#1 above)*.