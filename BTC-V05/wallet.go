package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
	"fmt"
	"bytes"
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
		fmt.Println("创建私钥失败,err:", err)
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

	//公钥hash
	pubKeyHash := getPubKeyHashFromPubKey(pubKey)

	//拼接一个字节
	payload := append([]byte{byte(0)}, pubKeyHash...)

	//生成校验码
	checksum := checkSum(payload)

	//将数据拼起来共==》25字节数据
	payload = append(payload, checksum...)

	//进行base58编码
	address := base58.Encode(payload)

	//返回数据
	return address
}

//通过给定公钥，得到公钥哈希值
func getPubKeyHashFromPubKey(pubKey []byte) []byte {
	hash1 := sha256.Sum256(pubKey)

	//hash160处理
	hasher := ripemd160.New()
	_, err := hasher.Write(hash1[:])
	if err != nil {
		fmt.Println("hasher.Write err:", err)
		return nil
	}

	//公钥hash
	pubKeyHash := hasher.Sum(nil)

	return pubKeyHash
}

//得到四字节验证码
func checkSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	return second[:4]
}

/*锁定output的是公钥Hash值，不是地址，所以要写一个函数，
  可以通过地址推出公钥Hash值 */
func getPubKeyHashFromAddress(address string) []byte {
	//base58解码
	decodeInfo := base58.Decode(address)

	//校验数据
	if len(decodeInfo) != 25 {
		fmt.Println("传入地址无效")
		return nil
	}

	//截取
	pubKeyHash := decodeInfo[1:len(decodeInfo)-4]

	return pubKeyHash
}

//校验地址
func isValidAddress(address string) bool {
	// 	解码，得到25字节数据
	decodeInfo := base58.Decode(address)

	if len(decodeInfo) != 25 {
		fmt.Println("isValidAddress, 传入地址长度无效")
		return false
	}

	// 截取前21字节payload，截取后四字节checksum1
	payload := decodeInfo[:len(decodeInfo)-4]   //21字节
	checksum1 := decodeInfo[len(decodeInfo)-4:] //4字节

	// 对palyload计算，得到checksum2，与checksum1对比，true校验成功，反之失败
	checksum2 := checkSum(payload)
	return bytes.Equal(checksum1, checksum2)
}
