package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Block struct {
	Data     string `json:"data" binding:"required"`
	PrevHash string
	Hash     string
	Index    int
}

type BlockChain struct {
	Blocks []*Block
}

var blockchain *BlockChain

func NewBlockchain() *BlockChain {
	return &BlockChain{Blocks: []*Block{GenesisBlock()}}
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, &Block{Data: "Genesis block"})
}

func main() {
	blockchain = NewBlockchain()
	go func() {
		for _, block := range blockchain.Blocks {
			fmt.Println("prev hash", block.PrevHash)
			fmt.Println("index", block.Index)
			fmt.Println("data", block.Data)
			fmt.Println("hash", block.Hash)
			fmt.Println()
		}
	}()
	engine := gin.New()
	engine.GET("/", GetBlockChain)
	engine.POST("/writeblock", WriteBlock)
	log.Println("listening on port:8080")
	engine.Run()
}

func GetBlockChain(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Blockchain": blockchain.Blocks,
	})
}

func WriteBlock(c *gin.Context) {
	var block *Block
	if err := c.BindJSON(&block); err != nil {
		log.Println(err, "unable to bind the details")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error() + "unable to bind the details",
		})
		return
	}
	blockchain.AddBlock(block)
	c.JSON(http.StatusOK, gin.H{
		"Succes": blockchain.Blocks[len(blockchain.Blocks)-1],
	})
}

func (bc *BlockChain) AddBlock(data *Block) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	block := CreateBlock(prevBlock, data)
	if validBlock(prevBlock, block) {
		bc.Blocks = append(bc.Blocks, block)
	}
}

func validBlock(prevBlock, block *Block) bool {
	if prevBlock.Hash != block.PrevHash {
		return false
	}
	if block.Index != prevBlock.Index+1 {
		return false
	}
	if !block.validHash(block.Hash) {
		return false
	}
	return true
}

func (b *Block) validHash(hash string) bool {
	b.GenerateHash()
	return b.Hash == hash
}

func CreateBlock(prevBlock, data *Block) *Block {
	data.Index = prevBlock.Index + 1
	data.PrevHash = prevBlock.Hash
	data.GenerateHash()
	return data
}

func (b *Block) GenerateHash() {
	bytes, _ := json.Marshal(b.Data)
	data := b.PrevHash + string(bytes)
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}
