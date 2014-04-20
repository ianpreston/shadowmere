package puzzle

import (
	"strings"
)


func parseLine(line string) (string, string, []string) {
	var command string
	var origin string
	var args []string

	line = strings.Trim(line, " \r\n")
	tokens := splitLine(line)

	if line[0] == ':' {
		// TODO - Check bounds. A malformed message could cause a panic
		// here.
		origin = strings.TrimPrefix(tokens[0], ":")
		command = tokens[1]
		args = tokens[2:]
	} else {
		command = strings.ToUpper(tokens[0])
		if len(tokens) > 1 {
			args = tokens[1:]
		}
	}
	
	return command, origin, args
}

func splitLine(line string) []string {
	tokens := strings.Split(line, " ")

	rightTokenIdx := -1
	for idx, t := range tokens {
		if strings.HasPrefix(t, ":") && idx != 0 {
			rightTokenIdx = idx
			break
		}
	}

	if rightTokenIdx != -1 {
		leftTokens := tokens[:rightTokenIdx]
		rightToken := strings.Join(tokens[rightTokenIdx:], " ")
		rightToken = strings.TrimPrefix(rightToken, ":")

		tokens = leftTokens
		tokens = append(tokens, rightToken)
	}

	return tokens
}