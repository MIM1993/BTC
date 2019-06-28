package main

import "fmt"

func (cli *CLI) addBlock(data string) {
	fmt.Println("添加区块链被调用")
	bc, err := GetBlockChainInstance()
	if err!=nil{
		fmt.Println("GetBlockChainInstance failed:", err)
		return
	}

	if err = bc.AddBlockToChain(data); err != nil {
		fmt.Println("addBlock failed:", err)
		return
	}
	fmt.Println("添加区块成功！")
}

func (cli *CLI)createBlockChain(){
	err :=CreateBlockChain()
	if err!=nil{
		fmt.Println("CreateBlockChain err :",err)
		return
	}
	fmt.Println("创建区块链成功")
}

func (cli *CLI)print(){
	//实例化区块链
	bc,err := GetBlockChainInstance()
	if err!=nil{
		fmt.Println("GetBlockChainInstance failed:", err)
		return
	}

	//生成迭代器
	it :=bc.NewIterator()
	for {
		//调用next
		block := it.Next()

		fmt.Printf("\n===============================================\n")
		fmt.Printf("PrevHash : %x\n", block.PrevHash)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Bits : %d\n", block.Bits)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", string(block.Data))

		//判断区块是否有效
		pow := NewProofofWork(block)
		fmt.Printf("Isalid:%v \n", pow.IsValid())

		//退出条件
		if block.PrevHash == nil {
			fmt.Println("区块链遍历结束!")
			break
		}
	}
}