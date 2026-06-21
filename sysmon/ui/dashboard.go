package ui

import (
	"fmt"
	"sysmon/monitor"

	tea "github.com/charmbracelet/bubbletea"
)

// 1. Define Custom Message
// create unique types so the Update function knows exactly what data just arrived.
type cpuMsg float64
type netMsg monitor.NetworkSpeed
type diskMsg monitor.DiskSpeed

// 2. The Model
// It holds both the state channels (the pipes) and the current state (the data) 
type Dashboard struct {

	cpuChan <-chan float64
	netChan <-chan monitor.NetworkSpeed
	diskChan <-chan monitor.DiskSpeed

	cpu float64
	net monitor.NetworkSpeed
	disk monitor.DiskSpeed 
} 

// pass the channels in from main.go
func InitialModel(c <-chan float64, n <-chan monitor.NetworkSpeed, d <-chan monitor.DiskSpeed) Dashboard {
	return Dashboard{
		cpuChan: c,
		netChan: n,
		diskChan: d,
	}
}

// 3. The Commands (Background Listener)
// These functions wait for data in the background and return it as a tea.Msg
func listenForCPU(sub <-chan float64) tea.Cmd {
	return func() tea.Msg {
		return cpuMsg(<-sub) // wait for data, then wrap it in the custom type
	}
}

func listenForNet(sub <-chan monitor.NetworkSpeed) tea.Cmd {
	return func() tea.Msg {
		return netMsg(<-sub)
	}
}

func listenForDisk(sub <-chan monitor.DiskSpeed) tea.Cmd {
	return func() tea.Msg {
		return diskMsg(<-sub)
	}
}



// 4. The Init Function
// This runs once. Use tea.Batch() to kick off all these listeners simultaneously
func (m Dashboard) Init() tea.Cmd {
	return tea.Batch(
		listenForCPU(m.cpuChan),
		listenForNet(m.netChan),
		listenForDisk(m.diskChan),
	)
}

// 5. The Update Function
func (m Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	// Handle CPU Data arriving!
	case cpuMsg:
		m.cpu = float64(msg) // save the new data to the state
		return m, listenForCPU(m.cpuChan)
	
	// Handle Net Data arriving!
	case netMsg:
		m.net = monitor.NetworkSpeed(msg)
		return m, listenForNet(m.netChan)
	
	// Handle Disk Data arriving!
	case diskMsg:
		m.disk = monitor.DiskSpeed(msg)
		return m, listenForDisk(m.diskChan)
	}

	return m, nil
}

// 6. The View Function
func (m Dashboard) View() string {
	ui := "🧋 System Monitor Dashboard\n\n"

	ui += fmt.Sprintf("CPU Usage:    %5.2f %%\n", m.cpu)
	ui += fmt.Sprintf("Network:      DL: %8.2f KB/s  |  UL: %8.2f KB/s\n", m.net.RxSpeedKB, m.net.TxSpeedKB)
	ui += fmt.Sprintf("Disk I/O:     Read: %6.2f KB/s  |  Write: %6.2f KB/s\n\n", m.disk.ReadSpeedKB, m.disk.WriteSpeedKB)

	ui += "Press 'q' or 'ctrl+c' to quit.\n"
	return ui
}