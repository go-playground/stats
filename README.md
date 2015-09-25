Package stats
=============

Package stats allows for gathering of statistics regarding your Go application and system it is running on and
sent them via UDP to a server where you can do whatever you wish to the stats; display, store in database or
send off to a logging service.

## NOTE: This is in Beta testing

###### We currently gather the following Go related information:

* # of Garbabage collects
* Last Garbage Collection
* Last Garbage Collection Pause Duration
* Memory Allocated
* Memory Heap Allocated
* Memory Heap System Allocated
* Go version
* Number of goroutines
* HTTP request logging; when implemented via middleware

###### And the following System Information:

* Host Information; hostname, OS....
* CPU Information; tpye, model, # of cores...
* Total CPU Timings
* Per Core CPU Timings
* Memory + Swap Information