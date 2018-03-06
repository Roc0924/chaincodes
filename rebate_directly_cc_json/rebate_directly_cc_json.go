/**
 * Create with GoLand
 * Author               : wangzhenpeng
 * Date                 : 2018/3/6
 * Time                 : 下午4:33
 * Description          :
 */
package rebate_directly_cc_json

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

type RebateChainCodeJSon struct{
}

//Init chaincode
func (chainCode *RebateChainCodeJSon) Init(stub shim.ChaincodeStubInterface) pb.Response{
	fmt.Println("Init chaincode rebate_directly_cc_json")
	if transientMap, err := stub.GetTransient(); nil == err {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)
}

//Invoke method
func (chainCode *RebateChainCodeJSon) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Query the method invoke of chaincode rebate_directly_cc\n")
	function, args := stub.GetFunctionAndParameters()
	if "invoke" != function {
		return shim.Error("unknown function call: " + function)
	}
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if "registerJson" == args[0] {
		if len(args) != 3 {
			return shim.Error("Call method registerJson error. Incorrect number of arguments. Expecting 3")
		}

		return chainCode.registerJson(stub, args[1], args[2])
	}

	if "queryJson" == args[0] {
		if len(args) != 2 {
			return shim.Error("Call method queryJson error. Incorrect number of arguments. Expecting 2")
		}

		return chainCode.queryJson(stub, args[1])
	}


	return shim.Success(nil)
}
func (chainCode *RebateChainCodeJSon) registerJson(stub shim.ChaincodeStubInterface, userId string, value string) pb.Response {
	fmt.Println("================== registerJson ====================")


	err := stub.PutState(userId, []byte(value))
	if nil != err {
		fmt.Errorf("put state error %s", err.Error())
		return shim.Error(err.Error())
	}

	fmt.Printf("register user: %s, value: %s\n", userId, value)
	return shim.Success(nil)
}
func (chainCode *RebateChainCodeJSon) queryJson(stub shim.ChaincodeStubInterface, userId string) pb.Response {
	fmt.Println("================== queryJson ====================")

	state, err := stub.GetState(userId)
	if nil != err {
		return shim.Error("query state by userId:" + userId + " error: " + err.Error())
	}

	return shim.Success(state)

}

func main() {
	err := shim.Start(new(RebateChainCodeJSon))

	if nil != err {
		fmt.Printf("Error starting chaincode rebate_directly_cc, error: %s\n", err)
	}
}

