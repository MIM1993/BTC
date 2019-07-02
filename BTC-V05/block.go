package main

import (
	"bytes"
	"time"
	"encoding/gob"
	"fmt"
	"crypto/sha256"
)

//定义区块结构
//实现基础字段  前区块 hash 数据
//补充字段   version  时间戳 难度值
type Block struct {
	//version
	Version uint64
	//前区块
	PrevHash []byte
	//交易的根哈希值
	MerkleRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值
	Bits uint64
	//随机数
	Nonce uint64
	//哈希
	Hash []byte

	//交易集合，一个区块可以有很多集合
	Transactions []*Transaction
}

//创建区块方法		区块       前hash值
func NewBlock(txs []*Transaction, PrevHash []byte) *Block {
	b := Block{
		Version:    0,
		PrevHash:   PrevHash,
		MerkleRoot: nil,
		TimeStamp:  uint64(time.Now().Unix()),

		Bits:  0, //随意写的
		Nonce: 0, //随意写的
		Hash:  nil,
		//Data:  []byte(data),
		Transactions: txs,
	}

	//添加简易梅克尔根
	b.HashTransactionMerkleRoot()

	//挖矿，计算Hash值
	pow := NewProofofWork(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce

	return &b
}

//绑定Serialize方法  将block转码为字节流
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(b)
	if err != nil {
		fmt.Printf("Encode err:", err)
		return nil
	}
	return buffer.Bytes()
}

//将字节流转为block
//反序列化，输入[]byte，返回block
func Deserialize(src []byte) *Block {
	var block Block

	//解码器
	decoder := gob.NewDecoder(bytes.NewReader(src))

	//解码
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Decode err:", err)
		return nil
	}
	return &block
}

//简易梅克尔根生成
func (block *Block) HashTransactionMerkleRoot() {
	var info [][]byte

	//遍历所有交易
	for _, tx := range block.Transactions {
		txhashinfo := tx.Txid
		info = append(info, txhashinfo)
	}

	value := bytes.Join(info, []byte{})

	hash := sha256.Sum256(value)

	block.MerkleRoot = hash[:]
}
