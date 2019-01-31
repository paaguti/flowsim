package common

// A dictionary to map DSCP labels to values
// CAVEAT: Multiply by 4 to get the full TOS byte value

import (
	"errors"
	"fmt"
)

var dscpDict = map[string]int{
	"CS0":  0,
	"CS1":  8,
	"CS2":  16,
	"CS3":  24,
	"CS4":  32,
	"CS5":  40,
	"CS6":  48,
	"CS7":  56,
	"EF":   46,
	"AF11": 10,
	"AF12": 12,
	"AF13": 14,
	"AF21": 18,
	"AF22": 20,
	"AF23": 22,
	"AF31": 26,
	"AF32": 28,
	"AF33": 30,
	"AF41": 34,
	"AF42": 36,
	"AF43": 38,
}

func Dscp(s string) (int, error) {
	var val int
	if val, ok := dscpDict[s]; ok {
		return val, nil
	}
	_, err := fmt.Sscanf(s, "%d", &val)
	if err != nil {
		return 0, errors.New("Unknown DSCP ID ")
	}
	if val < 0 || val > 64 {
		return 0, errors.New("Value out of range [0,64)")
	}
	return val, nil
}

func ToDscp(dscp int) (string, error) {
	for k, v := range dscpDict {
		if v == dscp {
			return k, nil
		}
	}
	return "Undefined", errors.New(fmt.Sprintf("Value %d not mapped to DSCP", dscp))
}
