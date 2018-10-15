package main

import (
	"fmt"

	"github.com/google/uuid"
)

var NewFlows int

func insertOrUpdateFlow(ft FiveTuple, pktInfo PacketInfo) {
	if flowEntry, ok := flowMap[ft]; ok {
		//exists
		flowEntry.PacketCount += 1
		flowEntry.Volume += pktInfo.Volume
		flowEntry.LastUpdated = pktInfo.Timestamp
		buildProfile(&flowEntry)
		flowMap[ft] = flowEntry
	} else {
		NewFlows += 1
		fe := FlowEntry{
			UUID:        uuid.New().String(),
			PacketCount: 1,
			Volume:      pktInfo.Volume,
			Created:     pktInfo.Timestamp,
			LastUpdated: pktInfo.Timestamp,
		}
		initProfile(&fe)
		flowMap[ft] = fe
	}
}

// func printStringMap() {
// 	for key, value := range flowMap {
// 		dnsNames := getDNSNames(key.SrcIP)
// 		if dnsNames != nil {
// 			fmt.Println(key, ": ", value.PacketCount, value.Volume, value.GetDuration(), dnsNames)
// 		}
// 	}
// 	fmt.Println("Total flows", len(flowMap))
// }

func printFlow(f Flow) {
	dnsNames := getDNSNames(f.FiveTuple.SrcIP)
	fmt.Println(f.FiveTuple, f.FlowEntry.PacketCount,
		f.FlowEntry.Volume, f.FlowEntry.GetDuration(), dnsNames)
}
