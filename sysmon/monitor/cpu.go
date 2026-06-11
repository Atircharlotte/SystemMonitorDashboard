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
type CPUUnit struct {
	usedUnits int
	idleUnits int
}

func ReadCPUUnit() CPUUnit {
	
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	var AllCPUUnits CPUUnit;
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// split the line into words
		fields := strings.Fields(line)

		if fields[0] == "cpu" {
			usedUnits := 0
			idleUnits := 0
			for j := 1; j < 11; j++ {
				if j != 5 {
					unitToAdd, _ := strconv.Atoi(fields[j])
					usedUnits += unitToAdd
				} else {
					idleUnit, _ := strconv.Atoi(fields[j])
					idleUnits += idleUnit
				}
				AllCPUUnits.usedUnits = usedUnits
				AllCPUUnits.idleUnits = idleUnits
			}
			break
		}
	}
	fmt.Printf("usedUnit: %d\n", AllCPUUnits.usedUnits);
	fmt.Printf("idleUnit: %d\n", AllCPUUnits.idleUnits);
	return AllCPUUnits
}


func StartCPUWorker() <- chan float64 {
	
	// 1. create channel (pipe)
	cpuChannel := make(chan float64)

	// launch the background worker (anonymous function)
	go func() {
		for {
			// TODO: read /proc/stat (snapshot 1)
			unitRecordBefore := ReadCPUUnit()
			
			// wait for 1 second
			time.Sleep(1 * time.Second)
			// TODO: read /proc/stat (snapshot 2)
			unitRecordAfter := ReadCPUUnit()

			totalDelta := unitRecordAfter.usedUnits - unitRecordBefore.usedUnits
			idleDelta := unitRecordAfter.idleUnits - unitRecordBefore.idleUnits

			fmt.Printf("totalDelta: %d\n", totalDelta);
			fmt.Printf("idleDelta: %d\n", idleDelta);
			// calculate 
			calculatedPercentage := (float64(totalDelta - idleDelta) / float64(totalDelta)) * 100 

			// send the result down the pipe
			cpuChannel <- calculatedPercentage
		}
	}()

	return cpuChannel
}