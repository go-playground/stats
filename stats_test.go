package stats

import (
	"os"
	"testing"
	"time"

	. "gopkg.in/bluesuncorp/assert.v1"
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
		Port:   3008,
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

func TestBadListenUDP(t *testing.T) {

	serverConfig := &ServerConfig{
		Domain: "",
		Port:   3008,
	}

	server, err := NewServer(serverConfig)
	Equal(t, err, nil)

	PanicMatches(t, func() { server.Run() }, "listen udp :3008: bind: address already in use")
}

func TestBadServerAndEncodingFailure(t *testing.T) {
	serverConfig := &ServerConfig{
		Domain: "",
		Port:   3012,
	}

	localConfig := &ClientConfig{
		Domain: "",
		Port:   3013,
	}

	client, err := NewClient(localConfig, serverConfig)
	Equal(t, err, nil)
	go client.Run()

	<-time.Tick(time.Second * 1)

	client2, err := NewClient(localConfig, serverConfig)
	Equal(t, err, nil)
	PanicMatches(t, func() {
		client2.Run()
	}, "dial udp :3013->:3012: bind: address already in use")
}

func TestBadAddrs(t *testing.T) {
	serverConfig := &ServerConfig{
		Domain: "werfewfewfewf",
		Port:   -1000,
	}

	server, err := NewServer(serverConfig)
	Equal(t, err, nil)

	PanicMatches(t, func() { server.Run() }, "lookup udp/-1000: nodename nor servname provided, or not known")

	localConfig := &ClientConfig{
		Domain: "erferfergergerg",
		Port:   -2000,
	}

	client, err := NewClient(localConfig, serverConfig)
	Equal(t, err, nil)
	PanicMatches(t, func() { client.Run() }, "lookup udp/-1000: nodename nor servname provided, or not known")

	// set good server, but bad local remains
	serverConfig.Domain = ""
	serverConfig.Port = 3011

	client, err = NewClient(localConfig, serverConfig)
	Equal(t, err, nil)
	PanicMatches(t, func() { client.Run() }, "lookup udp/-2000: nodename nor servname provided, or not known")
}

func TestClientSendingData(t *testing.T) {

	serverConfig := &ServerConfig{
		Domain: "",
		Port:   3008,
	}

	localConfig := &ClientConfig{
		Domain: "",
		Port:   3009,
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
