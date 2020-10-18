package assignment02IBC

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const miningReward = 100
const rootUser = "Satoshi"

type Block struct {
	Spender     map[string]int
	Receiver    map[string]int
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

func CalculateBalance(userName string, chainHead *Block) int {
	balance := 0
	for c := chainHead; c != nil; c = c.PrevPointer {
		balance += c.Receiver[userName]
		balance -= c.Spender[userName]
	}
	//fmt.Println(balance)
	return balance
}

func CalculateHash(inputBlock *Block) string {
	hasher := sha256.New()
	var t_bytes = []byte{}

	for key, value := range inputBlock.Receiver {
		bKey := []byte(key)
		bVal := byte(value)
		for i := 0; i < len(bKey); i++ {
			t_bytes = append(t_bytes, bKey[i])
		}
		t_bytes = append(t_bytes, bVal)
	}

	for key, value := range inputBlock.Spender {
		bKey := []byte(key)
		bVal := byte(value)
		for i := 0; i < len(bKey); i++ {
			t_bytes = append(t_bytes, bKey[i])
		}
		t_bytes = append(t_bytes, bVal)
	}
	hasher.Write(t_bytes)
	currentHash := hex.EncodeToString(hasher.Sum(nil))
	return currentHash
}

func InsertBlock(spendingUser string, receivingUser string, miner string, amount int, chainHead *Block) *Block {
	if miner != rootUser {
		fmt.Print("Invalid insertion: miner error\n")
		return chainHead
	}
	if spendingUser != rootUser {
		check_bal := CalculateBalance(spendingUser, chainHead)
		if amount > check_bal {
			fmt.Print("Invalid insertion: bal error\n")
			return chainHead
		}
	}

	if chainHead == nil {
		b := &Block{nil, nil, chainHead, "", ""}
		b.Receiver = make(map[string]int)
		b.Spender = make(map[string]int)
		b.Receiver[rootUser] += miningReward
		b.CurrentHash = CalculateHash(b)
		//fmt.Println("First block mined")
		//fmt.Print(b.Receiver)
		return b
	}

	b := &Block{nil, nil, chainHead, "", ""}
	b.PrevHash = CalculateHash(chainHead)
	check_bal := CalculateBalance(spendingUser, chainHead)
	b.Receiver = make(map[string]int)
	b.Spender = make(map[string]int)
	if amount <= check_bal {
		b.Spender[spendingUser] = amount
		b.Receiver[receivingUser] = amount
		b.Receiver[rootUser] += miningReward
		b.CurrentHash = CalculateHash(b)
		return b
	}

	return chainHead
}

func ListBlocks(chainHead *Block) {
	for c := chainHead; c != nil; c = c.PrevPointer {
		//print("Hash Current: ", c.currentHash, " ")
		fmt.Print("Received: ")
		for key, value := range c.Receiver {
			fmt.Print(key, " ", value)
			fmt.Print(", ")
		}
		fmt.Print("Spent: ")
		for key, value := range c.Spender {
			fmt.Print(key, " ", value)
			fmt.Print(", ")
		}
		fmt.Print(" ---> ")
	}
	fmt.Println("Blockchain Start")
}

func VerifyChain(chainHead *Block) {
	for c := chainHead; c != nil; c = c.PrevPointer {
		hash_c := CalculateHash(c)
		if c.PrevPointer != nil {
			hash_p := CalculateHash(c.PrevPointer)
			if hash_p != c.PrevHash || hash_c != c.CurrentHash {
				fmt.Println("Blockchain is compromised")
				return
			}
		}
		if hash_c != c.CurrentHash {
			fmt.Println("Blockchain is compromised")
			return
		}
	}
	fmt.Println("Blockchain Verified")
	return
}
