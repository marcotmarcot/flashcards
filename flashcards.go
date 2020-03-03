package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cs := cards()
	ps := probs(len(cs))
	for true {
		i := index(ps)
		fmt.Println(cs[i].front)
		var input string
		fmt.Scanln(&input)
		if input == "q" {
			break
		}
		fmt.Println(cs[i].back)
		if right() {
			ps[i] /= 2
		} else {
			ps[i] *= 2
		}
	}
	f, err := os.Create("probs.tsv")
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	for _, p := range ps {
		fmt.Fprintln(w, p)
	}
	w.Flush()
}

func cards() []*card {
	var cs []*card
	bytes, err := ioutil.ReadFile("cards.tsv")
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		fields := strings.Split(line, "\t")
		if len(fields) < 2 {
			continue
		}
		cs = append(cs, &card{fields[0], fields[1]})
	}
	return cs
}

type card struct {
	front, back string
}

func probs(size int) []float64 {
	var ps []float64
	bytes, err := ioutil.ReadFile("probs.tsv")
	if err != nil {
		ps = make([]float64, size)
		for i := range ps {
			ps[i] = 1
		}
		return ps
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		p, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		ps = append(ps, p)
	}
	if len(ps) != size {
		log.Fatal(fmt.Sprintf("%v != %v", len(ps), + size))
	}
	return ps
}

func index(ps []float64) int {
	total := 0.0
	for _, p := range ps {
		total += p
	}
	chosen := rand.Float64() * total
	for i, p := range ps {
		chosen -= p
		if chosen < 0 {
			return i
		}
	}
	return -1
}

func right() bool {
	var input string
	for true {
		fmt.Println("Right? [y/n]")
		fmt.Scanln(&input)
		if input == "y" || input == "n" {
			break
		}
	}
	return input == "y"
}
