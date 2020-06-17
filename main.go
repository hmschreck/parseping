package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
)


var pingRegex = regexp.MustCompile(`.*icmp_seq=([0-9]+) ttl=[0-9]+ time=([0-9\.]+) ms`)
var filename = flag.String("filename", "dummy", "filename to parse")

func main() {
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("couldn't open file %v", err)
	}
	defer file.Close()
	sequences := []int{}
	times := []float64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result := pingRegex.FindSubmatch([]byte(scanner.Text()))
		if len(result) < 2 {
			continue
		}
		icmpSeq, err := strconv.Atoi(string(result[1]))
		if err != nil {
			log.Errorf("couldn't convert to int %v", err)
		}
		time, err := strconv.ParseFloat(string(result[2]), 64)
		if err != nil {
			log.Errorf("couldn't convert to float64 %v", err)
			break
		}
		sequences = append(sequences, icmpSeq)
		times = append(times, time)
	}
	// last line should be the total attempts, and the number of actual successes should be the count of lines
	missedPackets := sequences[len(sequences)-1] - len(sequences)
	fmt.Println("Total missed packets: %d", missedPackets)
	misses := make(map[int]int)
	previousSequence := 0
	for _, sequence := range sequences {
		gap := sequence - previousSequence - 1
		previousSequence = sequence
		if _, ok := misses[gap]; !ok {
			misses[gap] = 1
		} else {
			misses[gap] += 1
		}
	}
	fmt.Printf("%v\n", misses)
}
