package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
)

// ------------------------------- Instructions -------------------------------

type Instruction interface {
    Delay() int
}

type AddX struct {
    Amount int
    delay int
}

type Noop struct {
    delay int
}

func NewAddx(amount int) AddX {
    return AddX{amount, 2}
}

func NewNoop() Noop {
    return Noop{1}
}

func (ins AddX) Delay() int {
    return ins.delay
}

func (ins Noop) Delay() int {
    return ins.delay
}

// ----------------------------------- CPU ------------------------------------

type Cpu struct {
    X int
    pc int
    program []Instruction
    delay int
}

func NewCpu(program []Instruction) *Cpu {
    return &Cpu{X: 1, pc: 0, program: program, delay: 0}
}

func (cpu *Cpu) Fetch() {
    cpu.delay = cpu.program[cpu.pc].Delay()
}

func (cpu *Cpu) Exec() {
    cpu.delay--

    if cpu.delay > 0 {
        return
    }

    var ins = cpu.program[cpu.pc]

    switch ins.(type) {
    case AddX:
        addx, _ := ins.(AddX)
        cpu.X += addx.Amount
    case Noop:
        ;
    default:
        panic("unreachable line")
    }

    cpu.pc++
}

func (cpu *Cpu) Processing() bool {
    return cpu.delay > 0
}

func (cpu *Cpu) Ready() bool {
    return cpu.pc == len(cpu.program)
}

// ----------------------------------- Main -----------------------------------

func EachLineDo(f func(string)) error {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        line := scanner.Text()
        f(line)
    }

    err := scanner.Err()

    return err
}

func main() {
    var program []Instruction

    err := EachLineDo(func(line string) {
        split := strings.Split(line, " ")
        switch split[0] {
        case "addx":
            amount, _ := strconv.Atoi(split[1])
            program = append(program, NewAddx(amount))
        case "noop":
            program = append(program, NewNoop())
        default:
            panic("unreachable line")
        }
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

    var cpu = NewCpu(program)
    var cycle = 1
    var checkpoints = [...]int{220, 180, 140, 100, 60, 20}
    var totalSignalStrength = 0

    for !cpu.Ready() {
        for _, checkpoint := range checkpoints {
            if cycle == checkpoint {
                totalSignalStrength += cycle * cpu.X
                break
            }
        }

        if !cpu.Processing() {
            cpu.Fetch()
            cpu.Exec()
        } else {
            cpu.Exec()
        }

        cycle++
    }

    fmt.Println("Sum of signal strengths:", totalSignalStrength)
}
