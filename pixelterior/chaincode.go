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
	Balance int            `json:"balance"`
	Address string         `json:"address"`
	Tokens  map[string]int `json:"tokens"`
}

type User struct {
	Name    string `json:"name"`
	Wallet  Wallet `json:"wallet"`
}

type Token struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	TotalSupply int `json:"totalSupply"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")

    token := Token{
        Name: "MyToken",
        Symbol: "MTK",
        Decimals: 18,
        TotalSupply: 1000000,
    }
    tokenJSON, _ := json.Marshal(token)

    stub.PutState("MTK", tokenJSON)

	user1 := User{
		Name: "user1",
		Wallet: Wallet{
			Address: "0x1234",
			Balance: 10000,
			Tokens: map[string]int{"MTK": 500000},
		},
	}
	user2 := User{
		Name: "user2",
		Wallet: Wallet{
			Address: "0x5678",
			Balance: 10000,
			Tokens: map[string]int{"MTK": 500000},
		},
	}


	user1JSON, _ := json.Marshal(user1)
	user2JSON, _ := json.Marshal(user2)

	stub.PutState("user1", user1JSON)
	stub.PutState("user2", user2JSON)

	fmt.Println("ex02 Initialized well")

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "transferToken":
		return t.transferToken(stub, args)
	case "createAsset":
		return t.createAsset(stub, args)
	case "createUser":
		return t.createUser(stub, args)
	case "deleteUser":
		return t.deleteUser(stub, args)
	case "deleteAllUsers":
		return t.deleteAllUsers(stub)
	case "trade":
		return t.trade(stub, args)
	case "updateAsset":
		return t.updateAsset(stub, args)
	case "deleteAsset":
		return t.deleteAsset(stub, args)
	case "deleteAllAssets":
		return t.deleteAllAssets(stub)
	case "queryAssets":
		return t.queryAssets(stub)
	case "queryAsset":
		return t.queryAsset(stub, args)
	case "queryUsers":
		return t.queryUsers(stub)
	case "queryUser":
		return t.queryUser(stub, args)
	default:
		return shim.Error("Invalid invoke function name")
	}
}

func (t *SimpleChaincode) transferToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments. Expecting 3: from, to, amount")
    }

    from := args[0]
    to := args[1]
    amount, err := strconv.Atoi(args[2])
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
    if fromUser.Tokens["MTK"] < amount {
        return shim.Error("Insufficient balance")
    }

    // 토큰 전송
    fromUser.Tokens["MTK"] -= amount
    if toUser.Tokens == nil {
        toUser.Tokens = make(map[string]int)
    }
    toUser.Tokens["MTK"] += amount

    // 상태 업데이트
    fromJSON, _ := json.Marshal(fromUser)
    toJSON, _ := json.Marshal(toUser)
    stub.PutState(from, fromJSON)
    stub.PutState(to, toJSON)

    return shim.Success([]byte("Token transfer successful"))
}

func (t *SimpleChaincode) createAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	name := args[0]
	assetKey := "asset_" + name

	assetBytes, err := stub.GetState(assetKey)
	if err != nil {
		return shim.Error("Failed to get asset")
	}

	if assetBytes != nil {
		return shim.Error("Inputed name of asset already exist")
	}

	price, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid price inputed. Expecting integer value")
	}

	totalAmount, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid totalAmount inputed. Expecting integer value")
	}

	asset := Asset{
		Name: name,
		Price: price,
		TotalStock: totalAmount,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return shim.Error("Failed to marshal asset")
	}

	err = stub.PutState(assetKey, assetJSON)
	if err != nil {
		return shim.Error("Failed to create asset")
	}

	return shim.Success([]byte("Asset created successfully"))
}

func (t *SimpleChaincode) updateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3: name, price, totalAmount")
	}

	name := args[0]
	assetKey := "asset_" + name
	price, _ := strconv.Atoi(args[1])
	totalAmount, _ := strconv.Atoi(args[2])

	assetBytes, err := stub.GetState(assetKey)
	if err != nil {
		return shim.Error("Failed to get asset")
	}

	if assetBytes == nil {
		return shim.Error("Asset not found")
	}

	asset := Asset{
		Name: name,
		Price: price,
		TotalStock: totalAmount,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return shim.Error("Failed to marshal asset")
	}

	err = stub.PutState(assetKey, assetJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to put state: %v", err))
	}

	return shim.Success([]byte("Asset updated successfully"))
}

func (t *SimpleChaincode) deleteAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	assetKey := "asset_" + name

	assetBytes, err := stub.GetState(assetKey)
	if err != nil {
		return shim.Error("Failed to get asset")
	}

	if assetBytes == nil {
		return shim.Error("Asset not found")
	}

	err = stub.DelState(assetKey)
	if err != nil {
		return shim.Error("Failed to delete asset")
	}

	return shim.Success([]byte("Asset deleted successfully"))
}

func (t *SimpleChaincode) deleteAllAssets(stub shim.ChaincodeStubInterface) pb.Response {
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

func (t *SimpleChaincode) queryAssets(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("asset_", "asset_~")

	if err != nil {
		return shim.Error("Failed to get assets")
	}
	defer iterator.Close()

	var assets []Asset
	for iterator.HasNext() {
		assetData, _ := iterator.Next()
		var asset Asset
		json.Unmarshal(assetData.Value, &asset)
		assets = append(assets, asset)
	}

	assetsBytes, _ := json.Marshal(assets)

	return shim.Success(assetsBytes)
}

func (t *SimpleChaincode) queryAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	assetKey := "asset_" + name

	assetBytes, err := stub.GetState(assetKey)
	if err != nil {
		return shim.Error("Failed to get asset")
	}

	if assetBytes == nil {
		return shim.Error("Asset not found")
	}

	return shim.Success(assetBytes)
}

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]

	existingUserBytes, _ := stub.GetState(name)

	if existingUserBytes != nil {
		return shim.Error("Username already exists")
	}

	newUser := User{
		Name: name, 
		Wallet: Wallet{Balance: 10000, Address: "0xasdf", Assets: []MyAsset{{Name: "asset1", Amount: 10}}}, 
		Tokens: make(map[string]int),
	}

	newUserBytes, _ := json.Marshal(newUser)
	err := stub.PutState(name, newUserBytes)

	if err != nil {
		return shim.Error("Failed to create user")
	}

	return shim.Success([]byte("User created successfully"))
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

func (t *SimpleChaincode) trade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	buyerId := args[0]
	sellerId := args[1]
	assetName := args[2]
	quantity, err := strconv.Atoi(args[3])

	if err != nil {
		return shim.Error("invalid quantity inputed. Expecting integer value")
	}

	buyerBytes, err := stub.GetState(buyerId)
	if err != nil || buyerBytes == nil {
		return shim.Error("Buyer Id not found")
	}

	sellerBytes, err := stub.GetState(sellerId)
	if err != nil || sellerBytes == nil {
		return shim.Error("Seller Id not found")
	}

	var buyer, seller User
	json.Unmarshal(buyerBytes, &buyer)
	json.Unmarshal(sellerBytes, &seller)

	assetBytes, err := stub.GetState("assets")
	var assets map[string]Asset
	json.Unmarshal(assetBytes, &assets)
	asset, exists := assets[assetName]
	if !exists {
		return shim.Error("Asset not found")
	}

	// Find seller's asset
	sellerAssetIndex := -1
	for i, a := range seller.Wallet.Assets {
		if a.Name == assetName {
			if a.Amount < quantity {
				return shim.Error("Seller does not have enough stock")
			}
			sellerAssetIndex = i
			break
		}
	}
	if sellerAssetIndex == -1 {
		return shim.Error("Seller does not have the asset")
	}

	totalPrice := asset.Price * quantity

	if buyer.Wallet.Balance < totalPrice {
		return shim.Error("Buyer does not have enough balance")
	}

	// Update balances
	buyer.Wallet.Balance -= totalPrice
	seller.Wallet.Balance += totalPrice

	// Update buyer's assets
	buyerAssetIndex := -1
	for i, a := range buyer.Wallet.Assets {
		if a.Name == assetName {
			buyer.Wallet.Assets[i].Amount += quantity
			buyerAssetIndex = i
			break
		}
	}
	if buyerAssetIndex == -1 {
		buyer.Wallet.Assets = append(buyer.Wallet.Assets, MyAsset{Name: assetName, Amount: quantity})
	}

	// Update seller's assets
	seller.Wallet.Assets[sellerAssetIndex].Amount -= quantity
	if seller.Wallet.Assets[sellerAssetIndex].Amount == 0 {
		seller.Wallet.Assets = append(seller.Wallet.Assets[:sellerAssetIndex], seller.Wallet.Assets[sellerAssetIndex+1:]...)
	}

	buyerBytes, _ = json.Marshal(buyer)
	sellerBytes, _ = json.Marshal(seller)
	stub.PutState(buyerId, buyerBytes)
	stub.PutState(sellerId, sellerBytes)

	return shim.Success([]byte("Trade completed"))
}

func (t *SimpleChaincode) queryUsers(stub shim.ChaincodeStubInterface) pb.Response {
	iterator, err := stub.GetStateByRange("", "")

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
	walletBytes, err := stub.GetState(userId)

	if err != nil || walletBytes == nil {
		return shim.Error("Wallet not found")
	}

	return shim.Success(walletBytes)
}

func (t *SimpleChaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name string

	name = args[0]

	userData, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get Data about inputed user: " + name)
	}

	if userData == nil {
		return shim.Error("Not Existing user name. Input other name")
	}

	return shim.Success(userData)
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

