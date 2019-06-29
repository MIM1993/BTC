package main

import (
	"math/big"
	"fmt"
	"bytes"
	"crypto/sha256"
)

//实现挖矿功能

//方法
//run计算
//功能：找到nonce，满足哈希值比目标值小

type ProofofWork struct {
	//区块
	block *Block
	//目标值 与计算出的哈希值作对比
	target *big.Int
}

//创建一个工作量证明
//block  用户提供
//target 系统提供
func NewProofofWork(block *Block) *ProofofWork {
	pow := ProofofWork{
		block: block,
	}

	//难度值  写死 后面补充
	targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	tmpBigInt := new(big.Int)
	//将我们的难度值赋值给bigint
	tmpBigInt.SetString(targetStr, 16)

	pow.target = tmpBigInt

	return &pow
}

//挖矿函数，不断变化nonce，使得sha256(数据+nonce) < 难度值
//返回：区块哈希，nonce
func (pow *ProofofWork) Run() ([]byte, uint64) {
	//定义随机数
	var nonce uint64
	var hash [32]byte
	fmt.Println("开始挖矿了。。。")

	for {
		//time.Sleep(time.Millisecond * 300)
		fmt.Printf("%x\r", hash[:])
		//拼接字符串 + nonce
		data := pow.PrepareData(nonce)
		//hash 值= sha256（data）
		hash = sha256.Sum256(data)

		tmpInt := new(big.Int)
		tmpInt.SetBytes(hash[:])

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//当前计算的哈希.Cmp(难度值)
		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功,hash :%x, nonce :%d\n", hash[:], nonce)
			break
		} else {
			//如果不小于难度值
			nonce++
		}
	} //for

	return hash[:], nonce
}

//拼接nonce和block数据
func (pow *ProofofWork) PrepareData(nonce uint64) []byte {
	b := pow.block

	tmp := [][]byte{
		uintToByte(b.Version),
		b.PrevHash,
		b.MerkleRoot,//计算所有的交易得出，用于Hash计算
		uintToByte(b.TimeStamp),
		uintToByte(b.Bits),
		uintToByte(nonce),
		//b.Hash,
		//b.Data,
	}
	data := bytes.Join(tmp, []byte{})

	return data
}

//判断生成的区块是否有效
func (pow *ProofofWork) IsValid() bool {
	//获取区块
	//拼装数据
	data := pow.PrepareData(pow.block.Nonce)
	//计算sha256
	hash := sha256.Sum256(data)

	//与难度值进行比较
	tmpInt := new(big.Int)
	tmpInt.SetBytes(hash[:])

	//满足条件 返回true
	return tmpInt.Cmp(pow.target) == -1
}
