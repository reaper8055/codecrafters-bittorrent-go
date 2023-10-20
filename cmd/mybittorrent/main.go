package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

type DecodedData struct {
	S string
	I int
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else {
		bencodedSlice := strings.Split(bencodedString, "")
		if bencodedSlice[0] == "l" {
			if len(bencodedSlice) == 2 && bencodedSlice[1] == "e" {
				return []interface{}{}, nil
			}
			idx, err := strconv.Atoi(bencodedSlice[1])
			idx += 3
			if err != nil {
				return nil, err
			}
			i, _ := strconv.Atoi(strings.Join(bencodedSlice[idx+1:len(bencodedSlice)-2], ""))
			s := strings.Join(bencodedSlice[3:idx], "")
			return []interface{}{s, i}, nil
		}
		var result string
		for _, v := range bencodedSlice[1 : len(bencodedSlice)-1] {
			result += v
		}
		i, err := strconv.Atoi(result)
		if err != nil {
			return nil, err
		}
		return i, nil
	}
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
