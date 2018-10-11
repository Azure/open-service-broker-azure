# Parent-Child Model for Multiple Layers Services

Services like MySQL and MS SQL are not single-layer services -- they have database management system (DBMS) and database. To satisfy the user requirement to create databases on a shared DBMS, OSBA brings in the _Parent-Child Model_. Typically, the DBMS services are in the parent service category, and the database services are in the child service category.

## Parent Service

For a parent service, `alias` is always one of the required provisioning parameters. It is the correlative key for a parent service instance to engage with child service instances. For easy memo, it is recommended to give the same value to the service instance name.

## Child Service

For a child service, `parentAlias` is always one of the required provisioning parameters. It should be specified an existing `alias` of a parent service instance. A child service instance can only be created successfully in this way.

## Notes

* A parent service instance can have multiple children. That is the way to achieve that, databases share a DBMS.

* A service module can contain several parent services. For example MS SQL service module, parent DBMS service instances can be created by OSBA, and you can bring your own Azure SQL servers to register them as parent service instances.
