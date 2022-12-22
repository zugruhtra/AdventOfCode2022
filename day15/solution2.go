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

func (p *Point) Chop(r1, r2 int) {
	switch {
	case p.X < r1:
		p.X = r1
	case p.X > r2:
		p.X = r2
	}
	switch {
	case p.Y < r1:
		p.Y = r1
	case p.Y > r2:
		p.Y = r2
	}
}

type Line struct {
	P1, P2 Point
}

func (l *Line) Includes(o Line) bool {
	if !(l.Horizontal() && o.Horizontal() && l.P1.Y == o.P1.Y) {
		panic("both lines need to be horizontal and on the same y-coordinate")
	}
	return l.IncludesPoint(o.P1) && l.IncludesPoint(o.P2)
}

func (l *Line) IncludesPoint(p Point) bool {
	return l.P1.Y == p.Y && l.P2.Y == p.Y && l.P1.X <= p.X && p.X <= l.P2.X
}

func (l *Line) Horizontal() bool {
	return l.P1.Y == l.P2.Y
}

func (l *Line) Length() int {
	return l.P2.X - l.P1.X + 1
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
	if distance == c.Radius {
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
	searchRange, _ := strconv.Atoi(os.Args[1])

	var sensors []*Sensor
	err := EachLineDo(func(line string) {
		sensor := GetSensorFromRawLine(line)
		sensors = append(sensors, sensor)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading stdin:", err)
		os.Exit(1)
	}

	for y := 0; y <= searchRange; y++ {
		var ranges []Line
		for _, sensor := range sensors {
			pos := sensor.Pos
			coverageArea := Circle{pos, sensor.DistanceToClosetstBeacon()}
			intersect1, intersect2 := coverageArea.IntersectionsWithHorizontalLine(y)

			beaconPos := sensor.ClosestBeacon.Pos
			if beaconPos.Y == y {
				ranges = append(ranges, Line{sensor.ClosestBeacon.Pos, sensor.ClosestBeacon.Pos})
			}

			switch {
			case intersect1 == nil && intersect2 == nil:
				// no intersection
			case intersect2 == nil:
				// tangent
				intersect1.Chop(0, searchRange)
				ranges = append(ranges, Line{*intersect1, *intersect1})
			default:
				// full intersection
				intersect1.Chop(0, searchRange)
				intersect2.Chop(0, searchRange)
				ranges = append(ranges, Line{*intersect1, *intersect2})
			}
		}

		sort.Slice(ranges, func(i, j int) bool {
			return ranges[i].P1.X < ranges[j].P1.X
		})

		included := GetAllIncluded(ranges)
		ranges = make([]Line, 0)
		for _, include := range included {
			ranges = append(ranges, include[0])
		}

		coverage := 0
		rangeLen := len(ranges)
		last := ranges[rangeLen-1]
		for i, line := range ranges[:rangeLen-1] {
			switch next := ranges[i+1]; {
			case next.P1.X < line.P2.X:
				coverage += line.Length() - line.P2.Diff(next.P1).X
			default:
				coverage += line.Length()
			}
		}
		coverage += last.Length()
		coverage -= rangeLen

		if coverage != searchRange-2 {
			continue
		}

		var p Point
		for i := 1; i < len(ranges); i++ {
			prev, cur := ranges[i-1], ranges[i]
			if d := cur.P1.Diff(prev.P2); d.X < 2 {
				continue
			}
			p = Point{prev.P2.X + 1, y}
			break
		}
		fmt.Println("tuning frequenzy:", p.X*4000000+p.Y)
		break
	}
}
