package main

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func Test_pdrParser_processPayload(t *testing.T) {
	in := []byte{231, 205, 24, 4} // must be ~1345

	//fmt.Println(binary.LittleEndian.Uint16(in))
	//fmt.Println(binary.LittleEndian.Uint32(in))
	//fmt.Println(binary.BigEndian.Uint16(in))
	//fmt.Println(binary.BigEndian.Uint32(in))

	fmt.Println(binary.LittleEndian.AppendUint32([]byte{}, 1345))
	fmt.Println(binary.LittleEndian.Uint32(in))
}
