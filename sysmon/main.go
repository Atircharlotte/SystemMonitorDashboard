package main

import (
	"fmt"
	"sysmon/monitor"
)

func main() {
	fmt.Println("Starting System Monitor... ")

	// fetched structured data
	memData := monitor.GetMemory()
	
	// The OS gives us kilobytes. Convert to megabytes
	totalMB := memData.Total / 1024
	availableMB := memData.Available / 1024
	usedMB := totalMB - availableMB

	fmt.Printf("Total Memory: %dMB\n", totalMB)
	fmt.Printf("Available Memory: %dMB\n", availableMB)
	fmt.Printf("Used Memory: %dMB\n", usedMB)
}