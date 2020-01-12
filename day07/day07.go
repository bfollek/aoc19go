package day07

import (
	"fmt"

	"github.com/gitchander/permutation"
)

// Part1 tries every combination of phase settings on the amplifiers.
// What is the highest signal that can be sent to the thrusters?
func Part1(fileName string) int {
	//program := intcode.LoadFromFile(fileName)
	combos := phaseSettings([]int{0, 1, 2, 3, 4})
	for _, combo := range combos {
		// func Part1(fileName string) int {
		// 	_, output := intcode.RunFromFile(fileName, []int{1})
		// 	fmt.Printf("output == %v\n", output)
		// 	return output[len(output)-1]
		// }
		fmt.Printf("%T %v\n", combo, combo)
	}
	return 0
}

func phaseSettings(sl []int) [][]int {
	combos := [][]int{}
	p := permutation.New(permutation.IntSlice(sl))
	for p.Next() {
		combos = append(combos, sl)
	}
	return combos
}
