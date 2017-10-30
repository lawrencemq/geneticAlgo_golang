package main

import (
	"sort"
	"fmt"
	"math"
	"./ga"
)




type Quad3D struct {
	x, y float64
}
func makeNewQuadEntry(newX, newY float64) Quad3D {
	return Quad3D{
		x: newX,
		y: newY,
	}
}

func calculate3D(entry Quad3D) float64 {
	return math.Pow(entry.x, 2)- 6 * entry.x + math.Pow(entry.y, 2)- 6 * entry.y + 2
}

type Quadratic3dGA struct {

}
func (l Quadratic3dGA) GenerateInitialPopulation(populationSize int) []interface{}{

	initialPopulation := make([]interface{}, 0, populationSize)
	for i:= 0; i < populationSize; i++ {
		initialPopulation = append(initialPopulation, makeNewQuadEntry(makeNewEntry(), makeNewEntry()))
	}

	return initialPopulation
}
func (l Quadratic3dGA) PerformCrossover(result1, result2 interface{}, _ int) interface{}{
	r1Entry, r2Entry := result1.(Quad3D), result2.(Quad3D)
	return makeNewQuadEntry(
		(r1Entry.x + r2Entry.x) / 2,
		(r1Entry.y + r2Entry.y) / 2,
	)
}
func (l Quadratic3dGA) PerformMutation(_ interface{}, _ int) interface{}{
	return makeNewQuadEntry(makeNewEntry(), makeNewEntry())
}
func (l Quadratic3dGA) Sort(population []interface{}){
	sort.Slice(population, func(i, j int) bool {
		return calculate3D(population[i].(Quad3D)) > calculate3D(population[j].(Quad3D))
	})
}



func quadratic3dMain(){
	settings := ga.GeneticAlgorithmSettings{
		PopulationSize: 100,
		MutationRate: 10,
		CrossoverRate: 100,
		NumGenerations: 20,
		KeepBestAcrossPopulation: true,
	}

	best, err := ga.Run(Quadratic3dGA{}, settings)
	entry := best.(Quad3D)

	if err != nil {
		println(err)
	}else{
		fmt.Printf("Best: x: %f  y: %f  z: %f\n", entry.x, entry.y, calculate3D(entry))
	}
}
