package calculator

import (
	"github.com/InsafMin/web_calculator/pkg/errors"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	expression = strings.Replace(expression, " ", "", -1)
	if strings.Count(expression, "(") > strings.Count(expression, ")") {
		return 0, errors.ErrExtraOpenBracket
	}
	if strings.Count(expression, "(") < strings.Count(expression, ")") {
		return 0, errors.ErrExtraCloseBracket
	}

	tokens, err := Tokenize(expression)
	if err != nil {
		return 0, err
	}

	result, err := Evaluate(tokens)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func Tokenize(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	var tokens []string
	var number string

	for _, r := range expression {
		if unicode.IsDigit(r) || r == '.' {
			number += string(r)
		} else if IsOperator(r) || r == '(' || r == ')' {
			if number != "" {
				tokens = append(tokens, number)
				number = ""
			}
			tokens = append(tokens, string(r))
		} else {
			return nil, errors.ErrUnacceptableSymbol
		}
	}

	if number != "" {
		tokens = append(tokens, number)
	}

	return tokens, nil
}

func IsOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/'
}

func Resolve(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, errors.ErrOperatorNotSupported
	}
}

// https://www.youtube.com/watch?v=Vk-tGND2bfc
func Evaluate(tokens []string) (float64, error) {
	var numbers []float64
	var operators []string

	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			numbers = append(numbers, num)
		} else if IsOperator(rune(token[0])) {
			for len(operators) > 0 && Priority(token) <= Priority(operators[len(operators)-1]) {
				if len(numbers) < 2 {
					return 0, errors.ErrExtraOperator
				}

				num2 := numbers[len(numbers)-1]
				numbers = numbers[:len(numbers)-1]

				num1 := numbers[len(numbers)-1]
				numbers = numbers[:len(numbers)-1]

				operator := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				res, err := Resolve(num1, num2, operator)
				if err != nil {
					return 0, err
				}

				numbers = append(numbers, res)
			}
			operators = append(operators, token)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				num2 := numbers[len(numbers)-1]
				numbers = numbers[:len(numbers)-1]

				num1 := numbers[len(numbers)-1]
				numbers = numbers[:len(numbers)-1]

				operator := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				res, err := Resolve(num1, num2, operator)
				if err != nil {
					return 0, err
				}

				numbers = append(numbers, res)
			}
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			}
		}
	}

	for len(operators) > 0 {
		if len(numbers) < 2 {
			return 0, errors.ErrExtraOperator
		}

		num2 := numbers[len(numbers)-1]
		numbers = numbers[:len(numbers)-1]

		num1 := numbers[len(numbers)-1]
		numbers = numbers[:len(numbers)-1]

		operator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		res, err := Resolve(num1, num2, operator)
		if err != nil {
			return 0, err
		}

		numbers = append(numbers, res)
	}

	if len(numbers) != 1 {
		return 0, errors.ErrInvalidExpression
	}

	return numbers[0], nil
}

func Priority(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}
