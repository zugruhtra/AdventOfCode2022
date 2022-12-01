package main

import (
    "os"
    "bufio"
    "fmt"
    "strconv"
    "sort"
)

func max(ns []int) int {
    ans := ns[0]
    for i := 1; i < len(ns); i++ {
        if ans < ns[i] {
            ans = ns[i]
        }
    }
    return ans
}

func sum(ns []int) int {
    total := 0
    for i := 0; i < len(ns); i++ {
        total += ns[i]
    }
    return total
}


func main() {
    elf, food := 0, 0
    elfs := make([][]int, 239)
    scanner := bufio.NewScanner(os.Stdin)


    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            elf++
            food = 0
        } else {
            calories, err := strconv.Atoi(line)
            if err != nil {
                fmt.Fprintln(os.Stderr, "can not convert", line, "to integer")
                continue
            }

            if len(elfs[elf]) == 0 {
                elfs[elf] = make([]int, 15)
            }
            elfs[elf][food] = calories
            food++
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading stdin:", err)
    }

    total_calories := make([]int, len(elfs))
    for i := 0; i < len(elfs); i++ {
        total_calories[i] = sum(elfs[i])
    }

    sort.Ints(total_calories)

    fmt.Printf("Sum of top three calories carrying elfs: %d calories", sum(total_calories[len(elfs)-3:]))
}
