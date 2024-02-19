package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type cel struct {
	bomba  bool
	aberto bool
	quanto int
}

type campo struct {
	bombas       [][]cel
	inicializado bool
	cursorx      int
	cursory      int
}

func (c *campo) Randomizar(quantidade, cx, cy int) {
	for i := range len(c.bombas) {
		for j := range len(c.bombas[0]) {
			c.bombas[i][j].bomba = false
		}
	}

	contador := 0
	for adicionado := 0; adicionado < quantidade; {
		i := (contador / 5) % 5
		j := contador % 5
		chance := rand.Intn(5)
		if chance == 0 && (cx != i || cy != j) {
			c.bombas[i][j].bomba = true
			adicionado += 1
		}
		contador += 1
	}
	c.inicializado = true
}

func initialModel() campo {
	campo := campo{
		bombas: make([][]cel, 5),
	}
	for i := range campo.bombas {
		campo.bombas[i] = make([]cel, 5)
	}
	campo.inicializado = false
	return campo
}

func (m campo) Init() tea.Cmd {
	return nil
}

func (m campo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursorx > 0 {
				m.cursorx--
			}

		case "down", "j":
			if m.cursorx < len(m.bombas)-1 {
				m.cursorx++
			}

		case "left", "h":
			if m.cursory > 0 {
				m.cursory--
			}

		case "right", "l":
			if m.cursory < len(m.bombas[0])-1 {
				m.cursory++
			}

		case "r":
			m.Randomizar(5, m.cursorx, m.cursory)

		case "enter", " ":
			m.bombas[m.cursorx][m.cursory].aberto = !m.bombas[m.cursorx][m.cursory].aberto
		}
	}

	return m, nil
}

func (m campo) View() string {
	s := "Campo Minado\n"

	for i := range m.bombas {
		for j := range m.bombas[0] {
			cursor := ""
			cnesse := m.cursorx == i && m.cursory == j
			if cnesse {
				cursor += "["
			} else {
				cursor += " "
			}
			if m.bombas[i][j].aberto {
				if m.bombas[i][j].bomba {
					cursor += "X"
				} else {
					cursor += strconv.FormatInt(int64(m.bombas[i][j].quanto), 10)
				}
			} else {
				cursor += "#"
			}
			if cnesse {
				cursor += "]"
			} else {
				cursor += " "
			}
			s += cursor
		}
		s += "\n"
	}

	s += "Pressione q para sair.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
