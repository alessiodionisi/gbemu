package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type OpCode struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Parameters     []string          `json:"parameters"`
	Flags          map[string]string `json:"flags"`
	Bits           int               `json:"bits"`
	CyclesBranch   int               `json:"cycles_branch"`
	CyclesNoBranch int               `json:"cycles_no_branch"`
}

type OpCodes struct {
	Unprefixed map[string]OpCode `json:"unprefixed"`
	CBPrefixed map[string]OpCode `json:"cb_prefixed"`
}

type InputOpCode struct {
	Name            string
	Group           string
	Flags           map[string]string
	TCyclesBranch   int
	TCyclesNoBranch int
}

type InputOpCodes struct {
	Unprefixed []InputOpCode
	CBPrefixed []InputOpCode
}

func parseInputOpCode(inputOpCode InputOpCode) OpCode {
	var name string
	params := make([]string, 0)

	nameSplitSpace := strings.Split(inputOpCode.Name, " ")

	if len(nameSplitSpace) > 1 {
		name = nameSplitSpace[0]
		nameSplitComma := strings.Split(nameSplitSpace[1], ",")
		params = nameSplitComma
	} else {
		name = inputOpCode.Name
	}

	bits := 8
	if strings.Contains(inputOpCode.Group, "x16") {
		bits = 16
	}

	fmt.Println(name, params)

	return OpCode{
		Name:           name,
		Description:    inputOpCode.Name,
		Parameters:     params,
		Flags:          inputOpCode.Flags,
		Bits:           bits,
		CyclesBranch:   inputOpCode.TCyclesBranch,
		CyclesNoBranch: inputOpCode.TCyclesNoBranch,
	}
}

func generateInstructionString(opCode string, opCodeData OpCode) string {
	parametersString := ""

	for _, parameter := range opCodeData.Parameters {
		parametersString += fmt.Sprintf("\"%s\",", parameter)
		/*if i+1 != len(opCodeData.Parameters) {
			parametersString += ","
		}*/
	}

	flagsString := ""

	for key, flag := range opCodeData.Flags {
		flagsString += fmt.Sprintf("\"%s\":\"%s\",", key, flag)
		/*if i+1 != len(opCodeData.Parameters) {
			parametersString += ","
		}*/
	}

	return fmt.Sprintf("%s: {OpCode: %s, Name: \"%s\", Description: \"%s\", Parameters: []string{%s}, Flags: map[string]string{%s}, Bits: %d, CyclesBranch: %d, CyclesNoBranch: %d }", opCode, opCode, opCodeData.Name, opCodeData.Description, parametersString, flagsString, opCodeData.Bits, opCodeData.CyclesBranch, opCodeData.CyclesNoBranch)
}

func main() {
	inputFile, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	inputFileStat, err := inputFile.Stat()
	if err != nil {
		panic(err)
	}

	inputFileBytes := make([]byte, inputFileStat.Size())

	_, err = inputFile.Read(inputFileBytes)
	if err != nil {
		panic(err)
	}

	err = inputFile.Close()
	if err != nil {
		panic(err)
	}

	var inputOpCodes InputOpCodes
	err = json.Unmarshal(inputFileBytes, &inputOpCodes)
	if err != nil {
		panic(err)
	}

	unprefixedOps := make(map[string]OpCode)
	for opCodeNumber, inputOpCode := range inputOpCodes.Unprefixed {
		opCodeHex := fmt.Sprintf("%#02x", opCodeNumber)
		fmt.Printf("%s; %s\n", opCodeHex, inputOpCode.Name)

		opCode := parseInputOpCode(inputOpCode)

		unprefixedOps[opCodeHex] = opCode
	}

	cbPrefixedOps := make(map[string]OpCode)
	for opCodeNumber, inputOpCode := range inputOpCodes.CBPrefixed {
		opCodeHex := fmt.Sprintf("%#02x", opCodeNumber)
		fmt.Printf("%s; %s\n", opCodeHex, inputOpCode.Name)

		opCode := parseInputOpCode(inputOpCode)

		cbPrefixedOps[opCodeHex] = opCode
	}

	opCodes := OpCodes{
		Unprefixed: unprefixedOps,
		CBPrefixed: cbPrefixedOps,
	}

	opCodesBytes, err := json.Marshal(opCodes)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(opCodesBytes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", opCodes)

	var goInstructions string
	var goPrefixedInstructions string

	for opCode, opCodeData := range opCodes.Unprefixed {
		inst := generateInstructionString(opCode, opCodeData)
		goInstructions += fmt.Sprintf("%s,\n", inst)
	}

	for opCode, opCodeData := range opCodes.CBPrefixed {
		inst := generateInstructionString(opCode, opCodeData)
		goPrefixedInstructions += fmt.Sprintf("%s,\n", inst)
	}

	goInstructionsFileContent := fmt.Sprintf("var Instructions = map[uint8]*Instruction{\n%s}\n\nvar PrefixedInstructions = map[uint8]*Instruction{\n%s}\n", goInstructions, goPrefixedInstructions)

	fmt.Println(goInstructionsFileContent)

	goInstructionsFile, err := os.Create("instructions_map")
	if err != nil {
		panic(err)
	}
	defer goInstructionsFile.Close()

	_, err = goInstructionsFile.WriteString(goInstructionsFileContent)
	if err != nil {
		panic(err)
	}
}
