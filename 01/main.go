package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

var number = regexp.MustCompile(`[0-9]`)

var replacements = map[string]string{
	"one":   "o1e",
	"two":   "t2o",
	"three": "t3e",
	"four":  "f4r",
	"five":  "f5e",
	"six":   "s6x",
	"seven": "s7n",
	"eight": "e8t",
	"nine":  "n9e",
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

	sum := 0
	arr := strings.Split(input, "\n")
	// for _, l := range arr {
	// 	all := number.FindAllStringSubmatch(l, -1)
	// 	fmt.Printf("all: %q\n", all)
	// 	first := all[0][0]
	// 	last := all[len(all)-1][0]
	// 	i, err := strconv.Atoi(first + last)
	// 	if err != nil {
	// 		fatalf("%v", err)
	// 	}
	// 	fmt.Printf("i: %d\n", i)
	// 	sum += i
	// }

	// fmt.Printf("sum: %d\n", sum)

	sum = 0
	for _, l := range arr {
		fmt.Printf("line: %s\n", l)
		for {
			earliestMatch, index := "", len(l)
			for k := range replacements {
				i := strings.Index(l, k)
				if i != -1 && i < index {
					earliestMatch = k
					index = i
				}
			}
			if index == len(l) {
				// fmt.Printf("break\n")
				break
			}
			l = strings.Replace(l, earliestMatch, replacements[earliestMatch], 1)
			// fmt.Printf("fixing l: %q\n", l)
		}

		fmt.Printf("new line: %s\n", l)
		all := number.FindAllStringSubmatch(l, -1)
		// fmt.Printf("all: %q\n", all)
		first := all[0][0]
		last := all[len(all)-1][0]
		i, err := strconv.Atoi(first + last)
		if err != nil {
			fatalf("%v", err)
		}
		fmt.Printf("i: %d\n", i)
		sum += i
	}

	fmt.Printf("sum: %d\n", sum)
}
