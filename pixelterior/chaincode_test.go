package pixelterior

import (
	"fmt"
	"testing"

	"github/com/hyperledger/fabric/core/chaincode/shim"
)

type Asset struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	TotalStock int    `json:"totalStock"`
}

func checkInit(t *testing.T, stub *shim.MockStub) {
	res = stub.MockInit()
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State(name)
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shin.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestPixelterior_Init(t *testing.T) {
	scc := new(SimpleChainCode)
	stub := shin.NewMockStub("ex02", scc)

	checkInit(t, stub)

	checkState(t, stub, "Asset1", []byte(Asset{Name: "asset1", Price: 100, TotalStock: 10000}))
}

func TestPixelterior_Invoke(t *testing.T) {
	scc := new(SimpleChainCode)
	stub := shim.NewMockStub("ex02", scc)

	// Init Asset1{Name: "asset1", Price: 100, TotalStock: 10000}, Asset2{Name: "asset2", Price: 200, TotalStock: 10000}
	checkInit(t, stub)

	// Invoke createAsset asset3 100 10000
	checkInvoke(t, stub, [][]byte{[]byte("createAsset"), []byte("Asset3"), []byte("100"), []byte("10000")})

	// Invoke updateAsset asset1 price 100 -> 200
	checkInvoke(t, stub, [][]byte{[]byte("updateAsset"), []byte("Asset1"), []byte("200"), []byte("10000")})

	// Invoke queryAsset asset1
	checkInvoke(t, stub, [][]byte{[]byte("queryAsset"), []byte("Asset1")})

	// Invoke createUser User{Name: "TestUser", Wallet: Wallet{Balance: 10000, Address: "0xasdf", Assets: []MyAsset{Name: "asset1", Amount: 10}}}, Tokens: make(map[string]int)}
	checkInvoke(t, stub, [][]byte{[]byte("createUser"), []byte("TestUser")})

	// Invoke trade User1 -> User2 / Asset1 50ea
	checkInvoke(t, stub, [][]byte{[]byte("trade"), []byte("User1"), []byte("User2"), []byte("Asset1"), []byte("50")})
}
