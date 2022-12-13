package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const INF = math.MaxInt32

type Position struct {
	x, y int
}

type Graph struct {
	edges    []*Edge
	vertices []int
}

type Edge struct {
	From, To, Weight int
}

func NewEdge(from, to, weight int) *Edge {
	return &Edge{from, to, weight}
}

func NewGraph(edges []*Edge, vertices []int) *Graph {
	return &Graph{edges, vertices}
}

func (g *Graph) BellmanFord(source int) ([]int, []int) {
	size := len(g.vertices)
	distances := make([]int, size)
	predecessors := make([]int, size)

	for _, v := range g.vertices {
		distances[v] = INF
	}
	distances[source] = 0

	for i, changes := 0, 0; i < size-1; i, changes = i+1, 0 {
		for _, edge := range g.edges {
			newDist := distances[edge.From] + edge.Weight
			if newDist < distances[edge.To] && edge.Weight < 2 {
				distances[edge.To] = newDist
				predecessors[edge.To] = edge.From
				changes++
			}
		}
		if changes == 0 {
			break
		}
	}

	return predecessors, distances
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
	var heightmap [][]int
	var start, end Position

	heightmap = make([][]int, 0)
	row := 0

	err := EachLineDo(func(line string) {
		heightmap = append(heightmap, make([]int, 0))
		for i, c := range line {
			switch c {
			case 'S':
				start = Position{i, row}
				heightmap[row] = append(heightmap[row], 0)
			case 'E':
				end = Position{i, row}
				heightmap[row] = append(heightmap[row], 25)
			default:
				heightmap[row] = append(heightmap[row], int(c-'a'))
			}
		}
		row++
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "reading stdin:", err)
		os.Exit(1)
	}

	nRows := len(heightmap)
	nCols := len(heightmap[0])
	nVertices := nRows * nCols

	vertices := make([]int, nVertices)
	for i := 0; i < nVertices; i++ {
		vertices[i] = i
	}

	edges := make([]*Edge, 0)
	for i := 0; i < nRows; i++ {
		for j := 0; j < nCols; j++ {
			p := j + i*nCols
			// left
			if j-1 >= 0 {
				edge := NewEdge(p, p-1, heightmap[i][j-1]-heightmap[i][j])
				edges = append(edges, edge)
			}
			// right
			if j+1 < nCols {
				edge := NewEdge(p, p+1, heightmap[i][j+1]-heightmap[i][j])
				edges = append(edges, edge)
			}
			// top
			if p-nCols >= 0 {
				edge := NewEdge(p, p-nCols, heightmap[i-1][j]-heightmap[i][j])
				edges = append(edges, edge)
			}
			// down
			if p+nCols < nVertices {
				edge := NewEdge(p, p+nCols, heightmap[i+1][j]-heightmap[i][j])
				edges = append(edges, edge)
			}
		}
	}

	source := start.x + start.y*nCols
	goal := end.x + end.y*nCols

	graph := NewGraph(edges, vertices)
	pred, _ := graph.BellmanFord(source)

	steps := 0
	prev := goal
	for prev != source {
		prev = pred[prev]
		steps++
	}

	fmt.Println("Num. of Steps:", steps)
}
