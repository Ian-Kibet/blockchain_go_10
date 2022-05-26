# Instructions

Update your terminal profile `alias blockchain=go run .` then close and reopen the terminal window

## Create Blockchain
1. Create a Node ID
    `export NODE_ID=9000`
2. Create a wallet in the node. You can provide an alternative alias to the wallet address.
    `blockchain createwallet -alias test.go10`
3. Create the blockchain where ADDRESS is the mining reward address
    `blockchain createblockchain -address ADDRESS`
4. 

## Start HTTP Server
1. Run `blockchain http -port PORT` where PORT is the web server port to listen for connections

```
$ blockchain http -port 8080
Current Node ID: 9001

Starting http server on port: 8080
2022/05/26 18:19:12 Listening on 8080
```
2. You can now explore the APIs

## APIs

`GET /blockchain` - Retrieves the blockchain, all the blocks available
```
[
    {
        "Hash":  "0000c13e9e5f641ebe0dd19bb718e1f60851091abbe20c67f6114c2956109a51",
        "PrevBlockHash": "",
        "Timestamp": 1653556027,
        "Transactions": [
            {
                "ID": "ac3c6b31696a66ba3f5c59c84f033ad4757c4bd9097a6bab5dbfd52e71b6190e",
                "Vin": [
                    {
                        "Txid": null,
                        "Vout": -1,
                        "Signature": null,
                        "PubKey": "VGhlIFRpbWVzIDAzL0phbi8yMDA5IENoYW5jZWxsb3Igb24gYnJpbmsgb2Ygc2Vjb25kIGJhaWxvdXQgZm9yIGJhbmtz"
                    }
                ],
                "Vout": [
                    {
                    "Value": 10,
                    "PubKeyHash": "FbrkgXC5Yj8QanFB5pmEgrD9xqw="
                    }
                ]
            }
        ],
        "Nonce": 32259,
        "Height": 0,
        "Pow": true
    }
]
```

`GET /blockchain/{hash}` - Retrieves a block using its hash
```
{
    "Hash":  "0000c13e9e5f641ebe0dd19bb718e1f60851091abbe20c67f6114c2956109a51",
    "PrevBlockHash": "",
    "Timestamp": 1653556027,
    "Transactions": [
        {
            "ID": "ac3c6b31696a66ba3f5c59c84f033ad4757c4bd9097a6bab5dbfd52e71b6190e",
            "Vin": [
                {
                    "Txid": null,
                    "Vout": -1,
                    "Signature": null,
                    "PubKey": "VGhlIFRpbWVzIDAzL0phbi8yMDA5IENoYW5jZWxsb3Igb24gYnJpbmsgb2Ygc2Vjb25kIGJhaWxvdXQgZm9yIGJhbmtz"
                }
            ],
            "Vout": [
                {
                "Value": 10,
                "PubKeyHash": "FbrkgXC5Yj8QanFB5pmEgrD9xqw="
                }
            ]
        }
    ],
    "Nonce": 32259,
    "Height": 0,
    "Pow": true
}
```

`GET /blockchain/current` - Retrieves the latest block
```
{
    "Hash":  "0000c13e9e5f641ebe0dd19bb718e1f60851091abbe20c67f6114c2956109a51",
    "PrevBlockHash": "",
    "Timestamp": 1653556027,
    "Transactions": [
        {
            "ID": "ac3c6b31696a66ba3f5c59c84f033ad4757c4bd9097a6bab5dbfd52e71b6190e",
            "Vin": [
                {
                    "Txid": null,
                    "Vout": -1,
                    "Signature": null,
                    "PubKey": "VGhlIFRpbWVzIDAzL0phbi8yMDA5IENoYW5jZWxsb3Igb24gYnJpbmsgb2Ygc2Vjb25kIGJhaWxvdXQgZm9yIGJhbmtz"
                }
            ],
            "Vout": [
                {
                "Value": 10,
                "PubKeyHash": "FbrkgXC5Yj8QanFB5pmEgrD9xqw="
                }
            ]
        }
    ],
    "Nonce": 32259,
    "Height": 0,
    "Pow": true
}
```

`GET /wallets` - Retrieves the wallets
```
[
    "17r7kCwrKCbtVSXJw5ptxD9mALEBzRDbtE",
    "1MMYtKYd6Bj6WZ5JXPLYYMvD4d9ZsTkzo5"
]
```

`GET /wallets/17r7kCwrKCbtVSXJw5ptxD9mALEBzRDbtE` - Retrieves the wallet data
```
{
    "publicKey": "4HeZIvlW8bmOr9bFq5UIvOlmDQffdMPxh261dQt/YySBTuGwWSpVq+Zmg6OvqlGmbpOr/hKRrQcGsd51LK6Gdw==",
    "address": "17r7kCwrKCbtVSXJw5ptxD9mALEBzRDbtE",
    "alias": "ian.go10"
}
```

`GET /wallets/17r7kCwrKCbtVSXJw5ptxD9mALEBzRDbtE/balance` - Retrieves the wallet balance
```
{
    "balance": 0
}
```

`POST /wallets/17r7kCwrKCbtVSXJw5ptxD9mALEBzRDbtE/send` - Send token
```
Body
{
    "amount": 1,
    "receiver": "ian9001.go10"
}
Response
{
    "status": "pending",
    "id": "17r7kCwrKCbtVSXJw5ptxD9mALEBzRDbtE"
}
```

`GET /transactions` - Gets transactions
```
{
    
}
```