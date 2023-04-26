# haproxy-runtime-api-client

Go package for accessing the HA-Proxy Runtime API (stats socket)

The basic functionality handles sending a raw commands and getting the raw response back on command at a time.

## maintenance state

The main reason for creating this module is handling maintenance of backend servers when using a probe on the servers is not sufficient.

The module therefore contains a set of higher level functions for placing a backend server into maintenance state and back to ready state again. When placing a backend server into maintenance state this is done using the following process:

First the backend server state is set to drain.

The backend server is then monitored for number of concurrent connections to determine when the backend server can be placed into maintenance state. The state change occures either with the number of concurrent connections reaches 0 or when a given period of time has elapsed.
