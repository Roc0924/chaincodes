/**
 * Create with GoLand
 * Author               : wangzhenpeng
 * Date                 : 2018/3/6
 * Time                 : 下午2:33
 * Description          : directly rebate chaincode
 */
package rebate_direct_cc

import (
	"fmt"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/orderer/multichain"
)
//define chaincode struct
type RebateChainCode struct {
}


//init chaincode
func (chainCode *RebateChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Init chaincode rebate_direct_cc\n")


	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}

	}

	return shim.Success(nil)
}


//invoke method
func (chainCode *RebateChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("query the method invoke of chaincode rebate_direct_cc\n")
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


	return shim.Success(nil)
}


//register a record
func (chainCode *RebateChainCode) register(stub shim.ChaincodeStubInterface, userId string, value string) pb.Response {


	return shim.Success(nil)
}
func (chainCode *RebateChainCode) query(stub shim.ChaincodeStubInterface, userId string) pb.Response {


	return shim.Success(nil)
}
