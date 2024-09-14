package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}

func infof(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

type Looped int

const (
	UNKNOWN Looped = iota
	INSIDE
	OUTSIDE
)

type Node struct {
	val    string
	dist   int
	inLoop Looped
}

func NewWorld() *World {
	return &World{
		w: map[int]map[int]*Node{},
	}
}

type World struct {
	w map[int]map[int]*Node
}

func (w *World) LenY() int {
	return len(w.w[0]) - 1
}

func (w *World) LenX() int {
	return len(w.w) - 1
}

func (w *World) Add(x, y int, v *Node) {
	if _, ok := w.w[x]; !ok {
		w.w[x] = map[int]*Node{}
	}
	w.w[x][y] = v
}

func (w *World) Get(x, y int) *Node {
	if _, ok := w.w[x]; !ok {
		w.w[x] = map[int]*Node{}
	}
	return w.w[x][y]
}

type NodeRef struct {
	n    *Node
	x, y int
}

func (w *World) Connected(x, y int) []*NodeRef {
	c := []*NodeRef{}
	if x > 0 {
		c = append(c, &NodeRef{n: w.w[x-1][y], x: x - 1, y: y})
	}
	if y > 0 {
		c = append(c, &NodeRef{n: w.w[x][y-1], x: x, y: y - 1})
	}
	if y < w.LenY() {
		c = append(c, &NodeRef{n: w.w[x][y+1], x: x, y: y + 1})
	}
	if x < w.LenX() {
		c = append(c, &NodeRef{n: w.w[x+1][y], x: x + 1, y: y})
	}
	return c
}

func (w *World) ConnectedUnvisited(x, y int) []*NodeRef {
	c := []*NodeRef{}
	if x > 0 {
		n := w.w[x-1][y]
		v := n.val
		if (v == "-" || v == "L" || v == "F") && n.dist == 0 {
			c = append(c, &NodeRef{n: w.w[x-1][y], x: x - 1, y: y})
		}
	}
	if y > 0 {
		n := w.w[x][y-1]
		v := n.val
		if (v == "|" || v == "7" || v == "F") && n.dist == 0 {
			c = append(c, &NodeRef{n: w.w[x][y-1], x: x, y: y - 1})
		}
	}
	if y < w.LenY() {
		n := w.w[x][y+1]
		v := n.val
		if (v == "|" || v == "L" || v == "J") && n.dist == 0 {
			c = append(c, &NodeRef{n: w.w[x][y+1], x: x, y: y + 1})
		}
	}
	if x < w.LenX() {
		n := w.w[x+1][y]
		v := n.val
		if (v == "-" || v == "7" || v == "J") && n.dist == 0 {
			c = append(c, &NodeRef{n: w.w[x+1][y], x: x + 1, y: y})
		}
	}
	return c
}

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		fatalf("Error: %v", err)
	}
	input := string(buff)
	input = strings.TrimSpace(input)

	arr := strings.Split(input, "\n")
	world := NewWorld()
	startX, startY := 0, 0
	for y, line := range arr {
		for x, s := range strings.Split(line, "") {
			n := &Node{val: s}
			if s == "S" {
				startX, startY = x, y
				n.dist = -1
			}
			world.Add(x, y, n)
		}
	}
	infof("start: %d, %d", startX, startY)

	highest := 0
	next := world.ConnectedUnvisited(startX, startY)
	step := 1
	for {
		if len(next) == 0 {
			break
		}
		more := []*NodeRef{}
		for _, n := range next {
			infof("at %d, %d", n.x, n.y)
			world.Get(n.x, n.y).dist = step
			highest = step
			more = append(more, world.ConnectedUnvisited(n.x, n.y)...)
		}
		next = more
		step++
	}

	infof("last step: %d", highest)

	for y := 0; y <= world.LenY(); y++ {
		for x := 0; x <= world.LenX(); x++ {
			n := world.Get(x, y)
			if n.val != "." || n.inLoop != UNKNOWN {
				continue
			}
			around := world.Connected(x, y)
			if len(around) < 4 {
				n.inLoop = OUTSIDE
			}
		}
	}
}

func castInt(s string) int {
	sign := 1
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = s[1:]
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		fatalf("castInt %q: %v", s, err)
	}
	return sign * i
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
