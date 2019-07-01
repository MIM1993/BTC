package main

import (
	"os"
	"fmt"
	"strconv"
)

//命令行处理
type CLI struct {
}

const Usage = `
正确的使用用法：
	./blockchain create                   				  "创建区块链"
	./blockchain addBlock <需要写入的数据>   				  "添加区块"
	./blockchain print                     				  "打印区块链"
	./blockchain getBalance <地址> 						  "获取余额"
	./blockchain send <FROM> <TO> <AMOUNT> <MINER> <DATA> "转账"
	./blockchain createwallet 							  "创建钱包"
    ./blockchain listwallet   							  "显示钱包中所有地址"
`

///负责解析命令
func (cli *CLI) Run() {
	cmds := os.Args

	//用户至少输入两个参数
	if len(cmds) < 2 {
		fmt.Println("输入的参数无效，请检查！")
		fmt.Println(Usage)
		return
	}

	switch cmds[1] {
	case "create":
		fmt.Println("创建区块链被调用")
		cli.createBlockChain()

	case "addBlock":
		if len(cmds) != 3 {
			fmt.Println("输入的参数无效，请检查！")
			fmt.Println(Usage)
			return
		}
		data := cmds[2]
		cli.addBlock(data)

	case "print":
		fmt.Println("打印区块链被调用")
		cli.print()

	case "getBalance":
		fmt.Println("获取余额")
		if len(cmds) != 3 {
			fmt.Println("输入的参数无效，请检查！")
			fmt.Println(Usage)
			return
		}
		address := cmds[2]
		cli.getBalance(address)
	case "send":
		fmt.Println("send 命令被调用")
		if len(cmds) != 7 {
			fmt.Println("send 参数无效")
			return
		}
		from := cmds[2]
		to := cmds[3]
		//这个是金额，float64，命令接收都是字符串，需要转换
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.send(from, to, amount, miner, data)
	case "createwallet":
		fmt.Println("createwallet 命令被调用")
		cli.createwallet()
	case "listwallet":
		fmt.Println("listwallet 命令被调用")
		cli.listwallet()
	default:
		fmt.Println("输入的参数无效，请检查！")
		fmt.Println(Usage)
	}
}
