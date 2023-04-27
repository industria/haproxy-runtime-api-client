package haproxy

import (
	"log"
	"testing"
	"time"
)

func TestClientDial(t *testing.T) {
	client, err := NewClient("tcp://localhost:9999")
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}

	resp, err := client.Execute("show servers conn")
	if err != nil {
		t.Fatalf("unable to execute show servers conn : %v", err)
	}

	log.Printf("RESPONSE: %s", resp)

}

func TestSetServerState(t *testing.T) {
	client, err := NewClient("tcp://localhost:9999")
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}

	err = client.SetServerState("indexws", "iws01", ServerStateDrain)
	if err != nil {
		t.Fatalf("drain failed : %v", err)
	}

	time.Sleep(10 * time.Second)

	err = client.SetServerState("indexws", "iws01", ServerStateMaint)
	if err != nil {
		t.Fatalf("maint failed : %v", err)
	}

	time.Sleep(10 * time.Second)

	err = client.SetServerState("indexws", "iws01", ServerStateReady)
	if err != nil {
		t.Fatalf("ready failed : %v", err)
	}
}

func TestShowServersState(t *testing.T) {
	client, err := NewClient("tcp://localhost:9999")
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}

	resp, err := client.ShowServersState()
	if err != nil {
		t.Fatalf("failed show state : %v", err)
	}
	//log.Printf("%+v", resp)

	for _, state := range resp {
		log.Printf("%s/%s %d %d", state.BeName, state.SrvName, state.SrvAdminState, state.SrvOpState)
	}

}
