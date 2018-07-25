package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	for {
		cpuUsage, err := calcCPUUsage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("CPUUsage:\t%6.3f %%\n", cpuUsage)

		mi, err := memInfo()
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range []string{"MemTotal", "MemFree", "SwapTotal", "SwapFree"} {
			fmt.Printf("%s:\t %s\n", s, mi[s])
		}
		fmt.Println("------------------------------------------------")

		// TODO probably should use a ticker
		time.Sleep(time.Second)
	}
}

var (
	prevIdleTime  uint64
	prevTotalTime uint64
)

// https://rosettacode.org/wiki/Linux_CPU_utilization#Go
func calcCPUUsage() (cpuUsage float64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return cpuUsage, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()[5:] // get rid of cpu plus 2 spaces
	file.Close()
	if err := scanner.Err(); err != nil {
		return cpuUsage, err
	}
	split := strings.Fields(firstLine)
	idleTime, _ := strconv.ParseUint(split[3], 10, 64)
	totalTime := uint64(0)
	for _, s := range split {
		u, _ := strconv.ParseUint(s, 10, 64)
		totalTime += u
	}
	if prevIdleTime != 0 && prevTotalTime != 0 {
		deltaIdleTime := idleTime - prevIdleTime
		deltaTotalTime := totalTime - prevTotalTime
		cpuUsage = (1.0 - float64(deltaIdleTime)/float64(deltaTotalTime)) * 100.0
	}
	prevIdleTime = idleTime
	prevTotalTime = totalTime

	return cpuUsage, nil
}

func memInfo() (map[string]string, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)

	mi := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		if kv := strings.SplitN(line, ":", 2); len(kv) > 1 {
			k := strings.TrimSpace(kv[0])
			mi[k] = strings.TrimSpace(kv[1])
		}
	}
	file.Close()
	return mi, nil
}
