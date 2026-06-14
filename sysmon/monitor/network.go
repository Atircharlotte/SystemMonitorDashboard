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

// Represents the raw bytes at a specific moment in time
type NetworkSnapshot struct {
	TotalRxBytes int // Total Downloaded
	TotalTxBytes int // Total Uploaded
}
// Represents the calculated speed (what u will send down to the channel) 
type NetworkSpeed struct {
	RxSpeedMB float64 // download spped in MB/s
	TxSpeedMB float64 // upload speed in MB/s
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
	for i := 0; i < 2; i++ {
		if !scanner.Scan() {
			break
		}
	}

	// 2. Tle colon trap
	// The Fix: Don't use Fields right away. Use strings.Split(line, ":") first.
	for scanner.Scan() {
		line := scanner.Text()
		devidedParts := strings.Split(line, ":")
		fields := strings.Fields(devidedParts[1])

		// The Indices: 
		// Download (RX Bytes): This is at index 0. 
		// Upload (TX Bytes): This is at index 8.
		rxBytes, _ := strconv.Atoi(fields[0])
		txBytes, _ := strconv.Atoi(fields[8])
		snapshot.TotalRxBytes += rxBytes
		snapshot.TotalTxBytes += txBytes
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("file parsing error: ", err)
	}
	fmt.Printf("totalRxBytes: %d\n", snapshot.TotalRxBytes)
	fmt.Printf("totalTxBytes: %d\n", snapshot.TotalTxBytes)
	return snapshot
}

func StartNetworkWorker() <- chan NetworkSpeed {

	// 1. create channel (pipe)
	networkChannel := make(chan NetworkSpeed)

	// launch the background worker
	go func() {
		for {
			var speed NetworkSpeed
			// 1. snapshot 1
			snap1 := ReadNetworkSnapshot()

			// wait for 1 second
			time.Sleep(1 * time.Second)

			// 2. snapshot 2
			snap2 := ReadNetworkSnapshot()

			// calculate the speed
			speed.RxSpeedMB = float64(snap2.TotalRxBytes - snap1.TotalRxBytes) / (1024*1024)
			speed.TxSpeedMB = float64(snap2.TotalTxBytes - snap1.TotalTxBytes) / (1024*1024)

			networkChannel <- speed
		}
	}()

	return networkChannel
}