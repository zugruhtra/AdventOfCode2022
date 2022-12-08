package main

import (
    "bufio"
    "fmt"
    "os"
)

type Tree struct {
    Height int
    Visible bool
}

func NewTree(height int) *Tree {
    return &Tree{height, false}
}

type Grid struct {
    grid [][]*Tree
    nRows int
}

func NewGrid() *Grid {
    return &Grid{
        grid: make([][]*Tree, 0), 
        nRows: 0,
    }
}

func (g *Grid) AddRow(row []*Tree) {
    var nRow = len(row)
    g.grid = append(g.grid, make([]*Tree, nRow))
    for i, tree := range row {
        g.grid[g.nRows][i] = tree
    }
    g.nRows++
}

func (g *Grid) LeftToRight(f func(*Tree, bool)) {
    for y := 0; y < g.nRows; y++ {
        f(g.grid[y][0], true)
        for x := 1; x < len(g.grid[y]); x++ {
            f(g.grid[y][x], false)
        }
    }
}

func (g *Grid) RightToLeft(f func(*Tree, bool)) {
    for y := 0; y < g.nRows; y++ {
        f(g.grid[y][len(g.grid[y])-1], true)
        for x := len(g.grid[y])-2;x >= 0; x-- {
            f(g.grid[y][x], false)
        }
    }
}

func (g *Grid) TopToBottom(f func(*Tree, bool)) {
    var rowLen = len(g.grid[0])
    for x := 0; x < rowLen; x++ {
        f(g.grid[0][x], true)
        for y := 1; y < g.nRows; y++ {
            f(g.grid[y][x], false)
        }
    }
}

func (g *Grid) BottomToTop(f func(*Tree, bool)) {
    var rowLen = len(g.grid[0])
    for x := 0; x < rowLen; x++ {
        f(g.grid[g.nRows-1][x], true)
        for y := g.nRows-2; y >= 0; y-- {
            f(g.grid[y][x], false)
        }
    }
}

func (g *Grid) ScenicScore(x, y int) int {
    if x == 0 || x == g.nRows-1 || y == 0 || y == len(g.grid[0])-1 {
        return 0
    }

    var score = 1
    var tree = g.grid[y][x]

    // left to right
    var reachedEdge = true
    for i := x+1; i < len(g.grid[y]); i++ {
        if tree.Height <= g.grid[y][i].Height {
            score *= i - x
            reachedEdge = false
            break
        }
    }
    if reachedEdge {
        score *= len(g.grid[y]) - x - 1
    }

    // right to left
    reachedEdge = true
    for i := x-1; i >= 0; i-- {
        if tree.Height <= g.grid[y][i].Height {
            score *= x - i
            reachedEdge = false
            break
        }
    }

    if reachedEdge {
        score *= x
    }
    // top to bottom
    reachedEdge = true
    for i := y+1; i < g.nRows; i++ {
        if tree.Height <= g.grid[i][x].Height {
            score *= i - y
            reachedEdge = false
            break
        }
    }

    if reachedEdge {
        score *= g.nRows - y - 1
    }
    // bottom to top
    reachedEdge = true
    for i := y-1; i >= 0; i-- {
        if tree.Height <= g.grid[i][x].Height {
            score *= y - i
            reachedEdge = false
            break
        }
    }

    if reachedEdge {
        score *= y
    }

    return score
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
    var grid = NewGrid()

    err := EachLineDo(func(line string) {
        var trees = make([]*Tree, 0)
        for _, c := range line {
            var height = int(c) - 48
            var tree = NewTree(height)
            trees = append(trees, tree)
        }
        grid.AddRow(trees)
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

    var cols = len(grid.grid)
    var rows = len(grid.grid[0])
    var maxScore = 0
    for x := 0; x < rows; x++ {
        for y := 0; y < cols; y++ {
            var score = grid.ScenicScore(x, y)
            if score > maxScore {
                maxScore = score
            }
        }
    }

    fmt.Println("Max. scenic score:", maxScore)
}
