package main

import (
	"fmt"

	"github.com/alijared/nr-log-parser/internal/cmd"
	"github.com/alijared/nr-log-parser/pkg/errors"
)

func main() {
	if err := cmd.Execute(); err != nil {
		err := err.(errors.CMDError)
		if err.Type() == errors.VALIDATION_ERROR {
			fmt.Println("Error:", err.Error())
			fmt.Println(err.Usage())
		} else {
			fmt.Println("Error executing command:", err.Error())
		}
	}
}
