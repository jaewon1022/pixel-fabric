package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/pixelterior_local"
)

func main() {
	err := shim.Start(new(pixelterior_local.SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %S", err)
	}
}
