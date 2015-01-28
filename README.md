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

Authors
-------

* Marcin Kostrzewa
* [Grzegorz Świrski](http://swirski.name)
