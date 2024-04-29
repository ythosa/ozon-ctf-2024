package main

import (
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

type pdrReader struct {
	ngreader *pcapgo.NgReader
	parser   *pdrParser
	limit    int
}

func newPdrReader(sourcePath string, parser *pdrParser, limit int) pdrReader {
	src, err := os.Open(sourcePath)
	if err != nil {
		log.Fatalf("failed to open pcapng: %v", err)
	}

	ngreader, err := pcapgo.NewNgReader(src, pcapgo.DefaultNgReaderOptions)
	if err != nil {
		log.Fatalf("failed create reader pcapng: %v", err)
	}

	return pdrReader{ngreader: ngreader, parser: parser, limit: limit}
}

func (pdrReader pdrReader) read() {
	packetSource := gopacket.NewPacketSource(pdrReader.ngreader, pdrReader.ngreader.LinkType())
	i := 0

	for packet := range packetSource.Packets() {
		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			udp := udpLayer.(*layers.UDP)
			pdrReader.parser.processUDP(udp)
		}

		i++
		if pdrReader.limit != 0 && i > pdrReader.limit {
			break
		}
	}
}
