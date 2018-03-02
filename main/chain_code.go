package main

import (
"fmt"
"strconv"
"encoding/json"

"github.com/hyperledger/fabric/core/chaincode/shim"
pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init initializes the chaincode state


type Response struct{
	message string
	status string
	data interface{}
}

type Record struct {
	userId string
	asset int64
	shadowAsset int64
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### chain_code Init ###########")
	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)

}

// Invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### chain_code Invoke ###########")
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0] == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	if args[0] == "anticipateRebate" {
		// Deletes an entity from its state
		return t.anticipateRebate(stub, args)
	}
	if args[0] == "rebateAccounted" {
		// Deletes an entity from its state
		return t.rebateAccounted(stub, args)
	}
	if args[0] == "rollbackRebate" {
		// Deletes an entity from its state
		return t.rollbackRebate(stub, args)
	}
	if args[0] == "register" {
		// Deletes an entity from its state
		return t.register(stub, args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', 'anticipateRebate', 'rebateAccounted', 'rollbackRebate',or 'register'")
}



func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	userId := args[0]
	value := args[1]
	shadowUserId := userId + "_s"

	////////////////////////////////////////
	// check whether ths user has been registered

	isRegistered, checkError := t.checkUserIsRegistered(stub, userId)
	if nil != checkError {
		return shim.Error("check register error " + checkError.Error())
	}
	if isRegistered {
		return shim.Error("user " + userId + " has registered")
	}

	////////////////////////////////////////
	// create the account

	putError := stub.PutState(userId, []byte(value))
	if nil != putError {
		return shim.Error("create account " + userId + "error" + putError.Error())
	}

	////////////////////////////////////////
	// create the shadow account
	shadowPutError := stub.PutState(shadowUserId, []byte("0"))
	if nil != shadowPutError {
		return shim.Error("create account " + shadowUserId + "error" + shadowPutError.Error())
	}

	valueInt64, parseValueError := strconv.ParseInt(value, 10, 64)
	if nil != parseValueError {
		return shim.Error("string convert error" + parseValueError.Error())
	}


	record := Record{userId:userId, asset:valueInt64, shadowAsset:0}

	recordBytes, jsonError := json.Marshal(record)
	if nil != jsonError {
		return shim.Error("parse json error " + jsonError.Error())
	}

	return shim.Success(recordBytes)
}


// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var userId string
	var err error

	response := &Response{}

	if len(args) != 2 {
		response.message = "Incorrect number of arguments"
		response.status = "failed"
		response.data = nil
		return shim.Error(string(response2bytes(response)))
	}


	userId = args[1]

	recordByte, err := t.getRecordByUserId(stub, userId)

	if nil != err {
		return shim.Error("query error " + err.Error())
	}
	if nil != recordByte {
		response.data = bytes2response(recordByte)
		response.status = "success"
		response.message = "query success"
	} else {
		response.message = "record is empty"
		response.status = "failed"
		response.data = nil
	}
	return shim.Success(response2bytes(response))

}


/*
	anticipate rebate
 */
func (t *SimpleChaincode) anticipateRebate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	///////////////////////////////////////
	// check arguments
	if 4 != len(args) {
		return shim.Error("Incorrect number of arguments, Expecting 4")
	}
	var source string = args[1]
	var shadowDestination string = args[2] + "_s"
	delta, parseDeltaError := strconv.ParseInt(args[3], 10, 64)
	if nil != parseDeltaError {
		return shim.Error("error:[" + parseDeltaError.Error() + "] arise when parse " + args[3] )
	}

	return t.move(stub, source, shadowDestination, delta)
}



func (chaincode *SimpleChaincode) rebateAccounted(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	///////////////////////////////////////
	// check arguments
	if 3 != len(args) {
		return shim.Error("Incorrect number of arguments, Expecting 4")
	}
	var userId string = args[1]
	var shadowUserId string = userId + "_s"
	delta, parseDeltaError := strconv.ParseInt(args[2], 10, 64)
	if nil != parseDeltaError {
		return shim.Error("error:[" + parseDeltaError.Error() + "] arise when parse " + args[3] )
	}
	return chaincode.move(stub, userId, shadowUserId, delta)
}





/*
	rollback rebate
 */
func (chaincode *SimpleChaincode) rollbackRebate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	///////////////////////////////////////
	// check arguments
	if 4 != len(args) {
		return shim.Error("Incorrect number of arguments, Expecting 4")
	}

	var shadowDestination = args[1] + "s"
	var source = args[2]
	delta, parseDeltaError := strconv.ParseInt(args[3], 10, 64)
	if nil != parseDeltaError {
		return shim.Error("error:[" + parseDeltaError.Error() + "] arise when parse " + args[3] )
	}
	return chaincode.move(stub, shadowDestination, source, delta)
}




func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, source string, destination string, delta int64) pb.Response {
	///////////////////////////////////////
	// get state from ledger
	sourceStateBytes, sourceStateBytesError := stub.GetState(source)
	if nil != sourceStateBytesError {
		return shim.Error("error " + sourceStateBytesError.Error() + " arise when getting " + source + " state from ledger")
	}

	destinationStateBytes, destinationStateBytesError := stub.GetState(destination)
	if nil != destinationStateBytesError {
		return shim.Error("error " + destinationStateBytesError.Error() + " arise when getting " + destination + " state from ledger")
	}

	sourceState, sourceStateError := strconv.ParseInt(string(sourceStateBytes), 10, 64)
	if nil != sourceStateError {
		return shim.Error("error:[" + sourceStateError.Error() + "] arise when parse " + string(sourceStateBytes) )
	}
	destinationState, destinationStateError := strconv.ParseInt(string(destinationStateBytes), 10, 64)
	if nil != destinationStateError {
		return shim.Error("error:[" + destinationStateError.Error() + "] arise when parse " + string(destinationStateBytes) )
	}


	///////////////////////////////////////
	// move rebate
	sourceState -= delta
	destinationState += delta



	///////////////////////////////////////
	// put new state to ledger

	sourcePutErr :=stub.PutState(source, []byte(string(sourceState)))
	if nil != sourcePutErr {
		return shim.Error("put " + source + " error: " + sourcePutErr.Error())
	}

	destinationPutErr :=stub.PutState(destination, []byte(string(destinationState)))
	if nil != destinationPutErr {
		return shim.Error("put " + destination + " error" + destinationPutErr.Error())
	}

	response := &Response{}
	response.message = "move rebate from " + destination + " to " + source + " success, amount: " + string(delta)
	response.status = "success"
	response.data = nil

	return shim.Success(response2bytes(response))
}

/*
get record by user id
 */
func (t *SimpleChaincode) getRecordByUserId(stub shim.ChaincodeStubInterface, userId string) ([]byte, error) {
	////////////////////////////////////////
	// generate shadow user id
	shadowUserId := userId + "s"


	storeValue, storeError := stub.GetState(userId)
	var record Record
	if nil != storeError {
		shim.Error("get user error " + storeError.Error())
		return nil, storeError
	}
	if nil != storeValue {
		shim.Success([]byte("user " + userId + "has registered"))

		// get asset
		storeValueInt64, parseValueError := strconv.ParseInt(string(storeValue), 10, 64)
		if nil != parseValueError {
			shim.Error("string convert value error" + parseValueError.Error())
			return nil, parseValueError
		}

		// get shadow asset
		storeShadowValue, storeShadowError := stub.GetState(shadowUserId)

		if nil != storeShadowError {
			shim.Error("get shadow user error " + storeShadowError.Error())
			return nil, storeShadowError
		}

		if nil != storeShadowValue {
			storeShadowValueInt64, parseShadowValueError := strconv.ParseInt(string(storeShadowValue), 10, 64)
			if nil != parseShadowValueError {
				shim.Error("string convert shadow value error" + parseShadowValueError.Error())
				return nil, parseShadowValueError
			}

			record = Record{userId:userId, asset:storeValueInt64, shadowAsset:storeShadowValueInt64}

			recordByte, recordError := json.Marshal(record)
			if nil != recordError {
				return nil, recordError
			}

			return recordByte, nil
		}
	}
	return nil, nil
}




// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userId := args[1]
	shadowUserId := userId + "_s"

	// Delete the key from the state in ledger
	err := stub.DelState(userId)
	if err != nil {
		return shim.Error("Failed to delete state")
	}
	shadowErr := stub.DelState(shadowUserId)
	if shadowErr != nil {
		return shim.Error("Failed to delete shadow state")
	}

	return shim.Success(nil)
}


func (t *SimpleChaincode) checkUserIsRegistered(stub shim.ChaincodeStubInterface, userId string) (bool, error) {

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



func response2bytes(response *Response) []byte {

	responseBytes, responseError := json.Marshal(response)
	if nil != responseError {
		fmt.Printf("error %s arise when marshal response", responseError.Error())
		return nil
	}
	return responseBytes

}


func bytes2response(data []byte) Response {
	var response Response
	responseError := json.Unmarshal(data, &response)
	if nil != responseError {
		fmt.Printf("error %s arise when marshal response", responseError.Error())
	}
	return response

}






func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
