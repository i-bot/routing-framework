package errorHandler

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
