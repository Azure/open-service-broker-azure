# OSBA Versions and Module Stability

OSBA has an overall stability, defined by its release version, which is
based on [semver](https://semver.org/), but each of its
[services](/README.md#supported-services) also has a stability label.

Note that OSBA services refer to the catalog offering, and not necessarily one
specific Azure service.

# Overall Stability and OSBA Version

We version OSBA generally via its release version to indicate its overall
stability. Generally, this includes how well it's tested, its stability "in the wild,"
how many people using it, bug reports, and more.

But when we increase a semver MAJOR, MINOR or PATCH version, we also take into
account if any of the OSBA services change in stability. Read on for how this works.

# Service Stability

We indicate a stability for each service in the OSBA catalog and docs to
indicate how mature it is. Each service may have a different stability:

* `experimental` - We have an idea for a new service in the catalog that we want to support,
  but don't have a good idea how it should look. We are essentially
  “throwing something against the wall” to see if it'll stick. Experimental
  services may be radically changed or removed at any time
* `preview` - We have a better understanding how a service should look, but we
  don't yet guarantee backward compatibility. We do guarantee that this service
  won't go back to `experimental`, so we are committing to making it exist in
  some shape or form
* `stable` - We now understand usage patterns of the service very well and we
  guarantee full backward compatibility. We will not promote a service to stable
  if the Azure service it provisions is not GA.

We've added these stability labels so that services can move freely and
independently across semver releases, which is important so we can improve
an individual service, or subset of service at one time.

# Service Stability and OSBA Version

We correlate service stability changes to OSBA version changes according to the
following rules:

* If a service stability goes from `experimental` to `preview`, a semver MINOR
  or MAJOR release must happen
* If a service stability goes from `preview` to `stable`, a semver MINOR
  or MAJOR release must happen
* If any service's stability goes "down" (from `stable` to `preview` or
  `experimental`, or from `preview` to `experimental`), a MAJOR release must
  happen (this would be a breaking change)
  * Downgrading a service's stability will be rare
