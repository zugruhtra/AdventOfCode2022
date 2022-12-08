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

    var height int
    var markVisibleTree = func(tree *Tree, newRow bool) {
        if newRow {
            height = -1
        }
        if tree.Height > height {
            height = tree.Height
            tree.Visible = true
        }
    }

    grid.LeftToRight(markVisibleTree)
    grid.RightToLeft(markVisibleTree)
    grid.TopToBottom(markVisibleTree)
    grid.BottomToTop(markVisibleTree)

    var nVisible = 0
    grid.LeftToRight(func(tree *Tree, newRow bool) {
        if tree.Visible {
            nVisible++
        }
    })

    fmt.Println("Num. of visible trees:", nVisible)
}
