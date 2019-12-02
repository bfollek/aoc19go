package day02

import (
	"errors"
	"log"

	"github.com/bfollek/advent2019go/util"
)

const add = 1
const multiply = 2
const halt = 99

const moonLanding = 19690720

// Part1 "What value is left at position 0 after the program halts?"
func Part1(fileName string) int {
	program := loadProgram(fileName)
	// "...before running the program, replace position 1 with the value 12 and replace
	// position 2 with the value 2.
	program[1] = 12
	program[2] = 2
	program = runProgram(program)
	return program[0]
}

// Part2 "...you need to determine what pair of inputs produces the output 19690720."
func Part2(fileName string) (int, error) {
	cleanMemory := loadProgram(fileName)
	for i := 0; i < 99; i++ {
		for j := 0; j < 99; j++ {
			program := resetMemory(cleanMemory)
			program[1] = i
			program[2] = j
			program = runProgram(program)
			if program[0] == moonLanding {
				return 100*i + j, nil
			}
		}
	}
	return -1, errors.New("No solution found")
}

/*
 * "Each time you try a pair of inputs, make sure you first reset the computer's
 * memory to the values in the program (your puzzle input) - in other words, don't
 * reuse memory from a previous attempt."
 *
 * We can't just assign slice to slice, or we'll get two references to the same slice
 * and we'll overwrite our clean memory. So we build a new slice and copy all values.
 */
func resetMemory(cleanMemory []int) []int {
	m := []int{}
	for _, i := range cleanMemory {
		m = append(m, i)
	}
	return m
}

func loadProgram(fileName string) []int {
	ss := util.MustLoadStringSlice(fileName, ",")
	program := []int{}
	for _, s := range ss {
		i := util.MustAtoi(s)
		program = append(program, i)
	}
	return program
}

// "The three integers immediately after the opcode tell you these three positions -
// the first two indicate the positions from which you should read the input values, // // and the third indicates the position at which the output should be stored."
func runProgram(program []int) []int {
	opCodeIndex := 0
	for {
		switch opCode := program[opCodeIndex]; opCode {
		case add, multiply:
			op1 := program[program[opCodeIndex+1]]
			op2 := program[program[opCodeIndex+2]]
			var value int
			if opCode == add {
				value = op1 + op2
			} else {
				value = op1 * op2
			}
			program[program[opCodeIndex+3]] = value
		case halt:
			return program
		default:
			log.Fatalf("Unexpected op code: %d", opCode)
		}
		// "Once you're done processing an opcode, move to the next one
		// by stepping forward 4 positions."
		opCodeIndex += 4
	}
}