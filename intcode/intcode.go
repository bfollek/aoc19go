package intcode

import "log"

// NoInput is a convenience for clients.
var NoInput = []int{}

type computer struct {
	memory []int
	iP     int // Instruction pointer
	input  []int
	inP    int // Input pointer
	output []int
}

// ------------------------------------------------------------------
// Modes
// ------------------------------------------------------------------

// Your ship computer already understands parameter mode 0, position mode, which
// causes the parameter to be interpreted as a position - if the parameter is 50,
// its value is the value stored at address 50 in memory. Until now, all parameters
//  have been in position mode.

// Now, your ship computer will also need to handle parameters in mode 1,
// immediate mode. In immediate mode, a parameter is interpreted as a value -
// if the parameter is 50, its value is simply 50.

const (
	positionMode = iota
	immediateMode
)

// ------------------------------------------------------------------
// Opcodes
// ------------------------------------------------------------------

// Opcode 1 adds together numbers read from two positions and stores the result
// in a third position. The three integers immediately after the opcode tell you
// these three positions - the first two indicate the positions from which you
// should read the input values, and the third indicates the position at which
// the output should be stored.

// Opcode 2 works exactly like opcode 1, except it multiplies the two inputs
// instead of adding them. Again, the three integers after the opcode indicate
// where the inputs and outputs are, not their values.

// Opcode 3 takes a single integer as input and saves it to the position given
// by its only parameter. For example, the instruction 3,50 would take an input
// value and store it at address 50.

// Opcode 4 outputs the value of its only parameter. For example, the instruction
// 4,50 would output the value at address 50.

// Opcode 99 means that the program is finished and should immediately halt.
// The instruction 99 contains only an opcode and has no parameters.

// Encountering an unknown opcode means something went wrong.

const (
	add = iota + 1
	multiply
	input
	output
	halt = 99
)

// opcode => number of params
var opCodeNumParams = map[int]int{add: 3, multiply: 3, input: 1, output: 1, halt: 0}

// Run executes an intcode program.
// The first param, `program`, is the program code.
// The second param, `input`, is any input the program needs.
// The first return value is memory after the program runs.
// The second return value is the program's output.
//
// The program may modify the memory it's in as it runs.
// This means that the program may be self-modifying.
func Run(program []int, input []int) ([]int, []int) {
	vm := load(program, input)
	var opCode int
	for {
		switch opCode = vm.memory[vm.iP]; opCode {
		case add, multiply:
			op1 := vm.memory[vm.memory[vm.iP+1]]
			op2 := vm.memory[vm.memory[vm.iP+2]]
			var value int
			if opCode == add {
				value = op1 + op2
			} else {
				value = op1 * op2
			}
			vm.memory[vm.memory[vm.iP+3]] = value
		case halt:
			return vm.memory, vm.output
		default:
			log.Fatalf("Unexpected op code: %d", opCode)
		}
		vm.iP += (opCodeNumParams[opCode] + 1)
	}
}

// load creates the vm and loads the program into it.
func load(program []int, input []int) *computer {
	vm := new(computer)
	vm.memory = make([]int, len(program))
	copy(vm.memory, program)
	vm.iP = 0
	vm.input = input
	vm.inP = 0
	vm.output = []int{}
	return vm
}