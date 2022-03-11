package common

func VisitOne(args []string) (string, []string) {
	if len(args) > 0 {
		return args[0], args[1:]
	}
	return "", []string{}
}
