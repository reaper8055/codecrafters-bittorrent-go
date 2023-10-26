package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func trimListDelimiters(bencodedString string) (string, int) {
	l := len(bencodedString)
	i := 0
	for l > 1 {
		if bencodedString[i] == 'l' && bencodedString[l-1] == 'e' {
			i++
			l--
			continue
		}
		break
	}
	// fmt.Println("26: ", bencodedString[i:l])
	return bencodedString[i:l], i - 1
}

func format(s interface{}, level int) interface{} {
	if level <= 0 {
		return s
	}
	result := []interface{}{format(s, level-1)}
	return result
}

func printType(v interface{}) {
	fmt.Println(reflect.TypeOf(v))
}

func decodeString(bencodedString string) (string, error) {
	l := len(bencodedString)
	colonIdx := 0
	for idx, ch := range bencodedString {
		if ch != ':' {
			continue
		}
		colonIdx = idx
		break
	}
	return bencodedString[colonIdx+1 : l], nil
}

func decodeInteger(bencodedString string) (int, error) {
	l := len(bencodedString)
	decodedInteger, err := strconv.Atoi(bencodedString[1 : l-1])
	if err != nil {
		return 0, err
	}
	return decodedInteger, nil
}

func getColonIndex(bencodedString string) (int, error) {
	var colonIdx int
	for idx, char := range bencodedString {
		if char != ':' {
			continue
		}
		colonIdx = idx
		return colonIdx, nil
	}
	return 0, errors.New("colon not found in the given string")
}

func decodeList(bencodedString string) (interface{}, error) {
	l := len(bencodedString)
	var colonIdx, decodedInteger int
	var decodedString string
	firstChar := bencodedString[0]
	if firstChar == 'i' {
		colonIdx, _ = getColonIndex(bencodedString)
		decodedString, _ = decodeString(bencodedString[colonIdx-1 : l])
		decodedInteger, _ = decodeInteger(bencodedString[:colonIdx-1])
		return []interface{}{
			decodedInteger,
			decodedString,
		}, nil
	}
	colonIdx, _ = getColonIndex(bencodedString)
	lengthOfString, _ := strconv.Atoi(bencodedString[:colonIdx])
	decodedString, _ = decodeString(bencodedString[:lengthOfString+colonIdx+1])
	decodedInteger, _ = decodeInteger(bencodedString[lengthOfString+colonIdx+1:])
	return []interface{}{
		decodedString,
		decodedInteger,
	}, nil
}

func decodeBencode(bencodedString string) (interface{}, error) {
	l := len(bencodedString)
	var level int

	firstChar := bencodedString[0]
	lastChar := bencodedString[l-1]

	if firstChar == 'l' && lastChar == 'e' {
		bencodedString, level = trimListDelimiters(bencodedString)
		if len(bencodedString) == 0 {
			return []interface{}{}, nil
		}
		result, _ := decodeList(bencodedString)
		r := format(result, level)
		return r, nil
	} else if firstChar == 'i' && lastChar == 'e' {
		result, _ := decodeInteger(bencodedString)
		return result, nil
	}
	result, _ := decodeString(bencodedString)
	return result, nil
}

func main() {
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
