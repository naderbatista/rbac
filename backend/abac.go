package main

import (
	"net"
	"strings"
	"time"
)

func evaluatePolicy(p Policy, clientIP string, now time.Time) bool {
	switch p.Type {
	case "horario":
		return evaluateTimeRange(p.Value, now)
	case "ip":
		return evaluateIPList(p.Value, clientIP)
	default:
		return true
	}
}

func evaluateTimeRange(value string, now time.Time) bool {
	parts := strings.SplitN(value, "-", 2)
	if len(parts) != 2 {
		return true
	}
	start, err1 := time.Parse("15:04", strings.TrimSpace(parts[0]))
	end, err2 := time.Parse("15:04", strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return true
	}

	cur := now.Hour()*60 + now.Minute()
	s := start.Hour()*60 + start.Minute()
	e := end.Hour()*60 + end.Minute()

	return cur >= s && cur <= e
}

func evaluateIPList(value, clientIP string) bool {
	for _, entry := range strings.Split(value, ",") {
		entry = strings.TrimSpace(entry)
		if entry == clientIP {
			return true
		}
		_, cidr, err := net.ParseCIDR(entry)
		if err == nil && cidr.Contains(net.ParseIP(clientIP)) {
			return true
		}
	}
	return false
}
