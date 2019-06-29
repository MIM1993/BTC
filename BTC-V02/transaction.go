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

//获取指定地址的金额，实现遍历账本的函数
func (bc *BlockChain) FindMyUTXO(address string) []TXOutput {
	var utxos []TXOutput
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

					//找到属于目标地址的output
					utxos = append(utxos, output)

				} //if

			} //for

			//---------------------------遍历input------------------------------------------
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
	return utxos
}

//使用Transaction改写程序
//获取挖矿人的金额
//创建普通交易
//转账
