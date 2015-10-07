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
