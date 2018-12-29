package utils

import (
	"fmt"
)

// Err - Err
func Err(context string, err error) bool {
	if err != nil {
		fmt.Println(context + ": " + err.Error())
		return true
	}
	return false
}
