package main

import "time"

type PktLvlSummary struct {
	PacketCount   int
	IPCount       int
	TCPCount      int
	UDPCount      int
	DNSCount      int
	Volume        int
	FirstPacketTS time.Time
	LastPacketTS  time.Time
}

func (p PktLvlSummary) GetDuration() time.Duration {
	return p.LastPacketTS.Sub(p.FirstPacketTS)
}

type FiveTuple struct {
	SrcIP, DstIP, SrcPort, DstPort, Protocol string
}

type Flow struct {
	FiveTuple FiveTuple
	FlowEntry FlowEntry
}

type FlowEntry struct {
	UUID        string
	Created     time.Time
	LastUpdated time.Time
	Volume      int
	PacketCount int
	Profile     *FlowProfile
}

func (f FlowEntry) GetDuration() time.Duration {
	return f.LastUpdated.Sub(f.Created)
}

type FlowProfile struct {
	Created      time.Time
	CurrentIndex int
	NEntries     int
	Snapshots    []*Snapshot
}

type Snapshot struct {
	Index, Volume, PacketCount int
}

type PacketInfo struct {
	Volume    int
	Timestamp time.Time
}

type DNSData struct {
	Timestamp time.Time
	IP        string
	Name      string
}

var flowMap map[FiveTuple]FlowEntry
var dnsMap map[string][]string
var dnsSlice []DNSData

func initMaps() {
	flowMap = make(map[FiveTuple]FlowEntry)
	dnsMap = make(map[string][]string)
	dnsSlice = nil
}
