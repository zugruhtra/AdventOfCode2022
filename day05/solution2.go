package main

import (
    "fmt"
    "errors"
    "bufio"
    "os"
    "strings"
    "strconv"
)

// ----------------------------------- Main -----------------------------------

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    
    stacks := make([]Stack, 9)
    cargoPhase := true
    movePhase := false

    var moves []*Move

    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { cargoPhase = false; movePhase = true; continue }

        if cargoPhase { ParseCrateLine(line, stacks) }
        if movePhase { moves = append(moves, ParseCommandLine(line)) }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

    for _, stack := range stacks { stack.Reverse() }
    for _, move := range moves { MakeMove(*move, stacks) }

    for _, stack := range stacks {
        fmt.Print(string(stack.Head().name))
    }

}

// ---------------------------------- Model -----------------------------------

//// --------------------------------- Crate ----------------------------------

type Crate struct {
    name string
}

//// --------------------------------- Stack ----------------------------------

type Stack struct {
    data []Crate
}

func NewStack() *Stack {
    return &Stack{}
}

func (s *Stack) Push(c Crate) {
    s.data = append(s.data, c)
}

func (s *Stack) Pop() error {
    if len(s.data) == 0 {
        return errors.New("stack is emtpy")
    }

    s.data = s.data[:len(s.data)-1]
    return nil
}

func (s *Stack) Head() Crate {
    if len(s.data) == 0 { 
        panic("stack is empty")
    }

    return s.data[len(s.data)-1]
}

func (s *Stack) Reverse() {
    // from: https://golangcookbook.com/chapters/arrays/reverse/
    for i, j := 0, len(s.data)-1; i < j; i, j = i+1, j-1 {
        s.data[i], s.data[j] = s.data[j], s.data[i]
    }
}

func (s *Stack) Print() {
    fmt.Println(s.data)
}

//// -------------------------------- Commands --------------------------------

type Move struct {
    amount int
    from int
    to int
}

func NewMove(amount, from, to int) *Move {
    return &Move{amount, from, to}
}


// --------------------------------- Parsing ----------------------------------

func ParseCrateLine(line string, stacks []Stack) {
    parens := false
    for idx, char := range line {
        switch char {
        case '[':
            parens = true
        case ']':
            parens = false
        default:
            if parens {
                crate := Crate{string(char)}
                stacks[(idx - 1) / 4].Push(crate)
            }
        }
    }
}

func ParseCommandLine(line string) *Move {
    split := strings.Split(line, " ")
    amount, _ := strconv.Atoi(split[1])
    from, _ := strconv.Atoi(split[3])
    to, _ := strconv.Atoi(split[5])
    return NewMove(amount, from-1, to-1)
}



func MakeMove(move Move, stacks []Stack) {
    dest := &stacks[move.to]
    src := &stacks[move.from]
    tmp := NewStack()

    for i := move.amount; i > 0; i-- {
        tmp.Push(src.Head())
        src.Pop()
    }

    for i := move.amount; i > 0; i-- {
        dest.Push(tmp.Head())
        tmp.Pop()
    }
}
