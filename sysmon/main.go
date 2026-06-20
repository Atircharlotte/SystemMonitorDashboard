package main

import (
	"fmt"
	"os"
	"sysmon/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// fmt.Println("Starting System Monitor... ")

	// // fetched structured data
	// memData := monitor.GetMemory()
	
	// // The OS gives us kilobytes. Convert to megabytes
	// totalMB := memData.Total / 1024
	// availableMB := memData.Available / 1024
	// usedMB := totalMB - availableMB

	// fmt.Printf("Total Memory: %dMB\n", totalMB)
	// fmt.Printf("Available Memory: %dMB\n", availableMB)
	// fmt.Printf("Used Memory: %dMB\n", usedMB)



	// fmt.Println("Calculating CPU Usage...")
	// // get the pipe
	// cpuStream := monitor.StartCPUWorker()

	// fmt.Println("Waiting for CPU data...")

	// // listen to pipe forever
	// for {
	// 	currentCPU := <-cpuStream
	// 	fmt.Printf("Current CPU Usage: %.2f%%\n", currentCPU)
	// }

	// fmt.Println("Calculating network speed...")
	// networkStream := monitor.StartNetworkWorker()
	// for {
	// 	currentNetworkSpeed := <- networkStream
	// 	fmt.Printf("Current Network RxSpeed speed: %.2f KB/s\n", currentNetworkSpeed.RxSpeedKB)
	// 	fmt.Printf("Current Network TxSpeed speed: %.2f KB/s\n", currentNetworkSpeed.TxSpeedKB)
	// }

	// fmt.Println("Calculating disk soeed...")
	// diskStream := monitor.StartDiskWorker()
	// for {
	// 	currentDiskSpeed := <- diskStream
	// 	fmt.Printf("Current Disk Reading speed: %.2f KB/s\n", currentDiskSpeed.ReadSpeedKB)
	// 	fmt.Printf("Current Disk Writing speed: %.2f KB/s\n", currentDiskSpeed.WriteSpeedKB)
	// }


	// 1. Create the initial dashboard state
	dashboard := ui.InitialModel()
	
	// 2. Initialize the Bubble Tea program with the dashboard
	p := tea.NewProgram(dashboard)

	// 3. run the program (this will take over the erminal screen)
	if _, err := p.Run(); err != nil {
		fmt.Println("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	

}