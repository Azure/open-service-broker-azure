# Azure SQL Database Examples

The examples in this directory show different ways to use the Azure SQL Database catalog. The catalog provides three services:

| Service Name | Description |
|--------------|-------------|
| `azure-sql` | Database Management System (DBMS) and Database |
| `azure-sql-dbms-only` |  DBMS-only |
| `azure-sql-database-only` |  Database-only |

This directory contains example Kubernetes manifests to exercise these services.

## Basic Usage

The easiest way to use the Azure SQL Database service is to use the `azure-sql` service to provision both a DBMS and a new database.

The `sql-instance.yaml` manifest will create an instance using the `azure-sql` service using the `basic` plan. This will result in a new Azure SQL Database that includes both the DBMS and the database. The `sql-binding.yaml` will create a binding for the to this new database, ultimately resulting in a new Kubernetes secret named `example-sql-secret`. Once created, you can use this secret in an application to connect to the new Azure SQL Database instance.

## Advanced Usage

The `advanced` directory contains manifests that you can use to provision the `azure-sql-dbms-only` and `azure-sql-database-only` services. These services allow you to independently provision the Azure SQL Database DBMS and the database itself for more advanced use cases, such as running multiple databases on a single DBMS.

The `sql-dbms-instance.yaml` manifest will provision an instance of the `azure-sql-dbms-only` service using the `sql-dbms-only` plan. This service is not bindable, so there is no corresponding binding manifest. An important element of this manifest is the `alias` parameter. This is used when provisioning an instance of the `azure-sql-database-only` service.

The `sql-database-instance.yaml` manifest will provision an instance of the `azure-sql-database-only` service using the `basic` plan. This service *requires* a parameter called `parentAlias`. The value of this parameter matches the `alias` parameter,  which is defined in the `sql-dbms-instance.yaml` manifest.

The `sql-database-binding.yaml` manifest can then be used to create a service binding to the database only service instance created above.
