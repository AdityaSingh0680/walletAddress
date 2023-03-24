package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <wallet address>")
		return
	}

	walletAddress := common.HexToAddress(os.Args[1])

	// Connect to Polygon Mumbai Testnet
	client, err := ethclient.Dial("https://rpc-mumbai.matic.today")
	if err != nil {
		fmt.Println("Failed to connect to the Ethereum client:", err)
		return
	}
	defer client.Close()

	// Load the contract instance
	contractAddress := common.HexToAddress("0x16581f93797e33fd2b1a3497822adf1762ee36e2")
	contractInstance, err := NewMyToken(contractAddress, client)
	if err != nil {
		fmt.Println("Failed to load contract instance:", err)
		return
	}

	// Get the balance of the wallet
	balance, err := contractInstance.BalanceOf(&bind.CallOpts{}, walletAddress)
	if err != nil {
		fmt.Println("Failed to get balance:", err)
		return
	}
	fmt.Printf("Total number of tokens for address %s: %s\n", walletAddress.Hex(), balance.String())

	// Get the token ID of the first token owned by the wallet
	tokenID, err := contractInstance.TokenOfOwnerByIndex(&bind.CallOpts{}, walletAddress, big.NewInt(0))
	if err != nil {
		fmt.Println("Failed to get token ID:", err)
		return
	}

	// Check if the token metadata is available
	tokenURI, err := contractInstance.TokenURI(&bind.CallOpts{}, tokenID)
	if err != nil {
		if err.Error() == "abi: cannot find matching function signature" {
			fmt.Println("No tokens in wallet.")
			return
		}
		fmt.Println("Failed to get token metadata:", err)
		return
	}

	fmt.Printf("Token metadata for token ID %s: %s\n", tokenID.String(), tokenURI)
}
