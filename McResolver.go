package main

import (
	"fmt"
	"net"
)

func mcResolvPostProcess(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		item = normalize_addr(item)
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ResolvMcAddr(addr string) (string, []string, error) {
	var aliases []string
	is_port_explicit := true
	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		aliases = append(aliases, addr)
		host = addr
		port = "25565"
		is_port_explicit = false
	}

	if ip := net.ParseIP(host); ip != nil {
		aliases = append(aliases, net.JoinHostPort(host, port))
		return net.JoinHostPort(host, port), mcResolvPostProcess(aliases), nil
	}

	if is_port_explicit == false {
		ret_addr := ""
		_, srvs, _ := net.LookupSRV("minecraft", "tcp", host)

		for _, srv := range srvs {
			if net.ParseIP(srv.Target) == nil {
				aliases = append(aliases, net.JoinHostPort(srv.Target[:len(srv.Target)-1], fmt.Sprintf("%d", srv.Port)))
				addresses, _ := net.LookupHost(srv.Target)
				for _, res_addr := range addresses {
					ret_addr = net.JoinHostPort(res_addr, fmt.Sprintf("%d", srv.Port))
					aliases = append(aliases, ret_addr)
				}
			} else {
				ret_addr = net.JoinHostPort(srv.Target, fmt.Sprintf("%d", srv.Port))
				aliases = append(aliases, ret_addr)
			}
		}

		if len(ret_addr) != 0 {
			return ret_addr, mcResolvPostProcess(aliases), nil
		}
	}

	aliases = append(aliases, net.JoinHostPort(host, port))
	addresses, _ := net.LookupHost(host)
	if len(addresses) == 0 {
		return "", mcResolvPostProcess(aliases), fmt.Errorf("no address for `%s`", host)
	}

	for _, address := range addresses {
		aliases = append(aliases, net.JoinHostPort(address, port))
	}

	// seems like net.LookupHost returns IPv4 after IPv6
	return net.JoinHostPort(addresses[len(addresses)-1], port), mcResolvPostProcess(aliases), nil
}
