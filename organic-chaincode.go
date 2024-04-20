// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
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

/* Define Tuna structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type Organic struct {
	Shipment string `json:"shipment"`
	Timestamp string `json:"timestamp"`
	Location  string `json:"location"`
	Holder  string `json:"holder"`
}

/*
 * The Init method *
 called when the Smart Contract "tuna-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "tuna-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryOrganic" {
		return s.queryOrganic(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordOrganic" {
		return s.recordOrganic(APIstub, args)
	} else if function == "queryAllOrganic" {
		return s.queryAllOrganic(APIstub)
	} else if function == "changeOrganicHolder" {
		return s.changeOrganicHolder(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryTuna method *
Used to view the records of one particular tuna
It takes one argument -- the key for the tuna in question
 */
func (s *SmartContract) queryOrganic(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	organicAsBytes, _ := APIstub.GetState(args[0])
	if organicAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	return shim.Success(organicAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 tuna catches)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	organic := []Organic{
		Organic{Shipment: "923F", Location: "67.0006, -70.5476", Timestamp: "1504054225", Holder: "Miriam"},
		Organic{Shipment: "M83T", Location: "91.2395, -49.4594", Timestamp: "1504057825", Holder: "Dave"},
		Organic{Shipment: "T012", Location: "58.0148, 59.01391", Timestamp: "1493517025", Holder: "Igor"},
		Organic{Shipment: "P490", Location: "-45.0945, 0.7949", Timestamp: "1496105425", Holder: "Amalea"},
		Organic{Shipment: "S439", Location: "-107.6043, 19.5003", Timestamp: "1493512301", Holder: "Rafa"},
		Organic{Shipment: "J205", Location: "-155.2304, -15.8723", Timestamp: "1494117101", Holder: "Shen"},
		Organic{Shipment: "S22L", Location: "103.8842, 22.1277", Timestamp: "1496104301", Holder: "Leila"},
		Organic{Shipment: "EI89", Location: "-132.3207, -34.0983", Timestamp: "1485066691", Holder: "Yuan"},
		Organic{Shipment: "129R", Location: "153.0054, 12.6429", Timestamp: "1485153091", Holder: "Carlo"},
		Organic{Shipment: "49W4", Location: "51.9435, 8.2735", Timestamp: "1487745091", Holder: "Fatima"},
	}

	i := 0
	for i < len(organic) {
		fmt.Println("i is ", i)
		organicAsBytes, _ := json.Marshal(organic[i])
		APIstub.PutState(strconv.Itoa(i+1), organicAsBytes)
		fmt.Println("Added", organic[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordTuna method *
Fisherman like Sarah would use to record each of her tuna catches. 
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordOrganic(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var organic = Organic{ Shipment: args[1], Location: args[2], Timestamp: args[3], Holder: args[4] }

	organicAsBytes, _ := json.Marshal(organic)
	err := APIstub.PutState(args[0], organicAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record product record: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllTuna method *
allows for assessing all the records added to the ledger(all tuna catches)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllOrganic(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

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
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllOrganic:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeTunaHolder method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, tuna id and new holder name. 
 */
func (s *SmartContract) changeOrganicHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	organicAsBytes, _ := APIstub.GetState(args[0])
	if organicAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	organic := Organic{}

	json.Unmarshal(organicAsBytes, &organic)
	// Normally check that the specified argument is a valid holder of tuna
	// we are skipping this check for this example
	organic.Holder = args[1]

	organicAsBytes, _ = json.Marshal(organic)
	err := APIstub.PutState(args[0], organicAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change organic holder: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
