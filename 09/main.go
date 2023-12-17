package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
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
	lines := [][]int{}
	for _, line := range arr {
		strs := strings.Split(strings.TrimSpace(line), " ")
		line := []int{}
		for _, s := range strs {
			line = append(line, castInt(s))
		}
		lines = append(lines, line)
	}
	// infof("lines:\n%+v", lines)

	firstSum, lastSum := 0, 0
	for _, l := range lines {
		line := l
		infof("line: %+v", line)
		firsts := []int{line[0]}
		lasts := []int{line[len(line)-1]}
		for {
			diffs := map[int]int{}
			newLine := []int{}
			for i := 0; i < len(line)-1; i++ {
				// infof("i: %d", i)
				// infof("line: %+v", line)
				// infof("v: %d", line[i])
				a, b := line[i], line[i+1]
				d := b - a
				diffs[d] += 1
				newLine = append(newLine, d)
			}
			// infof("diffs: %+v", newLine)
			firsts = append(firsts, newLine[0])
			lasts = append(lasts, newLine[len(newLine)-1])
			if len(diffs) == 1 && diffs[0] > 0 {
				break
			}
			line = newLine
		}
		// infof("lasts: %+v", lasts)
		finalValue := 0
		slices.Reverse(lasts)
		for _, v := range lasts {
			finalValue += v
		}
		infof("finalValue: %d", finalValue)
		lastSum += finalValue

		infof("firsts: %+v", firsts)
		firstValue := 0
		slices.Reverse(firsts)
		for _, v := range firsts {
			firstValue = v - firstValue
		}
		infof("firstValue: %d", firstValue)
		firstSum += firstValue
	}

	infof("firstSum: %d", firstSum)
	infof("lastSum: %d", lastSum)
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
