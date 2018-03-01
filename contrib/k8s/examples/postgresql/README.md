# Azure Database for PostgreSQL Examples

The examples in this directory show different ways to use the Azure Database for PostgreSQL catalog. The catalog provides three services:

| Service Name | Description |
|--------------|-------------|
| `azure-postgresql` | Database Management System (DBMS) and Database |
| `azure-postgresql-dbms` | DBMS-only |
| `azure-postgresql-database` | Database-only |

This directory contains example Kubernetes manifests to exercise these services.

## Basic Usage

The easiest way to use the Azure Database for PostgreSQL service is to use the `azure-postgresql` service to provision both a DBMS and a new database.

The `postgresql-instance.yaml` manifest will create an instance using the `azure-postgresql` service using the `basic50` plan. This will result in a new Azure Database for PostgreSQL that includes both the DBMS and the database. The `postgresql-binding.yaml` manifest will create a binding for the to this new database, ultimately resulting in a new Kubernetes secret named `example-postgresql-secret`. Once created, you can use this secret in an application to connect to the new Azure Database for PostgreSQL instance.

## Advanced Usage

The `advanced` directory contains manifests that you can use to provision the `azure-postgresql-dbms` and `azure-postgresql-database` services. These services allow you to independently provision the Azure SQL Database DBMS and the database itself for more advanced use cases, such as running multiple databases on a single DBMS.

The `postgresql-dbms-instance.yaml` manifest will provision an instance of the `azure-postgresql-dbms` service using the `basic50` plan. This service is not bindable, so there is no corresponding binding manifest. An important element of this manifest is the `alias` parameter. This is used when provisioning an instance of the `azure-postgresql-database` service.

The `postgresql-database-instance.yaml` manifest will provision an instance of the `azure-mypostgresqlql-database` service using the `database` plan. This service *requires* a parameter called `parentAlias`. The value of this parameter matches the `alias` parameter,  which is defined in the `postgresql-dbms-instance.yaml` manifest.

The `postgresql-database-binding.yaml` manifest can then be used to create a service binding to the database only service instance created above.