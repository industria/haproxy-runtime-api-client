package haproxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/industria/haproxy-runtime-api-client/state"
)

type RuntimeClient struct {
	network string
	address string
}

// create new RuntimeClient for stats socket with address
// where the address is expressed as a URI and can
// be either unix://path or tcp://address:port
func NewClient(uri string) (*RuntimeClient, error) {
	var network string
	var address string
	if strings.HasPrefix(uri, "unix://") {
		network = "unix"
		address = uri[7:]
	} else if strings.HasPrefix(uri, "tcp://") {
		network = "tcp"
		address = uri[6:]
	} else {
		return nil, fmt.Errorf("address [%s] must start with unix:// or tcp://", address)
	}

	return &RuntimeClient{
		network: network,
		address: address,
	}, nil
}

// execute a HA-Proxy runtime api command
// https://cbonte.github.io/haproxy-dconv/2.5/management.html#9.3
// result is the raw response from the API
func (rc *RuntimeClient) Execute(command string) ([]byte, error) {
	conn, err := net.Dial(rc.network, rc.address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cmd := []byte(command + "\n")

	log.Printf("%s", cmd)

	if _, err := conn.Write(cmd); err != nil {
		return nil, fmt.Errorf("unable to send command: %v", cmd)
	}

	resp, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %v", err)
	}

	return resp, nil
}

type ServerState string

const (
	ServerStateDrain ServerState = "drain"
	ServerStateReady ServerState = "ready"
	ServerStateMaint ServerState = "maint"
)

// execuie a server state change command
// set server <backend>/<server> state [ ready | drain | maint ]
func (rc *RuntimeClient) SetServerState(backend, server string, state ServerState) error {
	command := fmt.Sprintf("set server %s/%s state %s", backend, server, state)
	resp, err := rc.Execute(command)
	if err != nil {
		return err
	}

	// a response other than a line feed would be an error on these state changes
	if len(resp) != 1 && resp[0] != 10 {
		return fmt.Errorf("changing %s/%s to %s failed with: %s", backend, server, state, resp)
	}

	return nil
}

// show servers state [<backend>]
func (rc *RuntimeClient) ShowServersState() ([]state.ServerState, error) {
	//command := fmt.Sprintf("set server %s/%s state %s", backend, server, state)
	command := fmt.Sprintf("show servers state")

	resp, err := rc.Execute(command)
	if err != nil {
		return nil, err
	}

	return state.ParseShowServersState(resp)
}
