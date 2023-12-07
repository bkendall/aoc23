package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func infof(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

var numbers = regexp.MustCompile(`[0-9]+`)
var onlyNumbers = regexp.MustCompile(`^[0-9]+$`)

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		fatalf("Error: %v", err)
	}
	input := string(buff)
	input = strings.TrimSpace(input)

	sum := 0
	arr := strings.Split(input, "\n")
	cardCounts := map[int]int{}
	for i := range arr {
		cardCounts[i] = 1
	}
	for i, c := range arr {
		nCards := cardCounts[i]
		data := strings.Split(c, ":")
		ns := strings.Split(strings.TrimSpace(data[1]), "|")
		winners := strings.Split(strings.TrimSpace(ns[0]), " ")
		numbers := strings.Split(strings.TrimSpace(ns[1]), " ")

		nMap := map[string]bool{}
		for _, n := range numbers {
			if strings.TrimSpace(n) == "" {
				continue
			}
			nMap[strings.TrimSpace(n)] = true
		}

		wMap := map[string]bool{}
		for _, w := range winners {
			if strings.TrimSpace(w) == "" {
				continue
			}
			wMap[strings.TrimSpace(w)] = true
		}

		have := 0
		haveCount := 0
		for n := range nMap {
			if wMap[n] {
				haveCount++
				if have == 0 {
					have = 1
				} else {
					have *= 2
				}
			}
		}
		for x := i + 1; x < i+1+haveCount; x++ {
			cardCounts[x] += nCards
		}
		sum += have
	}

	infof("sum: %d", sum)

	cards := 0
	for _, t := range cardCounts {
		cards += t
	}
	infof("card total: %d", cards)
}

var notNumberOrDot = regexp.MustCompile(`[^0-9\.]`)

func hasSymbol(l string, x1, x2 int) bool {
	sub := l[x1:x2]
	loc := notNumberOrDot.FindStringIndex(sub)
	return loc != nil
}

func findNumbers(arr []string, x, y int) []int {
	ns := []int{}
	if y > 0 {
		l := arr[y-1]
		one, two := false, false
		if x > 0 {
			if onlyNumbers.MatchString(l[x-1 : x]) {
				ns = append(ns, wholeNumber(l, x-1))
				one = true
			}
		}
		if onlyNumbers.MatchString(l[x : x+1]) {
			if !one {
				ns = append(ns, wholeNumber(l, x))
			}
			two = true
		}
		if x < len(l)-1 {
			if onlyNumbers.MatchString(l[x+1 : x+2]) {
				if !two {
					ns = append(ns, wholeNumber(l, x+1))
				}
			}
		}
	}
	l := arr[y]
	one, two := false, false
	if x > 0 {
		if onlyNumbers.MatchString(l[x-1 : x]) {
			ns = append(ns, wholeNumber(l, x-1))
			one = true
		}
	}
	if onlyNumbers.MatchString(l[x : x+1]) {
		if !one {
			ns = append(ns, wholeNumber(l, x))
		}
		two = true
	}
	if x < len(l)-1 {
		if onlyNumbers.MatchString(l[x+1 : x+2]) {
			if !two {
				ns = append(ns, wholeNumber(l, x+1))
			}
		}
	}
	if y < len(arr)-1 {
		l := arr[y+1]
		one, two := false, false
		if x > 0 {
			if onlyNumbers.MatchString(l[x-1 : x]) {
				ns = append(ns, wholeNumber(l, x-1))
				one = true
			}
		}
		if onlyNumbers.MatchString(l[x : x+1]) {
			if !one {
				ns = append(ns, wholeNumber(l, x))
			}
			two = true
		}
		if x < len(l)-1 {
			if onlyNumbers.MatchString(l[x+1 : x+2]) {
				if !two {
					ns = append(ns, wholeNumber(l, x+1))
				}
			}
		}
	}
	return ns
}

func wholeNumber(l string, x int) int {
	n1, n2 := x, x+1
	for {
		if n1 == 0 {
			break
		}
		if onlyNumbers.MatchString(l[n1-1 : n2]) {
			n1--
		} else {
			break
		}
	}
	for {
		if n2 == len(l) {
			break
		}
		if onlyNumbers.MatchString(l[n1 : n2+1]) {
			n2++
		} else {
			break
		}
	}
	n, _ := strconv.Atoi(l[n1:n2])
	return n
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
