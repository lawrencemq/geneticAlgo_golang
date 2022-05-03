package main

import (
	"fmt"
	"lawrencemq/geneticalgo_golang/ga"
	"math"
	"sort"
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
	return math.Pow(entry.x, 2) - 6*entry.x + math.Pow(entry.y, 2) - 6*entry.y + 2
}

type Quadratic3dGA struct {
}

func (l Quadratic3dGA) GenerateInitialPopulation(populationSize int) []Quad3D {

	initialPopulation := make([]Quad3D, 0, populationSize)
	for i := 0; i < populationSize; i++ {
		initialPopulation = append(initialPopulation, makeNewQuadEntry(makeNewEntry(), makeNewEntry()))
	}

	return initialPopulation
}
func (l Quadratic3dGA) PerformCrossover(result1, result2 Quad3D, _ int) Quad3D {
	return makeNewQuadEntry(
		(result1.x+result2.x)/2,
		(result1.y+result2.y)/2,
	)
}
func (l Quadratic3dGA) PerformMutation(_ Quad3D, _ int) Quad3D {
	return makeNewQuadEntry(makeNewEntry(), makeNewEntry())
}
func (l Quadratic3dGA) Sort(population []Quad3D) {
	sort.Slice(population, func(i, j int) bool {
		return calculate3D(population[i]) > calculate3D(population[j])
	})
}

func quadratic3dMain() {
	settings := ga.GeneticAlgorithmSettings{
		PopulationSize:           100,
		MutationRate:             10,
		CrossoverRate:            100,
		NumGenerations:           20,
		KeepBestAcrossPopulation: true,
	}

	best, err := ga.Run[Quad3D](Quadratic3dGA{}, settings)

	if err != nil {
		println(err)
	} else {
		fmt.Printf("Best: x: %f  y: %f  z: %f\n", best.x, best.y, calculate3D(best))
	}
}
