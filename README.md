Randy - Cluster-Wide Random ID Generator
========================================

Randy is a microservice to generate truly random ID across multiple hosts,
or clusters. It does not use any database to store generated numbers, instead it
generates IDs algorithmically.

Randy ensures...
-------------------

1. **High-Performance** - ables to generate thousands of random IDs per second.

2. **High-Scalability** - allows to run multiple Randy instances on single host
   (using different ports), multiple hosts and multiple clusters.

3. **No centralized database** - common solutions are based on having some kind of
   centralized database (e.g. Redis). As database could be an another SPOF or point
   that is not so easy to scale or maintain, the goal was to avoid it.
   Randy stores it's state in process memory and files (used only in process initialization).

4. **Uniqueness of generated IDs** - Randy ensures that every generated ID is unique
   due to the algorithm and pattern of the generated IDs.
   
Requirements
------------

* [**Golang**](https://golang.org/doc/install)
* [**HAProxy**](http://www.haproxy.org/) - recommended to load-balance across multiple Randy instances

Running Randy
-------------

Running Randy is as easy as:

```
$ git clone git@github.com:scalebeat/randy.git
$ cd randy
$ go build
$ ./randy -port 8080 -hash raw

2016/12/11 13:38:08 Randy initialized, running on port 8080
```

**Note:** You can choose hash algo of: raw, md5, sha256 or sha512 - it will be used 
to hash the final id.

Keep in mind that you should **never use "raw" in production**
as it will expose your hostname and port on which Randy is running.

Now you can try it by running:

```
$ curl -i http://localhost:8080/

HTTP/1.1 200 OK
Date: Sun, 11 Dec 2016 12:37:15 GMT
Content-Length: 48
Content-Type: text/plain; charset=utf-8

1481535453344192624:16:17:localhost:8080
```

Running performance test
------------------------

You can easily run performance test by using [Apache Benchmark](http://httpd.apache.org/docs/2.4/programs/ab.html).
Sample test result from a very low-end setup (VPS with 1 vCore 2,4 GHz 2 GB RAM):

```
$ ab -n 5000 -c 800 http://localhost:8080/

...
Requests per second:    5422.84 [#/sec] (mean)
Time per request:       147.524 [ms] (mean)
Time per request:       0.184 [ms] (mean, across all concurrent requests)

```

License
-------

Copyright (c) 2016, Antoni Orfin - [Scalebeat](https://scalebeat.com).
Licensed under the [BSD 3-Clause](LICENSE)
