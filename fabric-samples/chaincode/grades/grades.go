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
type Grade struct {
	Course      string `json:"course"`
	StudentName string `json:"studentName"`
	Practice    string `json:"practice"`
	Theory      string `json:"theory"`
	Send		string `json:"send"`
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
	if function == "queryGrade" {
		return s.queryGrade(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createGrade" {
		return s.createGrade(APIstub, args)
	} else if function == "queryAllGrades" {
		return s.queryAllGrades(APIstub)
	} else if function == "changeGrades" {
		return s.changeGrades(APIstub, args)
	}else if function == "getHistory" {
		return s.getHistory(APIstub, args)
	}


	return shim.Error("Nombre de la función no es valida.")
}

func (s *SmartContract) queryGrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Error se esperaban 1 argumento.")
	}

	gradeAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(gradeAsBytes)
}
// Inicialización datos del contrato.
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	grades := []Grade{
		Grade{Course: "CC",   StudentName: "Juan Jose Lopez",   Practice: "8.4", Theory: "8.5", Send:"0" },
		//Grade{Course: "SIGE", StudentName: "Juan Jose Lopez",   Practice: "7.5", Theory: "5" },
		//Grade{Course: "SSWB", StudentName: "Rodrigo Perez",     Practice: "3",  Theory: "7" },
		//Grade{Course: "PGPI", StudentName: "Adrian Rodriguez",  Practice: "8.5", Theory: "7.5"},
	}

	i := 0
	for i < len(grades) {
		fmt.Println("i is ", i)
		gradeAsBytes, _ := json.Marshal(grades[i])
		APIstub.PutState("NOTA"+strconv.Itoa(i), gradeAsBytes)
		fmt.Println("Añadida", grades[i])
		i = i + 1
	}

	return shim.Success(nil)
}

//Crear un registro academico de notas.
func (s *SmartContract) createGrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Numero incorrectos de argumentos se esperaban 5.")
	}

	var grade = Grade{Course: args[1], StudentName: args[2], Theory: args[3], Practice: args[4], Send: "0"}

	gradeAsBytes, _ := json.Marshal(grade)
	APIstub.PutState(args[0], gradeAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllGrades(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Añadir la coma.
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		//se escribe el registro como un objeto Json
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllGrades:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeGrades(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Numero incorrecto de parametros, se esperaban 3")
	}

	gradeAsBytes, _ := APIstub.GetState(args[0])
	grade := Grade{}

	json.Unmarshal(gradeAsBytes, &grade)
	grade.Practice = args[1]
	grade.Theory = args[2]

	gradeAsBytes, _ = json.Marshal(grade)
	APIstub.PutState(args[0], gradeAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	type GradeHistory struct {
		TxId    string   	 `json:"txId"`
		Value   Grade  		 `json:"value"`
	}
	var history []GradeHistory;
	var grade Grade

	if len(args) != 1 {
		return shim.Error("Numero incorrecto de argumentos se esperaban 1")
	}

	gradeId := args[0]
	fmt.Printf(" Historial del registro academico. %s\n", gradeId)

	// Obtener el historial academico del registro.
	resultsIterator, err := APIstub.GetHistoryForKey(gradeId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx GradeHistory
		tx.TxId = historyData.TxId                     //Se copia el id de la transaccion.
		json.Unmarshal(historyData.Value, &grade)     //Se convierte el historial a json
		if historyData.Value == nil {                  //Si el registro Grade ha sido borrado.
			var emptyGrade Grade
			tx.Value = emptyGrade                 //Se copa el nodo vacio
		} else {
			json.Unmarshal(historyData.Value, &grade) //Se convierte el historial a json
			tx.Value = grade                      //se copia el registros
		}
		history = append(history, tx)              //Se añade la transacción a la lista.
	}
	fmt.Printf("- Historial del registro academico :\n%s", history)

	//se conbierte a arreglo a bytes
	historyAsBytes, _ := json.Marshal(history)     
	return shim.Success(historyAsBytes)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}

}
