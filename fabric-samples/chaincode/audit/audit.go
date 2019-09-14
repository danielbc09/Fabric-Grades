/*
 * Smart contract for the notes aplication.
 * Autor: Jairo Daniel Bautista Castro.
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Grades structure, with 4 properties.  Structure tags are used by encoding/json library
type AuditGrade struct {
	Course      	  string `json:"course"`
	StudentName 	  string `json:"studentName"`
	Practice    	  string `json:"practice"`
	Theory      	  string `json:"theory"`
	Department        string `json:"department"`
	Validated		  string `json:"validated"`
}

/*
 * The Init method is called when the Smart Contract "Grades" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "grades
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryAuditGrade" {
		return s.queryAuditGrade(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createAuditGrade" {
		return s.createAuditGrade(APIstub, args)
	} else if function == "queryAllAuditGrades" {
		return s.queryAllAuditGrades(APIstub)
	} else if function == "changeStatus" {
		return s.changeStatusGrades(APIstub, args)
	} else if function == "getHistory" {
		return s.getHistory(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryAuditGrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	auditAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(auditAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	
	auditGrades := []AuditGrade{
		AuditGrade{Course: "CC", StudentName: "Juan Jose Lopez", Practice: "8.4", Theory: "8.5", Department:"decsai", Validated:"0"},
		AuditGrade{Course: "SIGE", StudentName: "Juan Jose Lopez", Practice: "7.5", Theory: "5", Department:"atc", Validated:"0"},
	}

	i := 0
	for i < len(auditGrades) {
		fmt.Println("i is ", i)
		auditAsBytes, _ := json.Marshal(auditGrades[i])
		APIstub.PutState("AUDITNOTA"+strconv.Itoa(i), auditAsBytes)
		fmt.Println("Added", auditGrades[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createAuditGrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var audit = AuditGrade{Course: args[1], StudentName: args[2], Theory: args[3], Practice: args[4], Department: args[5], Validated: args[6] }

	auditAsBytes, _ := json.Marshal(audit)
	APIstub.PutState(args[0], auditAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllAuditGrades(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey:=""
	endKey:=""
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllGrades:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeStatusGrades(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	auditGradeAsBytes, _ := APIstub.GetState(args[0])
	auditGrade := AuditGrade{}

	json.Unmarshal(auditGradeAsBytes, &auditGrade)
	auditGrade.Validated = args[1]

	auditGradeAsBytes, _ = json.Marshal(auditGrade)
	APIstub.PutState(args[0], auditGradeAsBytes)

	return shim.Success(nil)
}


func (s *SmartContract) getHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	type AuditHistory struct {
		TxId    string   	 `json:"txId"`
		Value   AuditGrade   `json:"value"`
	}
	var history []AuditHistory;
	var audit AuditGrade

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	auditId := args[0]
	fmt.Printf("- start getHistoryForMarble: %s\n", auditId)

	// Get History
	resultsIterator, err := APIstub.GetHistoryForKey(auditId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = historyData.TxId                     //copy transaction id over
		json.Unmarshal(historyData.Value, &audit)     //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //marble has been deleted
			var emptyAudit AuditGrade
			tx.Value = emptyAudit                 //copy nil marble
		} else {
			json.Unmarshal(historyData.Value, &audit) //un stringify it aka JSON.parse()
			tx.Value = audit                      //copy marble over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForMarble returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}



func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

