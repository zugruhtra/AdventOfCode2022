package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
)

type Range struct {
    start int
    end int
}

func (r *Range) Contains(other *Range) bool {
    return r.start <= other.start && r.end >= other.end
}

func StringToRange(arr string) *Range {
    split := strings.Split(arr, "-")
    start, _ := strconv.Atoi(split[0])
    end, _ := strconv.Atoi(split[1])
    return &Range{start, end}
}



func main() {
    scanner := bufio.NewScanner(os.Stdin)
    count := 0

    for scanner.Scan() {
        line := scanner.Text()

        rawPair := strings.Split(line, ",")

        first := StringToRange(rawPair[0])
        second := StringToRange(rawPair[1])

        if first.Contains(second) || second.Contains(first) {
            count++
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

    fmt.Println("Total num of fully containing pairs:", count)
}
