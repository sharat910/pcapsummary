package main

import (
	"fmt"
	"sort"
)

func getFlows() []Flow {
	var flows []Flow
	empty := FiveTuple{}
	for key, value := range flowMap {
		//Don't know why this happens
		if key == empty {
			fmt.Println("!")
			continue
		}
		flows = append(flows, Flow{key, value})
	}
	return flows
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func passingFilter(flow Flow) bool {
	ft := flow.FiveTuple
	if ft.Protocol == Config.Filters.Protocol {
		for _, port := range Config.Filters.Ports {
			if ft.SrcPort == port || ft.DstPort == port {
				return true
			}
		}

	}
	return false
}

func Summarize(s PktLvlSummary) {
	fmt.Printf("%+v\n", s)
	writeSummaryToCSV(s)
	writeDNSToCSV()
	n := Config.TopN
	nFlows := len(flowMap)
	end := min(n, nFlows)

	flows := getFlows()

	// Sort by volume
	sort.Slice(flows, func(i, j int) bool {
		return flows[i].FlowEntry.Volume > flows[j].FlowEntry.Volume
	})
	fmt.Println("Top flows by volume")
	for i := 0; i < end; i++ {
		printFlow(flows[i])
	}
	writeFlowSummariesToCSV(flows[:end], "volume")

	sort.Slice(flows, func(i, j int) bool {
		return flows[i].FlowEntry.GetDuration() > flows[j].FlowEntry.GetDuration()
	})
	fmt.Println("Top flows by Duration")
	for i := 0; i < end; i++ {
		printFlow(flows[i])
	}

	writeFlowSummariesToCSV(flows[:end], "duration")

	if Config.Filters.Enabled {
		var filteredFlows []Flow
		for _, flow := range flows {
			if passingFilter(flow) {
				filteredFlows = append(filteredFlows, flow)
			}
		}
		writeFlowSummariesToCSV(filteredFlows, "filtered")
	}

	if Config.Timeseries.Enabled {
		for _, flow := range flows {
			if Config.Timeseries.ApplyFilters {
				if passingFilter(flow) {
					writeTimeseriesToCSV(flow)
				}
			} else {
				writeTimeseriesToCSV(flow)
			}

		}
	}

}
