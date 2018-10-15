package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FW struct {
	File   *os.File
	Writer *csv.Writer
}

var fileMap map[string]FW

var flowStatsWriter *csv.Writer
var flowStatsFile *os.File

func createDirs(filePath string) {
	//Creating directories
	directorystring := filepath.Dir(filePath)
	err := os.MkdirAll(directorystring, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func createFile(pcapFile string, filename string) *os.File {
	//Creating file
	filePath := filepath.Join(Config.OutputDir,
		strings.TrimSuffix(filepath.Base(pcapFile), filepath.Ext(pcapFile)),
		filename)
	createDirs(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func fillFWIntoMap(pcapFile string, filename string) {
	file := createFile(pcapFile, filename)
	fileMap[filename] = FW{
		File:   file,
		Writer: csv.NewWriter(file),
	}
}

func initCSVWriters(pcapFile string) {
	fileMap = make(map[string]FW)
	fillFWIntoMap(pcapFile, Config.Files.Summary)
	fillFWIntoMap(pcapFile, Config.Files.TopVolume)
	fillFWIntoMap(pcapFile, Config.Files.TopDuration)
	fillFWIntoMap(pcapFile, Config.Files.DNS)

	fileMap[Config.Files.Summary].Writer.Write([]string{"Packets", "Volume (bytes)",
		"Duration", "Flows", "IP Packets", "TCP Packets", "UDP Packets", "DNS Packets"})

	fileMap[Config.Files.TopVolume].Writer.Write([]string{"FlowID", "SrcIP", "DstIP",
		"SrcPort", "DstPort", "Protocol", "Volume (bytes)", "Packet Count", "Duration"})

	fileMap[Config.Files.TopDuration].Writer.Write([]string{"FlowID", "SrcIP", "DstIP",
		"SrcPort", "DstPort", "Protocol", "Volume (bytes)", "Packet Count", "Duration"})

	fileMap[Config.Files.DNS].Writer.Write([]string{"Timestamp", "IP", "Name"})

	if Config.Timeseries.Enabled {
		fillFWIntoMap(pcapFile, Config.Files.FlowTimeseries)
		fileMap[Config.Files.FlowTimeseries].Writer.Write([]string{"FlowID", "SrcIP", "DstIP",
			"SrcPort", "DstPort", "Protocol", "StatTimestamp", "Volume (bytes)",
			"Packet Count", "Duration"})
	}

	if Config.Filters.Enabled {
		fillFWIntoMap(pcapFile, Config.Files.FilteredFlows)
		fileMap[Config.Files.FilteredFlows].Writer.Write([]string{"FlowID", "SrcIP", "DstIP",
			"SrcPort", "DstPort", "Protocol", "Volume (bytes)", "Packet Count", "Duration"})
	}
}

func closeCSVFiles() {
	for _, value := range fileMap {
		value.File.Close()
		value.Writer.Flush()
	}
}

func writeTimeseriesToCSV(f Flow) {

	var rows [][]string

	common := []string{
		f.FlowEntry.UUID,
		f.FiveTuple.SrcIP,
		f.FiveTuple.DstIP,
		f.FiveTuple.SrcPort,
		f.FiveTuple.DstPort,
		f.FiveTuple.Protocol,
	}

	last_snapshot := f.FlowEntry.Profile.Snapshots[0]
	next_index := last_snapshot.Index
	for _, s := range f.FlowEntry.Profile.Snapshots {
		created := f.FlowEntry.Profile.Created
		for i := next_index; i <= s.Index; i++ {
			row := common
			dur := time.Duration(Config.Timeseries.Interval*int(i+1)) * time.Millisecond
			ts := created.Add(dur)
			vol := last_snapshot.Volume
			pktcount := last_snapshot.PacketCount
			if i == s.Index {
				vol = s.Volume
				pktcount = s.PacketCount
			}
			row = append(row,
				fmt.Sprint(ts),
				fmt.Sprint(vol),
				fmt.Sprint(pktcount),
				fmt.Sprint(dur))
			rows = append(rows, row)
		}
		next_index = s.Index + 1
		last_snapshot = s
	}
	fileMap[Config.Files.FlowTimeseries].Writer.WriteAll(rows)
	fileMap[Config.Files.FlowTimeseries].Writer.Flush()
}

func writeSummaryToCSV(s PktLvlSummary) {
	row := []string{
		fmt.Sprint(s.PacketCount),
		fmt.Sprint(s.Volume),
		fmt.Sprint(s.GetDuration()),
		fmt.Sprint(len(flowMap)),
		fmt.Sprint(s.IPCount),
		fmt.Sprint(s.TCPCount),
		fmt.Sprint(s.UDPCount),
		fmt.Sprint(s.DNSCount),
	}
	fileMap[Config.Files.Summary].Writer.Write(row)
	fileMap[Config.Files.Summary].Writer.Flush()
}

func writeFlowSummariesToCSV(flows []Flow, sortby string) {
	var rows [][]string
	for i := 0; i < len(flows); i++ {
		f := flows[i]
		row := []string{
			f.FlowEntry.UUID,
			f.FiveTuple.SrcIP,
			f.FiveTuple.DstIP,
			f.FiveTuple.SrcPort,
			f.FiveTuple.DstPort,
			f.FiveTuple.Protocol,
			fmt.Sprint(f.FlowEntry.Volume),
			fmt.Sprint(f.FlowEntry.PacketCount),
			fmt.Sprint(f.FlowEntry.GetDuration()),
		}
		rows = append(rows, row)
	}
	var writer *csv.Writer
	if sortby == "volume" {
		writer = fileMap[Config.Files.TopVolume].Writer
	} else if sortby == "duration" {
		writer = fileMap[Config.Files.TopDuration].Writer
	} else if sortby == "filtered" {
		writer = fileMap[Config.Files.FilteredFlows].Writer
	}
	writer.WriteAll(rows)
	writer.Flush()
}

func writeDNSToCSV() {
	var rows [][]string
	for _, dnsData := range dnsSlice {
		row := []string{
			dnsData.Timestamp.String(),
			dnsData.IP,
			dnsData.Name,
		}
		rows = append(rows, row)
	}
	writer := fileMap[Config.Files.DNS].Writer
	writer.WriteAll(rows)
	writer.Flush()

}
