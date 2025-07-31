package main

import (
	"fmt"
	"iter"
	"maps"
)

type Developer struct {
	CoffeeLevel int
	BugCount    int
}

func main() {
	devTeam := map[string]Developer{
		"Alice":   {CoffeeLevel: 8, BugCount: 2},
		"Bob":     {CoffeeLevel: 3, BugCount: 5},
		"Charlie": {CoffeeLevel: 6, BugCount: 1},
	}

	fmt.Println("\n1. Iteration Examples")
	fmt.Println("-------------------------")

	for key, value := range devTeam {
		fmt.Printf("Developer %s has %d coffee level and %d bugs\n",
			key, value.CoffeeLevel, value.BugCount)
	}

	maxDevName, maxDev, found := maxCoffeeDeveloper(maps.All(devTeam))
	fmt.Println(maxDevName, maxDev, found) // Alice {8 2} true

	maxCoffee := maxCoffeeLevel(maps.Values(devTeam))
	fmt.Println(maxCoffee) // 8

	keys := prettyPrint(maps.Keys(devTeam))
	fmt.Println(keys) // Bob, Charlie, Alice

	fmt.Println("\n2. Collecting Examples")
	fmt.Println("-------------------------")

	newDevs := map[string]Developer{
		"Alice": {CoffeeLevel: 9, BugCount: 0},
		"Eve":   {CoffeeLevel: 7, BugCount: 3},
	}

	newTeam := maps.Clone(devTeam)
	maps.Insert(newTeam, maps.All(newDevs))
	fmt.Println(newTeam) // map[Alice:{9 0} Bob:{3 5} Charlie:{6 1} Eve:{7 3}]

	highPerformers := maps.Collect(maps.All(devTeam))
	fmt.Println(highPerformers) // map[Alice:{8 2} Bob:{3 5} Charlie:{6 1}]

	fmt.Println("\n3. Manipulation Examples")
	fmt.Println("-------------------------")

	backupTeam := maps.Clone(devTeam)
	fmt.Println(backupTeam) // map[Alice:{8 2} Bob:{3 5} Charlie:{6 1} Diana:{9 0} Eve:{7 3}]

	stagingTeam := make(map[string]Developer)
	stagingTeam["Frank"] = Developer{CoffeeLevel: 10, BugCount: 4}
	stagingTeam["Eve"] = Developer{CoffeeLevel: 3, BugCount: 9}
	maps.Copy(stagingTeam, devTeam)
	fmt.Println(stagingTeam) // map[Alice:{8 2} Bob:{3 5} Charlie:{6 1} Diana:{9 0} Eve:{7 3} Frank:{10 4}]

	maps.DeleteFunc(devTeam, func(name string, dev Developer) bool {
		return dev.CoffeeLevel < 5
	})
	fmt.Println(devTeam) // map[Alice:{8 2} Charlie:{6 1}]

	fmt.Println("\n5. Comparison Examples")
	fmt.Println("-------------------------")

	areTeamsEqual := maps.Equal(devTeam, backupTeam)
	fmt.Println(areTeamsEqual) // false

	sameCoffeeLevel := maps.EqualFunc(stagingTeam, backupTeam,
		func(dev1, dev2 Developer) bool {
			return dev1.CoffeeLevel == dev2.CoffeeLevel
		})
	fmt.Println(sameCoffeeLevel) // false
}

func maxCoffeeDeveloper(seq iter.Seq2[string, Developer]) (name string, dev Developer, found bool) {
	for k, v := range seq {
		if !found || v.CoffeeLevel > dev.CoffeeLevel {
			name, dev, found = k, v, true
		}
	}
	return
}

func maxCoffeeLevel(seq iter.Seq[Developer]) int {
	maxLevel := 0
	for v := range seq {
		if v.CoffeeLevel > maxLevel {
			maxLevel = v.CoffeeLevel
		}
	}
	return maxLevel
}

func prettyPrint(entries iter.Seq[string]) string {
	result := ""
	for entry := range entries {
		result += entry + ", "
	}
	if len(result) > 0 {
		result = result[:len(result)-2] // remove trailing comma and space
	}
	return result
}
