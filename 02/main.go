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

var number = regexp.MustCompile(`[0-9]+`)

func main() {
	sourcePtr := flag.String("source", "./sample.txt", "input file")
	flag.Parse()

	buff, err := os.ReadFile(*sourcePtr)
	if err != nil {
		fatalf("Error: %v", err)
	}
	input := string(buff)
	input = strings.TrimSpace(input)

	// maxs := map[string]int{
	// 	"red":   12,
	// 	"green": 13,
	// 	"blue":  14,
	// }

	sum := 0
	arr := strings.Split(input, "\n")
	for _, l := range arr {
		g := strings.Split(l, ":")
		gameNum, _ := strconv.Atoi(number.FindStringSubmatch(g[0])[0])
		fmt.Printf("gameNum: %d\n", gameNum)
		maxs := map[string]int{}
		for _, r := range strings.Split(strings.TrimSpace(g[1]), ";") {
			game := map[string]int{}
			for _, set := range strings.Split(strings.TrimSpace(r), ",") {
				rnd := strings.Split(strings.TrimSpace(set), " ")
				c, _ := strconv.Atoi(number.FindStringSubmatch(rnd[0])[0])
				game[rnd[1]] += c
			}
			for color, v := range game {
				if maxs[color] < v {
					maxs[color] = v
				}
			}
			fmt.Printf("game %d: %+v\n", gameNum, game)
		}
		sum += (maxs["red"] * maxs["green"] * maxs["blue"])
	}

	fmt.Printf("sum: %d\n", sum)
}
