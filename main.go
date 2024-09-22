package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	input := `{
  "number_1": {
    "N": "1.50"
  },
  "string_1": {
    "S": "784498 "
  },
  "string_2": {
    "S": "2014-07-16T20:55:46Z"
  },
  "map_1": {
    "M": {
      "bool_1": {
        "BOOL": "truthy"
      },
      "null_1": {
        "NULL ": "true"
      },
      "list_1": {
        "L": [
          {
            "S": ""
          },
          {
            "N": "011"
          },
          {
            "N": "5215s"
          },
          {
            "BOOL": "f"
          },
          {
            "NULL": "0"
          }
        ]
      }
    }
  },
  "list_2": {
    "L": "noop"
  },
  "list_3": {
    "L": [
      "noop"
    ]
  },
  "": {
    "S": "noop"
  }
}`

	var inputData map[string]any
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		fmt.Println("error parsing input:", err)
		return
	}

	resultList := []map[string]interface{}{}
	result := map[string]interface{}{}

	for key, value := range inputData {
		resultKey := trimString(key)
		if resultKey != "" {
			resultValue := transformValue(value.(map[string]any))
			if resultValue != nil {
				result[resultKey] = resultValue
			}
		}
	}

	resultList = append(resultList, result)

	output, err := json.MarshalIndent(resultList, "", "  ")
	if err != nil {
		fmt.Println("Error generating output:", err)
		return
	}

	fmt.Println("Result: ")
	fmt.Println(string(output))
}

func trimString(s string) string {
	return strings.TrimSpace(s)
}

func transformValue(input map[string]any) any {
	for key, value := range input {
		key = trimString(key)
		switch key {
		case "S":
			val := trimString(value.(string))
			if val == "" {
				return nil
			}

			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				fmt.Println("not a timestamp...it is string value")
				return val
			}

			return t.Unix()

		case "N":
			val := trimString(value.(string))
			trimmedVal := strings.TrimLeft(val, "0")
			num, err := strconv.ParseFloat(trimmedVal, 64)
			if err != nil {
				fmt.Println("invalid integer...")
				return nil
			}

			return num

		case "BOOL":
			str := trimString(value.(string))
			if strings.ToLower(str) == "true" || str == "1" || strings.ToLower(str) == "t" {
				return true
			} else if strings.ToLower(str) == "false" || str == "0" || strings.ToLower(str) == "f" {
				return false
			}

			return nil

		case "NULL":
			str := trimString(value.(string))
			if strings.ToLower(str) == "true" || str == "1" || strings.ToLower(str) == "t" {
				return "null"
			}

			return nil

		case "L":
			list, ok := value.([]any)
			if !ok {
				fmt.Println("invalid list...")
				return nil
			}

			resultList := []interface{}{}
			for _, item := range list {
				listItem, ok := item.(map[string]any)
				if ok {
					resultItem := transformValue(listItem)
					if resultItem != nil {
						resultList = append(resultList, resultItem)
					}
				}

			}

			if len(resultList) > 0 {
				return resultList
			}

			return nil

		case "M":
			mapValues, ok := value.(map[string]any)
			if !ok {
				fmt.Println("invalid map...")
				return nil
			}

			resultMap := map[string]interface{}{}
			for k, val := range mapValues {
				resultKey := trimString(k)
				mapValue, ok := val.(map[string]any)
				if ok {
					resultVal := transformValue(mapValue)
					if resultVal == "null" {
						resultMap[resultKey] = nil
					} else if resultVal != nil {
						resultMap[resultKey] = resultVal
					}
				}
			}

			if len(resultMap) > 0 {
				return resultMap
			}

			return nil

		default:
			fmt.Println("invalid key...")
			return nil
		}
	}
	return nil
}
