# SQL Examples

This repository contains various tests/spikes with SQL databases in Go.

## Connect

The idea we're chasing here is to use be able to
connect to a MySQL instance regardless where the
code is running (local, dockerized, GCP, AWS, Azure,
you name it) and without additional infrastructure
components.

1. You change the connection string only.
2. The code doesn't change.
3. The infrastructure doesn't change. E.g. you don't
   run a CloudSQL Proxy on GCP. Everything is self-contained.

Notice that this doesn't work in a Dockerized environment yet.
We need to [wait for this issue](https://github.com/google/go-cloud/issues/2644). 

## Nanos

The code will create a `people` table in a database with a `Created`
column that has nanosecond precision.

# License

MIT.
