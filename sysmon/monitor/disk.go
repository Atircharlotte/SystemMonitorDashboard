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

// Helper function to check if a device is a primary physical disk
func isPhysicalDisk(name string) bool {
	// Skip loopback virtual environments
	if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") {
		return false
	}

	// If it's a standard SCSI/SATA drive (sda, sdb), it must be exactly 3 characters long
	// This will catch "sda" but drop "sda1" or "sda2"
	if strings.HasPrefix(name, "sda") && len(name) == 3 {
		return true
	} 

	// if it's an NVMe drive (nvme0n1), I want the drive but not the partion (nvme0n1p1)
	if strings.HasPrefix(name, "nvme") && !strings.Contains(name, "p") {
		return true
	}
	return false
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
		fields := strings.Fields(scanner.Text())

		// SAFETY CHECK: Ensure the line has the standard number of fields
		if len(fields) < 14 {
			continue
		}

		deviceName := fields[2]

		// only accumulate stats if it's a root physical drive
		if isPhysicalDisk(deviceName) {
			// fields[2]: Device name
			// fields[5]: Sectors Read
			// fields[9]: Sections Written
			readSectors, _ := strconv.ParseUint(fields[5], 10, 64)
			writeSectors, _ := strconv.ParseUint(fields[9], 10, 64)

			snapshot.TotalReadSectors += readSectors
			snapshot.TotalWriteSectors += writeSectors
		}
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

			// Optimized Math: Sectors / 2.0 is exactly equal to (Sectors * 512) / 1024
			speed.ReadSpeedKB = float64(snapshot2.TotalReadSectors - snapshot1.TotalReadSectors) / 2.0
			speed.WriteSpeedKB = float64(snapshot2.TotalWriteSectors - snapshot1.TotalWriteSectors) / 2.0

			diskChannel <- speed
		}
	}()
	return diskChannel
}