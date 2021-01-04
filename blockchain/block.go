package blockchain

import (
	"bytes"
	"encoding/gob"
	. "yu/common"
	"yu/event"
	"yu/trie"
	"yu/txn"
)

type Block struct {
	header *Header
	txns   []*txn.Txn
}

func NewBlock(preHeader *Header, txns []*txn.Txn, stateRoot Hash) (*Block, error) {
	preHash := preHeader.txnRoot
	blocknum := preHeader.number + 1

	txnsBytes := make([]Hash, 0)
	for _, tx := range txns {
		hash, err := tx.Hash()
		if err != nil {
			return nil, err
		}
		txnsBytes = append(txnsBytes, hash)
	}
	mTree := trie.NewMerkleTree(txnsBytes)
	txnRoot := mTree.RootNode.Data

	header := NewHeader(preHash, blocknum, txnRoot, stateRoot)
	return &Block{
		header,
		txns,
	},nil
}

func (b *Block) Head() *Header {
	return b.header
}

func (b *Block) Txns() []*txn.Txn {
	return b.txns
}

func (b *Block) Hash() Hash {
	return b.header.TxnRoot()
}

func (b *Block) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (b *Block) Events() []event.Event {
	allEvents := make([]event.Event, 0)
	for _, tx := range b.txns {
		events := tx.Events()
		allEvents = append(allEvents, events...)
	}
	return allEvents
}
