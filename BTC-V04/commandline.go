package main

import "fmt"

//添加区块
func (cli *CLI) addBlock(data string) {
	fmt.Println("添加区块链被调用")
	//bc, err := GetBlockChainInstance()
	//if err != nil {
	//	fmt.Println("GetBlockChainInstance failed:", err)
	//	return
	//}
	//
	//if err = bc.AddBlockToChain(data); err != nil {
	//	fmt.Println("addBlock failed:", err)
	//	return
	//}
	//fmt.Println("添加区块成功！")
}

//创建区块链
func (cli *CLI) createBlockChain() {
	err := CreateBlockChain()
	if err != nil {
		fmt.Println("CreateBlockChain err :", err)
		return
	}
	fmt.Println("创建区块链成功")
}

//打印区块
func (cli *CLI) print() {
	//实例化区块链
	bc, err := GetBlockChainInstance()
	if err != nil {
		fmt.Println("print GetBlockChainInstance failed:", err)
		return
	}

	defer bc.db.Close()

	//生成迭代器
	it := bc.NewIterator()
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
		fmt.Printf("Data : %s\n", block.Transactions[0].TXInputs[0].ScriptSig) //旷工写入的数据

		//判断区块是否有效
		pow := NewProofofWork(block)
		fmt.Printf("Isalid:%v \n", pow.IsValid())

		//退出条件
		if block.PrevHash == nil {
			fmt.Println("区块链遍历结束!")
			break
		}

	} //for
}

//获取余额
func (cil *CLI) getBalance(address string) {

	bc, err := GetBlockChainInstance()
	if err != nil {
		fmt.Println("getBalance GetBlockChainInstance failed:", err)
		return
	}

	defer bc.db.Close()

	//获取相关的utxos
	utxos := bc.FindMyUTXO(address)

	total := 0.0

	//循环
	for _, utxo := range utxos {
		total += utxo.Output.Value
	}

	fmt.Printf("'%s'的金额是%f/n", address, total)
}

//创建交易并添加区块
func (cil *CLI) send(from, to string, amount float64, miner, data string) {
	fmt.Println("from:", from)
	fmt.Println("to:", to)
	fmt.Println("amount:", amount)
	fmt.Println("miner:", miner)
	fmt.Println("data:", data)

	//每次send添加一个数组

	//获取区块链
	bc, err := GetBlockChainInstance()
	if err != nil {
		fmt.Println("send getBalance GetBlockChainInstance failed:", err)
		return
	}

	defer bc.db.Close()

	//创建挖矿交易
	coinbaseTx := NewCoinbaseTx(miner, data)

	//创建交易池
	txs := []*Transaction{coinbaseTx}

	//创建普通交易
	tx := NewTransaction(from, to, amount, bc)
	if tx != nil {
		fmt.Println("找到一笔有效的转账交易!")
		txs = append(txs, tx)
	} else {
		fmt.Println("注意，找到一笔无效的转账交易, 不添加到区块!")
	}

	//调用addblock
	err = bc.AddBlockToChain(txs)
	if err != nil {
		fmt.Println("添加区块失败，转账失败!")
	}

	fmt.Println("添加区块失败，转账成功!")
}

//创建钱包
func (cli *CLI) createwallet() {
	//创建管理钱包句柄
	wm := NewWalletManager()
	if wm == nil {
		fmt.Println("获取walletManager失败")
		return
	}

	//生成地址
	address := wm.createwallet()

	//进行错误判断，防止程序崩溃
	if address == "" {
		fmt.Println("创建钱包失败")
		return
	}

	fmt.Println("新钱包的地址是：", address)
}

//展示钱包中的所有地址
func (cli *CLI) listwallet() {
	//创建管理钱包句柄
	wm := NewWalletManager()

	if wm == nil {
		fmt.Println("钱包数据读取失败")
		return
	}

	//展示所有数据
	addresses := wm.listAdderss()

	//循环打印
	for _, address := range addresses {
		fmt.Println(address)
	}
}
