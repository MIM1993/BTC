package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"io/ioutil"
	"sort"
)

type WalletManager struct {
	//定义一个map来管理所有的钱包
	//key:地址
	//value:wallet结构(公钥，私钥)
	Wallets map[string]*wallet
}

//创建WalletManager结构
func NewWalletManager() *WalletManager {
	//创建一个, Wallets map[string]*wallet
	var wm WalletManager
	//分配空间
	wm.Wallets = make(map[string]*wallet)

	//从本地加载数据，写如map中
	if !wm.loadFile() {
		return nil
	}

	return &wm
}

//创建钱包地址
func (wm *WalletManager) createwallet() string {
	//创建密钥对
	w := newWalletKeyPair()
	if w == nil {
		fmt.Println("newWalletKeyPair 失败")
		return ""
	}

	//根据私钥生成地址
	address := w.getAddress()
	if address == "" {
		fmt.Println("getAddress 失败")
		return ""
	}

	//把地址写入map中
	wm.Wallets[address] = w

	//将秘钥写入磁盘
	if !wm.saveFile() {
		return ""
	}

	//返回地址
	return address
}

const walletFile = "wallet.dat"

//保存私钥
func (wm *WalletManager) saveFile() bool {
	//使用gob进行编码
	var buffer bytes.Buffer

	//注册接口，用于解码
	gob.Register(elliptic.P256())
	//创建编码器
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wm)
	if err != nil {
		fmt.Println("Encode err：", err)
		return false
	}

	//保存到磁盘中
	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println("WriteFile err:", err)
		return false
	}

	return true
}

//读取文件wallet.dat,将文件内容存到walletManager中
func (wm *WalletManager) loadFile() bool {
	//判断文件是否存在
	if !IsFileExist(walletFile) {
		fmt.Println("文件不存在，无需加载")
		return true //必须是true，不然会陷入死循环，注意
	}

	//读取文件
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return false
	}

	//将文件解码到wm中
	//注册接口到解码函数中
	gob.Register(elliptic.P256())
	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(content))
	//解码  直接将数据解码到wm中的map里
	err = decoder.Decode(wm)
	if err != nil {
		fmt.Println("Decode err :", err)
		return false
	}

	return true
}

//将wallet.dat中的数据打印出来
func (wm *WalletManager) listAdderss() []string {
	//定义容器
	addresses := []string{}

	//循环读取wm中的map
	for address := range wm.Wallets {
		addresses = append(addresses, address)
	}

	//排序
	sort.Strings(addresses)

	return addresses
}
