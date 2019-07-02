package main

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"os"
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

//判断文件是否存在
func IsFileExist(filename string) bool {
	//查看文件状态
	_, err := os.Stat(filename)

	//查看文件是否不存在  这个函数正确性更高，不要用os.IsExist
	if os.IsNotExist(err) {
		return false
	}

	return true
}


