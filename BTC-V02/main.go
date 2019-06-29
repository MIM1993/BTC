package main

//
////定义区块结构
////实现基础字段
////补充字段
//type Block struct {
//	//前区块
//	PrevHash []byte
//	//哈希
//	Hash []byte
//	//数据
//	Data []byte
//}
//
////创建区块方法		区块       前hash值
//func NewBlock(data string, PrevHash []byte) *Block {
//	b := Block{
//		PrevHash: PrevHash,
//		Hash:     nil,
//		Data:     []byte(data),
//	}
//
//	//计算hash值
//	b.setHash()
//
//	return &b
//}
//
////提供计算区块hash值的方法
//func (b *Block) setHash() {
//	//data 是block各个字节流进行拼接
//	//二维切片
//	temp := [][]byte{
//		b.PrevHash,
//		b.Hash,
//		b.Data,
//	}
//	//Join函数
//	data := bytes.Join(temp, []byte{})
//
//	//比特币  sha256
//	hash := sha256.Sum256(data)
//
//	//赋值
//	b.Hash = hash[:]
//}
//
////定义区块连结构(使用数组模拟)
//type BlockChain struct {
//	Blocks []*Block //区块链
//}
//
//const genesisinfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
//
////提供一个创建区块链方法
//func NewBlockChain() *BlockChain {
//	genesisBlock := NewBlock(genesisinfo, nil)
//
//	bc := BlockChain{
//		Blocks: []*Block{genesisBlock},
//	}
//
//	return &bc
//
//}
//
////提供一个添加区块到链中的方法
//func (blockchain *BlockChain) AddBlockToChain(block *Block) {
//
//	blockchain.Blocks = append(blockchain.Blocks, block)
//}

//打印区块
func main() {
	////创建区块链，并且实例化
	//err := CreateBlockChain()
	//if err != nil {
	//	fmt.Println("CreateBlockChain err:", err)
	//	return
	//}
	//
	////获取区块链实例
	//bc, err := GetBlockChainInstance()
	//if err != nil {
	//	fmt.Println("GetBlockChainInstance err :", err)
	//	return
	//}
	//
	////添加区块
	//err = bc.AddBlockToChain("hello world !")
	//if err != nil {
	//	fmt.Println("AddBlockToChain err :", err)
	//	return
	//}
	//err = bc.AddBlockToChain("财富自由，就看今天！")
	//if err != nil {
	//	fmt.Println("AddBlockToChain err :", err)
	//	return
	//}
	//
	////调用跌倒器 输出blockchain
	//it := bc.NewIterator()
	//for {
	//	//调用next
	//	block := it.Next()
	//
	//	fmt.Printf("\n===============================================\n")
	//	fmt.Printf("PrevHash : %x\n", block.PrevHash)
	//	fmt.Printf("Version : %d\n", block.Version)
	//	fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
	//	fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
	//	fmt.Printf("Bits : %d\n", block.Bits)
	//	fmt.Printf("Nonce : %d\n", block.Nonce)
	//	fmt.Printf("Hash : %x\n", block.Hash)
	//	fmt.Printf("Data : %s\n", string(block.Data))
	//
	//	//判断区块是否有效
	//	pow := NewProofofWork(block)
	//	fmt.Printf("Isalid:%v \n", pow.IsValid())
	//
	//	//退出条件
	//	if block.PrevHash == nil {
	//		fmt.Println("打印接受")
	//		break
	//	}
	//}

	cli := CLI{}
	cli.Run()
}
