package stat

import (
	"bytes"
	"encoding/csv"
	"testing"
)

func TestElementsToStatCounters(t *testing.T) {
	line := "indexws,iws01,0,0,0,1,,8592,1678803,9633529,,0,,4,0,14,0,UP,1,1,0,33,16,10916,15450,,1,4,1,,8578,,2,0,,1,L7OK,200,47,0,8554,0,0,18,0,,,,8572,0,1,,,,,4,,,0,23,31,4590,,,,Layer7 check passed,,2,3,4,,,,172.24.21.40:8080,,http,,,,,,,,0,7992,600,,,7953,,0,175,172,5272,0,7858,95,4294959343,4294959345,1,,,,-,0,0,0,,,,,,,,,,,,,,,,,,,,,,"
	r := csv.NewReader(bytes.NewReader([]byte(line)))
	r.Comma = ','   // white space separated colums in the runtime API responses
	r.Comment = '#' // lines with # is comments in the runtime API responses

	lines, err := r.ReadAll()
	if err != nil {
		t.Fatalf("unable to read CSV : %v", err)
	}
	c := elementsToStatCounters(lines[0])

	if c.PxName != "indexws" {
		t.Fatalf("PxName not indexws")
	}
	if c.SvName != "iws01" {
		t.Fatalf("SvName not iws01")
	}
	if c.Qcur != 0 {
		t.Fatalf("Qcur not 0")
	}
	if c.Qmax != 0 {
		t.Fatalf("Qmax not 0")
	}
	if c.Scur != 0 {
		t.Fatalf("Qmax not 0")
	}
	if c.Smax != 1 {
		t.Fatalf("Smax not 1")
	}
	if c.Slim != 0 {
		t.Fatalf("Slim not 0")
	}
	if c.Stot != 8592 {
		t.Fatalf("Stot not 8592")
	}
	if c.Bin != 1678803 {
		t.Fatalf("Bin not 1678803")
	}
	if c.Bout != 9633529 {
		t.Fatalf("Bout not 9633529")
	}
	if c.Dreg != 0 {
		t.Fatalf("Dreg not 0")
	}
	if c.Dresp != 0 {
		t.Fatalf("Dresp not 0")
	}
	if c.Ereg != 0 {
		t.Fatalf("Ereg not 0")
	}
	if c.Econ != 4 {
		t.Fatalf("Econ not 4")
	}
	if c.Eresp != 0 {
		t.Fatalf("Eresp not 0")
	}
	if c.Wretr != 14 {
		t.Fatalf("Wretr not 14")
	}
	if c.Wredis != 0 {
		t.Fatalf("Wredis not 0")
	}
	if c.Status != "UP" {
		t.Fatalf("Status not UP")
	}
	if c.Weight != 1 {
		t.Fatalf("Weight not 1")
	}
	if c.Act != 1 {
		t.Fatalf("Act not 1")
	}
	if c.Bck != 0 {
		t.Fatalf("Bck not 0")
	}
	if c.ChkFail != 33 {
		t.Fatalf("ChkFail not 33")
	}
	if c.ChkDown != 16 {
		t.Fatalf("ChkDown not 16")
	}
	if c.LastChg != 10916 {
		t.Fatalf("LastChg not 10916")
	}
	if c.Downtime != 15450 {
		t.Fatalf("Downtime not 15450")
	}
	if c.Qlimit != 0 {
		t.Fatalf("Qlimit not 0")
	}
	if c.Pid != 1 {
		t.Fatalf("Pid not 1")
	}
	if c.Iid != 4 {
		t.Fatalf("Iid not 4")
	}
	if c.Sid != 1 {
		t.Fatalf("Sid not 1")
	}
	if c.Throttle != 0 {
		t.Fatalf("Throttle not 0")
	}
	if c.Lbtot != 8578 {
		t.Fatalf("Lbtot not 8578")
	}
	if c.Tracked != 0 {
		t.Fatalf("Tracked not 0")
	}
	if c.Type != 2 {
		t.Fatalf("Type not 2")
	}
	if c.Rate != 0 {
		t.Fatalf("Rate not 0")
	}
	if c.RateLim != 0 {
		t.Fatalf("RateLim not 0")
	}
	if c.RateMax != 1 {
		t.Fatalf("RateMax not 1")
	}
	if c.CheckStatus != "L7OK" {
		t.Fatalf("CheckStatus not L7OK")
	}
	if c.CheckCode != 200 {
		t.Fatalf("CheckCode not 200")
	}
	if c.CheckDuration != 47 {
		t.Fatalf("CheckDuration not 47")
	}
	if c.Hrsp1xx != 0 {
		t.Fatalf("Hrsp1xx not 0")
	}
	if c.Hrsp2xx != 8554 {
		t.Fatalf("Hrsp2xx not 8554")
	}
	if c.Hrsp3xx != 0 {
		t.Fatalf("Hrsp3xx not 0")
	}
	if c.Hrsp4xx != 0 {
		t.Fatalf("Hrsp4xx not 0")
	}
	if c.Hrsp5xx != 18 {
		t.Fatalf("Hrsp5xx not 18")
	}
	if c.HrspOther != 0 {
		t.Fatalf("HrspOther not 0")
	}
	if c.HanaFail != 0 {
		t.Fatalf("HanaFail not 0")
	}
	if c.ReqRate != 0 {
		t.Fatalf("ReqRate not 0")
	}
	if c.ReqRateMax != 0 {
		t.Fatalf("ReqRateMax not 0")
	}
	if c.ReqTot != 8572 {
		t.Fatalf("ReqTot not 8572")
	}
	if c.CliAbrt != 0 {
		t.Fatalf("CliAbrt not 0")
	}
	if c.SrvAbrt != 1 {
		t.Fatalf("SrvAbrt not 1")
	}
	if c.CompIn != 0 {
		t.Fatalf("CompIn not 0")
	}
	if c.CompOut != 0 {
		t.Fatalf("CompOut not 0")
	}
	if c.CompByp != 0 {
		t.Fatalf("CompByp not 0")
	}
	if c.CompRsp != 0 {
		t.Fatalf("CompRsp not 0")
	}
	if c.LastSess != 4 {
		t.Fatalf("LastSess not 4")
	}
	if c.LastChk != "" {
		t.Fatalf("LastChk not blank")
	}
	if c.LastAgt != "" {
		t.Fatalf("LastAgt not blank")
	}
	if c.Qtime != 0 {
		t.Fatalf("Qtime not 0")
	}
	if c.Ctime != 23 {
		t.Fatalf("Ctime not 23")
	}
	if c.Rtime != 31 {
		t.Fatalf("Rtime not 31")
	}
	if c.Ttime != 4590 {
		t.Fatalf("Ttime not 4590")
	}
	if c.AgentStatus != "" {
		t.Fatalf("AgentStatus not blank")
	}
	if c.AgentCode != 0 {
		t.Fatalf("AgentCode not 0")
	}
	if c.AgentDuration != 0 {
		t.Fatalf("AgentDuration not 0")
	}
	if c.CheckDesc != "Layer7 check passed" {
		t.Fatalf("CheckDesc not Layer7 check passed")
	}
	if c.AgentDesc != "" {
		t.Fatalf("AgentDesc not blank")
	}
	if c.CheckRise != 2 {
		t.Fatalf("CheckRise not 2")
	}
	if c.CheckFall != 3 {
		t.Fatalf("CheckFall not 3")
	}
	if c.CheckHealth != 4 {
		t.Fatalf("CheckHealth not 4")
	}
	if c.AgentRise != 0 {
		t.Fatalf("AgentRise not 0")
	}
	if c.AgentFall != 0 {
		t.Fatalf("AgentFall not 0")
	}
	if c.AgentHealth != 0 {
		t.Fatalf("AgentHealth not 0")
	}
	if c.Addr != "172.24.21.40:8080" {
		t.Fatalf("Addr not 172.24.21.40:8080")
	}
	if c.Cookie != "" {
		t.Fatalf("Cookie not blank")
	}
	if c.Mode != "http" {
		t.Fatalf("Mode not http")
	}
	if c.Algo != "" {
		t.Fatalf("Algo not blank")
	}
	if c.ConnRate != 0 {
		t.Fatalf("ConnRate not 0")
	}
	if c.ConnRateMax != 0 {
		t.Fatalf("ConnRateMax not 0")
	}
	if c.ConnTot != 0 {
		t.Fatalf("ConnTot not 0")
	}
	if c.Intercepted != 0 {
		t.Fatalf("Intercepted not 0")
	}
	if c.Dcon != 0 {
		t.Fatalf("Dcon not 0")
	}
	if c.Dses != 0 {
		t.Fatalf("Dses not 0")
	}
	if c.Wrew != 0 {
		t.Fatalf("Wrew not 0")
	}
	if c.Connect != 7992 {
		t.Fatalf("Connect not 7992")
	}
	if c.Reuse != 600 {
		t.Fatalf("Reuse not 600")
	}
	if c.CacheLookups != 0 {
		t.Fatalf("CacheLookups not 0")
	}
	if c.CacheHits != 0 {
		t.Fatalf("CacheHits not 0")
	}
	if c.SrvIcur != 7953 {
		t.Fatalf("SrvIcur not 7953")
	}
	if c.SrcIlim != 0 {
		t.Fatalf("SrcIlim not 0")
	}
	if c.QtimeMax != 0 {
		t.Fatalf("QtimeMax not 0")
	}
	if c.CtimeMax != 175 {
		t.Fatalf("CtimeMax not 175")
	}
	if c.RtimeMax != 172 {
		t.Fatalf("RtimeMax not 172")
	}
	if c.TtimeMax != 5272 {
		t.Fatalf("TtimeMax not 5272")
	}
	if c.Eint != 0 {
		t.Fatalf("Eint not 0")
	}
	if c.IdleConnCur != 7858 {
		t.Fatalf("IdleConnCur not 7858")
	}
	if c.SafeConnCur != 95 {
		t.Fatalf("SafeConnCur not 95")
	}
	if c.UsedConnCur != 4294959343 {
		t.Fatalf("UsedConnCur not 4294959343")
	}
	if c.NeedConnEst != 4294959345 {
		t.Fatalf("NeedConnEst not 4294959345")
	}
	if c.Uweight != 1 {
		t.Fatalf("Uweight not 1")
	}

	// blank,blank,blank,-,0,0,0,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank,blank

}
