package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Collidable interface {
	Collision(Point) bool
}

type Point struct {
	X, Y int
}

func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) Diff(other Point) Point {
	return Point{p.X - other.X, p.Y - other.Y}
}

func (p Point) Between(o1, o2 Point) bool {
	if o1.Equal(o2) {
		return p.Equal(o1)
	}
	if o2.X*o2.Y < o1.X*o1.Y {
		o1, o2 = o2, o1
	}
	diff := o1.Diff(o2)
	return ((diff.Y == 0 && p.Y == o1.Y && o1.X <= p.X && p.X <= o2.X) ||
		(diff.X == 0 && p.X == o1.X && o1.Y <= p.Y && p.Y <= o2.Y))
}

type RockPath struct {
	rocks  []Point
	Lowest int
}

func NewRockPath() *RockPath {
	return &RockPath{rocks: make([]Point, 0), Lowest: -1}
}

func (rp *RockPath) AddPoint(p Point) {
	rp.rocks = append(rp.rocks, p)
	rp.setLowest()
}

func (rp *RockPath) Collision(p Point) bool {
	rpLen := len(rp.rocks)
	for i := 0; i < rpLen-1; i++ {
		if p.Between(rp.rocks[i], rp.rocks[i+1]) {
			return true
		}
	}
	return false
}

func (rp *RockPath) setLowest() {
	min := 0
	for _, rock := range rp.rocks {
		if rock.Y > min {
			min = rock.Y
		}
	}
	rp.Lowest = min
}

type Sand struct {
	Pos     Point
	falling bool
}

func NewSand(p Point) *Sand {
	return &Sand{p, true}
}

func (s *Sand) Collision(p Point) bool {
	return s.Pos.Equal(p)
}

func (s *Sand) OutOfStructure(rps *[]*RockPath) bool {
	count := 0
	for _, rp := range *rps {
		if s.Pos.Y > rp.Lowest {
			count++
		}
	}
	return count == len(*rps)
}

func (s *Sand) PossibleMoves() [3]Point {
	return [3]Point{{s.Pos.X, s.Pos.Y + 1},
		{s.Pos.X - 1, s.Pos.Y + 1},
		{s.Pos.X + 1, s.Pos.Y + 1}}
}

func (s *Sand) Falling() bool {
	return s.falling
}

func (s *Sand) Update(obstacles *[]Collidable) {
	var collision bool
	for _, move := range s.PossibleMoves() {
		collision = false
		for _, obstacle := range *obstacles {
			if obstacle.Collision(move) {
				collision = true
				break
			}
		}
		if !collision {
			s.Pos = move
			return
		}
	}
	s.falling = false
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
	var rockPaths []*RockPath
	err := EachLineDo(func(line string) {
		split := strings.Split(line, "->")
		rockPath := NewRockPath()
		for _, s := range split {
			s = strings.Trim(s, " ")
			sPos := strings.Split(s, ",")
			x, _ := strconv.Atoi(sPos[0])
			y, _ := strconv.Atoi(sPos[1])
			point := Point{x, y}
			rockPath.AddPoint(point)
		}
		rockPaths = append(rockPaths, rockPath)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "reading stdin: %s", err)
		os.Exit(1)
	}

	var obstacles []Collidable
	for _, rockPath := range rockPaths {
		obstacles = append(obstacles, rockPath)
	}

	resting := 0
	for {
		sand := NewSand(Point{500, 0})
		for sand.Falling() && !sand.OutOfStructure(&rockPaths) {
			sand.Update(&obstacles)
		}
		if !sand.Falling() {
			obstacles = append(obstacles, sand)
			resting++
		}
		if sand.OutOfStructure(&rockPaths) && sand.Falling() {
			break
		}
	}
	fmt.Println("No. of resting sand grains:", resting)
}
