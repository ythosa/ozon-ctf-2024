package main

import (
	"encoding/binary"
	"hash/crc32"
	"sort"

	"github.com/google/gopacket/layers"
)

const (
	payloadMagicOffset = 0
	payloadMagicSize   = 4

	payloadIdxOffset = 4
	payloadIdxSize   = 4

	payloadSizeOffset = 8
	payloadSizeSize   = 4

	payloadDataOffset = 12

	payloadCrc32Size = 4
)

type block struct {
	idx  int
	data []byte
}

type pdrParser struct {
	encryptedImage []block
	keyImage       []block
}

func newPdrParser() *pdrParser {
	return &pdrParser{encryptedImage: make([]block, 0), keyImage: make([]block, 0)}
}

func (ib *pdrParser) processUDP(layer *layers.UDP) {
	b, ok := ib.processPayload(layer.Payload)
	if !ok {
		return
	}

	if layer.DstPort%2 == 0 {
		ib.keyImage = append(ib.keyImage, b)
	} else {
		ib.encryptedImage = append(ib.encryptedImage, b)
	}
}

func (ib *pdrParser) processPayload(payload []byte) (block, bool) {
	magic := string(payload[payloadMagicOffset : payloadMagicOffset+payloadMagicSize])
	if magic != "PDR\x00" {
		return block{}, false
	}

	idx := binary.LittleEndian.Uint32(payload[payloadIdxOffset : payloadIdxOffset+payloadIdxSize]) // todo: fix me
	size := binary.LittleEndian.Uint32(payload[payloadSizeOffset : payloadSizeOffset+payloadSizeSize])
	data := payload[payloadDataOffset : len(payload)-payloadCrc32Size]

	crc := binary.LittleEndian.Uint32(payload[len(payload)-payloadCrc32Size:])
	if crc != crc32.ChecksumIEEE(data) {
		return block{}, false
	}

	if int(size) != len(data) {
		return block{}, false
	}

	return block{int(idx), data}, true
}

func (ib *pdrParser) fetchEncryptedImage() []byte {
	return ib.dataFromBlocks(ib.encryptedImage)
}

func (ib *pdrParser) fetchKeyImage() []byte {
	return ib.dataFromBlocks(ib.keyImage)
}

func (ib *pdrParser) dataFromBlocks(blocks []block) []byte {
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].idx < blocks[j].idx
	})

	lastIdx := -1

	var result []byte
	for _, b := range blocks {
		if b.idx != lastIdx {
			result = append(result, b.data...)
			lastIdx = b.idx
		}
	}

	return result
}
