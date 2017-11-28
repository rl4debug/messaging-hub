package blockchain

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
)

type Block struct {
	Hash         string
	PreviousHash string
	Data         []byte
}

func (b *Block) Seriallize() []byte {
	var hashBytes = []byte(b.Hash)
	var previousHashBytes = []byte(b.PreviousHash)

	var lenHash = uint32(len(hashBytes))
	var lenPrevioushHash = uint32(len(previousHashBytes))
	var lenData = uint32(len(b.Data))

	var data = append(intToBytes(lenHash), hashBytes...)

	data = append(data, intToBytes(lenPrevioushHash)...)
	data = append(data, previousHashBytes...)

	data = append(data, intToBytes(lenData)...)
	data = append(data, b.Data...)
	return data
}

func intToBytes(a uint32) []byte {
	var data = make([]byte, 4)
	binary.LittleEndian.PutUint32(data, a)
	return data
}

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) GenerateFirstBlock() *Block {
	var md5 = md5.New()
	io.WriteString(md5, "random string")

	var firstHash = fmt.Sprintf("%x", md5.Sum(nil))
	var block = &Block{
		PreviousHash: firstHash,
		Hash:         firstHash,
		Data:         nil,
	}

	return block
}

func (bc *BlockChain) CreateNewBlock(data []byte) *Block {
	var l = len(bc.Blocks)
	if l < 1 {
		panic("Invalid operation")
	}

	var lastBlock = bc.Blocks[l-1]

	var md5 = md5.New()
	var hash = fmt.Sprintf("%x", md5.Sum(append([]byte(lastBlock.Hash), data...)))

	var newBlock = &Block{
		PreviousHash: lastBlock.Hash,
		Hash:         hash,
		Data:         data,
	}

	return newBlock
}
