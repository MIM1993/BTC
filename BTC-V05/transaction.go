package main

import (
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"fmt"
	"time"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
	"strings"
)

//定义交易结构
type Transaction struct {
	//交易id
	Txid []byte
	//多个输入
	TXInputs []TXInput
	//多个输出
	TXOutputs []TXOutput
	//交易时间  时间戳
	TimeStamp uint64
}

//input
type TXInput struct {
	//这个input所引用的output所在的交易id
	Txid []byte
	//这个input所引用的output所在的交易中的索引
	Index int

	//付款人对当前交易的签名（新交易）
	//ScriptSig string

	//付款人对当前交易的签名
	ScriptSig []byte

	//公钥
	PubKey []byte
}

//output
type TXOutput struct {
	//收款人的公钥哈希值
	//ScriptPubk string
	ScriptPubKeyHash []byte

	//转账金额
	Value float64
}

//没办法直接赋值，需要直接创建output方法
func newOutput(address string, amount float64) TXOutput {
	output := TXOutput{Value: amount}

	//获取公钥hash
	pubKeyHash := getPubKeyHashFromAddress(address)

	output.ScriptPubKeyHash = pubKeyHash

	return output
}

//获取交易id（对交易做哈希）
func (tx *Transaction) setHash() error {
	//定义tx做gob编码得到字节流  做sha256   赋值给TXID
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(tx)
	if err != nil {
		fmt.Println("encode.Encode err :", err)
		return err
	}

	//进行sha256计算，得出的Hash值作物tx.Txid
	hash := sha256.Sum256(buffer.Bytes())

	//tx字节流的Hash值作为交易id
	tx.Txid = hash[:]

	return nil
}

//挖矿奖励
var reward = 12.5

//创建挖矿交易（input特殊）
func NewCoinbaseTx(miner /*矿工*/ string, data string) *Transaction {
	/*没有输入 只有输出  得到挖矿交易
	//挖矿交需要能够识别出来
	//挖矿交易不需要签名，所以签名可以任意写*/

	//输入交易
	input := TXInput{Txid: nil, Index: -1, ScriptSig: nil, PubKey: []byte(data)}

	//输出交易
	//output := TXOutput{Value: reward, ScriptPubk: miner}
	output := newOutput(miner, reward)

	//时间戳
	timeStamp := time.Now().Unix()

	//填写信息
	tx := Transaction{
		Txid:      nil,
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
		TimeStamp: uint64(timeStamp),
	}

	//计算tx.Txid
	tx.setHash()

	//返回
	return &tx
}

//判断是否是挖矿交易
func (tx *Transaction) isCoinbaseTx() bool {
	inputs := tx.TXInputs

	if len(inputs) == 1 && inputs[0].Txid == nil && inputs[0].Index == -1 {
		return true
	}

	return false
}

//
////定义一个结构，包含output的详情：output本身，位置信息
//type UTXOInfo struct {
//	//交易id
//	Txid []byte
//	//索引值
//	Index int64
//	//output
//	Output TXOutput
//}
//
////获取指定地址的金额，实现遍历账本的函数
//func (bc *BlockChain) FindMyUTXO(pubKeyHash []byte) []UTXOInfo {
//	var utxoinfos []UTXOInfo
//	//定义一个已经消耗过的utxo集合
//	spendUtxos := make(map[string][]int)
//
//	it := bc.NewIterator()
//
//	for {
//		//遍历区块
//		block := it.Next()
//		//遍历交易
//		for _, tx := range block.Transactions {
//
//		LABEL:
//		//遍历output
//			for outputIndex, output := range tx.TXOutputs {
//
//				//打印
//				fmt.Println("outputIndex :", outputIndex)
//
//				if bytes.Equal(output.ScriptPubKeyHash, pubKeyHash) {
//
//					//开始过滤
//					//大年甲乙
//					currentTxid := string(tx.Txid)
//
//					//去spentUtxo中查看
//					indexArray := spendUtxos[currentTxid]
//
//					//篮子中有数据
//					if len(indexArray) != 0 {
//						for _, spendIndex := range indexArray {
//							//判断下标
//							if outputIndex == spendIndex {
//								//跳过
//								//goto LABEL
//								continue LABEL
//							}
//
//						} //for
//					} //if
//
//					utxoinfo := UTXOInfo{Txid: tx.Txid, Index: int64(outputIndex), Output: output}
//
//					//找到属于目标地址的output
//					utxoinfos = append(utxoinfos, utxoinfo)
//
//				} //if
//
//			} //for
//
//			//---------------------------遍历input------------------------------------------
//
//			if tx.isCoinbaseTx() {
//				fmt.Println("发现挖矿交易，无需遍历inputs")
//				continue
//			}
//
//			for _, input := range tx.TXInputs {
//				hash := getPubKeyHashFromPubKey(input.PubKey)
//
//				if bytes.Equal(hash, pubKeyHash) {
//					//map         key：交易地址  value：数组下标
//					//map [string][]int{    "ox333":{0,1}  }
//					spentKey := string(input.Txid)
//					//向篮子添加已经消耗的outpt           []int             int
//					spendUtxos[spentKey] = append(spendUtxos[spentKey], input.Index)
//
//				} //if
//
//			} //for
//
//		} //for
//
//		//退出for循环条件
//		if len(block.PrevHash) == 0 {
//			break
//		}
//	} //for
//
//	//返回
//	return utxoinfos
//}

//使用Transaction改写程序
//获取挖矿人的金额
//创建普通交易

//转账

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//使用钱包  from
	wm := NewWalletManager()
	if wm == nil {

	}

	//找到对应的wallet key是地址
	wallet, ok := wm.Wallets[from]
	if !ok {
		return nil
	}

	//拿到公钥私钥
	pubKey := wallet.PruKey
	priKey := wallet.PriKey

	//获取公钥hash，所有的output的定为，是由公钥hash确定的
	pubKeyHash := getPubKeyHashFromPubKey(pubKey)

	//遍历账本，找到满足条件的utxo集合，返回总金额
	var spentUTXO = make(map[string][]int64) //包含utxo
	var retValue float64                     //总金额

	//遍历账本，找到能够使用的utxo，以及金额
	spentUTXO, retValue = bc.findNeedUTXO(pubKeyHash, amount)

	//如果金额不足，失败
	if retValue < amount {
		fmt.Println("金额不足")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//拼接input   遍历utxo，转化为output
	for txid, indexArray := range spentUTXO {
		for _, i := range indexArray {
			input := TXInput{Txid: []byte(txid), Index: int(i), ScriptSig: nil, PubKey: pubKey}
			inputs = append(inputs, input)
		}
	}

	//拼接outputs   创建属于收款方的output  找零
	//创建一个属于to的output
	//output1 := TXOutput{to, amount}
	output1 := newOutput(to, amount)
	//添加进入outputs
	outputs = append(outputs, output1)

	//> 如果总金额大于需要转账的金额，进行找零：给from创建一个output
	if retValue > amount {
		//output2 := TXOutput{from, retValue - amount}
		output2 := newOutput(from, retValue-amount)
		outputs = append(outputs, output2)
	}

	//时间戳
	timeStamp := time.Now().Unix()

	//设置哈希 返回
	tx := Transaction{nil, inputs, outputs, uint64(timeStamp)}
	//设置txid
	tx.setHash()

	if !bc.signTransaction(&tx, priKey) {
		fmt.Println("交易签名失败")
		return nil
	}

	return &tx
}

//实现具体的签名动作（copy 质控  签名）
func (tx *Transaction) sign(priKey *ecdsa.PrivateKey, prvTxs map[string]*Transaction) bool {
	fmt.Println("对交易进行具体的签名sign...")
	//判断是否是挖矿交易，挖矿交易没用inputs
	if !tx.isCoinbaseTx() {
		fmt.Println("找到挖矿交易，无需签名")
		return true
	}

	//获取交易copy 同时置空
	txCopy := tx.trimmedCpopy()

	//遍历copy交易的input
	for i, input := range txCopy.TXInputs {
		fmt.Printf("开始对input[%s]进行签名", i)

		//遍历发送过来的账本
		prevTx := prvTxs[string(input.Txid)]
		if prevTx == nil {
			return false
		}

		//获取引用的input引用的output
		output := prevTx.TXOutputs[input.Index]

		//获取output的公钥hash，赋值给txcopy中对应的input
		txCopy.TXInputs[i].PubKey = output.ScriptPubKeyHash

		//对copy进行签名 ,得到需要的hash值
		txCopy.setHash()

		//将pubKey设置为空
		txCopy.TXInputs[i].PubKey = nil

		//签名的具体数据
		hashData := txCopy.Txid

		//开始签名
		r, s, err := ecdsa.Sign(rand.Reader, priKey, hashData)
		if err != nil {
			fmt.Println("签名失败")
			return false
		}

		//将r,s追加到一起
		signature := append(r.Bytes(), s.Bytes()...)

		//将数子签名赋值给原始交易
		tx.TXInputs[i].ScriptSig = signature

	}

	fmt.Println("交易签名成功！")
	return true
}

//trim修剪  设置为空
func (tx *Transaction) trimmedCpopy() *Transaction {

	var (
		inputs  []TXInput
		outputs []TXOutput
	)

	//创建一个交易副本 其中数据设置为空
	for _, input := range tx.TXInputs {
		input := TXInput{
			Txid:      input.Txid,
			Index:     input.Index,
			ScriptSig: nil,
			PubKey:    nil,
		}
		inputs = append(inputs, input)
	}

	//output是一样的
	outputs = tx.TXOutputs

	//定义返回的副本交易
	txCopy := Transaction{
		Txid:      tx.Txid,
		TXInputs:  inputs,
		TXOutputs: outputs,
		TimeStamp: tx.TimeStamp,
	}

	//返回
	return &txCopy
}

//具体校验函数
func (tx *Transaction) verify(prevTxs map[string]*Transaction) bool {
	//获得交易副本
	txCopy := tx.trimmedCpopy()
	//遍历交易
	for i, input := range tx.TXInputs {
		prevTx := prevTxs[string(input.Txid)]
		if prevTx == nil {
			return false
		}

		//还原数据
		output := prevTx.TXOutputs[input.Index]
		txCopy.TXInputs[i].PubKey = output.ScriptPubKeyHash
		txCopy.setHash()

		//清零数据
		txCopy.TXInputs[i].PubKey = nil

		//具体还原的数据签名哈希值（数据 ）
		hashData := txCopy.Txid

		//签名
		signature := input.ScriptSig

		//储存公钥信息的数据，用于还原公钥
		pubKey := input.PubKey

		//开始校验
		var r, s, x, y big.Int
		//r,s从signarture中截取出来
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//x, y 从pubkey截取除来，还原为公钥本身
		x.SetBytes(pubKey[:len(pubKey)/2])
		y.SetBytes(pubKey[len(pubKey):])
		//曲线
		curve := elliptic.P256()
		pubKeyRaw := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		//进行校验
		res := ecdsa.Verify(&pubKeyRaw, hashData, &r, &s)
		if !res {
			fmt.Println("发现校验失败的input")
			return false
		}

	} //for

	fmt.Println("交易校验成功")

	return true
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.Txid))

	for i, input := range tx.TXInputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.ScriptSig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.TXOutputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.ScriptPubKeyHash))
	}

	return strings.Join(lines, "\n")
}
