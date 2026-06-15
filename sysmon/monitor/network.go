package monitor

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Represents the raw bytes at a specific moment in time
type NetworkSnapshot struct {
	TotalRxBytes uint64 // Total Downloaded (changed to uint 64 for safety)
	TotalTxBytes uint64 // Total Uploaded
}
// Represents the calculated speed (what u will send down to the channel) 
type NetworkSpeed struct {
	RxSpeedKB float64 // download spped in MB/s (changed name to KB)
	TxSpeedKB float64 // upload speed in MB/s
}

func ReadNetworkSnapshot() NetworkSnapshot {

	file, err := os.Open("/proc/net/dev")
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close()
	
	var snapshot NetworkSnapshot
	scanner := bufio.NewScanner(file)

	// 1. The header skip: skip two header lines
	scanner.Scan()
	scanner.Scan()

	// 2. Tle colon trap
	// The Fix: Don't use Fields right away. Use strings.Split(line, ":") first.
	for scanner.Scan() {
		line := scanner.Text()

		dividedParts := strings.Split(line, ":")

		// SAFETY CHECK: Ensure the line actually had a colon
		if len(dividedParts) < 2 {
			continue
		}
		fields := strings.Fields(dividedParts[1])

		// The Indices: 
		// Download (RX Bytes): This is at index 0. 
		// Upload (TX Bytes): This is at index 8.
		// SAFETY CHECK: ensure there is enough data 
		if len(fields) >= 9 {
			// ParseUnit is the standard for parsing large OS byte counters
			rxBytes, _ := strconv.ParseUint(fields[0], 10, 64)
			txBytes, _ := strconv.ParseUint(fields[8], 10, 64)

			snapshot.TotalRxBytes += rxBytes
			snapshot.TotalTxBytes += txBytes
		}
	}
	return snapshot
}

func StartNetworkWorker() <- chan NetworkSpeed {

	// 1. create channel (pipe)
	networkChannel := make(chan NetworkSpeed)

	// launch the background worker
	go func() {
		for {
			// 1. snapshot 1
			snap1 := ReadNetworkSnapshot()

			// wait for 1 second
			time.Sleep(1 * time.Second)

			// 2. snapshot 2
			snap2 := ReadNetworkSnapshot()

			var speed NetworkSpeed

			// calculate the speed
			// calculate KB/s: (divided by 1024.0)
			speed.RxSpeedKB = float64(snap2.TotalRxBytes - snap1.TotalRxBytes) / (1024)
			speed.TxSpeedKB = float64(snap2.TotalTxBytes - snap1.TotalTxBytes) / (1024)

			networkChannel <- speed
		}
	}()

	return networkChannel
}