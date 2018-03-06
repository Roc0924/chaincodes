/**
 * Create with GoLand
 * Author               : wangzhenpeng
 * Date                 : 2018/3/6
 * Time                 : 下午2:33
 * Description          : directly rebate chaincode
 */
package rebate_directly_cc

import (
	"fmt"
	"strconv"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
//define chaincode struct
type RebateChainCode struct {
}


//init chaincode
func (chainCode *RebateChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Init chaincode rebate_directly_cc\n")


	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}

	}

	return shim.Success(nil)
}


//invoke method
func (chainCode *RebateChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("query the method invoke of chaincode rebate_directly_cc\n")
	function, args := stub.GetFunctionAndParameters()
	if "invoke" != function {
		return shim.Error("unknown function call: " + function)
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if "register" == args[0] {
		if len(args) != 3 {
			return shim.Error("Call method register error. Incorrect number of arguments. Expecting 3")
		}

		return chainCode.register(stub, args[1], args[2])
	}

	if "query" == args[0] {
		if len(args) != 2 {
			return shim.Error("Call method query error. Incorrect number of arguments. Expecting 2")
		}
		return chainCode.query(stub, args[1])
	}

	if "rebateDirectly" == args[0] {
		if len(args) != 2 {
			return shim.Error("Call method rebateDirectly error. Incorrect number of arguments. Expecting 4")
		}
		return chainCode.rebateDirectly(stub, args[1], args[2], args[3])
	}


	return shim.Success(nil)
}


//register a record
func (chainCode *RebateChainCode) register(stub shim.ChaincodeStubInterface, userId string, value string) pb.Response {
	fmt.Println("================== register ====================")


	err := stub.PutState(userId, []byte(value))
	if nil != err {
		fmt.Errorf("put state error %s", err.Error())
		return shim.Error(err.Error())
	}

	fmt.Printf("register user: %s, value: %s\n", userId, value)
	return shim.Success(nil)
}

//query record by userId
func (chainCode *RebateChainCode) query(stub shim.ChaincodeStubInterface, userId string) pb.Response {
	fmt.Println("================== query ====================")

	state, err := stub.GetState(userId)
	if nil != err {
		return shim.Error("query state by userId:" + userId + " error: " + err.Error())
	}

	return shim.Success(state)
}

//rebate directly
func (chainCode *RebateChainCode) rebateDirectly(stub shim.ChaincodeStubInterface, source string, destination string, delta string) pb.Response {
	fmt.Println("================== rebateDirectly ====================")

	//load state from ledger
	sourceState, sourceStateErr := stub.GetState(source)
	if nil != sourceStateErr {
		return shim.Error("load source state error key:" + source + ", error:" + sourceStateErr.Error())
	}

	destinationState, destinationStateErr := stub.GetState(destination)
	if nil != destinationStateErr {
		return shim.Error("load destination state error key:" + destination + ", error:" + destinationStateErr.Error())
	}

	//parse state
	sourceInt, sourceIntErr := strconv.Atoi(string(sourceState))
	if nil != sourceIntErr {
		return shim.Error("conv sourceState error:" + sourceIntErr.Error())
	}
	destinationInt, destinationIntErr := strconv.Atoi(string(destinationState))
	if nil != destinationIntErr {
		return shim.Error("conv destinationState error:" + destinationIntErr.Error())
	}

	deltaInt, deltaIntErr := strconv.Atoi(delta)
	if nil != deltaIntErr {
		return shim.Error("conv delta:" + delta + " error:" + deltaIntErr.Error())
	}

	//rebate
	sourceInt -= deltaInt
	destinationInt += deltaInt


	//write state to ledger
	sourcePutErr := stub.PutState(source, []byte(string(sourceInt)))
	if nil != sourcePutErr {
		return shim.Error("put source " + source + " back to ledger error:" + sourcePutErr.Error())
	}
	destinationPutErr := stub.PutState(destination, []byte(string(destinationInt)))
	if nil != destinationPutErr {
		return shim.Error("put source " + destination + " back to ledger error:" + destinationPutErr.Error())
	}

	return shim.Success(nil)
}

func main()  {
	err := shim.Start(new(RebateChainCode))

	if nil != err {
		fmt.Printf("Error starting chaincode rebate_directly_cc, error: %s\n", err)
	}
}