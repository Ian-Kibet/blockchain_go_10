package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func MakeRouter() http.Handler {
	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(LogRequest)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Welcome to Go10 Blockchain","status": "OK"}`)
	}).Methods("GET")

	router.HandleFunc("/wallets", getWallets).Methods("GET")
	router.HandleFunc("/wallets/{address}", getWallet).Methods("GET")
	router.HandleFunc("/wallets/{address}/balance", getWalletBalance).Methods("GET")

	router.HandleFunc("/wallets/{address}/send", sendTokens).Methods("POST")

	router.HandleFunc("/transactions", getTransactions).Methods("GET")
	router.HandleFunc("/transactions/{id}", getTransaction).Methods("GET")

	router.HandleFunc("/blockchain", getBlockchain).Methods("GET")

	router.HandleFunc("/blockchain/current", getCurrentBlock).Methods("GET")
	router.HandleFunc("/blockchain/{hash}", getBlock).Methods("GET")

	router.HandleFunc("/nodes", getNodes).Methods("GET")

	return router
}

func LogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v] %v %v", r.RemoteAddr, r.Method, r.URL)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func ConvertAliases(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: translate aliases and save in session
		// address := mux.Vars(r)["address"]
		// if strings.Contains(address, ".go10") {
		// 	nodeID := os.Getenv("NODE_ID")
		// 	wallets, err := NewWallets(nodeID)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// 	wa, err := wallets.GetWalletByAlias(address)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusNotFound)
		// 		return
		// 	}
		// }
		h.ServeHTTP(w, r)
	})
}
