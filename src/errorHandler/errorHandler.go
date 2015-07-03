package errorHandler

import (
	"os"
)

func HandleError(err error) {
	if err != nil {
		panic(err)
		os.Exit(1)
	}
}
