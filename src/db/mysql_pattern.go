package db

import ()

var (
	SELECT       = fillRequest(2, 3, fillSelect)
	CREATE_TABLE = fillRequest(3, -1, fillCreate_Table)
	INSERT_INTO = fillRequest(3, -1, fillInsert_Into)
)

type fill func(values []string) (request string)

type MySQL_Request struct {
	min_args, max_args int
}

func fillRequest(min_args, max_args int, requestFiller fill) func([]string) string {
	return func(values []string) string {
		if (len(values) <= max_args || max_args < 0) && len(values) >= min_args {
			return requestFiller(values)
		}
		return ""
	}
}

func fillSelect(values []string) (request string) {
	switch len(values) {
	case 3:
		request = request + " WHERE " + values[2]
		fallthrough
	case 2:
		request = "SELECT " + values[0] + " FROM " + values[1] + request + ";"
	}
	return
}

func fillCreate_Table(values []string) (request string) {
	switch size := len(values); {
	case size > 3:
		for i := 2; i < size - 1; i++ {
			request += ", " + values[i]
		}
		fallthrough
	case size == 3:
		request = "CREATE TABLE IF NOT EXISTS " + values[0] + "(" + values[1] + request + ") ENGINE=" + values[size-1] + ";"
	}
	return
}

func fillInsert_Into(values []string) (request string) {
	switch size := len(values); {
	case size > 3:
		for i := 3; i < size; i++ {
			request += ", (" + values[i] + ")"
		}
		fallthrough
	case size == 3:
		request = "INSERT INTO " + values[0] + "(" + values[1] + ") VALUES(" + values[2] + ")" + request + ";"
	}
	return
}
