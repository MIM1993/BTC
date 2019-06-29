package main

import (
	"bolt"
	"fmt"
	"errors"
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
func CreateBlockChain() error {
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
			coinbase := NewCoinbaseTx("中本聪", genesisinfo)
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
	var lastHash []byte //内存中最后一个区块的哈希值

	// 打开数据库
	db, err := bolt.Open(blockchainDBFile, 0400, nil) //rwx  0100 => 4
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
//func (bc *BlockChain) AddBlockToChain(data string) error {
//
//	//区块链中最后一个去区块的Hash值
//	lastBlockHash := bc.tail
//
//	//创建区块
//	newblock := NewBlock(data, lastBlockHash)
//
//	//写入数据库
//	err := bc.db.Update(func(tx *bolt.Tx) error {
//		bucket := tx.Bucket([]byte(bucketBlock))
//		if bucket == nil {
//			return errors.New("AddBlock时Bucket不应为空")
//		}
//
//		//key是新区块的哈希值， value是这个区块的字节流
//		bucket.Put(newblock.Hash, newblock.Serialize())
//
//		//更新lastBlockHashKey
//		bucket.Put([]byte(lastBlockHashKey), newblock.Hash)
//
//		//更新bc的tail值,用于以后追加
//		bc.tail = newblock.Hash
//
//		return nil
//	})
//
//	return err
//}

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
