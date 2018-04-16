# FAQ
_What to do when things go sideways_

# Troubleshooting
Here are some issues that people have run into in the past:

## OSBA Pod in CrashLoopBackOff, can't lookup osba-redis

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

## I don't see all the services

The helm chart is configured to set the minimum stability level of the broker to `preview`,
meaning that only services that are marked `preview` or `stable` will be available from the broker.
This is currently Azure Database for MySQL, Azure Database for PostgreSQL and Azure SQL Database.
If you would like to use other services, you will need to add an additional flag to your helm install command:

```console
helm install azure/open-service-broker-azure --name osba --namespace osba \
  --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
  --set azure.tenantId=$AZURE_TENANT_ID \
  --set azure.clientId=$AZURE_CLIENT_ID \
  --set azure.clientSecret=$AZURE_CLIENT_SECRET \
  --set modules.minStability=EXPERIMENTAL
```