package pixelterior

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type Wallet struct {
	Address string         `json:"address"`
	Tokens  map[string]int `json:"tokens"`
}

type User struct {
	Name   string `json:"name"`
	Wallet Wallet `json:"wallet"`
}

type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply int    `json:"totalSupply"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")

	token := Token{
		Name:        "MyToken",
		Symbol:      "MTK",
		TotalSupply: 1000000,
	}
	tokenJSON, _ := json.Marshal(token)

	stub.PutState("MTK", tokenJSON)

	user1 := User{
		Name: "user1",
		Wallet: Wallet{
			Address: "0x1234",
			Tokens:  map[string]int{"MTK": 500000},
		},
	}
	user2 := User{
		Name: "user2",
		Wallet: Wallet{
			Address: "0x5678",
			Tokens:  map[string]int{"MTK": 500000},
		},
	}

	user1JSON, _ := json.Marshal(user1)
	user2JSON, _ := json.Marshal(user2)

	stub.PutState("user_user1", user1JSON)
	stub.PutState("user_user2", user2JSON)

	fmt.Println("ex02 Initialized well")

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "mint":
		return t.mint(stub, args)
	case "createUser":
		return t.createUser(stub, args)
	case "deleteUser":
		return t.deleteUser(stub, args)
	case "deleteAllUsers":
		return t.deleteAllUsers(stub)
	case "transfer":
		return t.transfer(stub, args)
	case "deleteAllTokens":
		return t.deleteAllTokens(stub)
	case "queryTokens":
		return t.queryTokens(stub)
	case "queryToken":
		return t.queryToken(stub, args)
	case "queryUsers":
		return t.queryUsers(stub)
	case "queryUser":
		return t.queryUser(stub, args)
	default:
		return shim.Error("Invalid function name invoked")
	}
}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3: from, to, amount")
	}

	from := args[0]
	to := args[1]
	tokenSymbol := args[2]
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid amount: " + err.Error())
	}

	// 송신자와 수신자의 상태 가져오기
	fromBytes, err := stub.GetState(from)
	if err != nil {
		return shim.Error("Failed to get sender: " + err.Error())
	}
	if fromBytes == nil {
		return shim.Error("Sender not found")
	}

	toBytes, err := stub.GetState(to)
	if err != nil {
		return shim.Error("Failed to get recipient: " + err.Error())
	}
	if toBytes == nil {
		return shim.Error("Recipient not found")
	}

	var fromUser, toUser User
	json.Unmarshal(fromBytes, &fromUser)
	json.Unmarshal(toBytes, &toUser)

	// 잔액 확인
	if fromUser.Wallet.Tokens[tokenSymbol] < amount {
		return shim.Error("Insufficient balance")
	}

	// 토큰 전송
	fromUser.Wallet.Tokens[tokenSymbol] -= amount
	toUser.Wallet.Tokens[tokenSymbol] += amount

	// 상태 업데이트
	fromJSON, _ := json.Marshal(fromUser)
	toJSON, _ := json.Marshal(toUser)
	stub.PutState(from, fromJSON)
	stub.PutState(to, toJSON)

	return shim.Success([]byte(nil))
}

func (t *SimpleChaincode) mint(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	name := args[0]
	symbol := args[1]
	totalAmount := args[2]
	to := args[3]

	tokenKey := "token_" + symbol

	symbolBytes, err := stub.GetState(tokenKey)
	if err != nil {
		return shim.Error("Failed to get tokens")
	}

	var token Token
	var totalSupply int

	totalSupply, err = strconv.Atoi(totalAmount)
	if err != nil {
		return shim.Error("Invalid totalAmount inputed. Expecting integer value")
	}

	// 이미 토큰이 존재할 경우 총 발행량을 더하고, 존재하지 않을 경우 토큰을 새로 발행함
	if symbolBytes != nil {
		json.Unmarshal(symbolBytes, &token)

		token.TotalSupply += totalSupply
	} else {
		token = Token{
			Name:        name,
			Symbol:      symbol,
			TotalSupply: totalSupply,
		}
	}

	toBytes, err := stub.GetState(to)
	if err != nil {
		return shim.Error("Failed to get recipient")
	}
	if toBytes == nil {
		return shim.Error("Recipient not found")
	}

	// to 에게 토큰을 발행
	var toUser User
	json.Unmarshal(toBytes, &toUser)

	toUser.Wallet.Tokens[symbol] += totalSupply

	toJson, err := json.Marshal(toUser)
	if err != nil {
		return shim.Error("Failed to marshal User")
	}

	err = stub.PutState(to, toJson)
	if err != nil {
		return shim.Error("Failed to mint token to User")
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return shim.Error("Failed to marshal Token")
	}

	err = stub.PutState(tokenKey, tokenJSON)
	if err != nil {
		return shim.Error("Failed to create Token")
	}

	return shim.Success([]byte(nil))
}

func (t *SimpleChaincode) deleteAllTokens(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("asset_", "asset_~")

	if err != nil {
		return shim.Error("Failed to get assets")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		assetData, _ := iterator.Next()
		assetKey := assetData.Key
		stub.DelState(assetKey)
	}

	return shim.Success([]byte("All assets deleted successfully"))
}

func (t *SimpleChaincode) deleteAllUsers(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("", "")

	if err != nil {
		return shim.Error("Failed to get users")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		userData, _ := iterator.Next()
		userKey := userData.Key
		stub.DelState(userKey)
	}

	return shim.Success([]byte("All users deleted successfully"))
}

func (t *SimpleChaincode) queryTokens(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("token_", "token_~")

	if err != nil {
		return shim.Error("Failed to get tokens")
	}
	defer iterator.Close()

	var tokens []Token
	for iterator.HasNext() {
		tokenData, _ := iterator.Next()
		var token Token
		json.Unmarshal(tokenData.Value, &token)
		tokens = append(tokens, token)
	}

	tokensBytes, _ := json.Marshal(tokens)
	return shim.Success(tokensBytes)
}

func (t *SimpleChaincode) queryToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	tokenKey := "token_" + name

	tokenBytes, err := stub.GetState(tokenKey)
	if err != nil {
		return shim.Error("Failed to get token")
	}

	if tokenBytes == nil {
		return shim.Error("Token not found")
	}

	return shim.Success(tokenBytes)
}

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userId := args[0]
	userKey := "user_" + userId

	existingUserBytes, _ := stub.GetState(userId)

	if existingUserBytes != nil {
		return shim.Error("Username already exists")
	}

	newUser := User{
		Name:   userId,
		Wallet: Wallet{Address: "0x", Tokens: make(map[string]int)},
	}

	newUserBytes, _ := json.Marshal(newUser)
	err := stub.PutState(userKey, newUserBytes)

	if err != nil {
		return shim.Error("Failed to create user")
	}

	return shim.Success([]byte(nil))
}

func (t *SimpleChaincode) deleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]

	existingUserBytes, _ := stub.GetState(name)

	if existingUserBytes == nil {
		return shim.Error("User not found")
	}

	err := stub.DelState(name)

	if err != nil {
		return shim.Error("Failed to delete user")
	}

	return shim.Success([]byte("User deleted successfully"))
}

func (t *SimpleChaincode) queryUsers(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("user_", "user_~")

	if err != nil {
		return shim.Error("Failed to get users")
	}
	defer iterator.Close()

	var users []User
	for iterator.HasNext() {
		userData, _ := iterator.Next()
		var user User
		json.Unmarshal(userData.Value, &user)
		users = append(users, user)
	}

	usersBytes, _ := json.Marshal(users)
	return shim.Success(usersBytes)
}

func (t *SimpleChaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userId := args[0]
	userKey := "user_" + userId

	userBytes, err := stub.GetState(userKey)

	if err != nil {
		return shim.Error("Failed to get user")
	}

	if userBytes == nil {
		return shim.Error("User not found")
	}

	return shim.Success(userBytes)
}

func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string
	var Aval, Bval int
	var X int
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}
