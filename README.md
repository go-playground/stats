Package stats
=============

[![Build Status](https://semaphoreci.com/api/v1/projects/b5c5d4f1-ec31-441f-a6d2-88f4e73105df/563578/badge.svg)](https://semaphoreci.com/joeybloggs/stats)
[![GoDoc](https://godoc.org/gopkg.in/go-playground/stats.v1?status.svg)](https://godoc.org/gopkg.in/go-playground/stats.v1)

Package stats allows for gathering of statistics regarding your Go application and system it is running on and
sent them via UDP to a server where you can do whatever you wish to the stats; display, store in database or
send off to a logging service.

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
* CPU Information; type, model, # of cores...
* Total CPU Timings
* Per Core CPU Timings
* Memory + Swap Information

Installation
------------

Use go get.

	go get gopkg.in/go-playground/stats.v1

or to update

	go get -u gopkg.in/go-playground/stats.v1

Then import the validator package into your own code.

	import "gopkg.in/go-playground/stats.v1"

#### Example
Server
```go
package main

import (
	"fmt"

	"gopkg.in/go-playground/stats.v1"
)

func main() {

	config := &stats.ServerConfig{
		Domain: "",
		Port:   3008,
		Debug:  false,
	}

	server, err := stats.NewServer(config)
	if err != nil {
		panic(err)
	}

	for stat := range server.Run() {

		// calculate CPU times
		// totalCPUTimes := stat.CalculateTotalCPUTimes()
		// perCoreCPUTimes := stat.CalculateCPUTimes()

		// Do whatever you want with the data
		// * Save to database
		// * Stream elsewhere
		// * Print to console
		//

		fmt.Println(stat)
	}
}
```

Client
```go
package main

import (
	"fmt"
	"net/http"
	"runtime"

	"gopkg.in/go-playground/stats.v1"
)

var statsClient *stats.ClientStats

func main() {

	serverConfig := &stats.ServerConfig{
		Domain: "remoteserver",
		Port:   3008,
		Debug:  false,
	}

	clientConfig := &stats.ClientConfig{
		Domain:           "",
		Port:             3009,
		PollInterval:     1000,
		Debug:            false,
		LogHostInfo:      true,
		LogCPUInfo:       true,
		LogTotalCPUTimes: true,
		LogPerCPUTimes:   true,
		LogMemory:        true,
		LogGoMemory:      true,
	}

	client, err := stats.NewClient(clientConfig, serverConfig)
	if err != nil {
		panic(err)
	}

	go client.Run()

	// if you want to capture HTTP requests in a middleware you'll have
	// to have access to the client.
	// see below for middleware example
	statsClient = client
}

// LoggingRecoveryHandler ...
//
//
// Middleware example for capturing HTTP Request info
// NOTE: this is standard go middlware, but could be adapted to other types/styles easily
//
func LoggingRecoveryHandler(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		// log incoming request
		logReq := statsClient.NewHTTPRequest(w, r)

		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)

				// log failure
				logReq.Failure(fmt.Sprintf("%s\n%s", err, trace[:n]))

				http.Error(w, "Friendly error message", 500)
				return
			}
		}()

		next.ServeHTTP(logReq.Writer(), r)

		// log completion
		logReq.Complete()
	}

	return http.HandlerFunc(fn)
}
```

License
------
Distributed under MIT License, please see license file in code for more details.
