package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ---------------------------------- Point -----------------------------------

type Point struct {
	X, Y int
}

func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) Diff(other Point) Point {
	return Point{p.X - other.X, p.Y - other.Y}
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
}

// ------------------------------- ObstacleMap --------------------------------

type ObstacleMap struct {
	Width, Height int
	buffer        []bool
}

func NewObstacleMap(width, height int) *ObstacleMap {
	return &ObstacleMap{width, height, make([]bool, width*height)}
}

func (obm *ObstacleMap) AddObstacle(p Point) {
	obm.buffer[p.X+p.Y*obm.Width] = true
}

func (obm *ObstacleMap) AddLine(p1, p2 Point) {
	if p2.X*p2.Y < p1.X*p1.Y {
		p1, p2 = p2, p1
	}
	diff := p1.Diff(p2)
	if diff.X == 0 {
		for i := 0; i <= p2.Y-p1.Y; i++ {
			obm.buffer[p1.X+(p1.Y+i)*obm.Width] = true
		}
	} else if diff.Y == 0 {
		for i := 0; i <= p2.X-p1.X; i++ {
			obm.buffer[p1.X+p1.Y*obm.Width+i] = true
		}
	} else {
		panic("not a horizontal or vertical line")
	}
}

func (obm *ObstacleMap) Collision(p Point) (bool, error) {
	var idx = p.X + p.Y*obm.Width
	if obm.OutOfMap(p) {
		return false, errors.New("Out of map")
	}
	return obm.buffer[idx], nil
}

func (obm *ObstacleMap) OutOfMap(p Point) bool {
	var result = p.X >= obm.Width || p.Y >= obm.Height
	return result
}

func (obm *ObstacleMap) String() string {
	var sb strings.Builder
	for i := 0; i < obm.Height; i++ {
		for j := 0; j < obm.Width; j++ {
			if obm.buffer[j+i*obm.Width] {
				sb.WriteRune('#')
			} else {
				sb.WriteRune('.')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

// --------------------------------- RockPath ---------------------------------

type RockPath struct {
	rocks            []Point
	Lowest, Furthest int
}

func NewRockPath() *RockPath {
	return &RockPath{rocks: make([]Point, 0), Lowest: -1, Furthest: -1}
}

func (rp *RockPath) AddPoint(p Point) {
	rp.rocks = append(rp.rocks, p)
	rp.setLowest()
	rp.setFurthest()
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

func (rp *RockPath) setFurthest() {
	min := 0
	for _, rock := range rp.rocks {
		if rock.X > min {
			min = rock.X
		}
	}
	rp.Furthest = min
}

// ----------------------------------- Sand -----------------------------------

type Sand struct {
	Pos               Point
	falling, outOfMap bool
}

func NewSand(p Point) *Sand {
	return &Sand{p, true, false}
}

func (s *Sand) PossibleMoves() [3]Point {
	return [3]Point{{s.Pos.X, s.Pos.Y + 1},
		{s.Pos.X - 1, s.Pos.Y + 1},
		{s.Pos.X + 1, s.Pos.Y + 1}}
}

func (s *Sand) Falling() bool {
	return s.falling
}

func (s *Sand) OutOfMap() bool {
	return s.outOfMap
}

func (s *Sand) Update(obstacles *ObstacleMap) {
	for _, move := range s.PossibleMoves() {
		if result, err := obstacles.Collision(move); !result {
			if err != nil {
				s.outOfMap = true
			}
			s.Pos = move
			return
		}
	}
	s.falling = false
}

// ----------------------------------- Misc -----------------------------------

func GetWidthAndHeight(rockPaths []*RockPath) (int, int) {
	var width, height int
	var lowest, furthest []int
	for _, rockPath := range rockPaths {
		lowest = append(lowest, rockPath.Lowest)
		furthest = append(furthest, rockPath.Furthest)
	}
	sort.Ints(lowest)
	sort.Ints(furthest)
	height = lowest[len(lowest)-1] + 1
	width = furthest[len(furthest)-1] + 1
	return width, height
}

func PopulateObstacleMap(obm *ObstacleMap, rp *RockPath) {
	rpLen := len(rp.rocks)
	if rpLen == 1 {
		obm.AddObstacle(rp.rocks[0])
		return
	}
	for i := 0; i < rpLen-1; i++ {
		r1, r2 := rp.rocks[i], rp.rocks[i+1]
		obm.AddLine(r1, r2)
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

// ----------------------------------- Main -----------------------------------

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

	width, height := GetWidthAndHeight(rockPaths)
	width, height = width*2, height+2
	obstacles := NewObstacleMap(width, height)
	for _, rockPath := range rockPaths {
		PopulateObstacleMap(obstacles, rockPath)
	}
	for i := 0; i < width; i++ {
		obstacles.AddObstacle(Point{i, height - 1})
	}

	resting := 0
	for {
		sand := NewSand(Point{500, 0})
		for sand.Falling() && !sand.OutOfMap() {
			sand.Update(obstacles)
		}
		if !sand.Falling() {
			obstacles.AddObstacle(sand.Pos)
			resting++
		}
		if sand.Pos.Equal(Point{500, 0}) {
			break
		}
	}
	fmt.Println("No. of resting sand grains:", resting)
}
