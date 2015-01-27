Projekt
=======

Nothing interesting. Move on.

Installation
------------
```
git clone https://github.com/gswirski/sieci_projekt.git $GOPATH/src/sieci
cd $GOPATH/src/sieci
./install
```

Usage
---------------

* `$GOPATH/bin/master :2000 :2001` starts master process
* `$GOPATH/bin/worker :2000` starts worker process and connects to master
* `$GOPATH/bin/client :2001` starts client process and connects to master

Authors
-------

* Marcin Kostrzewa
* [Grzegorz Świrski](http://swirski.name)
