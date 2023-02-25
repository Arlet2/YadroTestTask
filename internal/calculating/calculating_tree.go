package calculating

import (
	"errors"
	"strconv"
	"test_task/internal/format"
	"test_task/internal/operations"
	"test_task/internal/parsing"
)

type TreeCreatingError error
type FormulaCycleError error
type calculatingTree struct {
	nodes map[string][]string
}

func CreateTree(csv format.Csv) (calculatingTree, error) {
	nodes := make(map[string][]string, 0)
	for index, line := range csv.Data {
		for jndex, cell := range line {
			if parsing.IsFormula(cell) {
				link, err := csv.GetLinkByIndexes(jndex, index)
				if err != nil {
					panic(err)
				}
				formula := parsing.ParseFormula(cell)

				// если формулу можно посчитать сразу, то вычисляем на месте
				if !formula.FirstOperand.IsLink() && !formula.SecondOperand.IsLink() {
					value, err := formula.Action(formula.FirstOperand.GetConstant(),
						formula.SecondOperand.GetConstant())

					if err != nil {
						return calculatingTree{}, err.(operations.CalculatingError)
					}
					csv.Data[index][jndex] = strconv.FormatInt(int64(value), 10)
					continue
				}

				if _, ok := nodes[link]; !ok {
					nodes[link] = make([]string, 0)
				}

				if formula.FirstOperand.IsLink() {
					if !csv.IsLinkExist(formula.FirstOperand.GetLink()) {
						return calculatingTree{},
							errors.New("ячейки " + formula.FirstOperand.GetLink() + " не существует").(TreeCreatingError)
					}

					nodes[link] = append(nodes[link], formula.FirstOperand.GetLink())
				}

				if formula.SecondOperand.IsLink() {
					if !csv.IsLinkExist(formula.SecondOperand.GetLink()) {
						return calculatingTree{},
							errors.New("ячейки " + formula.SecondOperand.GetLink() + " не существует").(TreeCreatingError)
					}

					nodes[link] = append(nodes[link], formula.SecondOperand.GetLink())
				}
			}
		}
	}

	return calculatingTree{nodes: nodes}, nil
}

// (!) сортировка не детерминирована из-за недетерминированности ключей в tree.nodes
func (tree calculatingTree) SortTree() ([]string, error){
	
	nodesState := make(map[string]int, 0)
	sortedNodes := make([]string, 0)

	for key := range tree.nodes {
		nodesState[key] = 0
	}

	for key := range tree.nodes {
		err := tree.dfc(key, &nodesState, &sortedNodes)
		if err != nil {
			return nil, err
		}
	}

	// reverse
	for i, j := 0, len(sortedNodes)-1; i < j; i, j = i+1, j-1 {
		sortedNodes[i], sortedNodes[j] = sortedNodes[j], sortedNodes[i]
	}
	
	return sortedNodes, nil
}

func (tree calculatingTree) dfc(currentNode string, nodesState *map[string]int, sortedNodes *[]string) (error) {

	if (*nodesState)[currentNode] == 1 {
		return errors.New("обнаружена циклическая зависимость у формул").(FormulaCycleError)
	}
	if (*nodesState)[currentNode] == 2 {
		return nil
	}

	(*nodesState)[currentNode] = 1

	for _, value := range tree.nodes[currentNode] {
		err := tree.dfc(value, nodesState, sortedNodes)
		if err != nil {
			return err
		}
	}

	(*nodesState)[currentNode] = 2
	*sortedNodes = append(*sortedNodes, currentNode)

	return nil
}

func (tree calculatingTree) Calculate(csv format.Csv) (error) {
	return nil
}