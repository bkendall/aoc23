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

var values = map[string]int{
	"A": 13,
	"K": 12,
	"Q": 11,
	// "J": 10,
	"T": 9,
	"9": 8,
	"8": 7,
	"7": 6,
	"6": 5,
	"5": 4,
	"4": 3,
	"3": 2,
	"2": 1,
	"J": 0,
}

type HandType int

const (
	highCard HandType = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

type Hand struct {
	cards      string
	cardCounts map[string]int
	handType   HandType
	bid        int
}

func NewHand(str string, bid int) *Hand {
	h := &Hand{
		cards:      str,
		bid:        bid,
		cardCounts: make(map[string]int),
	}
	for _, c := range str {
		h.cardCounts[string(c)] += 1
	}

	h.handType = detectHandType(h)

	return h
}

func detectHandType(h *Hand) HandType {
	switch {
	case isFiveOfAKind(h):
		infof("hand %v is isFiveOfAKind", h.cards)
		return fiveOfAKind
	case isFourOfAKind(h):
		infof("hand %v is isFourOfAKind", h.cards)
		return fourOfAKind
	case isFullHouse(h):
		infof("hand %v is isFullHouse", h.cards)
		return fullHouse
	case isThreeOfAKind(h):
		infof("hand %v is isThreeOfAKind", h.cards)
		return threeOfAKind
	case isTwoPair(h):
		infof("hand %v is isTwoPair", h.cards)
		return twoPair
	case isPair(h):
		infof("hand %v is isPair", h.cards)
		return onePair
	default:
		infof("hand %v is isHighCard", h.cards)
		return highCard
	}
}

func isFiveOfAKind(h *Hand) bool {
	if len(h.cardCounts) == 1 {
		return true
	}
	return mapHas(h.cardCounts, 4) && h.cardCounts["J"] == 1 ||
		mapHas(h.cardCounts, 1) && h.cardCounts["J"] == 4 ||
		mapHas(h.cardCounts, 3) && h.cardCounts["J"] == 2 ||
		mapHas(h.cardCounts, 2) && h.cardCounts["J"] == 3
}

func isFourOfAKind(h *Hand) bool {
	if mapHas(h.cardCounts, 4) && mapHas(h.cardCounts, 1) {
		return true
	}
	if h.cardCounts["J"] == 2 {
		for k, v := range h.cardCounts {
			if k != "J" && v == 2 {
				return true
			}
		}
	}
	return h.cardCounts["J"] == 3 ||
		mapHas(h.cardCounts, 3) && h.cardCounts["J"] == 1

}

func isFullHouse(h *Hand) bool {
	if mapHas(h.cardCounts, 2) && mapHas(h.cardCounts, 3) {
		return true
	}
	if h.cardCounts["J"] == 1 {
		for k, v := range h.cardCounts {
			if k != "J" && v != 2 {
				break
			}
		}
		if mapHas(h.cardCounts, 3) {
			return true
		}
	}
	pairCount := 0
	for k, v := range h.cardCounts {
		if k != "J" && v == 2 {
			pairCount++
		}
	}
	return pairCount == 2 && h.cardCounts["J"] == 1
}

func isThreeOfAKind(h *Hand) bool {
	for _, v := range h.cardCounts {
		if v == 3 {
			return true
		}
	}
	return mapHas(h.cardCounts, 2) && h.cardCounts["J"] == 1 ||
		h.cardCounts["J"] == 2
}

func isTwoPair(h *Hand) bool {
	ones, twos := 0, 0
	for _, v := range h.cardCounts {
		switch v {
		case 1:
			ones++
		case 2:
			twos++
		}
	}
	if ones == 1 && twos == 2 {
		return true
	}
	return mapHas(h.cardCounts, 2) && h.cardCounts["J"] == 1
}

func isPair(h *Hand) bool {
	ones, twos := 0, 0
	for _, v := range h.cardCounts {
		switch v {
		case 1:
			ones++
		case 2:
			twos++
		}
	}
	if ones == 3 && twos == 1 {
		return true
	}
	return h.cardCounts["J"] == 1
}

func mapHas(m map[string]int, v int) bool {
	for _, val := range m {
		if v == val {
			return true
		}
	}
	return false
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
	hands := []*Hand{}
	for _, l := range arr {
		a := strings.Split(l, " ")
		handStr, bid := a[0], castInt(a[1])
		hand := NewHand(handStr, bid)
		// infof("hand: %+v", hand)
		hands = append(hands, hand)
	}

	sort.SliceStable(hands, func(i, j int) bool {
		a, b := hands[i], hands[j]
		if a.handType != b.handType {
			return a.handType < b.handType
		}
		for i := 0; i < len(a.cards); i++ {
			av, bv := values[string(a.cards[i])], values[string(b.cards[i])]
			if av != bv {
				return av < bv
			}
		}
		return false
	})

	winnings := 0
	for i, h := range hands {
		infof("rank %d: %+v", i, h)
		winnings += (i + 1) * h.bid
	}
	infof("winnings: %d", winnings)
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
