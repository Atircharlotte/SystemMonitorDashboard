package monitor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// type CPUUnit struct {
// 	usedUnits int
// 	idleUnits int
// }

// Renamed to Snapshot because it represents a single moment in time
type CPUSnapshot struct {
	Total int
	Idle  int
}

func ReadCPUSnapshot() CPUSnapshot {
	
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close() // CRITICAL: This prevents the memory/file leak!

	// var AllCPUUnits CPUUnit;
	var snapshot CPUSnapshot
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// split the line into words
		fields := strings.Fields(line)

		if len(fields) > 0 && fields[0] == "cpu" {
			total := 0
			idle := 0

			// loop through all 10 time columns
			for j := 1; j < len(fields); j++ {
				val, _ := strconv.Atoi(fields[j])
				total += val
				if j == 4 || j == 5 {
					idle += val
				}
			}
			snapshot.Total = total
			snapshot.Idle = idle
			break // only care about the first "cpu" line
		}
	}
	return snapshot
}


func StartCPUWorker() <- chan float64 {
	
	// 1. create channel (pipe)
	cpuChannel := make(chan float64)

	// launch the background worker (anonymous function)
	go func() {
		for {
			// TODO: read /proc/stat (snapshot 1)
			snap1 := ReadCPUSnapshot()
			
			// wait for 1 second
			time.Sleep(1 * time.Second)
			// TODO: read /proc/stat (snapshot 2)
			snap2 := ReadCPUSnapshot()

			totalDelta := snap2.Total - snap1.Total
			idleDelta := snap2.Idle - snap1.Idle

			fmt.Printf("totalDelta: %d\n", totalDelta);
			fmt.Printf("idleDelta: %d\n", idleDelta);
			// calculate 
			calculatedPercentage := (float64(totalDelta - idleDelta) / float64(totalDelta)) * 100.0 

			// send the result down the pipe
			cpuChannel <- calculatedPercentage
		}
	}()

	return cpuChannel
}