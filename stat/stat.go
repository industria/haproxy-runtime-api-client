// package for working with server stat counters
package stat

import (
	"bytes"
	"encoding/csv"
	"strconv"
)

// represents http://docs.haproxy.org/2.6/management.html#9.1
// The letters in brackets are L (Listeners), F (Frontends), B (Backends), and S (Servers) indicating the type which may have a value
// below is made for the following response header - not all are mapped - last included is uweight
// pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime,agent_status,agent_code,agent_duration,check_desc,agent_desc,check_rise,check_fall,check_health,agent_rise,agent_fall,agent_health,addr,cookie,mode,algo,conn_rate,conn_rate_max,conn_tot,intercepted,dcon,dses,wrew,connect,reuse,cache_lookups,cache_hits,srv_icur,src_ilim,qtime_max,ctime_max,rtime_max,ttime_max,eint,idle_conn_cur,safe_conn_cur,used_conn_cur,need_conn_est,uweight,agg_server_status,agg_server_check_status,agg_check_status,-,ssl_sess,ssl_reused_sess,ssl_failed_handshake,h2_headers_rcvd,h2_data_rcvd,h2_settings_rcvd,h2_rst_stream_rcvd,h2_goaway_rcvd,h2_detected_conn_protocol_errors,h2_detected_strm_protocol_errors,h2_rst_stream_resp,h2_goaway_resp,h2_open_connections,h2_backend_open_streams,h2_total_connections,h2_backend_total_streams,h1_open_connections,h1_open_streams,h1_total_connections,h1_total_streams,h1_bytes_in,h1_bytes_out,h1_spliced_bytes_in,h1_spliced_bytes_out,
type StatCounters struct {
	PxName        string // 0: Proxy name [LFBS]
	SvName        string // 1: service name (FRONTEND for frontend, BACKEND for backend, any name for server/listener) [LFBS]
	Qcur          uint32 // 2: current queued requests. For the backend this reports the number queued without a server assigned. [..BS]
	Qmax          uint32 // 3: max value of qcur [..BS]
	Scur          uint32 // 4: current sessions [LFBS]
	Smax          uint32 // 5: max sessions [LFBS]
	Slim          uint32 // 6: configured session limit [LFBS]
	Stot          uint64 // 7: cumulative number of sessions [LFBS]
	Bin           uint64 // 8: bytes in [LFBS]
	Bout          uint64 // 9: bytes out [LFBS]
	Dreg          uint64 // 10: requests denied because of security concerns. [LFB.]
	Dresp         uint64 // 11: responses denied because of security concerns. [LFBS]
	Ereg          uint64 // 12: request errors. [LF..]
	Econ          uint64 // 13:  number of requests that encountered an error trying to connect to a backend server. [..BS] The backend stat is the sum of the stat for all servers of that backend, plus any connection errors not associated with a particular server (such as the backend having no active servers)
	Eresp         uint64 // 14: response errors. [..BS]
	Wretr         uint64 // 15: number of times a connection to a server was retried. [..BS]
	Wredis        uint64 // 16: number of times a request was redispatched to another server. The server value counts the number of times that server was switched away from. [..BS]
	Status        string // 17: status (UP/DOWN/NOLB/MAINT/MAINT(via)/MAINT(resolution)...) [LFBS]
	Weight        uint32 // 18: total effective weight (backend), effective weight (server) [..BS]
	Act           uint32 // 19: number of active servers (backend), server is active (server) [..BS]
	Bck           uint32 // 20: number of backup servers (backend), server is backup (server) [..BS]
	ChkFail       uint64 // 21: number of failed checks. (Only counts checks failed when the server is up.) [...S]
	ChkDown       uint64 // 22: number of UP->DOWN transitions. The backend counter counts transitions to the whole backend being down, rather than the sum of the counters for each server. [..BS]
	LastChg       uint32 // 23: number of seconds since the last UP<->DOWN transition [..BS]
	Downtime      uint32 // 24: total downtime (in seconds). The value for the backend is the downtime for the whole backend, not the sum of the server downtime. [..BS]
	Qlimit        uint64 // 25: configured maxqueue for the server, or nothing in the value is 0 (default, meaning no limit) [...S]
	Pid           uint32 // 26: process id (0 for first instance, 1 for second, ...) [LFBS]
	Iid           uint32 // 27: unique proxy id [LFBS]
	Sid           uint32 // 28: server id (unique inside a proxy)[L..S]
	Throttle      uint64 // 29: current throttle percentage for the server, when slowstart is active, or no value if not in slowstart.
	Lbtot         uint64 // 30: total number of times a server was selected, either for new sessions, or when re-dispatching. The server counter is the number of times that server was selected. [..BS]
	Tracked       uint32 // 31: id of proxy/server if tracking is enabled. [...S]
	Type          uint32 // 32: (0=frontend, 1=backend, 2=server, 3=socket/listener) [LFBS]
	Rate          uint32 // 33: number of sessions per second over last elapsed second [.FBS]
	RateLim       uint32 // 34: configured limit on new sessions per second [.F..]
	RateMax       uint32 // 35: max number of new sessions per second [.FBS]
	CheckStatus   string // 36: status of last health check. [...S] Notice: If a check is currently running, the last known status will be reported, prefixed with "* ". e. g. "* L7OK".
	CheckCode     uint32 // 37: layer5-7 code, if available HTTP/SMTP/LDAP status code reported by the latest server health check [...S]
	CheckDuration uint64 // 38: time in ms took to finish last health check [...S] Total duration of the latest server health check, in milliseconds
	Hrsp1xx       uint64 // 39: http responses with 1xx code [.FBS]
	Hrsp2xx       uint64 // 40: http responses with 2xx code [.FBS]
	Hrsp3xx       uint64 // 41: http responses with 3xx code [.FBS]
	Hrsp4xx       uint64 // 42: http responses with 4xx code [.FBS]
	Hrsp5xx       uint64 // 43: http responses with 5xx code [.FBS]
	HrspOther     uint64 // 44: http responses with other codes (protocol error) [.FBS]
	HanaFail      uint64 // 45: failed health checks details [...S]
	ReqRate       uint32 // 46: HTTP requests per second over last elapsed second [.F..]
	ReqRateMax    uint32 // 47: max number of HTTP requests per second observed [.F..]
	ReqTot        uint64 // 48: total number of HTTP requests received [.FB.]
	CliAbrt       uint64 // 49: number of data transfers aborted by the client [..BS]
	SrvAbrt       uint64 // 50: number of data transfers aborted by the server (inc. in eresp) [..BS]
	CompIn        uint64 // 51: number of HTTP response bytes fed to the compressor [.FB.]
	CompOut       uint64 // 52: number of HTTP response bytes emitted by the compressor [.FB.]
	CompByp       uint64 // 53: number of bytes that bypassed the HTTP compressor (CPU/BW limit) [.FB.]
	CompRsp       uint64 // 54: number of HTTP responses that were compressed [.FB.]
	LastSess      int    // 55: number of seconds since last session assigned to server/backend [..BS]
	LastChk       string // 56: last health check contents or textual error [...S]
	LastAgt       string // 57: last agent check contents or textual error [...S]
	Qtime         uint32 // 58: the average queue time in ms over the 1024 last requests [..BS]
	Ctime         uint32 // 59: the average connect time in ms over the 1024 last requests [..BS]
	Rtime         uint32 // 60: the average response time in ms over the 1024 last requests (0 for TCP) [..BS]
	Ttime         uint32 // 61: the average total session time in ms over the 1024 last requests [..BS]
	AgentStatus   string // 62: status of last agent check [...S]
	AgentCode     uint32 // 63: numeric code reported by agent if any (unused for now) [...S]
	AgentDuration uint64 // 64: time in ms taken to finish last check [...S]
	CheckDesc     string // 65: short human-readable description of check_status [...S]
	AgentDesc     string // 66: short human-readable description of agent_status [...S]
	CheckRise     uint32 // 67: server's "rise" parameter used by checks, number of successful health checks before declaring a server UP (server 'rise' setting) [...S]
	CheckFall     uint32 // 68: server's "fall" parameter used by checks, number of failed health checks before declaring a server DOWN (server 'fall' setting) [...S]
	CheckHealth   uint32 // 69: server's health check value between 0 and rise+fall-1, current server health check level (0..fall-1=DOWN, fall..rise-1=UP) [...S]
	AgentRise     uint32 // 70: agent's "rise" parameter, normally 1 [...S]
	AgentFall     uint32 // 71: agent's "fall" parameter, normally 1 [...S]
	AgentHealth   uint32 // 72: agent's health parameter, between 0 and rise+fall-1 [...S]
	Addr          string // 73: address:port or "unix". IPv6 has brackets around the address. [L..S]
	Cookie        string // 74: server's cookie value or backend's cookie name [..BS]
	Mode          string // 75: proxy mode (tcp, http, health, unknown) [LFBS]
	Algo          string // 76: load balancing algorithm [..B.]
	ConnRate      uint32 // 77: number of connections over the last elapsed second [.F..]
	ConnRateMax   uint32 // 78: highest known conn_rate [.F..]
	ConnTot       uint64 // 79: cumulative number of connections [.F..]
	Intercepted   uint64 // 80: Total number of HTTP requests intercepted on the frontend (redirects/stats/services) since the worker process started [.FB.]
	Dcon          uint64 // 81: requests denied by "tcp-request connection" rules [LF..]
	Dses          uint64 // 82: requests denied by "tcp-request session" rules  [LF..]
	Wrew          uint64 // 83: cumulative number of failed header rewriting warnings [LFBS]
	Connect       uint64 // 84: cumulative number of connection establishment attempts [..BS]
	Reuse         uint64 // 85: cumulative number of connection reuses [..BS]
	CacheLookups  uint64 // 86: cumulative number of cache lookups [.FB.]
	CacheHits     uint64 // 87: cumulative number of cache hits [.FB.]
	SrvIcur       uint32 // 88: current number of idle connections available for reuse [...S]
	SrcIlim       uint32 // 89: limit on the number of available idle connections [...S]
	QtimeMax      uint32 // 90: the maximum observed queue time in ms [..BS]
	CtimeMax      uint32 // 91: the maximum observed connect time in ms [..BS]
	RtimeMax      uint32 // 92: the maximum observed response time in ms (0 for TCP) [..BS]
	TtimeMax      uint32 // 93: the maximum observed total session time in ms [..BS]
	Eint          uint64 // 94: cumulative number of internal errors [LFBS]
	IdleConnCur   uint32 // 95: current number of unsafe idle connections [...S]
	SafeConnCur   uint32 // 96: current number of safe idle connections [...S]
	UsedConnCur   uint32 // 97: current number of connections in use [...S]
	NeedConnEst   uint32 // 98: estimated needed number of connections [...S]
	Uweight       uint32 // 99: total user weight (backend), server user weight (server) [..BS]
}

// parse the response of the command show stat from the Runtime API
// using the CSV format as defined http://docs.haproxy.org/2.6/management.html#9.1
func ParseShowStat(response []byte) ([]StatCounters, error) {
	r := csv.NewReader(bytes.NewReader(response))
	r.Comma = ','   // comma separated colums in the runtime API responses
	r.Comment = '#' // lines with # is comments in the runtime API responses

	all, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	stats := make([]StatCounters, 0, len(all))
	for _, row := range all {
		stats = append(stats, elementsToStatCounters(row))
	}

	return stats, nil
}

// map the elements from a response line to the StatCounters struct. If the elements can not be mapped it will panic.
func elementsToStatCounters(elements []string) StatCounters {
	return StatCounters{
		PxName:        elements[0],
		SvName:        elements[1],
		Qcur:          u32(elements[2]),
		Qmax:          u32(elements[3]),
		Scur:          u32(elements[4]),
		Smax:          u32(elements[5]),
		Slim:          u32(elements[6]),
		Stot:          u64(elements[7]),
		Bin:           u64(elements[8]),
		Bout:          u64(elements[9]),
		Dreg:          u64(elements[10]),
		Dresp:         u64(elements[11]),
		Ereg:          u64(elements[12]),
		Econ:          u64(elements[13]),
		Eresp:         u64(elements[14]),
		Wretr:         u64(elements[15]),
		Wredis:        u64(elements[16]),
		Status:        elements[17],
		Weight:        u32(elements[18]),
		Act:           u32(elements[19]),
		Bck:           u32(elements[20]),
		ChkFail:       u64(elements[21]),
		ChkDown:       u64(elements[22]),
		LastChg:       u32(elements[23]),
		Downtime:      u32(elements[24]),
		Qlimit:        u64(elements[25]),
		Pid:           u32(elements[26]),
		Iid:           u32(elements[27]),
		Sid:           u32(elements[28]),
		Throttle:      u64(elements[29]),
		Lbtot:         u64(elements[30]),
		Tracked:       u32(elements[31]),
		Type:          u32(elements[32]),
		Rate:          u32(elements[33]),
		RateLim:       u32(elements[34]),
		RateMax:       u32(elements[35]),
		CheckStatus:   elements[36],
		CheckCode:     u32(elements[37]),
		CheckDuration: u64(elements[38]),
		Hrsp1xx:       u64(elements[39]),
		Hrsp2xx:       u64(elements[40]),
		Hrsp3xx:       u64(elements[41]),
		Hrsp4xx:       u64(elements[42]),
		Hrsp5xx:       u64(elements[43]),
		HrspOther:     u64(elements[44]),
		HanaFail:      u64(elements[45]),
		ReqRate:       u32(elements[46]),
		ReqRateMax:    u32(elements[47]),
		ReqTot:        u64(elements[48]),
		CliAbrt:       u64(elements[49]),
		SrvAbrt:       u64(elements[50]),
		CompIn:        u64(elements[51]),
		CompOut:       u64(elements[52]),
		CompByp:       u64(elements[53]),
		CompRsp:       u64(elements[54]),
		LastSess:      atoi(elements[55]),
		LastChk:       elements[56],
		LastAgt:       elements[57],
		Qtime:         u32(elements[58]),
		Ctime:         u32(elements[59]),
		Rtime:         u32(elements[60]),
		Ttime:         u32(elements[61]),
		AgentStatus:   elements[62],
		AgentCode:     u32(elements[63]),
		AgentDuration: u64(elements[64]),
		CheckDesc:     elements[65],
		AgentDesc:     elements[66],
		CheckRise:     u32(elements[67]),
		CheckFall:     u32(elements[68]),
		CheckHealth:   u32(elements[69]),
		AgentRise:     u32(elements[70]),
		AgentFall:     u32(elements[71]),
		AgentHealth:   u32(elements[72]),
		Addr:          elements[73],
		Cookie:        elements[74],
		Mode:          elements[75],
		Algo:          elements[76],
		ConnRate:      u32(elements[77]),
		ConnRateMax:   u32(elements[78]),
		ConnTot:       u64(elements[79]),
		Intercepted:   u64(elements[80]),
		Dcon:          u64(elements[81]),
		Dses:          u64(elements[82]),
		Wrew:          u64(elements[83]),
		Connect:       u64(elements[84]),
		Reuse:         u64(elements[85]),
		CacheLookups:  u64(elements[86]),
		CacheHits:     u64(elements[87]),
		SrvIcur:       u32(elements[88]),
		SrcIlim:       u32(elements[89]),
		QtimeMax:      u32(elements[90]),
		CtimeMax:      u32(elements[91]),
		RtimeMax:      u32(elements[92]),
		TtimeMax:      u32(elements[93]),
		Eint:          u64(elements[94]),
		IdleConnCur:   u32(elements[95]),
		SafeConnCur:   u32(elements[96]),
		UsedConnCur:   u32(elements[97]),
		NeedConnEst:   u32(elements[98]),
		Uweight:       u32(elements[99]),
	}
}
func atoi(s string) int {
	if len(s) == 0 {
		return 0
	}
	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return x
}

func u64(s string) uint64 {
	if len(s) == 0 {
		return 0
	}
	x, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return x
}

func u32(s string) uint32 {
	if len(s) == 0 {
		return 0
	}
	x, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(x)
}

func atob(s string) bool {
	return atoi(s) != 0
}
