package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// 1. the Model
// This struct holds all the state/data for the application.
// Will add CPU, Memory, and Disk stats in here
type Dashboard struct {
	loading bool
} 

// InitialModel is a helper function we will call from main.go
func InitialModel() Dashboard {
	return Dashboard{
		loading: true,
	}
}

// 2. The Init function
// This runs only once when the program starts. It's used to kick off
// initial background workers (like a loading spinner or fetching an API)
// For now, we don't need it to do anything, so we return nil   
func (m Dashboard) Init() tea.Cmd{
	return nil
}

// 3. The Update Function
// This is the "brain" of the UI. Any time an event happens 
// (a key is pressed, your mouse moves, or data arrives from a channel), 
// this function runs.
// It returns an updated version of the Model!
func (m Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg:= msg.(type) {

		// Did the user press a key on their keyboard?
	case tea.KeyMsg:
		switch msg.String() {
			// if user press 'q' or 'ctrl + c', send the Quit command
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}
	return m, nil
}


// 4. The View Function
// This function strictly draws the UI. It takes whatever data is currently
// in your Model and returns a single formatted string. Bubble Tea takes this 
// string and paints it to the terminal 
func (m Dashboard) View() string {
	ui := "🧋 Welcome to the System Monitor!\n\n"
	ui += "We are currently building the interface...\n\n"
	ui += "Press 'q' or 'ctrl+c' to quit.\n"
	return ui
}