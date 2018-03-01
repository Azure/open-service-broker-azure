# Azure Database for MySQL Examples

The examples in this directory show different ways to use the Azure Database for MySQL catalog. The catalog provides three services:

| Service Name | Description |
|--------------|-------------|
| `azure-mysql` | Database Management System (DBMS) and Database |
| `azure-mysql-dbms` | DBMS-only |
| `azure-mysql-database` | Database-only |

This directory contains example Kubernetes manifests to exercise these services.

## Basic Usage

The easiest way to use the Azure Database for MySQL service is to use the `azure-mysql` service to provision both a DBMS and a new database.

The `mysql-instance.yaml` manifest will create an instance using the `azure-mysql` service using the `basic50` plan. This will result in a new Azure Database for MySQL that includes both the DBMS and the database. The `mysql-binding.yaml` will create a binding for the to this new database, ultimately resulting in a new Kubernetes secret named `example-mysql-secret`. Once created, you can use this secret in an application to connect to the new Azure Database for MySQL instance.

## Advanced Usage

The `advanced` directory contains manifests that you can use to provision the `azure-mysql-dbms` and `azure-mysql-databasey` services. These services allow you to independently provision the Azure SQL Database DBMS and the database itself for more advanced use cases, such as running multiple databases on a single DBMS.

The `mysql-dbms-instance.yaml` manifest will provision an instance of the `azure-mysql-dbms` service using the `basic50` plan. This service is not bindable, so there is no corresponding binding manifest. An important element of this manifest is the `alias` parameter. This is used when provisioning an instance of the `azure-mysql-database` service.

The `mysql-database-instance.yaml` manifest will provision an instance of the `azure-mysql-database` service using the `database` plan. This service *requires* a parameter called `parentAlias`. The value of this parameter matches the `alias` parameter,  which is defined in the `mysql-dbms-instance.yaml` manifest.

The `mysql-database-binding.yaml` manifest can then be used to create a service binding to the database only service instance created above.
