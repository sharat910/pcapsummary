package main

import (
	"fmt"
	"time"

	"github.com/google/gopacket/layers"
)

func HandleDNS(dnsLayer layers.DNS, ts time.Time) {
	for i := 0; i < int(dnsLayer.ANCount); i++ {
		if dnsRec := dnsLayer.Answers[i]; dnsRec.Type.String() == "A" {
			IP := dnsRec.IP.String()
			Name := string(dnsRec.Name)
			dnsSlice = append(dnsSlice, DNSData{ts, IP, Name})
			insertOrUpdateDNSMap(IP, Name)
		}
	}
}

func insertOrUpdateDNSMap(IP string, dnsName string) {
	if dnsEntry, ok := dnsMap[IP]; ok {
		dnsMap[IP] = append(dnsEntry, dnsName)
	} else {
		dnsMap[IP] = []string{dnsName}
	}
}

func getDNSNames(IP string) []string {
	if dnsEntry, ok := dnsMap[IP]; ok {
		return dnsEntry
	} else {
		return nil
	}
}

func printDNSMap() {
	//fmt.Println("DNS Map")
	n_names := 0
	for _, value := range dnsMap {
		n_names += len(value)
	}
	fmt.Println("Total DNS IPs", len(dnsMap))
	fmt.Println("Total DNS Names", len(dnsMap))
}
