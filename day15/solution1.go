package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// ---------------------------------- Sensor ----------------------------------

type Sensor struct {
	Pos           Point
	ClosestBeacon *Beacon
}

func NewSensor(p Point, b *Beacon) *Sensor {
	return &Sensor{p, b}
}

func (s *Sensor) DistanceToClosetstBeacon() int {
	return ManhattanDistance(s.Pos, s.ClosestBeacon.Pos)
}

// ---------------------------------- Beacon ----------------------------------

type Beacon struct {
	Pos Point
}

func NewBeacon(p Point) *Beacon {
	return &Beacon{p}
}

// ----------------------------------- Misc -----------------------------------

func GetSensorFromRawLine(line string) *Sensor {
	var (
		re         *regexp.Regexp
		matches    []string
		rawX, rawY string
		x, y       int
	)

	pattern := `Sensor at x=(-*\d+), y=(-*\d+): closest beacon is at x=(-*\d+), y=(-*\d+)`
	re = regexp.MustCompile(pattern)
	matches = re.FindStringSubmatch(line)

	rawX, rawY = matches[3], matches[4]
	x, _ = strconv.Atoi(rawX)
	y, _ = strconv.Atoi(rawY)
	beacon := NewBeacon(Point{x, y})

	rawX, rawY = matches[1], matches[2]
	x, _ = strconv.Atoi(rawX)
	y, _ = strconv.Atoi(rawY)
	sensor := NewSensor(Point{x, y}, beacon)

	return sensor
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

type Point struct {
	X, Y int
}

func (p *Point) Equal(o Point) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p *Point) Diff(o Point) Point {
	return Point{p.X - o.X, p.Y - o.Y}
}

type Line struct {
	O, P Point
}

func (l *Line) Includes(o Line) bool {
	if !(l.Horizontal() && o.Horizontal() && l.O.Y == o.O.Y) {
		panic("both lines need to be horizontal and on the same y-coordinate")
	}
	return l.IncludesPoint(o.O) && l.IncludesPoint(o.P)
}

func (l *Line) IncludesPoint(p Point) bool {
	return l.O.Y == p.Y && l.P.Y == p.Y && l.O.X <= p.X && p.X <= l.P.X
}

func (l *Line) Horizontal() bool {
	return l.O.Y == l.P.Y
}

func (l *Line) Length() int {
	return l.P.X - l.O.X
}

type Circle struct {
	Pos    Point
	Radius int
}

func (c *Circle) IntersectsWithHorizontalLine(y int) bool {
	return c.Pos.Y-c.Radius <= y && y <= c.Pos.Y+c.Radius
}

func (c *Circle) IntersectionsWithHorizontalLine(y int) (*Point, *Point) {
	if !c.IntersectsWithHorizontalLine(y) {
		return nil, nil
	}

	distance := Abs(y - c.Pos.Y)
	if distance == 0 || distance == 2*c.Radius {
		return &Point{c.Pos.X, y}, nil
	}

	dx := c.Radius - distance
	p1 := &Point{c.Pos.X - dx, y}
	p2 := &Point{c.Pos.X + dx, y}

	return p1, p2
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func ManhattanDistance(p1, p2 Point) int {
	diff := p1.Diff(p2)
	return Abs(diff.X) + Abs(diff.Y)
}

func GetAllIncluded(lines []Line) [][]Line {
	var included [][]Line
	alreadyIncluded := make(map[Line]bool)
	idx := 0
	for i, line := range lines {
		if _, prs := alreadyIncluded[line]; prs {
			continue
		}
		included = append(included, make([]Line, 0))
		included[idx] = append(included[idx], line)
		for _, other := range lines[i+1:] {
			if line.Includes(other) {
				alreadyIncluded[other] = true
				included[idx] = append(included[idx], other)
			}
		}
		idx++
	}
	return included
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go run solution1.go <NUMB>")
	os.Exit(1)
}

// ----------------------------------- Main -----------------------------------

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	y, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
	}

	var sensors []*Sensor
	err = EachLineDo(func(line string) {
		sensor := GetSensorFromRawLine(line)
		sensors = append(sensors, sensor)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "reading stdin:", err)
		os.Exit(1)
	}

	var ranges []Line
	for _, sensor := range sensors {
		coverageArea := Circle{sensor.Pos, sensor.DistanceToClosetstBeacon()}
		intersect1, intersect2 := coverageArea.IntersectionsWithHorizontalLine(y)

		switch {
		case intersect1 == nil && intersect2 == nil:

		case intersect2 == nil:
			ranges = append(ranges, Line{*intersect1, *intersect1})
		default:
			ranges = append(ranges, Line{*intersect1, *intersect2})
		}
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].O.X < ranges[j].O.X
	})

	included := GetAllIncluded(ranges)
	lenIncluded := len(included)
	score := 1
	for i := 0; i < lenIncluded-1; i++ {
		line, nextLine := included[i][0], included[i+1][0]
		if line.IncludesPoint(nextLine.O) {
			score += line.Length() - (line.P.X - nextLine.O.X)
		} else {
			score += line.Length()
		}
	}
	score += included[lenIncluded-1][0].Length()

	seenBeacon := make(map[Point]bool)
	for _, sensor := range sensors {
		for _, include := range included {
			if _, prs := seenBeacon[sensor.ClosestBeacon.Pos]; prs {
				continue
			}
			if include[0].IncludesPoint(sensor.ClosestBeacon.Pos) {
				seenBeacon[sensor.ClosestBeacon.Pos] = true
				score--
			}
		}
	}

	fmt.Println("No. of Position which can not contain a beacon:", score)

}
