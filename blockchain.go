package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

const dbFile="blockchain.db"  //数据库文件名当前目录下
const blockBucket="blocks"  //名称，
type  BlockChain struct {
	tip []byte  //二进制数据
	db *bolt.DB //数据库
}
type  BlockChainIterator struct {
	currentHash  []byte //当前的哈希
	db  *bolt.DB //数据库
}

//增加一个区块 挖矿验证区块并 序列化压入桶中 （由于区块作为整体记录）添加key=1;value=lasthash作为参数赋值 存入数据库中
func (block *BlockChain)AddBlock(data string){
	var  lastHash []byte //上一块哈希
	err:=block.db.View(func(tx *bolt.Tx) error {
		block :=tx.Bucket([]byte(blockBucket))//取得数据
		lastHash=block.Get([]byte("1"))//取得第一块 get key的值为1 所有的数据均写在数据库中  传递参数为last
		return nil
	})
	if err!=nil{
		log.Panic(err)//处理打开错误
	}
	newBlock:=NewBlock(data,lastHash)//创建一个新的区块
	err= block.db.Update(func (tx *bolt.Tx)error{  //写入数据函数
		bucket :=tx.Bucket([]byte(blockBucket))//按照表的名字取出表 数据库 表单调取
		err:=bucket.Put(newBlock.Hash,newBlock.Serialize())//压入数据 key值为最新区块的哈希值
		if err!=nil{
			log.Panic(err)//处理压入错误
		}
		err =bucket.Put([]byte("1"),newBlock.Hash)//压入数据
		if err!=nil{
			log.Panic(err)//处理压入错误
		}
		block.tip=newBlock.Hash //处理

		return nil
	})
}
//迭代器  遍历所有区块链
func (block *BlockChain)Iterator()*BlockChainIterator{
	bcit:=&BlockChainIterator{block.tip,block.db}
	return bcit //根据区块链创建区块链迭代器
}
//取得下一个区块
func (it *BlockChainIterator)next()*Block{
	var block *Block
	err:= it.db.View(func (tx *bolt.Tx)error{
		bucket:=tx.Bucket([]byte(blockBucket))
		encodedBlock:=bucket.Get(it.currentHash)//抓取二进制数据
		block=DeserializeBlock(encodedBlock)//解码
		return nil
	})
	if err!=nil{
		log.Panic(err)//处理压入错误
	}
	it.currentHash=block.PrevBlockHash//哈希赋值
	return block
}
//创世区块
func NewBlockChain()*BlockChain{
	var tip []byte  //存储区块链的二进制数据
	db,err:=bolt.Open(dbFile,0600,nil)//打开数据库
	if err!=nil{
		log.Panic(err)//处理数据库打开错误
	}
	//处理数据更新
	err=db.Update(func (tx *bolt.Tx)error {
		bucket:=tx.Bucket([]byte(blockBucket))//按照名称打开数据库的表格
		if  bucket==nil{
			fmt.Println("当前数据库没有区块链，没有创建一个新的")
			genesis:=NewGenesisBlock()//创建创世区块
			bucket,err:=tx.CreateBucket([]byte(blockBucket))//创建一个数据库的表格  桶是键值对的一个集合 桶中的键值必须唯一
			if  err!=nil{
				log.Panic(err) //处理创建错误
			}
			err =bucket.Put(genesis.Hash,genesis.Serialize())//存入数据
			if err!=nil{
				log.Panic(err)//处理存入错误
			}
			err =bucket.Put([]byte("1"),genesis.Hash)//存入数据
			if err!=nil{
				log.Panic(err)//处理存入错误
			}
			tip=genesis.Hash //取得哈希
		}else{
			tip=bucket.Get([]byte("1"))
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)//处理数据库更新错误
	}
	bc:=BlockChain{tip,db} //创建一个区块链
	return &bc
}




