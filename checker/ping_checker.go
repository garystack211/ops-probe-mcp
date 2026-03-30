package checker

import (
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type PingResult struct {
	Host         string  `json:"host"`
	PacketsSent  int     `json:"packets_sent"`
	PacketsRecv  int     `json:"packets_recv"`
	PacketLoss   float64 `json:"packet_loss"`
	MinRTT       float64 `json:"min_rtt_ms"`
	MaxRTT       float64 `json:"max_rtt_ms"`
	AvgRTT       float64 `json:"avg_rtt_ms"`
	Reachable    bool    `json:"reachable"`
	Error        string  `json:"error,omitempty"`
}

func PingCheck(host string, count int, timeout time.Duration) (*PingResult, error) {
	result := &PingResult{
		Host:        host,
		PacketsSent: count,
		Reachable:   false,
	}

	if count <= 0 {
		count = 4
	}
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	args := []string{"-c", strconv.Itoa(count), "-W", strconv.Itoa(int(timeout.Seconds())), host}
	cmd := exec.Command("ping", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	outputStr := string(output)

	// 解析丢包率
	lossRegex := regexp.MustCompile(`(\d+)% packet loss`)
	if matches := lossRegex.FindStringSubmatch(outputStr); len(matches) > 1 {
		loss, _ := strconv.ParseFloat(matches[1], 64)
		result.PacketLoss = loss
		result.PacketsRecv = count - int(float64(count)*loss/100)
	}

	// 解析 RTT
	rttRegex := regexp.MustCompile(`min/avg/max[^=]*=\s*([\d.]+)/([\d.]+)/([\d.]+)`)
	if matches := rttRegex.FindStringSubmatch(outputStr); len(matches) > 3 {
		result.MinRTT, _ = strconv.ParseFloat(matches[1], 64)
		result.AvgRTT, _ = strconv.ParseFloat(matches[2], 64)
		result.MaxRTT, _ = strconv.ParseFloat(matches[3], 64)
		result.Reachable = true
	}

	return result, nil
}
