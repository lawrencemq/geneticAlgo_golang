package main

import (
	"math"
	"sort"
	"fmt"
	"math/rand"
	"./ga"
)

const highRange = 100.0

func makeNewEntry() float64 {
	return highRange * rand.Float64()
}


func calculate(x float64) float64 {
	return  math.Pow(x, 2) - 6*x + 2 // minimum should be at x=3
}


type QuadraticGA struct {

}
func (l QuadraticGA) GenerateInitialPopulation(populationSize int) []interface{}{

	initialPopulation := make([]interface{}, 0, populationSize)
	for i:= 0; i < populationSize; i++ {
		initialPopulation = append(initialPopulation, makeNewEntry())
	}

	return initialPopulation
}
func (l QuadraticGA) PerformCrossover(result1, result2 interface{}, _ int) interface{}{
	return (result1.(float64) + result2.(float64)) / 2
}
func (l QuadraticGA) PerformMutation(_ interface{}, _ int) interface{}{
	return makeNewEntry()
}
func (l QuadraticGA) Sort(population []interface{}){
	sort.Slice(population, func(i, j int) bool {
		return calculate(population[i].(float64)) > calculate(population[j].(float64))
	})
}


func quadraticMain(){
	settings := ga.GeneticAlgorithmSettings{
		PopulationSize: 10,
		MutationRate: 10,
		CrossoverRate: 100,
		NumGenerations: 20,
		KeepBestAcrossPopulation: true,
	}

	best, err := ga.Run(QuadraticGA{}, settings)

	if err != nil {
		println(err)
	}else{
		fmt.Printf("Best: x: %f  y: %f\n", best, calculate(best.(float64)))
	}

}
