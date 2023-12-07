package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}

func infof(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
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
	times := []int{}
	for _, t := range strings.Split(strings.Split(arr[0], ":")[1], " ") {
		if t == "" {
			continue
		}
		times = append(times, castInt(t))
	}
	distances := []int{}
	for _, d := range strings.Split(strings.Split(arr[1], ":")[1], " ") {
		if d == "" {
			continue
		}
		distances = append(distances, castInt(d))
	}

	product := 1
	for i := 0; i < len(times); i++ {
		t, d := times[i], distances[i]
		infof("time: %d, dist: %d", t, d)
		// max, maxCount := 0, 0
		winWays := 0
		for hold := 0; hold < t; hold++ {
			speed := hold
			remTime := t - hold
			dist := remTime * speed
			// infof("hold for %d; speed will be %d; dist will be %d", hold, speed, dist)

			if dist < d {
				continue
			}

			switch {
			case dist > d:
				// infof("new max: %d", dist)
				// max = dist
				// maxCount = 1
				winWays++
			}
		}
		// infof("max: %d; maxCount: %d", max, maxCount)
		infof("winWays: %d", winWays)
		product *= winWays
	}
	infof("product: %d", product)
}

func NewRangeMap() *RangeMap {
	return &RangeMap{
		ranges: map[*Range]int{},
	}
}

type RangeMap struct {
	ranges map[*Range]int
	l      sync.RWMutex
}

func (m *RangeMap) Add(r *Range, o int) {
	m.l.Lock()
	defer m.l.Unlock()
	m.ranges[r] = o
}

func (m *RangeMap) Get(i int) (int, *Range, bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	for k, v := range m.ranges {
		if i >= k.start && i <= k.end {
			return v, k, true
		}
	}
	return 0, nil, false
}

func (m *RangeMap) GetOffsettedValue(v int) int {
	off, r, ok := m.Get(v)
	if !ok {
		return v
	}
	dist := v - r.start
	offsetValue := off + dist
	return offsetValue
}

func (m *RangeMap) AllKeys() (chan int, int) {
	ranges := []Range{}
	for k := range m.ranges {
		ranges = append(ranges, *k)
	}

	sort.SliceStable(ranges, func(i, j int) bool {
		a, b := ranges[i], ranges[j]
		if a.start == b.start {
			return a.end < b.end
		}
		return a.start < b.start
	})

	combinedRanges := []Range{}
	for i := 0; i < len(ranges); i++ {
		start := ranges[i].start
		end := ranges[i].end
		for {
			if i+1 == len(ranges) {
				break
			}
			if ranges[i+1].start-1 <= end {
				end = ranges[i+1].end
				i++
			} else {
				break
			}
		}
		combinedRanges = append(combinedRanges, Range{start: start, end: end})
	}

	sum := 0
	infof("ranges:")
	for _, r := range combinedRanges {
		sum += r.end - r.start + 1
		infof("\t%d-%d", r.start, r.end)
	}
	infof("")

	c := make(chan int, 500)
	count := 0
	go func() {
		defer close(c)
		for _, r := range combinedRanges {
			for i := r.start; i <= r.end; i++ {
				count++
				if count%1e6 == 0 {
					infof("AllKeys: %d of %d (%f percent)", count, sum, float64(count)/float64(sum)*100)
				}
				c <- i
			}
		}
	}()
	return c, sum
}

type Range struct {
	start int // inclusive.
	end   int // inclusive.
}

func castInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fatalf("castInt %q: %v", s, err)
	}
	return i
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
