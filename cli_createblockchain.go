package main

import (
	"fmt"
	"log"
	"os"
)

func (cli *CLI) createBlockchain(address, nodeID string) {
	if !ValidateAddress(address) {
		fmt.Println("ERROR: Address is not valid")
		os.Exit(1)
	}
	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	if !stringInSlice(address, wallets.GetAddresses()) {
		fmt.Println("ERROR: Wallet not found in Node")
		os.Exit(1)
	}

	bc := CreateBlockchain(address, nodeID)
	defer bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}
