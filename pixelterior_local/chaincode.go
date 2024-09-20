package pixelterior_local

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type User struct {
	Name   string         `json:"name"`
	Tokens map[string]int `json:"tokens"`
}

type Token struct {
	Symbol      string `json:"symbol"`
	TotalSupply int    `json:"totalSupply"`
	Remain      int    `json:"remain"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")

	token := Token{
		Symbol:      "MTK",
		TotalSupply: 2000000,
		Remain:      1000000,
	}
	tokenJSON, _ := json.Marshal(token)

	stub.PutState("MTK", tokenJSON)

	user1 := User{
		Name:   "user1",
		Tokens: map[string]int{"MTK": 500000},
	}
	user2 := User{
		Name:   "user2",
		Tokens: map[string]int{"MTK": 500000},
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
	case "allocateToken":
		return t.allocateToken(stub, args)
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

func (t *SimpleChaincode) mint(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	symbol := args[0]
	totalAmount := args[1]

	tokenKey := "token_" + symbol

	tokenBytes, err := stub.GetState(tokenKey)
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
	if tokenBytes != nil {
		json.Unmarshal(tokenBytes, &token)

		token.TotalSupply += totalSupply
	} else {
		token = Token{
			Symbol:      symbol,
			TotalSupply: totalSupply,
			Remain:      totalSupply,
		}
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

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userId := args[0]
	userKey := "user_" + userId

	existingUserBytes, _ := stub.GetState(userKey)

	if existingUserBytes != nil {
		return shim.Error("Username already exists")
	}

	newUser := User{
		Name:   userId,
		Tokens: make(map[string]int),
	}

	newUserBytes, _ := json.Marshal(newUser)
	err := stub.PutState(userKey, newUserBytes)

	if err != nil {
		return shim.Error("Failed to create user")
	}

	return shim.Success([]byte(nil))
}

func (t *SimpleChaincode) allocateToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3: userId, symbol, amount")
	}

	userId := args[0]
	symbol := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid amount: " + err.Error())
	}

	userKey := "user_" + userId
	tokenKey := "token_" + symbol

	userBytes, err := stub.GetState(userKey)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	}

	if userBytes == nil {
		return shim.Error("User not found")
	}

	tokenBytes, err := stub.GetState(tokenKey)
	if err != nil {
		return shim.Error("Failed to get token: " + err.Error())
	}

	if tokenBytes == nil {
		return shim.Error("Token not found")
	}

	var user User
	var token Token
	json.Unmarshal(userBytes, &user)
	json.Unmarshal(tokenBytes, &token)

	if token.Remain < amount {
		return shim.Error("Insufficient token")
	}

	user.Tokens[symbol] += amount
	token.Remain -= amount

	userJSON, _ := json.Marshal(user)
	tokenJSON, _ := json.Marshal(token)

	stub.PutState(userKey, userJSON)
	stub.PutState(tokenKey, tokenJSON)

	return shim.Success([]byte(nil))
}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4: from, to, symbol, amount")
	}

	from := args[0]
	to := args[1]
	fromKey := "user_" + from
	toKey := "user_" + to

	tokenSymbol := args[2]
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid amount: " + err.Error())
	}

	// 송신자와 수신자의 상태 가져오기
	fromBytes, err := stub.GetState(fromKey)
	if err != nil {
		return shim.Error("Failed to get sender: " + err.Error())
	}
	if fromBytes == nil {
		return shim.Error("Sender not found")
	}

	toBytes, err := stub.GetState(toKey)
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
	if fromUser.Tokens[tokenSymbol] < amount {
		return shim.Error("Insufficient balance")
	}

	// 토큰 전송
	fromUser.Tokens[tokenSymbol] -= amount
	toUser.Tokens[tokenSymbol] += amount

	// 상태 업데이트
	fromJSON, _ := json.Marshal(fromUser)
	toJSON, _ := json.Marshal(toUser)
	stub.PutState(from, fromJSON)
	stub.PutState(to, toJSON)

	return shim.Success([]byte(nil))
}

func (t *SimpleChaincode) deleteAllTokens(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("token_", "token_~")

	if err != nil {
		return shim.Error("Failed to get assets")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		assetData, _ := iterator.Next()
		assetKey := assetData.Key
		stub.DelState(assetKey)
	}

	return shim.Success([]byte("All token deleted"))
}

func (t *SimpleChaincode) deleteAllUsers(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("", "")

	if err != nil {
		return shim.Error("Failed to get states")
	}
	defer iterator.Close()

	for iterator.HasNext() {
		state, _ := iterator.Next()
		stateKey := state.Key
		stub.DelState(stateKey)
	}

	return shim.Success([]byte("All state deleted"))
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

func (t *SimpleChaincode) deleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userId := args[0]
	userKey := "user_" + userId

	existingUserBytes, _ := stub.GetState(userKey)

	if existingUserBytes == nil {
		return shim.Error("User not found")
	}

	err := stub.DelState(userKey)

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
