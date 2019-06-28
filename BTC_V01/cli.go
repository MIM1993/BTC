package main

import (
	"os"
	"fmt"
)

//命令行处理
type CLI struct {
}

const Usage = `
正确的使用用法：
	./blockchain create "创建区块链"
	./blockchain addBlock <需要写入的数据> "添加区块"
	./blockchain print "打印区块链"
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

	default:
		fmt.Println("输入的参数无效，请检查！")
		fmt.Println(Usage)
	}
}
