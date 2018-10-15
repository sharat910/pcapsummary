package main

import (
	"log"
	"time"
)

func getSlotIndex(createdAt time.Time, pktTimestamp time.Time) int {
	dur := pktTimestamp.Sub(createdAt)
	interval := time.Duration(Config.Timeseries.Interval) * time.Millisecond
	return int(dur / interval)

}

func initProfile(flowEntry *FlowEntry) {
	flowProfile := FlowProfile{}
	s := &Snapshot{
		Index:       0,
		Volume:      flowEntry.Volume,
		PacketCount: flowEntry.PacketCount,
	}
	flowProfile.Created = flowEntry.Created
	flowProfile.Snapshots = append(flowProfile.Snapshots, s)
	flowProfile.CurrentIndex = 0
	flowProfile.NEntries = 1
	flowEntry.Profile = &flowProfile
	//fmt.Println("Initializing")
	// fmt.Println(flowEntry)
}

func buildProfile(flowEntry *FlowEntry) {
	flowProfile := flowEntry.Profile
	createdAt := flowEntry.Created
	pktTimestamp := flowEntry.LastUpdated
	idx := getSlotIndex(createdAt, pktTimestamp)
	if flowProfile.CurrentIndex == idx && flowProfile.NEntries > 0 {
		CurrSnapshot := flowProfile.Snapshots[flowProfile.NEntries-1]
		CurrSnapshot.Volume = flowEntry.Volume
		CurrSnapshot.PacketCount = flowEntry.PacketCount
	} else if idx > flowProfile.CurrentIndex || flowProfile.NEntries == 0 {
		s := &Snapshot{
			Index:       idx,
			Volume:      flowEntry.Volume,
			PacketCount: flowEntry.PacketCount,
		}
		flowProfile.Snapshots = append(flowProfile.Snapshots, s)
		flowProfile.CurrentIndex = idx
		flowProfile.NEntries += 1
	} else {
		log.Fatal("Current idx > SlotIndex")
	}
	// fmt.Println("Building")
	// fmt.Println(flowEntry)
}

func resetProfile(flowEntry FlowEntry) {
	flowProfile := flowEntry.Profile
	flowProfile.NEntries = 0
	flowProfile.Snapshots = nil
}

// func exportProfile(flowEntry FlowEntry, fiveTuple FiveTuple, threshold int) {
// 	flowProfile := flowEntry.Profile
// 	flowStat := FlowStat{
// 		UUID:      flowEntry.UUID,
// 		Timestamp: time.Time,
// 		FiveTuple: fiveTuple,
// 		Threshold: threshold,
// 		Profile:   flowProfile,
// 	}
// 	FlowStatChannel <- flowStat
// 	resetProfile(flowEntry)
// }
