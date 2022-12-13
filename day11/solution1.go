package main

import (
    "container/list"
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
    "sort"
)

// ---------------------------------- Model -----------------------------------

type Operation interface {
    Exec(int, int) int
}

type Add struct {}
type Mul struct {}

func (a Add) Exec(x, y int) int {
    fmt.Printf("    New worry level: %d + %d = %d\n", x, y, x+y)
    return x + y
}

func (t Mul) Exec(x, y int) int {
    fmt.Printf("    New worry level: %d * %d = %d\n", x, y, x*y)
    return x * y
}

type Test interface {
    Test(int) bool
}

type Divisible struct {
    x int
}

func (d Divisible) Test(n int) bool {
    fmt.Printf("    Current worry level is not divisible by %d.\n", d.x)
    return n % d.x == 0
}

type Monkey struct {
    id int
    items *list.List
    operation Operation
    operands [2]string
    test Test
    targets [2]int
    monkies *[]*Monkey
    inspectedItems int
}

func NewMonkey(id int, items []int, operation Operation, operands [2]string, test Test, targets [2]int) *Monkey {
    lst := list.New()
    for _, item := range items {
        lst.PushBack(item)
    }
    return &Monkey{id, lst, operation, operands, test, targets, nil, 0}
}

func (m *Monkey) Init(others *[]*Monkey) {
    m.monkies = others
}

func (m *Monkey) Inspect() {
    fmt.Printf("Monkey %d:\n", m.id)
    if m.items.Len() == 0 {
        fmt.Println("  Monkey has no items.")
        return
    }
    
    var next *list.Element
    for e := m.items.Front(); e != nil; e = next {
        fmt.Printf("  Monkey inspectes an item with a worry level of %d\n", e.Value)
        m.inspectedItems++
        var val [2]int
        worryLevel,_ := e.Value.(int)

        for j, operand := range m.operands {
            if operand == "old" {
                val[j] = worryLevel
            } else {
                val[j],_ = strconv.Atoi(operand)
            }
        }

        worryLevel = m.operation.Exec(val[0], val[1])

        worryLevel = worryLevel / 3
        fmt.Printf("    Monkey gets bored with item. Worry level is divided by 3 to %d.\n", worryLevel)

        var target int
        if m.test.Test(worryLevel) {
            target = 0
        } else {
            target = 1
        }

        fmt.Printf("    Item with worry level %d is thrown to monkey %d.\n", worryLevel, m.targets[target])

        e.Value = worryLevel
        next = e.Next()
        (*m.monkies)[m.targets[target]].items.PushBack(e.Value)
        m.items.Remove(e)
    }
}

func (m *Monkey) PrintItems() {
    fmt.Print("Monkey ", m.id)
    fmt.Print(": ")
    for e := m.items.Front(); e != nil; e = e.Next() {
        val,_ := e.Value.(int)
        fmt.Print(val, " ")
    }
    fmt.Println()
}

type MonkeyBuilder struct {
    id int
    items []int
    operation Operation
    operands []string
    test Test
    targets []int
}

func NewMonkeyBuilder() *MonkeyBuilder {
    return &MonkeyBuilder{-1, make([]int, 0), nil, make([]string, 0), nil, make([]int, 0)}
}

func (mb *MonkeyBuilder) AddId(id int) {
    mb.id = id
}

func (mb *MonkeyBuilder) AddItem(item int) {
    mb.items = append(mb.items, item)
}

func (mb *MonkeyBuilder) AddOperation(operation Operation) {
    mb.operation = operation
}

func (mb *MonkeyBuilder) AddOperand(op string) {
    mb.operands = append(mb.operands, op)
}

func (mb *MonkeyBuilder) AddTest(test Test) {
    mb.test = test
}

func (mb *MonkeyBuilder) AddTarget(target int) {
    mb.targets = append(mb.targets, target)
}

func (mb *MonkeyBuilder) Build() *Monkey {
    var operands = [2]string{mb.operands[0], mb.operands[1]}
    var targets = [2]int{mb.targets[0], mb.targets[1]}
    return NewMonkey(mb.id, mb.items, mb.operation, operands, mb.test, targets)
}

func (mb *MonkeyBuilder) Clear() {
    mb.id = -1
    mb.items = nil
    mb.operation = nil
    mb.operands = nil
    mb.test = nil
    mb.targets = nil
}

// ---------------------------------- Parser ----------------------------------

type Parser struct {
    mb *MonkeyBuilder
    processing, ready bool
}

func NewParser() *Parser {
    return &Parser{mb: NewMonkeyBuilder(), processing: false, ready: false}
}

func (p *Parser) Feed(line string) {
    if strings.HasPrefix(line, "Monkey") {
        p.processing = true
        p.ready = false
        split := strings.Split(line, " ")
        id,_ := strconv.Atoi(split[1][:len(split[1])-1])
        p.mb.AddId(id)
        return
    }

    line = strings.Trim(line, " ")
    split := strings.Split(line, ":")

    switch split[0] {
    case "Starting items":
        p.parseItemStatement(split[1])
    case "Operation":
        p.parseOperationStatement(split[1])
    case "Test":
        p.parseTestStatement(split[1])
    case "If true":
        p.parseIfStatement(split[1])
    case "If false":
        p.parseIfStatement(split[1])
    case "":
        if p.processing {
            p.processing = false
            p.ready = true
        }
    default:
        panic("unreachable line")
    }
}
func (p *Parser) parseItemStatement(line string) {
    line = strings.Trim(line, " ")
    sItems := strings.Split(line,  ",")
    for _, sItem := range sItems {
        sItem = strings.Trim(sItem, " ")
        item,_ := strconv.Atoi(sItem)
        p.mb.AddItem(item)
    }
}

func (p *Parser) parseOperationStatement(line string) {
    line = strings.Trim(line, " ")
    equation := strings.Split(line, " ")
    var operation Operation
    switch equation[3] {
    case "+":
        operation = Add{}
    case "*":
        operation = Mul{}
    default:
        panic("unreachable line")
    }
    p.mb.AddOperation(operation)
    p.mb.AddOperand(equation[2])
    p.mb.AddOperand(equation[4])
}

func (p *Parser) parseTestStatement(line string) {
    line = strings.Trim(line, " ")
    statement := strings.Split(line, " ")
    var test Test
    switch statement[0] {
    case "divisible":
        numb,_ := strconv.Atoi(statement[2])
        test = Divisible{numb}
    default:
        panic("unreachable line")
    }
    p.mb.AddTest(test)
}

func (p *Parser) parseIfStatement(line string) {
    line = strings.Trim(line, " ")
    statement := strings.Split(line, " ")
    var target int
    switch statement[0] {
    case "throw":
        target,_ = strconv.Atoi(statement[3])
    default:
        panic("unreachable line")
    }
    p.mb.AddTarget(target)
}

func (p *Parser) Ready() bool {
    return p.ready
}

func (p *Parser) GenMonkey() *Monkey {
    var monkey = p.mb.Build()
    p.mb.Clear()
    return monkey
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
    var monkies []*Monkey
    parser := NewParser()

    err := EachLineDo(func(line string) {
        parser.Feed(line)
        if parser.Ready() {
            monkey := parser.GenMonkey()
            monkies = append(monkies, monkey)
        }
    })

    monkey := parser.GenMonkey()
    monkies = append(monkies, monkey)

    if err != nil {
        fmt.Fprintf(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

    for _, monkey := range monkies {
        monkey.Init(&monkies)
    }

    for i := 0; i < 20; i++ {
        for _, monkey := range monkies {
            monkey.Inspect()
        }
    }

    fmt.Println()

    var inspected []int
    for _, monkey := range monkies {
        inspected = append(inspected, monkey.inspectedItems)
    }
    sort.Ints(inspected)
    var nMonkies = len(monkies)

    fmt.Println("Level of monkey business:", inspected[nMonkies-1] * inspected[nMonkies-2])
}
