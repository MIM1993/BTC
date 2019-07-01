package main

import (
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"fmt"
	"time"
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

type TXInput struct {
	//这个input所引用的output所在的交易id
	Txid []byte
	//这个input所引用的output所在的交易中的索引
	Index int
	//付款人对当前交易的签名（新交易）
	ScriptSig string
}

type TXOutput struct {
	//收款人的公钥哈希值
	ScriptPubk string
	//转账金额
	Value float64
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
	input := TXInput{Txid: nil, Index: -1, ScriptSig: data}

	//输出交易
	output := TXOutput{Value: reward, ScriptPubk: miner}

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

//判断挖矿交易
func (tx *Transaction) isCoinbaseTx() bool {
	inputs := tx.TXInputs

	if len(inputs) == 1 && inputs[0].Txid == nil && inputs[0].Index == -1 {
		return true
	}

	return false
}

type UTXOInfo struct {
	//交易id
	Txid []byte
	//索引值
	Index int64
	//output
	Output TXOutput
}

//获取指定地址的金额，实现遍历账本的函数
func (bc *BlockChain) FindMyUTXO(address string) []UTXOInfo {
	var utxoinfos []UTXOInfo
	//定义一个已经消耗过的utxo集合
	spendUtxos := make(map[string][]int)

	it := bc.NewIterator()

	for {
		//遍历区块
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transactions {

		LABEL:
		//遍历output
			for outputIndex, output := range tx.TXOutputs {

				//打印
				fmt.Println("outputIndex :", outputIndex)

				if output.ScriptPubk == address {

					//开始过滤
					//大年甲乙
					currentTxid := string(tx.Txid)

					//去spentUtxo中查看
					indexArray := spendUtxos[currentTxid]

					//篮子中有数据
					if len(indexArray) != 0 {
						for _, spendIndex := range indexArray {
							//判断下标
							if outputIndex == spendIndex {
								//跳过
								//goto LABEL
								continue LABEL
							}

						} //for
					} //if

					utxoinfo := UTXOInfo{Txid: tx.Txid, Index: int64(outputIndex), Output: output}

					//找到属于目标地址的output
					utxoinfos = append(utxoinfos, utxoinfo)

				} //if

			} //for

			//---------------------------遍历input------------------------------------------

			if tx.isCoinbaseTx() {
				fmt.Println("发现挖矿交易，无需遍历inputs")
				continue
			}

			for _, input := range tx.TXInputs {
				if input.ScriptSig == address {
					//map         key：交易地址  value：数组下标
					//map [string][]int{    "ox333":{0,1}  }
					spentKey := string(input.Txid)
					//向篮子添加已经消耗的outpt           []int             int
					spendUtxos[spentKey] = append(spendUtxos[spentKey], input.Index)

				} //if

			} //for

		} //for

		//退出for循环条件
		if len(block.PrevHash) == 0 {
			break
		}
	} //for

	//返回
	return utxoinfos
}

//使用Transaction改写程序
//获取挖矿人的金额
//创建普通交易

//转账

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//遍历账本，找到满足条件的utxo集合，返回总金额
	var spentUTXO = make(map[string][]int64) //包含utxo
	var retValue float64                     //总金额

	//遍历账本，找到能够使用的utxo，以及金额
	spentUTXO, retValue = bc.findNeedUTXO(from, amount)

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
			input := TXInput{[]byte(txid), int(i), to}
			inputs = append(inputs, input)
		}
	}

	//拼接outputs   创建属于收款方的output  找零
	//创建一个属于to的output
	output1 := TXOutput{to, amount}
	//添加进入outputs
	outputs = append(outputs, output1)

	//> 如果总金额大于需要转账的金额，进行找零：给from创建一个output
	if retValue > amount {
		output2 := TXOutput{from, retValue - amount}
		outputs = append(outputs, output2)
	}

	//时间戳
	timeStamp := time.Now().Unix()

	//设置哈希 返回
	tx := Transaction{nil, inputs, outputs, uint64(timeStamp)}

	//设置txid
	tx.setHash()

	return &tx
}
