package db

func Router(postgres bool, newrow *Baserow, method string) error {
	if postgres {
		switch method {
		case "GET":
			return QueryGetURL(newrow)
		case "POST":
			return QueryPOSTorSelect(newrow)
		}
	} else {
		switch method {
		case "GET":
			return GetLocal(newrow)
		case "POST":
			return PostLocal(newrow)
		}
	}
	return nil
}
