# Azure SQL Database Examples

The examples in this directory show different ways to use the expanded Azure SQL Database catalog. The catalog now provides three services:
* azure-sqldb (All-in-One SQL Database Server VM and Database)
* azure-sqldb-vm-only (SQL Database Server VM Only)
* azure-sqldb-db-only (SQL Database only)

The use of the azure-sqldb-vm-only service and the  azure-sqldb-db-only service allow you to independantly provision the SQL Database VM and the Database itself. This also enables the creation of multiple databases on a single server VM.

This directory contains example Kubernetes manifests to exercise these services. 

## sqldb-all-in-one

The `sqldb-all-in-one-instance.yaml` manifest will create an all-in-one Azure SQL Database that includes both the server VM and the database. The `sqldb-all-in-one-binding.yaml` will create a binding for the all in one service instance.

## sqldb-vm-only

The `sqldb-vm-only-instance.yaml` manifest will provision an instance of the `azure-sqldb-vm-only` service. This service is not bindable, so there is no corresponding binding manifest. An important elemnt of this manifest is the `alias` parameter. This is used when provisioning an instance of the `azure-sqldb-db-only` service.

## sqldb-db-only

The `sqldb-db-only-instance.yaml` manifest will provision an instance of the `azure-sqldb-db-only` service. This service *requires* a parameter called `parentAlias`. The value of this parameter matches the `alias` paramter,  which is defined in the `sqldb-vm-only-instance.yaml` manifest. 

The `sqldb-db-only-instance-binding.yaml` manifest can then be used to create a service binding to the database only service instance created above.
 
