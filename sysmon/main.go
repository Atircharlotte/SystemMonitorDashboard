package main

import (
	"fmt"
	"os"
	"sysmon/monitor"
	"sysmon/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	
	// 1. Start he background OS Workers
	cpuStream := monitor.StartCPUWorker()
	networkStream := monitor.StartNetworkWorker()
	diskStream := monitor.StartDiskWorker()


	// 2. Initialize the dashboard with the stream
	dashboard := ui.InitialModel(cpuStream, networkStream, diskStream)

	// 3. Start the bubble tea program
	p := tea.NewProgram(dashboard)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}