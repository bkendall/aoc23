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

type Choice struct {
	Left, Right string
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
	instructions := arr[0]

	mapping := map[string]*Choice{}
	for _, line := range arr[2:] {
		pts := strings.Split(line, " = ")
		k := strings.TrimSpace(pts[0])
		cs := strings.Split(pts[1], ", ")
		left := strings.TrimPrefix(cs[0], "(")
		right := strings.TrimSuffix(cs[1], ")")
		mapping[k] = &Choice{Left: left, Right: right}
	}

	curr := map[int]string{}
	stepsPerStart := map[string]int{}
	for k := range mapping {
		if string(k[2]) == "A" {
			curr[len(curr)-1] = k
		}
	}
	steps := 0
	for i := 0; i < len(instructions); i++ {
		s := string(instructions[i])
		for i, v := range curr {
			// infof("@ %s, going %s", v, s)
			switch s {
			case "L":
				curr[i] = mapping[v].Left
			case "R":
				curr[i] = mapping[v].Right
			}
		}
		steps++
		for k, v := range curr {
			if string(v[2]) == "Z" {
				stepsPerStart[v] = steps
				delete(curr, k)
				infof("%+v", stepsPerStart)
			}
		}
		if len(curr) == 0 {
			break
		}
		if i == len(instructions)-1 {
			i = -1
		}
	}

	nums := []int{}
	for _, v := range stepsPerStart {
		nums = append(nums, v)
	}

	infof("%v", nums)
	infof("LCM: %d", LCM(nums[0], nums[1], nums[2:]...))
}

// Thanks https://go.dev/play/p/SmzvkDjYlb:
// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
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
