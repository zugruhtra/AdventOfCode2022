package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
    "math"
)

type Direction int

const (
    UP Direction = iota
    DOWN
    LEFT
    RIGHT
)

type Command struct {
    Direction Direction
    Steps int
}

func NewCommand(d Direction, steps int) *Command {
    return &Command{d, steps}
}

func StringToDirection(s string) Direction {
    switch s {
    case "U":
        return UP
    case "D":
        return DOWN
    case "L":
        return LEFT
    case "R":
        return RIGHT
    default:
        panic("unreachable line")
    }
}

type Position struct {
    X, Y int
}

func (p Position) Difference(o Position) Position {
    return Position{p.X - o.X, p.Y - o.Y}
}

type Head struct {
    pos Position
}

func NewHead(x, y int) *Head {
    return &Head{Position{x, y}}
}

func (h *Head) Move(direction Direction) {
    switch direction {
        case UP:
            h.pos.Y += 1
        case DOWN:
            h.pos.Y -= 1
        case LEFT:
            h.pos.X -= 1
        case RIGHT:
            h.pos.X += 1
        default:
            panic("unreachable line")
    }
}

func (h *Head) Position() Position {
    return h.pos
}

type Tail struct {
    pos Position
}

func NewTail(x, y int) *Tail {
    return &Tail{Position{x, y}}
}

func (t *Tail) Position() Position {
    return t.pos
}

func (t *Tail) FollowHead(head *Head) {
    var diff = t.pos.Difference(head.Position())

    if diff.X == 0 && diff.Y == 0 {
        return
    }

    absX := math.Abs(float64(diff.X)) 
    absY := math.Abs(float64(diff.Y))
    if absX*absX + absY*absY <= 2 {
        return
    }

    // Assumption: the tail needs only one step (straigth, diagonally) 
    //             to stay connected to the head.
    
    if (diff.X != 1 || diff.X != -1) && (diff.Y != 1 || diff.Y != -1) {
        // need to move
        if diff.X == 0 {
            // move straight
            if diff.Y > 0 {
                // move up
                t.pos.Y -= 1
            } else {
                // move down
                t.pos.Y += 1
            }
        } else if diff.Y == 0 {
            // move straight
            if diff.X > 0 {
                // move right
                t.pos.X -= 1
            } else {
                // move left
                t.pos.X += 1
            }
        } else {
            if diff.X > 0 && diff.Y > 0 {
                // move top right
                t.pos.X -= 1
                t.pos.Y -= 1
            } else if diff.X > 0 && diff.Y < 0 {
                // move down right
                t.pos.X -= 1
                t.pos.Y += 1
            } else if diff.X < 0 && diff.Y > 0 {
                // move top left
                t.pos.X += 1
                t.pos.Y -= 1
            } else if diff.X < 0 && diff.Y < 0 {
                // move down left
                t.pos.X += 1
                t.pos.Y += 1
            }
        }
    }
}

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
    var commands = make([]*Command, 0)
    err := EachLineDo(func(line string) {
        split := strings.Split(line, " ")
        direction := StringToDirection(split[0])
        steps, _ := strconv.Atoi(split[1])
        command := NewCommand(direction, steps)
        commands = append(commands, command)
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

    var (
        head = NewHead(0, 0)
        tail = NewTail(0, 0)
        visited = make(map[Position]int)
    )

    for _, cmd := range commands {
        for i := 0; i < cmd.Steps; i++ {
            head.Move(cmd.Direction)
            tail.FollowHead(head)
            visited[tail.Position()] = 1
        }
    }

    fmt.Println("Num. of visited positions by the tail:", len(visited))
}
