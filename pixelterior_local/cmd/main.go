package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/pixelterior"
)

func main() {
	err := shim.Start(new(pixelterior.SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %S", err)
	}
}
