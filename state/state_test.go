package state

import "testing"

func TestStateElementsToServerState(t *testing.T) {
	elements := []string{"4", "indexws", "1", "iws01", "172.24.21.40", "2", "0", "1", "1", "776", "15", "3", "4", "6", "0", "0", "0", "-", "8080", "-", "0", "0", "-", "-", "0"}
	state := stateElementsToServerState(elements)
	if state.BeId != 4 {
		t.Fatalf("BeId not 4")
	}
	if state.BeName != "indexws" {
		t.Fatalf("BeName not indexws")
	}
	if state.SrvId != 1 {
		t.Fatalf("SrvId not 1")
	}
	if state.SrvName != "iws01" {
		t.Fatalf("SrvName not iws01")
	}
	if state.SrvAddr != "172.24.21.40" {
		t.Fatalf("SrvAddr not 172.24.21.40")
	}
	if state.SrvOpState != OperationalStateRunning {
		t.Fatalf("SrvOpState not 2 (OperationalStateRunning)")
	}
	if state.SrvAdminState != 0 {
		t.Fatalf("SrvAdminState not 0")
	}
	if state.SrvUWeight != 1 {
		t.Fatalf("SrvUWeight not 1")

	}
	if state.SrvIWeight != 1 {
		t.Fatalf("SrvIWeight not 1")
	}
	if state.SrvTimeSinceLastChange != 776 {
		t.Fatalf("SrvTimeSinceLastChange not 776")
	}
	if state.SrvCheckStatus != 15 {
		t.Fatalf("SrvCheckStatus not 15")
	}
	if state.SrvCheckResult != CheckResultPassed {
		t.Fatalf("SrvCheckStatus not 3 (CheckResultPassed)")
	}
	if state.SrvCheckHealth != 4 {
		t.Fatalf("SrvCheckStatus not 4")
	}
	if state.SrvCheckState != 6 { // TODO: should check the bitmask
		t.Fatalf("SrvCheckState not 6")
	}
	if state.SrvAgentState != 0 {
		t.Fatalf("SrvAgentState not 0")
	}
	if state.BkFForcedId {
		t.Fatalf("BkFForcedId not false")
	}
	if state.SrvFForcedId {
		t.Fatalf("SrvFForcedId not false")
	}
	if state.SrvFQDN != "-" {
		t.Fatalf("SrvFQDN not -")
	}
	if state.SrvPort != "8080" {
		t.Fatalf("SrvPort not 8080")
	}
	if state.SrvRecord != "-" {
		t.Fatalf("SrvRecord not -")
	}
	if state.SrvUseSSL {
		t.Fatalf("SrvUseSSL not false")
	}
	if state.SrvCheckPort != "0" {
		t.Fatalf("SrvCheckPort not 0")
	}
	if state.SrvCheckAddr != "-" {
		t.Fatalf("SrvCheckAddr not -")
	}
	if state.SrvAgentAddr != "-" {
		t.Fatalf("SrvAgentAddr not -")
	}
	if state.SrvAgentPort != "0" {
		t.Fatalf("SrvAgentPort not 0")
	}
}
