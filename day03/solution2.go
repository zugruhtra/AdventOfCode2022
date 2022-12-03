package main

import (
    "bufio"
    "fmt"
    "os"
    "unicode/utf8"
    "strings"
)

func main() {
    var (
        priorities []int
        groups     [][]string
    )

    groupIdx := 0
    rucksackIdx := 0
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        if rucksackIdx == 0 {
            groups = append(groups, make([]string, 3))
        }

        rucksack := scanner.Text()
        groups[groupIdx][rucksackIdx] = rucksack

        rucksackIdx = (rucksackIdx + 1) % 3
        if rucksackIdx == 0 {
            groupIdx++
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "error reading stdin:", err)
        os.Exit(1)
    }

    for _, group := range groups {
        sameItem := findSameItemInGroup(group)
        priority := getPriority(sameItem)
        priorities = append(priorities, priority)
    }

    fmt.Println("Sum of priorities:", sum(priorities))
}

func getPriority(r rune) int {
    if r >= 97 {
        return int(r - 96)
    } else {
        return int(r - 38)
    }
}

func findSameItemInGroup(group []string) rune {
    var sameItems string

    sameItems = findSimiliarItemsInTwoRucksack(group[0], group[1])
    sameItems = findSimiliarItemsInTwoRucksack(sameItems, group[2])

    if len(sameItems) == 1 {
        runeValue, _ := utf8.DecodeRuneInString(sameItems)
        return runeValue
    }

    panic("More than one same item in group")
}

func findSimiliarItemsInTwoRucksack(a, b string) string {
    var (
        sb strings.Builder
        sameItems []rune
    )

    compareSet := make(map[rune]int)

    a = getUniqueItems(a)
    b = getUniqueItems(b)

    for _, runeValue := range a {
        compareSet[runeValue] = 0
    }

    for _, runeValue := range b {
        _, ok := compareSet[runeValue]
        if ok {
            compareSet[runeValue]++
        } else {
            compareSet[runeValue] = 0
        }
    }

    for runeValue, count := range compareSet {
        if count > 0 {
            sameItems = append(sameItems, runeValue)
        }
    }

    for _, item := range sameItems {
        sb.WriteRune(item)
    }

    return sb.String()
}

func getUniqueItems(rucksack string) string {
    var sb strings.Builder
    compareSet := make(map[rune]int)

    for _, runeValue := range rucksack {
        compareSet[runeValue] = 0
    }

    for runeValue := range compareSet {
        sb.WriteRune(runeValue)
    }

    return sb.String()
}

func sum(ns []int) int {
    total := 0
    for _,n := range ns {
        total += n
    }
    return total
}
