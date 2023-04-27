// package for working with server state
package state

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
)

type OperationalState int

const (
	OperationalStateStopped  OperationalState = 0 // 0 = SRV_ST_STOPPED (The server is down.)
	OperationalStateStarting OperationalState = 1 // 1 = SRV_ST_STARTING (The server is warming up (up but throttled).)
	OperationalStateRunning  OperationalState = 2 // 2 = SRV_ST_RUNNING (The server is fully up.)
	OperationalStateStopping OperationalState = 3 // 3 = SRV_ST_STOPPING (The server is up but soft-stopping (eg: 404).)
)

type AdminState uint

const (
	AdminStateForcedMaintenance     AdminState = 0x01 // 0x01 = SRV_ADMF_FMAINT (The server was explicitly forced into maintenance.)
	AdminStateInheritedMaintenance  AdminState = 0x02 // 0x02 = SRV_ADMF_IMAINT (The server has inherited the maintenance status from a tracked server.)
	AdminStateConfiguredMaintenance AdminState = 0x04 // 0x04 = SRV_ADMF_CMAINT (The server is in maintenance because of the configuration.)
	AdminStateForcedDrain           AdminState = 0x08 // 0x08 = SRV_ADMF_FDRAIN (The server was explicitly forced into drain state.)
	AdminStateInheritedDrain        AdminState = 0x10 // 0x10 = SRV_ADMF_IDRAIN (The server has inherited the drain status from a tracked server.)
	AdminStateResolutionMaintenance AdminState = 0x20 // 0x20 = SRV_ADMF_RMAINT (The server is in maintenance because of an IP address resolution failure.)
	AdminStateHostMaintenance       AdminState = 0x40 // 0x40 = SRV_ADMF_HMAINT (The server FQDN was set from stats socket.)
)

type CheckResult int

const (
	CheckResultUnknown  CheckResult = 0 // 0 = CHK_RES_UNKNOWN (Initialized to this by default.)
	CheckResultNeutral  CheckResult = 1 // 1 = CHK_RES_NEUTRAL (Valid check but no status information.)
	CheckResultFailed   CheckResult = 2 // 2 = CHK_RES_FAILED (Check failed.)
	CheckResultPassed   CheckResult = 3 // 3 = CHK_RES_PASSED (Check succeeded and server is fully up again.)
	CheckResultCondpass CheckResult = 4 // 4 = CHK_RES_CONDPASS (Check reports the server doesn't want new sessions.)
)

type CheckState uint

const (
	CheckStateInProgress CheckState = 0x01 // 0x01 = CHK_ST_INPROGRESS (A check is currently running.)
	CheckStateConfigured CheckState = 0x02 // 0x02 = CHK_ST_CONFIGURED (This check is configured and may be enabled.)
	CheckStateEnabled    CheckState = 0x04 // 0x04 = CHK_ST_ENABLED (This check is currently administratively enabled.)
	CheckStatePaused     CheckState = 0x08 // 0x08 = CHK_ST_PAUSED (Checks are paused because of maintenance (health only).)
	CheckStateAgent      CheckState = 0x10 // 0x10 = CHK_ST_AGENT (Check is an agent check (otherwise it's a health check).) this is specific for srv_agent_state
)

// Parse Runtime API response for the command command: show servers state [<backend>]
// Reference: http://docs.haproxy.org/2.6/management.html#9.3-show%20servers%20state
type ServerState struct {
	BeId                   int              // be_id:                       Backend unique id.
	BeName                 string           // be_name:                     Backend label.
	SrvId                  int              // srv_id:                      Server unique id (in the backend).
	SrvName                string           // srv_name:                    Server label.
	SrvAddr                string           // srv_addr:                    Server IP address.
	SrvOpState             OperationalState // srv_op_state:                Server operational state (UP/DOWN/...).
	SrvAdminState          AdminState       // srv_admin_state:             Server administrative state (MAINT/DRAIN/...). The state is actually a mask of values
	SrvUWeight             int              // srv_uweight:                 User visible server's weight.
	SrvIWeight             int              // srv_iweight:                 Server's initial weight.
	SrvTimeSinceLastChange int              // srv_time_since_last_change:  Time since last operational change.
	SrvCheckStatus         int              // srv_check_status:            Last health check status.
	SrvCheckResult         CheckResult      // srv_check_result (Last check result (FAILED/PASSED/...).)
	SrvCheckHealth         int              // srv_check_health (Checks rise / fall current counter.)
	SrvCheckState          CheckState       // srv_check_state (tate of the check (ENABLED/PAUSED/...). The state is actually a mask of values)
	SrvAgentState          CheckState       // srv_agent_state: State of the agent check (ENABLED/PAUSED/...). This state uses the same mask values as srv_check_state, adding this specific one
	BkFForcedId            bool             // bk_f_forced_id:              Flag to know if the backend ID is forced by configuration.
	SrvFForcedId           bool             // srv_f_forced_id:             Flag to know if the server's ID is forced by configuration.
	SrvFQDN                string           // srv_fqdn:                    Server FQDN.
	SrvPort                string           // srv_port:                    Server port.
	SrvRecord              string           // srvrecord:                   DNS SRV record associated to this SRV.
	SrvUseSSL              bool             // srv_use_ssl:                 use ssl for server connections.
	SrvCheckPort           string           // srv_check_port:              Server health check port.
	SrvCheckAddr           string           // srv_check_addr:              Server health check address.
	SrvAgentAddr           string           // srv_agent_addr:              Server health agent address.
	SrvAgentPort           string           // srv_agent_port:              Server health agent port.
}

// parse show servers state response as defined in http://docs.haproxy.org/2.6/management.html#9.3-show%20servers%20state
func ParseShowServersState(response []byte) ([]ServerState, error) {
	// first line is version and must be 1 (49)
	if response[0] != 49 {
		return nil, fmt.Errorf("show servers state version was not 1")
	}
	// remove the first version line [49 10] 1 and line feed
	response = response[2:]

	r := csv.NewReader(bytes.NewReader(response))
	r.Comma = ' '   // white space separated colums in the runtime API responses
	r.Comment = '#' // lines with # is comments in the runtime API responses

	all, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	states := make([]ServerState, 0, len(all))
	for _, row := range all {
		states = append(states, stateElementsToServerState(row))
	}

	return states, nil
}

// map the elements from a response line to the ServerState struct. If the elements can not be mapped it will panic.
func stateElementsToServerState(elements []string) ServerState {
	return ServerState{
		BeId:                   atoi(elements[0]),
		BeName:                 elements[1],
		SrvId:                  atoi(elements[2]),
		SrvName:                elements[3],
		SrvAddr:                elements[4],
		SrvOpState:             OperationalState(atoi(elements[5])),
		SrvAdminState:          AdminState(atoui(elements[6])),
		SrvUWeight:             atoi(elements[7]),
		SrvIWeight:             atoi(elements[8]),
		SrvTimeSinceLastChange: atoi(elements[9]),
		SrvCheckStatus:         atoi(elements[10]),
		SrvCheckResult:         CheckResult(atoi(elements[11])),
		SrvCheckHealth:         atoi(elements[12]),
		SrvCheckState:          CheckState(atoui(elements[13])),
		SrvAgentState:          CheckState(atoui(elements[14])),
		BkFForcedId:            atob(elements[15]),
		SrvFForcedId:           atob(elements[16]),
		SrvFQDN:                elements[17],
		SrvPort:                elements[18],
		SrvRecord:              elements[19],
		SrvUseSSL:              atob(elements[20]),
		SrvCheckPort:           elements[21],
		SrvCheckAddr:           elements[22],
		SrvAgentAddr:           elements[23],
		SrvAgentPort:           elements[24],
	}
}

func atoi(s string) int {
	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return x
}

func atoui(s string) uint {
	x, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(x)
}

func atob(s string) bool {
	return atoi(s) != 0
}
