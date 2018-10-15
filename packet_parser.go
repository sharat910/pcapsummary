package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func Analyse(pcapFile string) PktLvlSummary {
	var (
		// Will reuse these for each packet
		ethLayer layers.Ethernet
		ipLayer  layers.IPv4
		tcpLayer layers.TCP
		udpLayer layers.UDP
		dnsLayer layers.DNS
	)
	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet,
		&ethLayer,
		&ipLayer,
		&tcpLayer,
		&udpLayer,
		&dnsLayer,
	)
	// Open file instead of device
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.DecodeOptions.Lazy = true
	packetSource.DecodeOptions.NoCopy = true

	// Initialize variables
	initMaps()
	s := PktLvlSummary{}
	firstPacketSeen := false

	var packet gopacket.Packet
	for packet = range packetSource.Packets() {
		//fmt.Println(packet)
		if !firstPacketSeen {
			s.FirstPacketTS = packet.Metadata().Timestamp
			firstPacketSeen = true
		}
		s.PacketCount += 1
		s.Volume += packet.Metadata().Length
		foundLayerTypes := []gopacket.LayerType{}

		_ = parser.DecodeLayers(packet.Data(), &foundLayerTypes)
		var fiveTup FiveTuple
		ip_flag, tcp_flag, udp_flag := false, false, false

		for _, layerType := range foundLayerTypes {
			switch layerType {
			case layers.LayerTypeIPv4:
				s.IPCount += 1
				ip_flag = true

				fiveTup.SrcIP = ipLayer.SrcIP.String()
				fiveTup.DstIP = ipLayer.DstIP.String()
				fiveTup.Protocol = ipLayer.Protocol.String()

			case layers.LayerTypeTCP:
				s.TCPCount += 1
				tcp_flag = true

				fiveTup.SrcPort = tcpLayer.SrcPort.String()
				fiveTup.DstPort = tcpLayer.DstPort.String()

			case layers.LayerTypeUDP:
				s.UDPCount += 1
				udp_flag = true

				fiveTup.SrcPort = udpLayer.SrcPort.String()
				fiveTup.DstPort = udpLayer.DstPort.String()

			case layers.LayerTypeDNS:
				s.DNSCount += 1
				HandleDNS(dnsLayer, packet.Metadata().Timestamp)
			}
		}

		// received_five_tuple := ip_flag && (tcp_flag || udp_flag)
		// empty := FiveTuple{}
		// if fiveTup == empty && received_five_tuple {
		// 	fmt.Println(packet)
		// }

		received_five_tuple := ip_flag && (tcp_flag || udp_flag)
		if received_five_tuple {
			//insertOrUpdate(five_tuple)
			pi := PacketInfo{packet.Metadata().Length, packet.Metadata().Timestamp}
			insertOrUpdateFlow(fiveTup, pi)
		}
	}
	s.LastPacketTS = packet.Metadata().Timestamp

	return s
}
