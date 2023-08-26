package target

import "strings"

type Name string

func (b Name) Key() string {
	if strings.Contains(string(b), "=") {
		parts := strings.SplitN(string(b), "=", 2)
		return parts[0]
	}
	return string(b)
}

func (b Name) Name() string {
	if strings.Contains(string(b), "=") {
		parts := strings.SplitN(string(b), "=", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return string(b)
}
