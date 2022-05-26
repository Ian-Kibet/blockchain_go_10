package main

import (
	"fmt"
	"os"
	"strings"
)

func (cli *CLI) createWallet(nodeID string, alias string) {
	if alias != "" && !strings.HasSuffix(alias, ".go10") {
		fmt.Println("ERROR: Invalid alias provided. Aliases should end in .go10. For Example; ian.go10")
		os.Exit(1)
	}
	wallets, _ := NewWallets(nodeID)
	address := wallets.CreateWallet(alias)
	wallets.SaveToFile(nodeID)

	fmt.Printf("Your new address: %s alias %v\n ", address, alias)
}
