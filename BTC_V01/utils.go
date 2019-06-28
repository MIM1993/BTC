package main

import (
	"encoding/binary"
	"bytes"
	"fmt"
)

//实现uint64转换为[]byte

func uintToByte(num uint64) []byte {
	var buffer bytes.Buffer
	//使用二进制编码
	if err := binary.Write(&buffer, binary.LittleEndian, &num); err != nil {
		fmt.Println("binary.Write err：", err)
		return nil
	}

	return buffer.Bytes()
}
