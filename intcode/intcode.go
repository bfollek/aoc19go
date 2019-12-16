package intcode

import (
	"log"

	"github.com/bfollek/advent2019go/util"
	"github.com/golang-collections/collections/stack"
)

// NoInput is a convenience for clients.
var NoInput = []int{}

type computer struct {
	memory         []int
	iP             int          // Instruction pointer
	input          *stack.Stack // Input buffer
	output         []int        // Output buffer
	parameterModes *stack.Stack
}

// ------------------------------------------------------------------
// Modes
// ------------------------------------------------------------------

const (
	positionMode = iota
	immediateMode
)

// ------------------------------------------------------------------
// Opcodes
// ------------------------------------------------------------------

type opCodeAttribs struct {
	numParams int
	exec      func(opCodeAttribs, *computer)
}

const (
	opAdd = iota + 1
	opMultiply
	opInput
	opOutput
	opJumpIfTrue
	opJumpIfFalse
	opLessThan
	opEquals
	opHalt = 99
)

var opCodes = map[int]opCodeAttribs{
	opAdd:         {3, add},
	opMultiply:    {3, multiply},
	opInput:       {1, in},
	opOutput:      {1, out},
	opJumpIfTrue:  {2, jumpIfTrue},
	opJumpIfFalse: {2, jumpIfFalse},
	opLessThan:    {3, lessThan},
	opEquals:      {3, equals},
}

// RunFromFile reads an intcode program from a file, then executes it.
func RunFromFile(fileName string, input []int) ([]int, []int) {
	ss := util.MustLoadStringSlice(fileName, ",")
	program := []int{}
	for _, s := range ss {
		i := util.MustAtoi(s)
		program = append(program, i)
	}
	return Run(program, input)
}

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
	for {
		// Parameter modes are stored in the same value as the instruction's opcode.
		// The opcode is a two-digit number based only on the ones and tens digit
		// of the value, that is, the opcode is the rightmost two digits of the
		// first value in an instruction.
		rawOpCode := vm.memory[vm.iP]
		setParameterModes(rawOpCode, vm)
		opCode := rawOpCode % 100
		// Opcode 99 means that the program is finished and should immediately halt.
		// The instruction 99 contains only an opcode and has no parameters.
		if opCode == opHalt {
			return vm.memory, vm.output
		}
		oca, ok := opCodes[opCode]
		if !ok {
			// Encountering an unknown opcode means something went wrong.
			log.Fatalf("Unexpected op code: %d", opCode)
		}
		oca.exec(oca, vm)
	}
}

// Parameter modes are single digits, one per parameter, read right-to-left
// from the opcode: the first parameter's mode is in the hundreds digit,
// the second parameter's mode is in the thousands digit, the third parameter's
// mode is in the ten-thousands digit, and so on. Any missing modes are 0.
func setParameterModes(remaining int, vm *computer) {
	// No harm setting more modes than the actual number of parameters.
	// Anything extra just won't get popped.
	var mode int
	vm.parameterModes = stack.New()
	for _, n := range []int{10000, 1000, 100} {
		if remaining > n {
			remaining -= n
			mode = 1
		} else {
			mode = 0
		}
		vm.parameterModes.Push(mode)
	}
}

// add (Opcode 1) - adds together numbers read from two positions and stores
// the result in a third position. The three integers immediately after the
// opcode tell you these three positions - the first two indicate the positions
// from which you should read the input values, and the third indicates the
// position at which the output should be stored.
func add(oca opCodeAttribs, vm *computer) {
	op1, op2 := next2Params(vm)
	store(op1+op2, vm.memory[vm.iP+3], vm)
	advanceInstructionPointer(oca.numParams+1, vm)
}

// multiply (Opcode 2) - works exactly like opcode 1, except it multiplies
// the two inputs instead of adding them. Again, the three integers after
// the opcode indicate where the inputs and outputs are, not their values.
func multiply(oca opCodeAttribs, vm *computer) {
	op1, op2 := next2Params(vm)
	store(op1*op2, vm.memory[vm.iP+3], vm)
	advanceInstructionPointer(oca.numParams+1, vm)
}

// in (Opcode 3) - takes a single integer as input and saves it to the position given
// by its only parameter. For example, the instruction 3,50 would take an input
// value and store it at address 50.
func in(oca opCodeAttribs, vm *computer) {
	i := vm.input.Pop()
	store(i.(int), vm.memory[vm.iP+1], vm)
	advanceInstructionPointer(oca.numParams+1, vm)
}

// out (Opcode 4) - outputs the value of its only parameter. For example,
// the instruction 4,50 would output the value at address 50.
func out(oca opCodeAttribs, vm *computer) {
	i := fetch(vm.iP+1, vm)
	vm.output = append(vm.output, i)
	advanceInstructionPointer(oca.numParams+1, vm)
}

// jumpIfTrue (opCode 5) - if the first parameter is non-zero, it sets the instruction
// pointer to the value from the second parameter. Otherwise, it does nothing.
func jumpIfTrue(oca opCodeAttribs, vm *computer) {
	p1, p2 := next2Params(vm)
	jump(p1 != 0, p2, oca, vm)
}

// jumpIfFalse (Opcode 6) - if the first parameter is zero, it sets the instruction
// pointer to the value from the second parameter. Otherwise, it does nothing.
func jumpIfFalse(oca opCodeAttribs, vm *computer) {
	p1, p2 := next2Params(vm)
	jump(p1 == 0, p2, oca, vm)
}

func jump(jump bool, jumpTo int, oca opCodeAttribs, vm *computer) {
	if jump {
		setInstructionPointer(jumpTo, vm)
	} else {
		advanceInstructionPointer(oca.numParams+1, vm)
	}
}

// lessThan (Opcode 7) - if the first parameter is less than the second parameter, it
// stores 1 in the position given by the third parameter. Otherwise, it stores 0.
func lessThan(oca opCodeAttribs, vm *computer) {
	op1, op2 := next2Params(vm)
	comparison(op1 < op2, oca, vm)
}

// equals (Opcode 8) - if the first parameter is equal to the second parameter, it
// stores 1 in the position given by the third parameter. Otherwise, it stores 0.
func equals(oca opCodeAttribs, vm *computer) {
	op1, op2 := next2Params(vm)
	comparison(op1 == op2, oca, vm)
}

func comparison(satisfied bool, oca opCodeAttribs, vm *computer) {
	var i int
	if satisfied {
		i = 1
	} else {
		i = 0
	}
	store(i, vm.memory[vm.iP+3], vm)
	advanceInstructionPointer(oca.numParams+1, vm)
}

func next2Params(vm *computer) (int, int) {
	return fetch(vm.iP+1, vm), fetch(vm.iP+2, vm)
}

// fetch returns the value at the address, taking into consideration
// the parameter mode.
//
// Immediate mode interprets the `address` param as the address of the
// value to return. If the `address` param is 50, we return the value
// stored at address 50.
//
// Position mode adds a level of indirection. It interprets the `address`
// param as the address of an address. The second address is the address of
// the value to return. If the `address` param is 50, we get the value stored
// at address 50. Suppose that value is 100. We then get the value stored at
// address 100, and return it.
func fetch(address int, vm *computer) int {
	var i int
	switch mode := vm.parameterModes.Pop(); mode {
	case immediateMode:
		i = vm.memory[address]
	case positionMode:
		i = vm.memory[address]
		i = vm.memory[i]
	default:
		log.Fatalf("Unexpected parameter mode: %d", mode)
	}
	return i
}

// "Parameters that an instruction writes to will never be in immediate mode."
func store(value int, address int, vm *computer) {
	// Make sure there's room
	if address >= len(vm.memory) {
		more := address + 1 - len(vm.memory)
		vm.memory = append(vm.memory, make([]int, more)...)
	}
	vm.memory[address] = value
}

func advanceInstructionPointer(i int, vm *computer) {
	vm.iP += i
}

func setInstructionPointer(i int, vm *computer) {
	vm.iP = i
}

// load creates the vm and loads the program into it.
func load(program []int, input []int) *computer {
	vm := new(computer)
	vm.memory = make([]int, len(program))
	// Copy so that we don't overwrite the program, in case we it run again.
	copy(vm.memory, program)
	vm.iP = 0
	vm.input = stack.New()
	// Start at end so we can use a stack and simply pop as needed
	for i := len(input) - 1; i >= 0; i-- {
		vm.input.Push(input[i])
	}
	vm.output = []int{}
	return vm
}
