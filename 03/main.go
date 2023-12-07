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
	y := 0
	for ; y < len(arr); y++ {
		nbrArr := numbers.FindAllStringSubmatchIndex(arr[y], -1)
		// infof("nbrArr: %+v", nbrArr)

		for _, mtch := range nbrArr {
			start, end := mtch[0], mtch[1]
			use := false
			if y > 0 {
				use = use || hasSymbol(arr[y-1], max(0, start-1), min(len(arr[0]), end+1))
			}
			use = use || hasSymbol(arr[y], max(0, start-1), min(len(arr[0]), end+1))
			if y < len(arr)-1 {
				use = use || hasSymbol(arr[y+1], max(0, start-1), min(len(arr[0]), end+1))
			}

			n, _ := strconv.Atoi(arr[y][start:end])
			if use {
				// infof("using: %d", n)
				sum += n
			} else {
				// infof("** not using: %d", n)
			}
		}
	}

	infof("sum: %d", sum)

	// Gear ratios

	sum = 0
	y = 0
	for ; y < len(arr); y++ {
		l := arr[y]
		for x := 0; x < len(l); x++ {
			i := strings.Index(l[x:], "*")
			// infof("indx? %v", i)
			if i == -1 {
				break
			}
			x = x + i
			ns := findNumbers(arr, x, y)
			// infof("ns @ %d, %d: %+v", x, y, ns)
			if len(ns) == 2 {
				sum += ns[0] * ns[1]
			}
		}
	}

	infof("sum: %d", sum)
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
