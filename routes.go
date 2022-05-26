package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// TODO: Consolidate types
// TODO: Remove repeated lines

type BlockResponse struct {
	Hash, PrevBlockHash string
	Timestamp           int64
	Transactions        []*struct {
		ID   string
		Vin  []TXInput
		Vout []TXOutput
	}
	Nonce  int
	Height int
	Pow    bool
}

type WalletData struct {
	publicKey, address, alias string
	balance                   int
}

type TransactionResponse struct {
	Block BlockResponse
	ID    string
	Vin   []TXInput
	Vout  []TXOutput
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	nodeID := os.Getenv("NODE_ID")
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	bci := bc.Iterator()

	var blockchain []BlockResponse

	for {
		block := bci.Next()

		pow := NewProofOfWork(block)

		var transactions []*struct {
			ID   string
			Vin  []TXInput
			Vout []TXOutput
		}
		for _, tx := range block.Transactions {
			transactions = append(transactions, &struct {
				ID   string
				Vin  []TXInput
				Vout []TXOutput
			}{
				ID:   hex.EncodeToString(tx.ID),
				Vin:  tx.Vin,
				Vout: tx.Vout,
			})
		}

		blockchain = append(blockchain, BlockResponse{
			Hash:          hex.EncodeToString(block.Hash),
			PrevBlockHash: hex.EncodeToString(block.PrevBlockHash),
			Timestamp:     block.Timestamp,
			Transactions:  transactions,
			Nonce:         block.Nonce,
			Height:        block.Height,
			Pow:           pow.Validate(),
		})
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	bytes, err := json.Marshal(blockchain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func getWallets(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.Header().Set("Content-Type", "application/json")

	nodeID := os.Getenv("NODE_ID")

	wallets, err := NewWallets(nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}

	fmt.Printf("%T", addresses)

	// var walletData []WalletData

	// for _, w := range addresses {
	// 	wallet := wallets.GetWallet(w)
	// 	walletData = append(walletData, WalletData{
	// 		publicKey: hex.EncodeToString(wallet.PublicKey),
	// 		address:   string(wallet.GetAddress()),
	// 		balance:   wallet.GetBalance(nodeID),
	// 		alias:     wallet.Alias,
	// 	})
	// }

	bytes, err := json.Marshal(addresses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func getWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	address := params["address"]
	usingAlias := strings.Contains(address, ".go10")

	nodeID := os.Getenv("NODE_ID")

	if !ValidateAddress(address) {
		http.Error(w, `{"message":"Wallet not valid"}`, http.StatusBadRequest)
		return
	}

	wallets, err := NewWallets(nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var searchWallet bool

	if usingAlias {
		searchWallet = stringInSlice(address, wallets.GetAliases())
	} else {
		searchWallet = stringInSlice(address, wallets.GetAddresses())
	}

	if !searchWallet {
		http.Error(w, `{"message":"Wallet not found"}`, http.StatusNotFound)
		return
	}

	var wallet Wallet

	if usingAlias {
		wa, err := wallets.GetWalletByAlias(address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		wallet = *wa
	} else {
		wallet = wallets.GetWallet(address)
	}

	bytes, err := json.Marshal(&struct {
		PublicKey []byte `json:"publicKey"`
		Address   string `json:"address"`
		Alias     string `json:"alias"`
	}{
		Address:   string(wallet.GetAddress()),
		PublicKey: wallet.PublicKey,
		Alias:     wallet.Alias,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func getWalletBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	address := params["address"]
	usingAlias := strings.Contains(address, ".go10")
	nodeID := os.Getenv("NODE_ID")

	if !ValidateAddress(address) {
		http.Error(w, `{"message":"Wallet not valid"}`, http.StatusBadRequest)
		return
	}

	wallets, err := NewWallets(nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var searchWallet bool

	if usingAlias {
		searchWallet = stringInSlice(address, wallets.GetAliases())
	} else {
		searchWallet = stringInSlice(address, wallets.GetAddresses())
	}

	if !searchWallet {
		http.Error(w, `{"message":"Wallet not found"}`, http.StatusNotFound)
		return
	}

	var wallet Wallet

	if usingAlias {
		wa, err := wallets.GetWalletByAlias(address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		wallet = *wa
	} else {
		wallet = wallets.GetWallet(address)
	}

	balance := wallet.GetBalance(nodeID)

	bytes, err := json.Marshal(&struct {
		Balance int `json:"balance"`
	}{
		Balance: balance,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	nodeID := os.Getenv("NODE_ID")
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	bci := bc.Iterator()

	var transactions []TransactionResponse

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			transactions = append(transactions, TransactionResponse{
				Block: BlockResponse{
					Hash:          hex.EncodeToString(block.Hash),
					PrevBlockHash: hex.EncodeToString(block.PrevBlockHash),
					Height:        block.Height,
					Nonce:         block.Nonce,
				},
				ID:   hex.EncodeToString(tx.ID),
				Vin:  tx.Vin,
				Vout: tx.Vout,
			})
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

	bytes, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	nodeID := os.Getenv("NODE_ID")

	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	idBytes, err := hex.DecodeString(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tx, err := bc.FindTransaction(idBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(&struct {
		ID   string
		Vin  []TXInput
		Vout []TXOutput
	}{
		ID:   hex.EncodeToString(tx.ID),
		Vin:  tx.Vin,
		Vout: tx.Vout,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	hash := params["hash"]
	nodeID := os.Getenv("NODE_ID")

	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	block, err := bc.GetBlock(hashBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	pow := NewProofOfWork(&block)

	var transactions []*struct {
		ID   string
		Vin  []TXInput
		Vout []TXOutput
	}

	for _, tx := range block.Transactions {
		transactions = append(transactions, &struct {
			ID   string
			Vin  []TXInput
			Vout []TXOutput
		}{
			ID:   hex.EncodeToString(tx.ID),
			Vin:  tx.Vin,
			Vout: tx.Vout,
		})
	}

	bytes, err := json.Marshal(BlockResponse{
		Hash:          hex.EncodeToString(block.Hash),
		PrevBlockHash: hex.EncodeToString(block.PrevBlockHash),
		Timestamp:     block.Timestamp,
		Transactions:  transactions,
		Nonce:         block.Nonce,
		Height:        block.Height,
		Pow:           pow.Validate(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func getCurrentBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nodeID := os.Getenv("NODE_ID")

	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	lastBlock := bc.GetLastBlock()

	pow := NewProofOfWork(&lastBlock)

	var transactions []*struct {
		ID   string
		Vin  []TXInput
		Vout []TXOutput
	}

	for _, tx := range lastBlock.Transactions {
		transactions = append(transactions, &struct {
			ID   string
			Vin  []TXInput
			Vout []TXOutput
		}{
			ID:   hex.EncodeToString(tx.ID),
			Vin:  tx.Vin,
			Vout: tx.Vout,
		})
	}

	bytes, err := json.Marshal(BlockResponse{
		Hash:          hex.EncodeToString(lastBlock.Hash),
		PrevBlockHash: hex.EncodeToString(lastBlock.PrevBlockHash),
		Timestamp:     lastBlock.Timestamp,
		Nonce:         lastBlock.Nonce,
		Height:        lastBlock.Height,
		Pow:           pow.Validate(),
		Transactions:  transactions,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func sendTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nodeID := os.Getenv("NODE_ID")

	type SendTokenRequest struct {
		Amount   int    `json:"amount"`
		Receiver string `json:"receiver"`
	}

	var body SendTokenRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	sender := params["address"]
	receiver := body.Receiver
	amount := body.Amount

	usingAlias := strings.Contains(sender, ".go10")

	if !ValidateAddress(sender) {
		http.Error(w, `{"message":"Sender Wallet not valid"}`, http.StatusBadRequest)
		return
	}
	if !ValidateAddress(receiver) {
		http.Error(w, `{"message":"Receiver Wallet not valid"}`, http.StatusBadRequest)
		return
	}

	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	wallets, err := NewWallets(nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var wallet Wallet

	if usingAlias {
		wa, err := wallets.GetWalletByAlias(sender)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		wallet = *wa
	} else {
		wallet = wallets.GetWallet(sender)
	}

	tx, err := NewUTXOTransaction(&wallet, receiver, amount, &UTXOSet, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sendTx(knownNodes[0], tx)

	bytes, err := json.Marshal(&struct {
		Status string `json:"status"`
		ID     string `json:"id"`
	}{
		Status: "pending",
		ID:     hex.EncodeToString(tx.ID),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}

func getNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bytes, err := json.Marshal(knownNodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}
