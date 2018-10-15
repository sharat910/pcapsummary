package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/google/gopacket/pcap"
)

var (
	pcapFile string = "pcaps/wifi_calling_13june_nexus_tmobile.pcap"
	handle   *pcap.Handle
	err      error
)

func main() {
	_ = flag.String("config", "config.yaml", "filepath for yaml configuration file")
	flag.Parse()
	LoadConfig()

	files, _ := filepath.Glob(Config.PCAPGlob)
	for i := 0; i < len(files); i++ {
		fmt.Println("\n\n====", files[i], "\n")
		initCSVWriters(files[i])
		s := Analyse(files[i])
		Summarize(s)
		closeCSVFiles()
	}
}
