package main

import (
	"fmt"
	"iter"
	"slices"
)

type Developer struct {
	Name        string
	CoffeeLevel int
	BugCount    int
}

func main() {
	devTeam := []Developer{
		{Name: "Alice", CoffeeLevel: 8, BugCount: 2},
		{Name: "Bob", CoffeeLevel: 3, BugCount: 5},
		{Name: "Charlie", CoffeeLevel: 6, BugCount: 1},
		{Name: "Diana", CoffeeLevel: 9, BugCount: 0},
		{Name: "Eve", CoffeeLevel: 4, BugCount: 3},
	}

	fmt.Println("\n1. Iterator Methods")
	fmt.Println("-------------------")

	pretty := prettyPrint(slices.All(devTeam))
	fmt.Println(pretty) // 0: {Alice 8 2}, 1: {Bob 3 5}, 2: {Charlie 6 1}, 3: {Diana 9 0}, 4: {Eve 4 3}

	maxLevel := maxCoffeeLevel(slices.Values(devTeam))
	fmt.Println(maxLevel) // 9

	pretty = prettyPrint(slices.Backward(devTeam))
	fmt.Println(pretty) // 4: {Eve 4 3}, 3: {Diana 9 0}, 2: {Charlie 6 1}, 1: {Bob 3 5}, 0: {Alice 8 2}

	fmt.Println("\n2. Creating and Populating Slices")
	fmt.Println("----------------------------------")

	collectedDevs := slices.Collect(slices.Values(devTeam))
	fmt.Println(collectedDevs) // [{Alice 8 2} {Bob 3 5} {Charlie 6 1} {Diana 9 0} {Eve 4 3}]

	existingDevelopers := []Developer{
		{Name: "Frank", CoffeeLevel: 7, BugCount: 2},
	}
	appendedDevs := slices.AppendSeq(existingDevelopers, slices.Values(devTeam))
	fmt.Println(appendedDevs) // [{Frank 7 2} {Alice 8 2} {Bob 3 5} {Charlie 6 1} {Diana 9 0} {Eve 4 3}]

	defaultDev := Developer{Name: "NewHire", CoffeeLevel: 5, BugCount: 0}
	newHires := slices.Repeat([]Developer{defaultDev}, 3)
	fmt.Println(newHires) // [{NewHire 5 0} {NewHire 5 0} {NewHire 5 0}]

	fmt.Println("\n3. Copying and Cloning Slices")
	fmt.Println("------------------------------")

	teamCopy := slices.Clone(devTeam)
	fmt.Println(teamCopy) // [{Alice 8 2} {Bob 3 5} {Charlie 6 1} {Diana 9 0} {Eve 4 3}]

	team1 := []Developer{
		{Name: "Frank", CoffeeLevel: 7, BugCount: 2},
	}
	team2 := []Developer{
		{Name: "Grace", CoffeeLevel: 10, BugCount: 0},
	}
	combinedTeam := slices.Concat(team1, team2)
	fmt.Println(combinedTeam) // [{Frank 7 2} {Grace 10 0}]

	fmt.Println("\n4. Searching Methods")
	fmt.Println("--------------------")

	alice := Developer{Name: "Alice", CoffeeLevel: 8, BugCount: 2}
	isAliceInTeam := slices.Contains(devTeam, alice)
	fmt.Println(isAliceInTeam) // true

	hasDevsWithZeroBugs := slices.ContainsFunc(devTeam, func(dev Developer) bool {
		return dev.BugCount == 0
	})
	fmt.Println(hasDevsWithZeroBugs) // true

	alicePos := slices.Index(devTeam, alice)
	fmt.Println(alicePos) // 0

	highCoffeeIdx := slices.IndexFunc(devTeam, func(dev Developer) bool {
		return dev.CoffeeLevel >= 9
	})
	fmt.Println(highCoffeeIdx) // 3

	coffeeLevels := []int{3, 4, 6, 8, 9}
	index, found := slices.BinarySearch(coffeeLevels, 6)
	fmt.Println(index, found) // 2 true

	comparisonFn := func(a, b Developer) int {
		return a.BugCount - b.BugCount
	}
	sortedByBugs := slices.Clone(devTeam)
	slices.SortFunc(sortedByBugs, comparisonFn)
	targetBugs := 2
	bugIndex, bugFound := slices.BinarySearchFunc(sortedByBugs, Developer{BugCount: targetBugs}, comparisonFn)
	fmt.Println(bugIndex, bugFound) // 2 true

	fmt.Println("\n5. Comparing Methods")
	fmt.Println("--------------------")

	anotherTeam := []Developer{
		{Name: "Alice", CoffeeLevel: 8, BugCount: 2},
		{Name: "Bob", CoffeeLevel: 3, BugCount: 5},
		{Name: "Charlie", CoffeeLevel: 6, BugCount: 1},
		{Name: "Diana", CoffeeLevel: 9, BugCount: 0},
		{Name: "Eve", CoffeeLevel: 4, BugCount: 3},
	}
	teamsEqual := slices.Equal(devTeam, anotherTeam)
	fmt.Println(teamsEqual) // true

	sameNames := slices.EqualFunc(devTeam, anotherTeam, func(a, b Developer) bool {
		return a.Name == b.Name
	})
	fmt.Println(sameNames) // true

	numbers1 := []int{1, 2, 3}
	numbers2 := []int{1, 2, 4}
	comparison := slices.Compare(numbers1, numbers2)
	fmt.Println(comparison) // -1 (since 3 < 4)

	compResult := slices.CompareFunc(devTeam, anotherTeam, func(a, b Developer) int {
		if a.Name < b.Name {
			return -1
		} else if a.Name > b.Name {
			return 1
		}
		return 0
	})
	fmt.Println(compResult) // 0

	fmt.Println("\n6. Sorting Methods")
	fmt.Println("------------------")

	numbers := []int{9, 3, 6, 1, 8}
	slices.Sort(numbers)
	fmt.Println(numbers) // [1 3 6 8 9]

	coffeeTeam := slices.Clone(devTeam)
	slices.SortFunc(coffeeTeam, func(a, b Developer) int {
		return a.CoffeeLevel - b.CoffeeLevel
	})
	fmt.Println(coffeeTeam) // [{Bob 3 5} {Eve 4 3} {Charlie 6 1} {Alice 8 2} {Diana 9 0}]

	stableTeam := slices.Clone(devTeam)
	slices.SortStableFunc(stableTeam, func(a, b Developer) int {
		return a.BugCount - b.BugCount
	})
	fmt.Println(stableTeam) // [{Diana 9 0} {Charlie 6 1} {Alice 8 2} {Eve 4 3} {Bob 3 5}]

	isSliceSorted := slices.IsSorted([]int{1, 7, 3, 4, 5})
	fmt.Println(isSliceSorted) // false

	isTeamSorted := slices.IsSortedFunc(coffeeTeam, func(a, b Developer) int {
		return a.CoffeeLevel - b.CoffeeLevel
	})
	fmt.Println(isTeamSorted) // true

	fmt.Println("\n7. Mutation Methods")
	fmt.Println("-------------------")

	reverseNumbers := []int{1, 2, 3, 4, 5}
	slices.Reverse(reverseNumbers)
	fmt.Println(reverseNumbers) // [5 4 3 2 1]

	withoutMiddle := slices.Delete(slices.Clone(devTeam), 1, 3)
	fmt.Println(withoutMiddle) // [{Alice 8 2} {Diana 9 0} {Eve 4 3}]

	lowCoffeeTeam := slices.DeleteFunc(slices.Clone(devTeam), func(dev Developer) bool {
		return dev.CoffeeLevel < 5
	})
	fmt.Println(lowCoffeeTeam) // [{Alice 8 2} {Charlie 6 1} {Diana 9 0}]

	duplicates := []int{1, 1, 2, 2, 2, 3, 1, 1}
	compacted := slices.Compact(slices.Clone(duplicates))
	fmt.Println(compacted) // [1 2 3 1]

	devDuplicates := []Developer{
		{Name: "Alice", CoffeeLevel: 8, BugCount: 2},
		{Name: "Bob", CoffeeLevel: 8, BugCount: 3},
		{Name: "Charlie", CoffeeLevel: 6, BugCount: 1},
	}
	compactedDevs := slices.CompactFunc(devDuplicates, func(a, b Developer) bool {
		return a.CoffeeLevel == b.CoffeeLevel
	})
	fmt.Println(compactedDevs) // [{Alice 8 2} {Charlie 6 1}]

	replacement := []Developer{{Name: "Grace", CoffeeLevel: 10, BugCount: 0}}
	replaced := slices.Replace(slices.Clone(devTeam), 0, 1, replacement...)
	fmt.Println(replaced) // [{Grace 10 0} {Bob 3 5} {Charlie 6 1} {Diana 9 0} {Eve 4 3}]

	newTeam := slices.Insert(devTeam, 2, Developer{Name: "Frank", CoffeeLevel: 7, BugCount: 2})
	fmt.Println(newTeam) // [{Alice 8 2} {Bob 3 5} {Frank 7 2} {Charlie 6 1} {Diana 9 0} {Eve 4 3}]

	fmt.Println("\n8. Utility Methods")
	fmt.Println("------------------")

	testNumbers := []int{8, 3, 6, 9, 4}
	minNum := slices.Min(testNumbers)
	fmt.Println(minNum) // 3

	minCoffeeDev := slices.MinFunc(devTeam, func(a, b Developer) int {
		return a.CoffeeLevel - b.CoffeeLevel
	})
	fmt.Println(minCoffeeDev) // {Bob 3 5}

	maxNum := slices.Max(testNumbers)
	fmt.Println(maxNum) // 9

	maxCoffeeDev := slices.MaxFunc(devTeam, func(a, b Developer) int {
		return a.CoffeeLevel - b.CoffeeLevel
	})
	fmt.Println(maxCoffeeDev) // {Diana 9 0}

	largeSlice := make([]int, 5, 20)
	fmt.Println(len(largeSlice), cap(largeSlice)) // 5 20
	clippedSlice := slices.Clip(largeSlice)
	fmt.Println(len(clippedSlice), cap(clippedSlice)) // 5 5

	smallSlice := []int{1, 2, 3}
	fmt.Println(len(smallSlice), cap(smallSlice)) // 3 3
	grownSlice := slices.Grow(smallSlice, 10)
	fmt.Println(len(grownSlice), cap(grownSlice)) // 3 14
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

func prettyPrint(entries iter.Seq2[int, Developer]) string {
	result := ""
	for i, v := range entries {
		result += fmt.Sprintf("%d: %v, ", i, v)
	}
	if len(result) > 0 {
		result = result[:len(result)-2]
	}
	return result
}
