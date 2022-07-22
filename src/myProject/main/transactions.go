package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10

type Transaction struct {
	TxID []byte
	Vin  []TxInput
	Vout []TxOutput
}

type TxOutput struct {
	Value        int64
	ScriptPubkey string
}

type TxInput struct {
	TxID      []byte
	Vout      int
	ScriptSig string
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encrypt := gob.NewEncoder(&encoded)
	err := encrypt.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.TxID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{subsidy, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetId()

	return &tx
}
