package monitor

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type DiskSnapshot struct {
	TotalReadSectors uint64 
	TotalWriteSectors uint64
}
type DiskSpeed struct {
	ReadSpeedKB float64
	WriteSpeedKB float64
}

func ReadDiskSnapshot() DiskSnapshot {
	
	// 1. read file "/proc/diskstats"
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close()

	// 2. Read the content of the file
	scanner := bufio.NewScanner(file)
	var snapshot DiskSnapshot
	
	for scanner.Scan() {
		line := scanner.Text() // read each line of the file
		fields := strings.Fields(line)
		deviceName := fields[2]
		if (strings.HasPrefix(deviceName, "loop")) {
			continue
		}
		// fields[2]: Device name
		// fields[5]: Sectors Read
		// fields[9]: Sections Written
		readSectors, _ := strconv.ParseUint(fields[5], 10, 64)
		writeSectors, _ := strconv.ParseUint(fields[9], 10, 64)
		snapshot.TotalReadSectors += readSectors
		snapshot.TotalWriteSectors += writeSectors
	}
	return snapshot
}

func StartDiskWorker() <- chan DiskSpeed {

	diskChannel := make(chan DiskSpeed)

	// run background anonymous function
	go func() {
		for {
			// snapshot1
			snapshot1 := ReadDiskSnapshot()
			
			// wait for one second
			time.Sleep(1 * time.Second)

			// snapshot2
			snapshot2 := ReadDiskSnapshot()

			var speed DiskSpeed
			speed.ReadSpeedKB = float64(snapshot2.TotalReadSectors - snapshot1.TotalReadSectors)*512 / 1024
			speed.WriteSpeedKB = float64(snapshot2.TotalWriteSectors - snapshot1.TotalWriteSectors)*512 / 1024

			diskChannel <- speed
		}
	}()
	return diskChannel
}