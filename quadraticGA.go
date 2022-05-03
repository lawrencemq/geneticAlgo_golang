package main

import (
	"fmt"
	"lawrencemq/geneticalgo_golang/ga"
	"math"
	"math/rand"
	"sort"
)

const highRange = 100.0

func makeNewEntry() float64 {
	return highRange * rand.Float64()
}

func calculate(x float64) float64 {
	return math.Pow(x, 2) - 6*x + 2 // minimum should be at x=3
}

type QuadraticGA struct {
}

func (l QuadraticGA) GenerateInitialPopulation(populationSize int) []float64 {

	initialPopulation := make([]float64, 0, populationSize)
	for i := 0; i < populationSize; i++ {
		initialPopulation = append(initialPopulation, makeNewEntry())
	}

	return initialPopulation
}
func (l QuadraticGA) PerformCrossover(result1, result2 float64, _ int) float64 {
	return (result1 + result2) / 2
}
func (l QuadraticGA) PerformMutation(_ float64, _ int) float64 {
	return makeNewEntry()
}
func (l QuadraticGA) Sort(population []float64) {
	sort.Slice(population, func(i, j int) bool {
		return calculate(population[i]) > calculate(population[j])
	})
}

func quadraticMain() {
	settings := ga.GeneticAlgorithmSettings{
		PopulationSize:           10,
		MutationRate:             10,
		CrossoverRate:            100,
		NumGenerations:           20,
		KeepBestAcrossPopulation: true,
	}

	best, err := ga.Run[float64](QuadraticGA{}, settings)

	if err != nil {
		println(err)
	} else {
		fmt.Printf("Best: x: %f  y: %f\n", best, calculate(best))
	}

}
