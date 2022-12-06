package main

import (
    "bufio"
    "fmt"
    "os"
    "unicode/utf8"
)

const MEMORY = 4

func EachLineDo(f func(string)) error {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        line := scanner.Text()
        f(line)
    }

    err := scanner.Err()

    return err

}

func AllUnequal(rs []rune) bool {
    nrs := len(rs)
    for i, r := range rs[:nrs-1] {
        for _, rr := range rs[i+1:] {
            if r == rr {
                return false
            }
        }
    }
    return true
}

func main() {
    var (
        marker string
        idx int
    )

    ring := make([]rune, MEMORY)

    err := EachLineDo(func(line string) {
        w := 0
        for i := 0; i < MEMORY; i++ {
            runeValue, width := utf8.DecodeRuneInString(line[w:])
            w += width
            ring[i] = runeValue
        }

        idx = MEMORY

        if AllUnequal(ring) {
            marker = line[:4]
            fmt.Printf("Marker: %s at %d\n", marker, idx)
            return
        }
        
        for i, c := range line[4:] {
            idx = i + 4
            ring[i%MEMORY] = c

            if AllUnequal(ring) {
                marker = line[idx-3:idx+1]
                fmt.Printf("Marker: %s at %d\n", marker, idx+1)
                return
            }
        }
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

}

