# Roadmap

This is the official roadmap for the Open Service Broker for Azure (OSBA) project. OSBA is currently available as an alpha release and this roadmap represents our plan to deliver (OSBA) and the services it provides in a stable form.

This roadmap addressess stability in terms of two dimensions:

* Services/Plans
* Database Schemas

See below for more information concerining stability plans for each.

# Services and Plans

Services/Plans have three stability tiers:

- `experimental` - We have an idea for a new service that we want to support, but don't have a good idea how it should look. We are essentially “throwing something against the wall” to see if it'll stick. Experimental services may be radically changed or removed at any time
- `preview` - We have a better understanding how a service should look, but we don't yet guarantee backward compatibility. We do guarantee that this service won't go back to `experimental`, so we are committing to making it exist in some shape or form
- `stable` - We now understand usage patterns of the service very well and we guarantee full backward compatibility

Here are our timelines for service stabilities:

- SQL Server: `preview` by the end of February 2018
- MySQL: `preview` by the end of February 2018
- PostgreSQL: `preview` by the end of February 2018

We aren't yet setting timelines on other service stabilities.

# Database Schemas

There are two tiers of database schemas to consider:

1. General outline: the broker currently uses Redis to store all data and orchestrate asynchronous queues. General outline means the names and arrangement of the "tables" and layout of asynchronous queues
1. Individual services: the broker stores provisioning and binding data for each service/plan in a "table" in Redis. Individual services means the layout of each provisioning & binding record. When a service is marked `stable`, its individual service schema will be backward compatible.

Our goal is to have a few services into `preview` before we declare the General outline (#1 above) as backward compatible. Since we're declaring three of them (see above) in `preview` at the end of February 2018, *we are aiming at Mid-march 2018 for a stable general database schema (#1 above)*.