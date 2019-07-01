package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
	"fmt"
)

//定义结构
type wallet struct {
	PriKey *ecdsa.PrivateKey
	PruKey []byte
}

//创建密钥对
func newWalletKeyPair() *wallet {
	curve := elliptic.P256()

	//创建私钥
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("创建私钥失败,err:",err)
		return nil
	}

	//获取公钥
	pubKeyRaw := priKey.PublicKey

	//拼接x，y
	pubKey := append(pubKeyRaw.X.Bytes(), pubKeyRaw.Y.Bytes()...)

	//创建wallet结构体
	wallet := wallet{PriKey: priKey, PruKey: pubKey}

	return &wallet
}

//根据私钥生成地址
func (w *wallet) getAddress() string {
	//公钥
	pubKey := w.PruKey
	hash1 := sha256.Sum256(pubKey)

	//hash160处理
	hasher := ripemd160.New()
	_,err :=hasher.Write(hash1[:])
	if err!=nil{
		fmt.Println("hasher.Write err:",err)
		return ""
	}

	//公钥hash
	pubKeyHash := hasher.Sum(nil)

	//拼接
	payload := append([]byte{byte(0)}, pubKeyHash...)

	//生成校验码
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	//取校验码的前四字节
	checksum := second[:4]

	//将数据拼起来共==》25字节数据
	payload = append(payload, checksum...)

	//进行base58编码
	address := base58.Encode(payload)

	//返回数据
	return address
}
