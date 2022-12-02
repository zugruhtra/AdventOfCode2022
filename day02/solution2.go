package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

type Shape int

const (
    Rock    Shape = 1
    Paper   Shape = 2
    Scissor Shape = 3
)

type Result int

const (
    Loss    Result = 0
    Draw    Result = 3
    Win     Result = 6
)

func (s *Shape) play(o Shape) Result {
    switch *s {
    case Rock:
        switch o {
        case Rock:
            return Draw
        case Paper:
            return Loss
        case Scissor:
            return Win
        default:
            panic("error: can't find opponent's shape")
        }
    case Paper:
        switch o {
        case Rock:
            return Win
        case Paper:
            return Draw
        case Scissor:
            return Loss
        default:
            panic("error: can't find opponent's shape")
        }
    case Scissor:
        switch o {
        case Rock:
            return Loss
        case Paper:
            return Win
        case Scissor:
            return Draw
        default:
            panic("error: can't find opponent's shape")
        }
    default:
        panic("error: can't find your shape")
    }
}

func getOpponentShape(s string) Shape {
    switch s {
    case "A":
        return Rock
    case "B":
        return Paper
    case "C":
        return Scissor
    default:
        panic("error: can't find opponent's shape")
    }
}

func getOwnShape(opponent Shape, strategy string) Shape {
    switch opponent {
    case Rock:
        switch strategy {
        case "X":
            return Scissor
        case "Y":
            return Rock
        case "Z":
            return Paper
        default:
            panic("error: can't find your shape")
        }
    case Paper:
        switch strategy {
        case "X":
            return Rock
        case "Y":
            return Paper
        case "Z":
            return Scissor
        default:
            panic("error: can't find your shape")
        }
    case Scissor:
        switch strategy {
        case "X":
            return Paper
        case "Y":
            return Scissor
        case "Z":
            return Rock
        default:
            panic("error: can't find your shape")
        }
    default:
        panic("error: can't find opponent's shape")
    }
}

func sum(ns []int) int {
    total := 0
    for i := 0; i < len(ns); i++ {
        total += ns[i]
    }
    return total
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    var results []int

    for scanner.Scan() {
        line := scanner.Text()
        split := strings.Split(line, " ")
        opponent := getOpponentShape(split[0])
        own := getOwnShape(opponent, split[1])

        result := own.play(opponent)
        total := int(result) + int(own)
        results = append(results, total)
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading stdin:", err)
        return
    }

    fmt.Println("Total Score:", sum(results))
}
