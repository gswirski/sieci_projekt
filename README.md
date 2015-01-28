Projekt
=======

Nothing interesting. Move on.

Installation
------------

1. Install Go lang according to documentation on http://golang.org
2. Export `GOPATH` env variable
3. Run following commands:

    ```
    git clone https://github.com/gswirski/sieci_projekt.git $GOPATH/src/sieci
    cd $GOPATH/src/sieci
    ./install
    ```

Usage
---------------

### Master process
Master processes manage jobs to execute. You can have as many master processes as you want, e.g. to support failover in clients.

To run a master process, execute `$GOPATH/bin/master` with two arguments: TCP address for worker communication and TCP address for client communication. Example: `$GOPATH/bin/master :2000 :2001`.

### Worker process
Worker proceses execute received code and return its result. A single worker can respond to many masters to better balance workload. Run `$GOPATH/bin/worker` with a list of masters to start a worker process. Example: `$GOPATH/bin/worker :2000 :2002 :2004` will communicate with masters on ports 2000, 2002 and 2004.

### Client process
Client processes send code to execute on cluster. They connect with master, wait for a response and die. To start a client process, run `$GOPATH/bin/client master_address file_with_code.py`. Example: `$GOPATH/bin/client 2001 code.py`.

Protocols
---------

Projekt uses two very simple protocols. One to manage workload between master and workers, the other to send data over client-master and master-worker connections.

### Data exchange protocol

After TCP connection is initialized, uploading side sends the following message (there must be a newline after _terminator_)

```
ENDSEQ terminator
code
terminator
```

After the _code_ is received, server responds with `RECEIVED` or `ERROR` message and, when results are available, begins another transfer with the same syntax:

```
ENDSEQ terminator
results
terminator
```

Server does not wait for clients to acknowledge success.

### Worker orchestration

When master process opens a connection with client process, it tries to find available worker. It does so by sending `AVAILABLE` message to every worker that is not executing a job known to this master (it is expected that master does not query workers that execute its own jobs at the time). When a worker is available, it responds with `READY` message and locks on that master. Master selects one of the workers and initializes a transfer using data exchange protocol. The rest of the workers receive `ROLLBACK` message to unblock them.

Implementation
--------------

Authors
-------

* Marcin Kostrzewa
* [Grzegorz Świrski](http://swirski.name)
