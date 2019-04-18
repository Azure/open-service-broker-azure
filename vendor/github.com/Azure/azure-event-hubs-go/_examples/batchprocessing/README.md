# Batch Processing

To process batches of events at one time a handler needs to implement both an event handler function as well as the `CheckpointPersister` interface found in `github.com/Azure/azure-amqp-common-go/persist`. 

## Running this example

Set the following environment variables `EVENTHUB_PARTITIONID`, `EVENTHUB_CONSUMERGROUP`, `EVENTHUB_NAMESPACE`, `EVENTHUB_NAME`, `EVENTHUB_KEY_NAME` and `EVENTHUB_KEY_VALUE` and run the executable.

