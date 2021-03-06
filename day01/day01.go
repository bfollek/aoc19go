package day01

import "github.com/bfollek/aoc19go/util"

type fuelFunc func(int) int

// Part1 calculates the fuel requirements for all of the modules on my spacecraft.
func Part1(fileName string) int {
	return calcFuel(fileName, fuelRequired)
}

// Part2 calculates the fuel requirements for all of the modules on my spacecraft,
// plus the fuel costs of the fuel.
func Part2(fileName string) int {
	return calcFuel(fileName, fuelRequiredMeta)
}

func calcFuel(fileName string, f fuelFunc) int {
	sum := 0
	masses := util.MustLoadIntSlice(fileName, "\n")
	for _, mass := range masses {
		sum += f(mass)
	}
	return sum
}

// "...to find the fuel required for a module, take its mass,
// divide by three, round down, and subtract 2."
func fuelRequired(mass int) int {
	return mass/3 - 2
}

// "So, for each module mass, calculate its fuel and add it to the total. Then,
// treat the fuel amount you just calculated as the input mass and repeat the
// process, continuing until a fuel requirement is zero or negative."
func fuelRequiredMeta(mass int) int {
	total := 0
	for f := fuelRequired(mass); f >= 1; f = fuelRequired(f) {
		total += f
	}
	return total
}
