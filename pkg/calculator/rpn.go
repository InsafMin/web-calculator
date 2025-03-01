package calculator

import (
	"github.com/InsafMin/web_calculator/pkg/errors"
	"strconv"
)

func ToRPN(tokens []string) ([]string, error) {
	var output []string
	var operators []string
	skobaLevel := 0

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			output = append(output, token)
		} else if token == "(" {
			skobaLevel += 2
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, errors.ErrInvalidExpression
			}
			operators = operators[:len(operators)-1]
			skobaLevel -= 2
		} else if IsOperator(rune(token[0])) {
			priority := Priority(token) + skobaLevel

			for len(operators) > 0 && Priority(operators[len(operators)-1]) >= priority {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		} else {
			return nil, errors.ErrUnacceptableSymbol
		}
	}

	for len(operators) > 0 {
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}
