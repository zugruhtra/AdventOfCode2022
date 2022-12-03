package main

import (
    "bufio"
    "fmt"
    "os"
    "unicode/utf8"
)

func main() {
    var priorities []int

    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        rucksack := scanner.Text()
        sameItem := findSameItemInRucksack(rucksack)
        priority := getPriority(sameItem)
        priorities = append(priorities, priority)
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "error reading stdin:", err)
        os.Exit(1)
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

func findSameItemInRucksack(rucksack string) rune {
    firstCompartment := rucksack[:len(rucksack)/2]
    secondCompartment := rucksack[len(rucksack)/2:]

    compareSet := make(map[rune]int)

    for i, w := 0, 0; i < len(firstCompartment); i += w {
        runeValue, width := utf8.DecodeRuneInString(firstCompartment[i:])
        w = width
        compareSet[runeValue] = i
    }

    for i, w := 0, 0; i < len(secondCompartment); i += w {
        runeValue, width := utf8.DecodeRuneInString(secondCompartment[i:])
        w = width
        _, ok := compareSet[runeValue]
        if ok {
            return runeValue
        }
    }

    panic("No similiar item in both compartments")
}

func sum(ns []int) int {
    total := 0
    for _,n := range ns {
        total += n
    }
    return total
}
