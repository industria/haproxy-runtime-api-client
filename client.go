package haproxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/industria/haproxy-runtime-api-client/stat"
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

// execute a server state change command
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

// get the server state for all backend
// show servers state [<backend>]
func (rc *RuntimeClient) ShowServersState() ([]state.ServerState, error) {
	command := "show servers state"
	resp, err := rc.Execute(command)
	if err != nil {
		return nil, err
	}
	return state.ParseShowServersState(resp)
}

// get stat counters using the command show stat
func (rc *RuntimeClient) ShowStat() ([]stat.StatCounters, error) {
	command := "show stat"
	resp, err := rc.Execute(command)
	if err != nil {
		return nil, err
	}
	return stat.ParseShowStat(resp)
}

//	place server into maintenance state with a previous drain operation
//
// this funcation will place a server into maintenance state by first
// placing the server in drain state waiting for the currenct conntations
// going to 0 or a timeout is reached and the server is forced into
// maintenance state regardless of the number of connections.
// a context with timeout should be used to avoid waiting forever for the draining to complete
// the draining can take a long time if there is an active persistent connection
func (rc *RuntimeClient) ServerMaintenance(ctx context.Context, backend, server string) error {
	// start by setting the backend server to draining
	if err := rc.SetServerState(backend, server, ServerStateDrain); err != nil {
		return err
	}

	// time for allowing a pause between checking if draining the backend server is complete
	timer := time.NewTimer(time.Millisecond * 10)
	defer timer.Stop()

	// check for complete or timeout
	for {
		select {
		case <-ctx.Done():
			log.Println("timeout - force the server to maint state")
			return rc.SetServerState(backend, server, ServerStateMaint)
		case <-timer.C:
			completed, err := rc.drainingComplet(backend, server)
			if err != nil {
				return err
			}
			if completed {
				return rc.SetServerState(backend, server, ServerStateMaint)
			}
			timer.Reset(time.Millisecond * 10)
		}
	}
}

func (rc *RuntimeClient) drainingComplet(backend, server string) (bool, error) {
	cs, err := rc.ShowStat()
	if err != nil {
		return false, err
	}

	for _, c := range cs {
		if c.PxName == backend && c.SvName == server {
			log.Printf("connections for %s/%s %d", backend, server, c.Scur)
			return c.Scur == 0, nil

		}
	}
	return false, nil
}
