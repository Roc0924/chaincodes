package main_return_json


import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


var logger = shim.NewLogger("main")

type RebateDirectChainCode struct {

}

type Record struct {
	userId string
	state interface{}
}

func (chaincode *RebateDirectChainCode)Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.SetLevel(shim.LogInfo)
	logger.Infof("===========================================Init rebate_direct_cc===========================================")

	if transientMap, err := stub.GetTransient(); nil == err {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)

}

func (chaincode *RebateDirectChainCode)Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if "register" == args[0] {

		return chaincode.register(stub, args)

	}

	if "query" == args[0] {
		return chaincode.query(stub, args)
	}

	return shim.Success(nil)

}



func (chaincode *RebateDirectChainCode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if 3 != len(args) {
		logger.Errorf("Incorrect number of arguments, Expecting 3")
	}

	userId := args[1]
	recordStr := args[2]

	// check whether the user has registered
	isRegistered, checkError := chaincode.checkRegistrition(stub, userId)

	if nil != checkError {
		return shim.Error("check register error " + checkError.Error())
	}
	if isRegistered {
		return shim.Error("user " + userId + " has registered")
	}

	// create the account
	putError := stub.PutState(userId, []byte(recordStr))
	if nil != putError {
		return shim.Error("create account " + userId + "error" + putError.Error())
	}



	return shim.Success(nil)
}
func (chaincode *RebateDirectChainCode) query(stubInterface shim.ChaincodeStubInterface, args []string) pb.Response {



	return shim.Success(nil)
}
func (chainCode *RebateDirectChainCode) checkRegistrition(stub shim.ChaincodeStubInterface, userId string) (bool, error) {
	storeValue, storeError := stub.GetState(userId)
	if nil != storeError {
		shim.Error("get user error " + storeError.Error())
		return false, storeError
	}
	if nil != storeValue {
		return true, nil
	}

	return false, nil
}