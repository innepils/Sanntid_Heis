package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config represents a key-value storage for configuration parameters.
type Config map[string]string

// LoadConfig loads configuration from a file into a Config map.
func LoadConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file %s: %w", filePath, err)
	}
	defer file.Close()

	config := make(Config)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			parts := strings.Fields(line[2:])
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", filePath, err)
	}

	return config, nil
}

// GetValue retrieves a value from the Config and tries to parse it into the provided variable, based on the format specifier.
func (c Config) GetValue(key string, varPtr interface{}, fmtStr string) error {
	value, ok := c[key]
	if !ok {
		return fmt.Errorf("key %s not found in config", key)
	}

	_, err := fmt.Sscanf(value, fmtStr, varPtr)
	return err
}

// // Example usage of the above functions:

// // Assuming you have an enum-like type and a matching function to convert strings to this type's values.
// type MyEnumType int

// const (
// 	EnumValue1 MyEnumType = iota
// 	EnumValue2
// )

// // MatchStringToMyEnumType matches string values to MyEnumType constants.
// func MatchStringToMyEnumType(val string) (MyEnumType, bool) {
// 	switch strings.ToLower(val) {
// 	case "enumvalue1":
// 		return EnumValue1, true
// 	case "enumvalue2":
// 		return EnumValue2, true
// 	default:
// 		return 0, false
// 	}
// }

// // SetEnumValue reads a string from the config and uses a match function to set an enum variable.
// func SetEnumValue(c Config, key string, enumVar *MyEnumType, matchFunc func(string) (MyEnumType, bool)) error {
// 	valStr, ok := c[key]
// 	if !ok {
// 		return fmt.Errorf("key %s not found", key)
// 	}

// 	enumVal, matched := matchFunc(valStr)
// 	if !matched {
// 		return fmt.Errorf("value %s for key %s is not a valid enum value", valStr, key)
// 	}

// 	*enumVar = enumVal
// 	return nil
// }
