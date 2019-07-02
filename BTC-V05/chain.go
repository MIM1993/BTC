package main

import (
	"bolt"
	"fmt"
	"errors"
	"crypto/ecdsa"
	"bytes"
)

//定义区块连结构(使用数组模拟)
type BlockChain struct {
	db   *bolt.DB //用于存储数据 //区块链
	tail []byte   //最后一个区块的Hash值
}

//常量
const genesisinfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const blockchainDBFile = "blockchain.db"
const bucketBlock = "bucketBlock" //桶
const lastBlockHashKey = "lastBlockHashKey"

//提供一个创建区块链方法
//func NewBlockChain() (*BlockChain, error) {
//	var lastHash []byte //内存中最后一个区块的哈希值
//
//	//打开数据库
//	db, err := bolt.Open(blockchainDBFile, 0600, nil)
//	if err != nil {
//		fmt.Println("bolt Open err :", err)
//		return nil, err
//	}
//
//	//defer db.Close()
//
//	//存在bucket
//	db.Update(func(tx *bolt.Tx) error {
//		//更新，找到目标bucket
//		bucket := tx.Bucket([]byte(blockchainDBFile))
//
//		//如果bucket不存在组创建，写入创世快
//		if bucket == nil {
//			//创建一个bucket
//			bucket, err = tx.CreateBucket([]byte(blockchainDBFile))
//			if err != nil {
//				fmt.Println("tx.CreateBucket err :", err)
//				return err
//			}
//
//			//写入创世块
//			//创建BlockChain，同时添加一个创世块
//			genesisBlock := NewBlock(genesisinfo, nil)
//			//key是区块的哈希值，value是block的字节流//TODO
//			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()) //将Block序列化
//			//更新最后区块哈希值到数据库
//			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
//
//			//更新内存中的最后一个区块哈希值, 后续操作就可以基于这个哈希增加区块
//			lastHash = genesisBlock.Hash
//		} else {
//			//直接读取则定的key
//			lastHash = bucket.Get([]byte(lastBlockHashKey))
//		}
//
//		return nil
//
//	})
//
//	//拼成blockchain然后返回
//	bc := BlockChain{db, lastHash}
//	return &bc, nil
//}

//创建区块链，持久化到本地，从无到有：这个函数仅执行一次
func CreateBlockChain(address string) error {
	//判断区块链是否存在
	if IsFileExist(blockchainDBFile) {
		fmt.Println("区块链已经存在")
		return nil
	}

	// 1. 区块链不存在，创建
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
	if err != nil {
		return err
	}

	//最后关闭数据库
	defer db.Close()

	// 2. 开始创建
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在
		if bucket == nil {
			//创建bucket
			bucket, err := tx.CreateBucket([]byte(bucketBlock))
			if err != nil {
				return err
			}

			/*第一个区块*/
			//创建挖矿交易
			coinbase := NewCoinbaseTx(address, genesisinfo)
			//拼装txs，挖矿交易切片
			txs := []*Transaction{coinbase}
			//创建创世快
			genesisBlock := NewBlock(txs, nil)

			//key是区块的哈希值，value是block的字节流
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()) //将block序列化
			//更新最后区块哈希值到数据库
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
		}
		return nil
	})
	return err //nil
}

//获取区块链实例，用于后续操作, 每一次有业务时都会调用
func GetBlockChainInstance() (*BlockChain, error) {
	//判断区块链是否存在
	if !IsFileExist(blockchainDBFile) {
		fmt.Println("区块链不存在，请创建区块链")
		return nil, nil
	}

	var lastHash []byte //内存中最后一个区块的哈希值

	// 打开数据库
	db, err := bolt.Open(blockchainDBFile, 0600, nil) //rwx
	if err != nil {
		return nil, err
	}

	//不要db.Close，后续要使用这个句柄

	//区块链一定存在，直接读取
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在，返回错误
		if bucket == nil {
			return errors.New("bucket不应为nil")
		} else {
			//直接读取特定的key，得到最后一个区块的哈希值
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}

		return nil
	})

	//5. 拼成BlockChain然后返回
	bc := BlockChain{db, lastHash}
	return &bc, nil
}

//提供一个添加区块到链中的方法
func (bc *BlockChain) AddBlockToChain(txs []*Transaction) error {
	//校验交易
	fmt.Println("添加区块前，对交易进行校验...")
	//挖矿交易无需校验，过滤掉

	//区块链中最后一个去区块的Hash值
	lastBlockHash := bc.tail

	//创建区块
	newblock := NewBlock(txs, lastBlockHash)

	//写入数据库
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("AddBlock时Bucket不应为空")
		}

		//key是新区块的哈希值， value是这个区块的字节流
		bucket.Put(newblock.Hash, newblock.Serialize())

		//更新lastBlockHashKey
		bucket.Put([]byte(lastBlockHashKey), newblock.Hash)

		//更新bc的tail值,用于以后追加
		bc.tail = newblock.Hash

		return nil
	})

	return err

}

//----------------------------------------------------------------------------------
//定义迭代器
type Iterator struct {
	db          *bolt.DB
	currentHash []byte //不断移动的哈希值，由于访问所有区块
}

//生成迭代器
func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}
	return &it
}

//动起来  Next方法  1、返回指向区块 2、向左移动
func (it *Iterator) Next() (block *Block) {

	//block := &Block{}

	//读取bucket当前hash数据
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		if bucket == nil {
			return errors.New("不能为空！")
		}

		blockTmpInfo := bucket.Get(it.currentHash)

		block = Deserialize(blockTmpInfo)

		//指针向左移动
		it.currentHash = block.PrevHash

		return nil
	})

	if err != nil {
		fmt.Println("Iterator Next err :", err)
		return nil
	}

	return
}

//定义一个结构，包含output的详情：output本身，位置信息
type UTXOInfo struct {
	//交易id
	Txid []byte
	//索引值
	Index int64
	//output
	Output TXOutput
}

//获取指定地址的金额，实现遍历账本的函数
func (bc *BlockChain) FindMyUTXO(pubKeyHash []byte) []UTXOInfo {
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

				if bytes.Equal(output.ScriptPubKeyHash, pubKeyHash) {

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
				hash := getPubKeyHashFromPubKey(input.PubKey)

				if bytes.Equal(hash, pubKeyHash) {
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

//获取需要使用的utxo                  地址          金额
func (bc *BlockChain) findNeedUTXO(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {
	retMap := make(map[string][]int64)
	var retValue float64

	//遍历账本  找到所有utxo
	utxoinfos := bc.FindMyUTXO(pubKeyHash)

	//遍历utxo，统计当前总额
	for _, utxoinfo := range utxoinfos {
		//统计当前总额
		retValue += utxoinfo.Output.Value

		//统计将要返回的utxo    key是txid,value是int64切片
		key := string(utxoinfo.Txid)
		retMap[key] = append(retMap[key], utxoinfo.Index)

		//判断金额，与amount比较  大于就返回  小于继续
		if retValue >= amount {
			break
		}
	}

	//返回 utxo信息切片   总金额
	return retMap, retValue

}

//	签名函数
func (bc *BlockChain) signTransaction(tx *Transaction, priKey *ecdsa.PrivateKey) bool {
	//遍历所有的input，通过id获取对应的交易，  map  id==》key   交易==》value

	//定义容器，储存所需的前交易
	prevTxs := make(map[string]*Transaction)

	//遍历账本，通过交易tx,获取需要的前交易
	for _, input := range tx.TXInputs {
		//定义新函数通过input.txid获取前交易
		prevTx := bc.findTransaction(input.Txid)

		fmt.Println("找到了引用的交易")

		//将查询出来的交易放入容器中
		prevTxs[string(input.Txid)] = prevTx
	}

	//调用签名函数
	result := tx.sign(priKey, prevTxs)

	return result
}

//通过input.txid获取前交易
func (bc *BlockChain) findTransaction(txid []byte) *Transaction {
	//遍历账本，获取前交易

	//利用迭代器
	it := bc.NewIterator()

	for {
		block := it.Next()

		for _, transaction := range block.Transactions {
			if bytes.Equal(transaction.Txid, txid) {
				return transaction
			} //if

		} //for

		//返回条件
		if len(block.PrevHash) == 0 {
			break
		}

	} //for

	return nil
}

//校验单笔交易
func (bc *BlockChain) verifTransaction(tx *Transaction) bool {
	fmt.Println("交易校验开始了...")

	//挖矿交易无需校验，因为无input
	if tx.isCoinbaseTx() {
		fmt.Println("发现挖矿交易，无需校验！")
		return true
	}

	//根据传递进来的tx得到所有的所需的交易集合
	prevTxs := make(map[string]*Transaction)

	//遍历账本，找到所有交易集合
	for _, input := range tx.TXInputs {
		prevTx := bc.findTransaction(input.Txid)
		if prevTx == nil {
			fmt.Println("没有找到有效的引用交易")
			return false
		}

		fmt.Println("找到了有效的引用交易")

		prevTxs[string(input.Txid)] = prevTx
	}

	//调用具体的校验函数
	return tx.verify(prevTxs)
}
