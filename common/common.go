package common

func Must(args ...interface{}) {
	if len(args) > 0 {
		err, ok := args[len(args)-1].(error)
		if ok && err != nil {
			panic(err)
		}
	}
}
