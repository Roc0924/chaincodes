/**
 * Create with GoLand
 * Author               : wangzhenpeng
 * Date                 : 2018/3/9
 * Time                 : 下午4:50
 * Description          : rebate chain code
 */
package rebate_cc



import (
	"fmt"
	"strconv"
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Rebate chaincode implementation
type RebateChaincode struct {
}

type RebateAccount struct{
	AccountId string
	Amount,ExpectAmount  int64
	Status string // normal,frozen,stop
	Details string
	Memo string
}

func (chaincode *RebateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init chaincode rebate_cc")

	return shim.Success(nil)
}

func (chaincode *RebateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("call method Invoke")
	function, args := stub.GetFunctionAndParameters()
	if "invoke" != function {
		return shim.Error("unknown function invoke")
	}
	if args[0] == "createAccount" {
		// get one key's history records
		return chaincode.createAccount(stub,args)
	} else if args[0] == "queryAccount" {
		// get one key's history records
		return chaincode.queryAccount(stub,args)
	} else if args[0] == "deleteAccount" {
		// Deletes account from ledger state
		return chaincode.deleteAccount(stub, args)
	} else if args[0] == "queryHistory" {
		// get one key's history records
		return chaincode.queryHistory(stub,args)
	} else if args[0] == "createPlan" {
		// get one key's history records
		return chaincode.createPlan(stub,args)
	} else if args[0] == "queryPlan" {
		// the old "Query" is now implemtned in invoke
		return chaincode.queryPlan(stub, args)
	} else if args[0] == "addAmountFromBudget" {
		// get one key's history records
		return chaincode.addAmountFromBudget(stub,args)
	} else if args[0] == "addAmountFromExpect" {
		// get one key's history records
		return chaincode.addAmountFromExpect(stub,args)
	} /*else if function == "minusAmount" {
		// get one key's history records
		return chaincode.minusAmount(stub,args)
	} else if function == "addExpectAmount" {
		// get one key's history records
		return chaincode.addExpectAmount(stub,args)
	} else if function == "minusExpectAmount" {
		// get one key's history records
		return chaincode.minusExpectAmount(stub,args)
	}*/

	return shim.Error("Invalid invoke function name. Expecting 'createAccount'" +
		" 'queryAccount' 'deleteAccount' 'queryHistory' 'createPlan' 'queryPlan' 'addAmountFromBudget' 'query' 'query'")
}
func (chaincode *RebateChaincode) queryHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	key:=args[1]
	it,err:= stub.GetHistoryForKey(key)
	if err!=nil{
		return shim.Error(err.Error())
	}
	var result,_= getHistoryListResult(it)
	return shim.Success(result)
}
func getHistoryListResult(resultsIterator shim.HistoryQueryIteratorInterface) ([]byte,error){

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		item,_:= json.Marshal( queryResponse)
		buffer.Write(item)
		//      buffer.Write(queryResponse)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}


// create account in ledger
func (chaincode *RebateChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//var amount, expectAmount int64 // Asset holdings
	var err error
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	rebateAccountStr := args[1]
	var rebateAccount RebateAccount
	fmt.Println(rebateAccountStr);
	err = json.Unmarshal([]byte(rebateAccountStr),&rebateAccount)
	fmt.Printf("rebate account : accountId:%s, amount:%d, expectAmount:%d, status:%s, details:%s, memo:%s\n",
		rebateAccount.AccountId, rebateAccount.Amount, rebateAccount.ExpectAmount, rebateAccount.Status,
			rebateAccount.Details, rebateAccount.Memo)

	if nil != err {
		return shim.Error("unmarshal error :" + rebateAccountStr);
	}




	//account := args[1]
	//amount, err = strconv.ParseInt(args[2], 10, 64)
	//if err != nil {
	//	return shim.Error("Expecting integer value for asset holding")
	//}
	//expectAmount, err = strconv.ParseInt(args[2], 10, 64)
	//if err != nil {
	//	return shim.Error("Expecting integer value for asset holding")
	//}


	//whether the account has registered
	record, recordErr := stub.GetState(rebateAccount.AccountId)
	if nil != recordErr {
		return shim.Error("Inner Error" + recordErr.Error())
	}

	if nil != record {
		return shim.Error("Account has registered")
	}


	//rebateAccount := &RebateAccount{amount: amount, expectAmount: expectAmount,status:args[3],details:args[4],memo:args[5]}

	//byteObject,_ := json.Marshal(rebateAccount)
	// Put the key into the state in ledger
	putErr := stub.PutState(rebateAccount.AccountId,[]byte(rebateAccountStr))
	if putErr != nil {
		return shim.Error("Failed to put state, error:" + putErr.Error())
	}

	return shim.Success(nil)
}
// Deletes an entity from state
func (chaincode *RebateChaincode) createPlan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	A := args[1]
	val := args[2]
	A ="plan_"+A
	// Delete the key from the state in ledger
	err := stub.PutState(A,[]byte(val))
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}
// Deletes an entity from state
func (chaincode *RebateChaincode) deleteAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accountId := args[1]

	accountStateByte, getErr := stub.GetState(accountId)
	if nil != getErr {
		shim.Error("get state of " + accountId + " error:" + getErr.Error())
	}

	// Delete the key from the state in ledger
	err := stub.DelState(accountId)
	if err != nil {
		return shim.Error("Failed to delete state， error:" + err.Error())
	}

	return shim.Success(accountStateByte)
}

// query callback representing the query of a chaincode
func (chaincode *RebateChaincode) queryPlan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[1]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState("plan_"+A)
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
// query callback representing the query of a chaincode
func (chaincode *RebateChaincode) queryAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountId string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accountId = args[1]

	// Get the state from the ledger
	recordbytes, err := stub.GetState(accountId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + accountId + "\"}"
		return shim.Error(jsonResp)
	}

	if recordbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + accountId + "\"}"
		return shim.Error(jsonResp)
	}

	//jsonResp := "{\"Name\":\"" + accountId + "\",\"Amount\":\"" + string(recordbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", recordbytes)
	return shim.Success(recordbytes)
}

// query callback representing the query of a chaincode
func (chaincode *RebateChaincode) addAmountFromBudget(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var planId,accountId string
	var budgetVal,accountVal,accountByte []byte
	var budget,val int64
	var err error
	var account RebateAccount
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	planId = args[0]
	accountId = args[1]
	val, err = strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// Get the state from the ledger
	budgetVal, err = stub.GetState("plan_"+planId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + planId + "\"}"
		return shim.Error(jsonResp)
	}
	accountVal, err = stub.GetState(accountId)
	if err != nil{
		jsonResp :="{\"Error\":\"get account "+accountId +" err \"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal(accountVal,&account)
	if err != nil{
		jsonResp :="{\"Error\":\"account "+accountId +" unmarshal err \"}"
		return shim.Error(jsonResp)
	}

	budget,err = strconv.ParseInt(string(budgetVal), 10, 64)
	if err != nil {
		jsonResp := "{\"Error\":\"budget is not int \"}"
		return shim.Error(jsonResp)
	}
	budget = budget - val
	if budget < 0 {
		jsonResp := "{\"Error\":\"budget is not enough \"}"
		return shim.Error(jsonResp)
	}
	err = stub.PutState("plan_"+planId,[]byte(strconv.FormatInt(budget, 10)))
	if err != nil{
		return shim.Error(err.Error())
	}
	account.Amount = account.Amount + val
	accountByte,err = json.Marshal(account)
	if err != nil{
		jsonResp :="{\"Error\":\"account "+accountId +" format err \"}"
		return shim.Error(jsonResp)
	}
	err = stub.PutState(accountId,accountByte)
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func (chaincode *RebateChaincode) addExpectAmountFromBudget(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var planId,accountId string
	var budgetVal,accountVal,accountByte []byte
	var budget,val int64
	var err error
	var account RebateAccount
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	planId = args[0]
	accountId = args[1]
	val, err = strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// Get the state from the ledger
	budgetVal, err = stub.GetState("plan_"+planId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + planId + "\"}"
		return shim.Error(jsonResp)
	}
	accountVal, err = stub.GetState(accountId)
	if err != nil{
		jsonResp :="{\"Error\":\"get account "+accountId +" err \"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal(accountVal,&account)
	if err != nil{
		jsonResp :="{\"Error\":\"account "+accountId +" unmarshal err \"}"
		return shim.Error(jsonResp)
	}

	budget,err = strconv.ParseInt(string(budgetVal), 10, 64)
	if err != nil {
		jsonResp := "{\"Error\":\"budget is not int \"}"
		return shim.Error(jsonResp)
	}
	budget = budget - val
	if budget < 0 {
		jsonResp := "{\"Error\":\"budget is not enough \"}"
		return shim.Error(jsonResp)
	}
	err = stub.PutState("plan_"+planId,[]byte(strconv.FormatInt(budget, 10)))
	if err != nil{
		return shim.Error(err.Error())
	}
	account.ExpectAmount = account.ExpectAmount + val
	accountByte,err = json.Marshal(account)
	if err != nil{
		jsonResp :="{\"Error\":\"account "+accountId +" format err \"}"
		return shim.Error(jsonResp)
	}
	err = stub.PutState(accountId,accountByte)
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (chaincode *RebateChaincode) rollBackExpectRebate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var planId,accountId string
	var budgetVal,accountVal,accountByte []byte
	var budget,val int64
	var err error
	var account RebateAccount
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	planId = args[0]
	accountId = args[1]
	val, err = strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// Get the state from the ledger
	budgetVal, err = stub.GetState("plan_"+planId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + planId + "\"}"
		return shim.Error(jsonResp)
	}
	accountVal, err = stub.GetState(accountId)
	if err != nil{
		jsonResp :="{\"Error\":\"get account "+accountId +" err \"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal(accountVal,&account)
	if err != nil{
		jsonResp :="{\"Error\":\"account "+accountId +" unmarshal err \"}"
		return shim.Error(jsonResp)
	}

	budget,err = strconv.ParseInt(string(budgetVal), 10, 64)
	if err != nil {
		jsonResp := "{\"Error\":\"budget is not int \"}"
		return shim.Error(jsonResp)
	}
	budget = budget + val
	if budget < 0 {
		jsonResp := "{\"Error\":\"budget is not enough \"}"
		return shim.Error(jsonResp)
	}
	err = stub.PutState("plan_"+planId,[]byte(strconv.FormatInt(budget, 10)))
	if err != nil{
		return shim.Error(err.Error())
	}
	account.ExpectAmount = account.ExpectAmount - val
	accountByte,err = json.Marshal(account)
	if err != nil{
		jsonResp :="{\"Error\":\"account "+accountId +" format err \"}"
		return shim.Error(jsonResp)
	}
	err = stub.PutState(accountId,accountByte)
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}



func (chaincode *RebateChaincode) addAmountFromExpect(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}



func (chaincode *RebateChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}








func main() {
	err := shim.Start(new(RebateChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
