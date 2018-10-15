package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Files struct {
	OutputDir      string `yaml:"output_dir"`
	Summary        string `yaml:"summary"`
	DNS            string `yaml:"dns"`
	TopVolume      string `yaml:"topvolume"`
	TopDuration    string `yaml:"topduration"`
	FlowTimeseries string `yaml:"flowtimeseries"`
	FilteredFlows  string `yaml:"filteredflows"`
}

type Timeseries struct {
	Enabled      bool `yaml:"enabled"`
	ApplyFilters bool `yaml:"apply_filters"`
	Interval     int  `yaml:"interval"`
}

type Filters struct {
	Enabled  bool     `yaml:"enabled"`
	Protocol string   `yaml:"protocol"`
	Ports    []string `yaml:"ports"`
}

type ConfStruct struct {
	PCAPGlob string `yaml:"pcap_glob"`
	TopN     int    `yaml:"top_n"`
	Files
	Timeseries
	Filters
}

func (c *ConfStruct) loadConfig() *ConfStruct {
	log.Println("Loading config...")
	filepath := get_filepath_from_flag()
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Println(c)
	return c
}

func get_filepath_from_flag() string {
	return flag.Lookup("config").Value.(flag.Getter).Get().(string)
}

var Config ConfStruct

func LoadConfig() {
	Config.loadConfig()
}
