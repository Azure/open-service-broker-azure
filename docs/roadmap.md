# Roadmap

This is the official roadmap for the Open Service Broker for Azure (OSBA) project. OSBA is currently available as an alpha release and this roadmap represents our plan to deliver (OSBA) and the services it provides in a stable form.

This roadmap addresses stability in terms of two dimensions:

* Services/Plans
* Database Schemas

These aspects of the broker can mature out of step and need to be considered independently. For example, the broker currently uses Redis to store all data and orchestrate asynchronous queues. As the broker matures, we may make changes that require additional information to be persisted for provisioned services or change table or collection names. These will result in incompatibilities with previously populated Redis instances. As we make these changes to mature the broker and work toward a stable release, we are not able to guarantee backward compatibility or migration paths for previous broker installations. Once we declare services and our broker have reached preview status, we will begin to ensure backward compatibility for future releases.

See below for more information concerning stability plans for each.

# Services and Plans

As services are developed, the plans available and features of the service may change and are not guaranteed to exist in the final stable release. When a service is promoted to stable, we will ensure backward compatibility. Services/Plans have three stability tiers:

- `experimental` - We have an idea for a new service that we want to support, but don't have a good idea how it should look. We are essentially “throwing something against the wall” to see if it'll stick. Experimental services may be radically changed or removed at any time
- `preview` - We have a better understanding how a service should look, but we don't yet guarantee backward compatibility. We do guarantee that this service won't go back to `experimental`, so we are committing to making it exist in some shape or form
- `stable` - We now understand usage patterns of the service very well and we guarantee full backward compatibility

Here are our timelines for service stabilities:

- SQL Server: `preview` by the end of February 2018
- MySQL: `preview` by the end of February 2018
- PostgreSQL: `preview` by the end of February 2018

We aren't yet setting timelines on other service stabilities, but will continue to mature additional services once we promote the service modules above to preview. Promotion from preview to stable will depend on the maturity of the underlying service. For example, the [Azure Database for PostgreSQL](https://azure.microsoft.com/en-us/services/postgresql/) service offering is currently in `Preview`. We will not promote the OSBA PostgreSQL offering to `stable` until the service itself is promoted out of `Preview` status.

# Database Schemas

There are two tiers of database schemas to consider:

1. General outline: the broker currently uses Redis to store all data and orchestrate asynchronous queues. General outline means the names and arrangement of the "tables" and layout of asynchronous queues
1. Individual services: the broker stores provisioning and binding data for each service/plan in a "table" in Redis. Individual services means the layout of each provisioning & binding record. When a service is marked `stable`, its individual service schema will be backward compatible.

In both areas, maturation of the broker may result in changes to the schema that will be incompatible with previous broker and Redis deployments. Our goal is to have a few services into `preview` before we declare the General outline (#1 above) as backward compatible. Since we're declaring three of them (see above) in `preview` at the end of February 2018, *we are aiming at Mid-march 2018 for a stable general database schema (#1 above)*.