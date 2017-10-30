package ga

import (
	"math/rand"
)

type GeneticAlgorithmSettings struct {

	PopulationSize           int
	MutationRate             int
	CrossoverRate            int
	NumGenerations           int
	KeepBestAcrossPopulation bool

}

type GeneticAlgorithmRunner interface {

	GenerateInitialPopulation(populationSize int) []interface{}
	PerformCrossover(individual1, individual2 interface{}, crossoverRate int) interface{}
	PerformMutation(individual interface{}, mutationRate int) interface{}
	Sort([]interface{})

}


func createStochasticProbableListOfIndividuals(population []interface{}) []interface{} {

	totalCount, populationLength:= 0, len(population)
	for j:= 0; j < populationLength; j++ {
		totalCount += j
	}

	probableIndividuals := make([]interface{}, 0, totalCount)
	for index, individual := range population {
		for i:= 0; i < index; i++{
			probableIndividuals = append(probableIndividuals, individual)
		}
	}

	return probableIndividuals
}


func Run(geneticAlgoRunner GeneticAlgorithmRunner, settings GeneticAlgorithmSettings) (interface{}, error){

	population := geneticAlgoRunner.GenerateInitialPopulation(settings.PopulationSize)

	geneticAlgoRunner.Sort(population)

	bestSoFar := population[len(population) - 1]

	for i:= 0; i < settings.NumGenerations; i++ {

		newPopulation := make([]interface{}, 0, settings.PopulationSize)

		if settings.KeepBestAcrossPopulation {
			newPopulation = append(newPopulation, bestSoFar)
		}

		// perform crossovers with random selection
		probabilisticListOfPerformers := createStochasticProbableListOfIndividuals(population)

		newPopIndex := 0
		if settings.KeepBestAcrossPopulation{
			newPopIndex = 1
		}
		for ; newPopIndex < settings.PopulationSize; newPopIndex++ {
			indexSelection1 := rand.Int() % len(probabilisticListOfPerformers)
			indexSelection2 := rand.Int() % len(probabilisticListOfPerformers)

			// crossover
			newIndividual := geneticAlgoRunner.PerformCrossover(
				probabilisticListOfPerformers[indexSelection1],
				probabilisticListOfPerformers[indexSelection2], settings.CrossoverRate)

			// mutate
			if rand.Intn(101) < settings.MutationRate {
				newIndividual = geneticAlgoRunner.PerformMutation(newIndividual, settings.MutationRate)
			}

			newPopulation = append(newPopulation, newIndividual)
		}

		population = newPopulation

		// sort by performance
		geneticAlgoRunner.Sort(population)

		// keep the best so far
		bestSoFar = population[len(population) - 1]

	}

	return bestSoFar, nil

}