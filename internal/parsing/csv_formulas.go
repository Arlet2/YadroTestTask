package parsing

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"test_task/internal/calculating"
)

type operand interface{
	isLink() bool
	getConstant() int
	getLink() string
}

type link struct {
	value string
}

func (link) isLink() bool {
	return true
}

func (link) getConstant() int {
	return 0
}

func (link link) getLink() string {
	return link.value
}

type constant struct {
	value int
}

func (constant) isLink() bool {
	return false
}

func (constant constant) getConstant() int {
	return constant.value
}

func (constant) getLink() string {
	return ""
}

type formula struct {
	firstOperand operand
	secondOperand operand
	action calculating.Operation
}

type FormulaParseError error

func parseCell(cell string) (formula, error) {
	if cell[0] != '=' {
		return formula{}, errors.New(cell+" is not a formula. Formulas starts from =").(FormulaParseError)
	}
	cell = strings.ReplaceAll(cell, " ", "")
	cell = strings.ReplaceAll(cell, "\t", "")

	regexOperations := ""
	for key := range calculating.AllowedOperations {
		// экранируем для регулярок
		regexOperations += "\\"+key
	}

	regex := regexp.MustCompile(`=[A-z0-9]+[`+regexOperations+`][A-z0-9]+`)

	if regex.FindAllString(cell, 1) == nil {
		regexOperations = strings.ReplaceAll(regexOperations, "\\", "")
		return formula{}, errors.New("incorrect formula. Formula need to be in format =OP1 ["+regexOperations+"] OP2")
	}

	// Дальнейшие регулярки точно найдут совпадения, 
	// так как они являются лишь частями регулярки, которая проверялась до этого

	regex = regexp.MustCompile(`[A-z0-9]+`)

	notParsedOperands := regex.FindAllString(cell, 2)

	notParsedOperand := notParsedOperands[0]
	numericValue, err := strconv.ParseInt(notParsedOperand, 10, 64)

	var firstOperand operand
	// если не удалось спарсить число == это ссылка
	if err != nil {
		firstOperand = link{notParsedOperand}
	} else { // удалось спарсить == это константа
		firstOperand = constant{int(numericValue)}
	}

	notParsedOperand = notParsedOperands[1]
	numericValue, err = strconv.ParseInt(notParsedOperand, 10, 64)

	var secondOperand operand
	// если не удалось спарсить число == это ссылка
	if err != nil {
		secondOperand = link{notParsedOperand}
	} else { // удалось спарсить == это константа
		secondOperand = constant{int(numericValue)}
	}

	regex = regexp.MustCompile(`[+`+regexOperations+`+]`)

	action := calculating.AllowedOperations[regex.FindAllString(cell, 1)[0]]

	return formula{firstOperand: firstOperand, secondOperand: secondOperand, action: action}, nil
}