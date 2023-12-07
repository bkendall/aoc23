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

	arr := strings.Split(input, "\n\n")

	seedStrings := strings.Split(strings.TrimSpace(arr[0]), " ")[1:]
	seeds := []*Range{}
	// for _, s := range seedStrings {
	// 	seeds = append(seeds, castInt(s))
	// }
	seedRanges := NewRangeMap()
	for i := 0; i < len(seedStrings); i += 2 {
		start := castInt(seedStrings[i])
		r := castInt(seedStrings[i+1])
		seeds = append(seeds, &Range{start: start, end: start + r})
		seedRanges.Add(&Range{start: start, end: start + r}, 1)
	}

	seedToSoil := NewRangeMap()
	soilToFertilizer := NewRangeMap()
	fertilizerToWater := NewRangeMap()
	waterToLight := NewRangeMap()
	lightToTemperature := NewRangeMap()
	temperatureToHumidity := NewRangeMap()
	humidityToLocation := NewRangeMap()

	for i, m := range []*RangeMap{seedToSoil, soilToFertilizer, fertilizerToWater, waterToLight, lightToTemperature, temperatureToHumidity, humidityToLocation} {
		lines := strings.Split(strings.TrimSpace(arr[i+1]), "\n")[1:]
		for _, l := range lines {
			ns := strings.Split(strings.TrimSpace(l), " ")
			dest, sourceStart, length := castInt(ns[0]), castInt(ns[1]), castInt(ns[2])
			m.Add(&Range{start: sourceStart, end: sourceStart + length}, dest)
			// infof("[%d] %d %d %d (%d)", i, dest, sourceStart, length, sourceStart+length)
		}
	}

	lowest := math.MaxInt
	lowestLock := sync.Mutex{}
	wg := sync.WaitGroup{}

	workers := 32
	seedChan := make(chan int, 500)
	seedRanges.AllKeys(seedChan)

	worker := func(id int, seeds <-chan int) {
		defer wg.Done()
		count := 0
		localLowest := math.MaxInt64
		for s := range seeds {
			count++
			v := s
			for _, m := range []*RangeMap{seedToSoil, soilToFertilizer, fertilizerToWater, waterToLight, lightToTemperature, temperatureToHumidity, humidityToLocation} {
				v = m.GetOffsettedValue(v)
			}
			// infof("[%d] seed %d: location: %d", id, s, v)
			if v < localLowest {
				localLowest = v
			}
		}
		lowestLock.Lock()
		defer lowestLock.Unlock()
		if localLowest < lowest {
			lowest = localLowest
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i, seedChan)
	}
	wg.Wait()
	infof("lowest: %v", lowest)
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

func (m *RangeMap) AllKeys(c chan<- int) int {
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
	return sum
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
