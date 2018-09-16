package main

import (
	"fmt"
	"os"
	"strconv"
	"flag"
	"github.com/labstack/gommon/log"
)

//命令行接口
type CLI struct {
	blockchain *BlockChain
}
//用法 cli.printUsage 实例化
func (cli *CLI)printUsage(){
	fmt.Println("用法如下")
	fmt.Println("addblock 向区块链增加块")
	fmt.Println("showchain 显示区块链")
}
//不符合规范
func (cli *CLI)validateArgs(){
	if len(os.Args)<2{
		cli.printUsage()//显示用法
		os.Exit(1)
	}
}
// 命令行实例化
func (cli *CLI)addBlock(data string){
	cli.blockchain.AddBlock(data)//增加区块
	fmt.Println("区块增加成功")
}
//循环打印函数 不断的每一个区块都是以其hash值作为索引 本身的区块信息有包含前一区块的hash值 不断的赋上一区块的值 直到
//创世区块 1.从新到旧 2.顺序打印----问题后期数据库查找怎么办 查看每笔交易
func (cli *CLI)showBlockChain(){
	bci:=cli.blockchain.Iterator()//创建循环迭代器
	for{
		block:=bci.next()//取得下一个区块
		fmt.Printf("上一块哈希%x\n",block.PrevBlockHash)
		fmt.Printf("数据：%s\n",block.Data)
		fmt.Printf("当前哈希%x\n",block.Hash)
		pow:=NewProofOfWork(block)
		fmt.Printf("pow %s",strconv.FormatBool(pow.Validate()))
		fmt.Println("\n")

		if len(block.PrevBlockHash)==0{//遇到创世区块终止
			break
		}
	}
}
//命令行触发 adddata showchain 新加 显示
func (cli *CLI)Run(){
	cli.validateArgs()//校验
	//处理命令行参数
	//**.exe showchain  name表示参数
	addblockcmd:=flag.NewFlagSet("addblock",flag.ExitOnError)
	showchaincmd:=flag.NewFlagSet("showchain",flag.ExitOnError)
	//**.exe addblock -data ""
	addBlockData:=addblockcmd.String("data","","Block data")
	switch os.Args[1]{
	case "addblock":
		err:=addblockcmd.Parse(os.Args[2:])//解析参数
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"showchain":
		err:=showchaincmd.Parse(os.Args[2:])//解析参数
		if err!=nil{
			log.Panic(err)//处理错误
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addblockcmd.Parsed(){
		if *addBlockData=="" {
			addblockcmd.Usage()
			os.Exit(1)
		}else{
			cli.addBlock(*addBlockData)//增加区块
		}
	}
	if showchaincmd.Parsed(){
		cli.showBlockChain()//显示区块链
	}



}

