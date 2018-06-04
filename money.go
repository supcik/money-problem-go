package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type verbalSum struct {
	operands []string
	result   string
}

func verbToInt(verb string, mapping map[rune]int) int {
	if mapping[[]rune(verb)[0]] == 0 {
		return -1
	}
	result := 0
	for _, c := range []rune(verb) {
		result = result*10 + mapping[c]
	}
	return result
}

func withoutItem(list []int, i int) []int {
	res := make([]int, 0, len(list)-1)
	res = append(res, list[:i]...)
	res = append(res, list[i+1:]...)
	return res
}

func (v verbalSum) allLetters() []rune {
	letterSet := make(map[rune]bool)
	for _, arg := range v.operands {
		for _, l := range arg {
			letterSet[l] = true
		}
	}
	for _, l := range v.result {
		letterSet[l] = true
	}
	letters := make([]rune, len(letterSet))
	i := 0
	// The order of the letters is not deterministic, but it doesn't matter
	for k := range letterSet {
		letters[i] = k
		i++
	}
	return letters
}

func (v verbalSum) isValidSolution(mapping map[rune]int) bool {
	sum := 0
	for _, arg := range v.operands {
		n := verbToInt(arg, mapping)
		if n < 0 {
			return false
		}
		sum += n
	}
	n := verbToInt(v.result, mapping)
	if n < 0 {
		return false
	}
	return sum == n
}

func (v verbalSum) recursiveSolve(
	letters []rune, digits []int, mapping map[rune]int,
) (combinations int, solutions int) {
	if len(letters) == 0 {
		if v.isValidSolution(mapping) {
			fmt.Printf("Found   %v\n", v.solutionString(mapping))
			return 1, 1 // one solution tried, one solution found
		}
		return 1, 0 // one solution tried, zero solution found
	}
	letter := letters[0] // take the first letter
	combinations, solutions = 0, 0
	for i, digit := range digits { // try all digits for the letter "letter"
		mapping[letter] = digit
		c, s := v.recursiveSolve(letters[1:], withoutItem(digits, i), mapping)
		combinations += c
		solutions += s
	}
	delete(mapping, letter) // backtrack
	return
}

func (v verbalSum) solve() error {
	fmt.Printf("Solving %v\n", v)
	// Build a letters array with all letter from all operands and from the result
	letters := v.allLetters()
	// Build a numbers array with all digits from 0 to 9
	digits := make([]int, 10)
	for i := 0; i < 10; i++ {
		digits[i] = i
	}
	if len(letters) > len(digits) {
		return errors.Errorf("Too many letters")
	}

	mapping := make(map[rune]int)
	count, sols := v.recursiveSolve(letters, digits, mapping)
	fmt.Printf("I tried %v combinations\n", count)
	fmt.Printf("I found %v solution(s)\n", sols)
	return nil
}

func (v verbalSum) String() string {
	return strings.Join(v.operands, " + ") + " = " + v.result
}

func (v verbalSum) solutionString(mapping map[rune]int) string {
	args := make([]string, len(v.operands))
	for i, x := range v.operands {
		n := verbToInt(x, mapping)
		args[i] = strconv.Itoa(n)
	}
	n := verbToInt(v.result, mapping)
	return strings.Join(args, " + ") + " = " + strconv.Itoa(n)
}

func main() {
	challenge := verbalSum{
		operands: []string{"SEND", "MORE"},
		result:   "MONEY",
	}

	err := challenge.solve()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
