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
	log.Printf("%+v", resp)
}

/*
1
# be_id be_name srv_id srv_name srv_addr srv_op_state srv_admin_state srv_uweight srv_iweight srv_time_since_last_change srv_check_status srv_check_result srv_check_health srv_check_state srv_agent_state bk_f_forced_id srv_f_forced_id srv_fqdn srv_port srvrecord srv_use_ssl srv_check_port srv_check_addr srv_agent_addr srv_agent_port
4 indexws 1 iws01 172.24.21.40 2 0 1 1 96 15 3 4 6 0 0 0 - 8080 - 0 0 - - 0
4 indexws 2 iws02 172.24.21.32 2 0 1 1 96 15 3 4 6 0 0 0 - 8080 - 0 0 - - 0
4 indexws 3 iws03 172.24.21.33 2 0 1 1 96 15 3 4 6 0 0 0 - 8080 - 0 0 - - 0
5 solr 1 solr01 192.168.220.4 2 0 1 1 83 15 0 4 7 0 0 0 solr 8983 - 0 0 - - 0
*/
