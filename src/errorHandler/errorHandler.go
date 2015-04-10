package errorHandler

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
}
