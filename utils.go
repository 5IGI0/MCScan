package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func normalize_addr(addr string) string {
	if strings.HasSuffix(addr, ":25565") {
		return strings.ToLower(addr[:len(addr)-6])
	}
	return addr
}

func assert_err(err error) {
	if err != nil {
		panic(err)
	}
}

func get_timestamp() int64 {
	return time.Now().Unix()
}

var symbdlt = regexp.MustCompile(`§.`)

func delete_mc_symbols(text string) string {
	return symbdlt.ReplaceAllString(text, "")
}

func escape_like(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, "\\", "\\\\"), "_", "\\_"), "%", "\\%")
}

func timestamp_to_relative(timestamp uint32) string {
	var rel int64 = int64(timestamp) - get_timestamp()

	if rel >= -15 && rel <= 15 {
		return "now"
	}

	var sign bool = true
	var out string

	if rel < 0 {
		sign = false
		rel = -rel
	}

	if rel > 60*60*24 {
		out = fmt.Sprint(rel/(60*60*24)) + " day(s)"
	} else if rel > 60*60 {
		out = fmt.Sprint(rel/(60*60)) + " hour(s)"
	} else if rel > 60 {
		out = fmt.Sprint(rel/60) + " minute(s)"
	} else {
		out = fmt.Sprint(rel) + " second(s)"
	}

	if sign {
		return "in " + out
	}
	return out + " ago"
}
