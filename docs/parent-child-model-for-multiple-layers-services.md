# Parent-Child Model for Multiple Layers Services

Services like MySQL and MS SQL are not single-layer services -- they have database management system (DBMS) and database. To satisfy the user requirement to create databases on a shared DBMS, OSBA brings in the _Parent-Child Model_. Typically, the DBMS services are in the parent service category, and the database services are in the child service category.

## Parent Service

For a parent service, `alias` is always one of the required provisioning parameters. It is the correlative key for a parent service instance to engage with child service instances. For easy memo, it is recommended to set the `alias` with the service instance name.

## Child Service

For a child service, `parentAlias` is always one of the required provisioning parameters. A child service instance can only be created with a valid `parentAlias` of an existing parent service instance.

## Notes

* A parent service instance can have multiple children. For example, a DBMS instance can have several database instances as its children.

* A service module can contain several parent services. For example, for MS SQL service module, parent DBMS service instances can be created by OSBA. Or, you can bring your own Azure SQL Servers and register them as parent service instances.
