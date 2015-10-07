package stats

import (
	"os"
	"testing"
	"time"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

func TestMain(m *testing.M) {
	// setup
	//
	// setup server for tests
	config := &ServerConfig{
		Domain: "",
		Port:   3010,
	}

	server, err := NewServer(config)
	if err != nil {
		panic(err)
	}

	go server.Run()

	<-time.Tick(time.Second * 1)

	os.Exit(m.Run())

	// teardown
}

func TestClientSendingData(t *testing.T) {

	serverConfig := &ServerConfig{
		Domain: "",
		Port:   3010,
	}

	localConfig := &ClientConfig{
		Domain:           "",
		Port:             3011,
		PollInterval:     1000,
		LogHostInfo:      true,
		LogCPUInfo:       true,
		LogTotalCPUTimes: true,
		LogPerCPUTimes:   true,
		LogMemory:        true,
		LogGoMemory:      true,
	}

	client, err := NewClient(localConfig, serverConfig)
	Equal(t, err, nil)

	go client.Run()

	ticker := time.NewTicker(time.Second * 1)
	i := 0
	for range ticker.C {
		if i == 1 {
			client.Stop()
			return
		}
		i++
	}
}
