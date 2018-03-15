# OSBA Design

We can roughly split OSBA's architecture into three pieces:

1.  REST API frontend
2.  Stable Storage
3.  Asynchronous processing

The remainder of this document will overview each piece.

(Note: testing infrastructure is out of scope for this document)

# REST API Frontend

OSBA runs an HTTP server that provides a REST API that conforms to the
[Open Service Broker API specification](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md)

TODOs for this section:

* Specific version of the spec?
* More detail on the server & its handlers

# Stable Storage

Almost all the operations that OSBA exposes in the REST API (above) require
something to be stored in a database. OSBA provides a standard
[`Store`](./pkg/storage/store.go) interface, which you can implement with
any backend. Currently there are two implementations:

* Redis
  * Used by default on OSBA startup
* Memory
  * Used for testing

There are currently more storage backends planned, but none fully implemented
yet.

# Asynchronous Processing

OSBA supports several services that take a long time to provision and
deprovision, so it uses
[asynchronous operations](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#asynchronous-operations)
defined by the OSB API spec to do them in the background.

Behind the scenes, OSBA defines an asynchronous processing interface composed
of the following pieces:

* [Engine](./pkg/async/engine.go) - the top-level system that you can submit
  tasks (see below bullet point) to schedule for execution
* [Task](./pkg/async/task.go) - metadata about a specific job
  (see below bullet point) to run. The engine uses this metadata to do scheduling,
  restarts, and more
* [Job](./pkg/async/job.go) - the actual functionality to run asynchronously.
  Jobs are assigned a name and invoked by name in code that starts the async.
  work

The current, default asynchronous queue implementation uses Redis and the
[Redis reliable queue](https://danielkokott.wordpress.com/2015/02/14/redis-reliable-queue-pattern/)
pattern to implement reliable, restartable asynchronous jobs.
