package pixelterior

import (
	"fmt"
	"testing"

	"github/com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChainCode struct {
}

func checkInit(t *testing.T, stub *shim.MockStub) {
	res := stub.MockInit()
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
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestPixelterior_Init(t *testing.T) {
	scc := new(SimpleChainCode)
	stub := shim.NewMockStub("ex02", scc)

	checkInit(t, stub)

	checkState(t, stub, "MTK", `{"symbol":"MTK","totalSupply":2000000,"remain":1000000}`)
}

func TestPixelterior_Invoke(t *testing.T) {
	scc := new(SimpleChainCode)
	stub := shim.NewMockStub("ex02", scc)

	// Init Token and Users
	checkInit(t, stub)

	// Transfer 100000 MTK from user1 to user2
	checkInvoke(t, stub, [][]byte{[]byte("transfer"), []byte("user_1"), []byte("user_2"), []byte("MTK"), []byte("100000")})

	// allocate 100000 MTK to user1
	checkInvoke(t, stub, [][]byte{[]byte("allocateToken"), []byte("user_1"), []byte("MTK"), []byte("100000")})
}
