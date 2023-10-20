package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func trimIdentifiers(bencodedString string) string {
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
	// fmt.Println(bencodedString[i:l])
	return bencodedString[i:l]
}

func decodeBencode(bencodedString string) (interface{}, error) {
	bencodedString = trimIdentifiers(bencodedString)
	if len(bencodedString) == 0 {
		return []interface{}{}, nil
	}
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
		for idx, ch := range bencodedString {
			if ch == ':' {
				i, _ := strconv.Atoi(bencodedString[1 : idx-2])
				return []interface{}{
					bencodedString[idx+1:],
					i,
				}, nil
			}
		}
		i, _ := strconv.Atoi(bencodedString[1 : len(bencodedString)-1])
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
