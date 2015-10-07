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
