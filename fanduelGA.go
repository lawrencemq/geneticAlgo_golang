package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"lawrencemq/geneticalgo_golang/ga"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
)

const salaryCap = 60000

const filename = "FanDuel-NFL-2016-11-17-16937-players-list.csv"

var qbs, rbs, wrs, tes, ds, ks []Player

type Player struct {
	position  string
	firstName string
	lastName  string
	fppg      float64
	salary    int
	game      string
	team      string
	injured   bool
}

func makePlayer(position, firstName, lastName string, ffpg float64, salary int, game, team string, injured bool) Player {
	return Player{
		position,
		firstName,
		lastName,
		ffpg,
		salary,
		game,
		team,
		injured,
	}
}

type Lineup struct {
	qb, rb1, rb2, wr1, wr2, wr3, te, def, k Player
}

func lineupToSlice(lineup Lineup) []Player {
	return []Player{lineup.qb, lineup.rb1, lineup.rb2, lineup.wr1, lineup.wr2, lineup.wr3, lineup.te, lineup.def, lineup.k}
}
func sliceToLineup(players []Player) Lineup {
	return Lineup{
		players[0],
		players[1],
		players[2],
		players[3],
		players[4],
		players[5],
		players[6],
		players[7],
		players[8],
	}
}

func makeFanduelEntry(qb, rb1, rb2, wr1, wr2, wr3, te, def, k Player) Lineup {
	return Lineup{
		qb,
		rb1,
		rb2,
		wr1,
		wr2,
		wr3,
		te,
		def,
		k,
	}
}

func generateFanduelEntry() Lineup {
	var possibleEntry Lineup
	for {
		possibleEntry = makeFanduelEntry(
			qbs[rand.Intn(len(qbs))],
			rbs[rand.Intn(len(rbs))],
			rbs[rand.Intn(len(rbs))],
			wrs[rand.Intn(len(wrs))],
			wrs[rand.Intn(len(wrs))],
			wrs[rand.Intn(len(wrs))],
			tes[rand.Intn(len(tes))],
			ds[rand.Intn(len(ds))],
			ks[rand.Intn(len(ks))],
		)
		if isValidFanduelEntry(possibleEntry) {
			break
		}
	}
	return possibleEntry
}

func getProjectedPointsForLineup(entry Lineup) float64 {
	var total float64
	for _, player := range lineupToSlice(entry) {
		total += player.fppg
	}
	return total
}

func entryHasValidNumPlayers(entry Lineup) bool {
	playersSet := mapset.NewSet()
	for _, player := range lineupToSlice(entry) {
		playersSet.Add(player)
	}
	return len(playersSet.ToSlice()) == 9
}

func countTimeTeamSeen(entry Lineup, team string) int {
	var count int
	for _, player := range lineupToSlice(entry) {
		if player.team == team {
			count++
		}
	}
	return count
}

func findMaxSameTeam(entry Lineup) int {
	allTeamsSet := mapset.NewSet()
	for _, player := range lineupToSlice(entry) {
		allTeamsSet.Add(player.team)
	}

	var max int
	for team := range allTeamsSet.Iterator().C {
		timesSeen := countTimeTeamSeen(entry, team.(string))
		if timesSeen > max {
			max = timesSeen
		}
	}

	return max

}

func sumSalaryNeededForEntry(entry Lineup) int {
	var totalSalary int
	for _, player := range lineupToSlice(entry) {
		totalSalary += player.salary
	}
	return totalSalary
}

func isValidFanduelEntry(entry Lineup) bool {
	return entryHasValidNumPlayers(entry) &&
		findMaxSameTeam(entry) < 3 &&
		sumSalaryNeededForEntry(entry) < salaryCap
}

func readInData() ([]Player, []Player, []Player, []Player, []Player, []Player) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	qbs := make([]Player, 0)
	rbs := make([]Player, 0)
	wrs := make([]Player, 0)
	tes := make([]Player, 0)
	ds := make([]Player, 0)
	ks := make([]Player, 0)

	firstEntry := true
	r := csv.NewReader(strings.NewReader(string(dat)))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if firstEntry {
			firstEntry = false
			continue
		}

		ffpg, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			panic(err)
		}
		salary, err := strconv.Atoi(record[7])
		if err != nil {
			panic(err)
		}

		position := record[1]
		firstName := record[2]
		lastName := record[4]
		game := record[8]
		team := record[9]
		injured := record[11] != ""
		newEntry := makePlayer(position, firstName, lastName, ffpg, salary, game, team, injured)

	Position:
		switch newEntry.position {
		case "QB":
			qbs = append(qbs, newEntry)
			break Position
		case "RB":
			rbs = append(rbs, newEntry)
			break Position
		case "WR":
			wrs = append(wrs, newEntry)
			break Position
		case "TE":
			tes = append(tes, newEntry)
			break Position
		case "D":
			ds = append(ds, newEntry)
			break Position
		case "K":
			ks = append(ks, newEntry)
			break Position
		default:
			panic("Can't figure out position!")
		}
	}
	return qbs, rbs, wrs, tes, ds, ks
}

type FanDuelGA struct {
}

func (l FanDuelGA) GenerateInitialPopulation(populationSize int) []Lineup {

	initialPopulation := make([]Lineup, 0, populationSize)
	for i := 0; i < populationSize; i++ {
		initialPopulation = append(initialPopulation, generateFanduelEntry())
	}
	return initialPopulation
}
func (l FanDuelGA) PerformCrossover(result1, result2 Lineup, crossoverRate int) Lineup {
	players1, players2 := lineupToSlice(result1), lineupToSlice(result2)

	crossoverIndex := int(float64(len(players1)) * (float64(crossoverRate) / 100.0))
	newEntrySlice1 := players1[:crossoverIndex]
	newEntrySlice2 := players2[crossoverIndex:]
	newEntryPlayers := append(newEntrySlice1, newEntrySlice2...)

	return sliceToLineup(newEntryPlayers)
}

func getRandomPlayerForPosition(position string) Player {
	switch position {
	case "QB":
		return qbs[rand.Intn(len(qbs))]
	case "RB":
		return rbs[rand.Intn(len(rbs))]
	case "WR":
		return wrs[rand.Intn(len(wrs))]
	case "TE":
		return tes[rand.Intn(len(tes))]
	case "D":
		return ds[rand.Intn(len(ds))]
	case "K":
		return ks[rand.Intn(len(ks))]
	default:
		panic("eek!!!!")
	}
}

func (l FanDuelGA) PerformMutation(result Lineup, mutationRate int) Lineup {

	players := lineupToSlice(result)
	for index, player := range players {
		if 100-rand.Intn(100) >= mutationRate {
			continue
		}

		for {
			finalPlayers := append(players[:index], getRandomPlayerForPosition(player.position))
			finalPlayers = append(finalPlayers, players[index+1:]...)
			lineup := sliceToLineup(finalPlayers)
			//fmt.Println(lineup)
			if !isValidFanduelEntry(lineup) {
				break
			}
			return lineup
		}
	}
	return sliceToLineup(players) // no mutation happened
}
func (l FanDuelGA) Sort(population []Lineup) {
	sort.Slice(population, func(i, j int) bool {
		return getProjectedPointsForLineup(population[i]) < getProjectedPointsForLineup(population[j])
	})
}

func fanduelMain() {
	rand.Seed(time.Now().Unix())
	qbs, rbs, wrs, tes, ds, ks = readInData()

	settings := ga.GeneticAlgorithmSettings{
		PopulationSize:           100,
		MutationRate:             10,
		CrossoverRate:            50,
		NumGenerations:           25,
		KeepBestAcrossPopulation: true,
	}

	best, err := ga.Run[Lineup](FanDuelGA{}, settings)
	if err != nil {
		println(err)
	} else {
		fmt.Printf("Best: %f:, $%d\n", getProjectedPointsForLineup(best), sumSalaryNeededForEntry(best))
		for _, player := range lineupToSlice(best) {
			fmt.Printf("%s: %s %s - %f\n", player.position, player.firstName, player.lastName, player.fppg)
		}
	}
}
