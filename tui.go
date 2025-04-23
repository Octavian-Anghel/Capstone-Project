package main

import(
	"fmt"
	"log"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	// Panel 1 - Options
	selectedOption int
	options []string

	//Panel 2 - User input
	inputHash string
	isInputActive bool

	//Panel 3 - Metadata
	hashData map[string]map[string]string // hash metadata (timestamp, temp, humid)
	selectedHash string
}

const (
	viewHash = iota
	addHash
	checkHash
)

func main(){
	//Initialize with a default model
	p := tea.NewProgram(model{
		options: []string{"View Hash", "Add Hash", "Check Hash"},
		hashData: map[string]map[string]string{
			"hash1": {"timestamp": "2025-03-01T12:00:00Z", "temperature": "25°C", "humidity": "60%"},
			"hash2": {"timestamp": "2025-03-01T12:30:00Z", "temperature": "30°C", "humidity": "55%"},
		},
	})
	if err := p.Start(); err != nil{
		log.Fatalf("Error starting program: %v", err)
	}
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case m.isInputActive: // Handle input when user is typing a hash
            switch msg.String() {
            case "enter":
                // Handle user input after pressing Enter (like adding a new hash)
                if m.selectedOption == addHash {
                    m.hashData[m.inputHash] = map[string]string{
                        "timestamp": "2025-03-01T15:00:00Z",
                        "temperature": "28°C",
                        "humidity": "50%",
                    }
                }
                m.isInputActive = false
                m.inputHash = ""
            case "backspace":
                if len(m.inputHash) > 0 {
                    m.inputHash = m.inputHash[:len(m.inputHash)-1]
                }
            case "esc":
                m.isInputActive = false
            default:
                m.inputHash += msg.String()
            }

        case msg.String() == "up" && m.selectedOption > 0:
            m.selectedOption--
        case msg.String() == "down" && m.selectedOption < len(m.options)-1:
            m.selectedOption++
        case msg.String() == "enter":
            // Activate input mode or display data
            if m.selectedOption == addHash || m.selectedOption == checkHash {
                m.isInputActive = true
            }
        case msg.String() == "esc":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    // Panel 1 - Options (with the current selection)
    panel1Style := lipgloss.NewStyle().Padding(1)
    panel1Content := "Select an option:\n"
    for i, option := range m.options {
        if i == m.selectedOption {
            panel1Content += "  → " + option + "\n"
        } else {
            panel1Content += "    " + option + "\n"
        }
    }

    // Panel 2 - User input (when in input mode)
    panel2Style := lipgloss.NewStyle().Padding(1)
    panel2Content := ""
    if m.isInputActive {
        panel2Content = fmt.Sprintf("Enter hash: %s", m.inputHash)
    } else {
        panel2Content = fmt.Sprintf("Selected Hash: %s", m.selectedHash)
    }

    // Panel 3 - Metadata display (for the selected hash)
    panel3Style := lipgloss.NewStyle().Padding(1)
    panel3Content := ""
    if data, ok := m.hashData[m.selectedHash]; ok {
        panel3Content = fmt.Sprintf("Timestamp: %s\nTemperature: %s\nHumidity: %s", data["timestamp"], data["temperature"], data["humidity"])
    } else {
        panel3Content = "No data available for the selected hash."
    }

    // Combine the panels into a single view
    layout := lipgloss.NewStyle().
        Align(lipgloss.Center).
        Render(fmt.Sprintf(
            "%s\n\n%s\n\n%s",
            panel1Style.Render(panel1Content),
            panel2Style.Render(panel2Content),
            panel3Style.Render(panel3Content),
        ))

    return layout
}


