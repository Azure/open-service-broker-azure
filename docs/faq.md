# FAQ
_What to do when things go sideways_



## Troubleshooting

Here are some issues that people have run into in the past:

### OSBA Pod in CrashLoopBackOff, can't lookup osba-redis

```console
$ kubectl get pod -n osba
NAME                                              READY     STATUS             RESTARTS   AGE
osba-open-service-broker-azure-7d7f455b79-c69v8   0/1       CrashLoopBackOff   15         53m
osba-redis-59656f567c-497dl                       1/1       Running            0          16m

$ kubectl logs osba-open-service-broker-azure-7d7f455b79-c69v8 -n osba
time="2018-03-13T12:08:31Z" level=info msg="Setting log level" logLevel=INFO
time="2018-03-13T12:08:31Z" level=info msg="Open Service Broker for Azure starting" commit=61f415e version=devel
time="2018-03-13T12:08:31Z" level=info msg="Sensitive instance and binding details will be encrypted" encryptionScheme=AES256
time="2018-03-13T12:08:31Z" level=info msg="API server is listening" address="http://0.0.0.0:8080"
time="2018-03-13T12:08:31Z" level=fatal msg="async engine stopped: error sending heartbeat for worker c88dea9f-f55f-4f32-b141-5729e867f48a: dial tcp: lookup osba-redis on 10.0.0.10:53: server misbehaving"
```

After checking on the health of the kube-system pods, kube-dns was not working. This caused osba 
to be unable to connect to its Redis instance. The problem was resolved by deleting the
kube-dns pods, so that Kubernetes would recreate them.

### I don't see all the services

Open Service Broker for Azure provides a number of services and each of these services is implemented by a separate module. The stability of individual modules is independent of overall broker stability and is ranked on a scale of `experimental`, `preview`, and `stable`. The broker can be configured to only load modules at or above a specified stability threshold. By default, the helm chart configures the broker to only load modules that are marked as `preview` or `stable`. This currently includes Azure Database for MySQL, Azure Database for PostgreSQL and Azure SQL Database. If you would like to use other services, you will need to add an additional flag to your helm install command:

```console
helm install azure/open-service-broker-azure --name osba --namespace osba \
  --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
  --set azure.tenantId=$AZURE_TENANT_ID \
  --set azure.clientId=$AZURE_CLIENT_ID \
  --set azure.clientSecret=$AZURE_CLIENT_SECRET \
  --set modules.minStability=EXPERIMENTAL
```

## "osba" is forbidden: not yet ready to handle request

[Service Catalog](https://github.com/kubernetes-incubator/service-catalog) is the software component that is used to integrate Open Service Broker for Azure with Kubernetes. Service Catalog currently works with Kubernetes version 1.9.0 and higher, so you will need a Kubernetes cluster that is version 1.9.0 or higher. If you receive this error message when trying to install Open Service Broker For Azure, check the version of your Kubernetes cluster. If you are running a version less than 1.9.0, please upgrade. If you are making a new cluster with AKS, you can specify the version with the `--kubernetes-version` flag like so:

```console
az aks create --resource-group osba --name my-aks-cluster --generate-ssh-keys --kubernetes-version 1.9.6
```

To see available versions, you can use the `az aks get-versions` command:

```console
$ az aks get-versions -l eastus -o table
KubernetesVersion    Upgrades
-------------------  -------------------------------------------------------------------------
1.9.6                None available
1.9.2                1.9.6
1.9.1                1.9.2, 1.9.6
1.8.11               1.9.1, 1.9.2, 1.9.6
1.8.10               1.8.11, 1.9.1, 1.9.2, 1.9.6
1.8.7                1.8.10, 1.8.11, 1.9.1, 1.9.2, 1.9.6
1.8.6                1.8.7, 1.8.10, 1.8.11, 1.9.1, 1.9.2, 1.9.6
1.8.2                1.8.6, 1.8.7, 1.8.10, 1.8.11, 1.9.1, 1.9.2, 1.9.6
1.8.1                1.8.2, 1.8.6, 1.8.7, 1.8.10, 1.8.11, 1.9.1, 1.9.2, 1.9.6
1.7.16               1.8.1, 1.8.2, 1.8.6, 1.8.7, 1.8.10, 1.8.11
1.7.15               1.7.16, 1.8.1, 1.8.2, 1.8.6, 1.8.7, 1.8.10, 1.8.11
1.7.12               1.7.15, 1.7.16, 1.8.1, 1.8.2, 1.8.6, 1.8.7, 1.8.10, 1.8.11
1.7.9                1.7.12, 1.7.15, 1.7.16, 1.8.1, 1.8.2, 1.8.6, 1.8.7, 1.8.10, 1.8.11
1.7.7                1.7.9, 1.7.12, 1.7.15, 1.7.16, 1.8.1, 1.8.2, 1.8.6, 1.8.7, 1.8.10, 1.8.11
```
