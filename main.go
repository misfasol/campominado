package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type cel struct {
	bomba    bool
	aberto   bool
	bandeira bool
	quanto   int
}

type campo struct {
	bombas       [][]cel
	inicializado bool
	cursorx      int
	cursory      int
}

func (c *campo) Randomizar(quantidade int) {
	for i := range len(c.bombas) {
		for j := range len(c.bombas[0]) {
			c.bombas[i][j].bomba = false
			c.bombas[i][j].quanto = 0
			c.bombas[i][j].aberto = false
		}
	}

	contador := 0
	for adicionado := 0; adicionado < quantidade; {
		i := (contador / len(c.bombas)) % len(c.bombas[0])
		j := contador % len(c.bombas[0])
		chance := rand.Intn(5)
		if chance == 0 && (c.cursorx != i || c.cursory != j) {
			c.bombas[i][j].bomba = true
			adicionado += 1
		}
		contador += 1
	}

	for i := range len(c.bombas) {
		for j := range len(c.bombas[0]) {
			for a := range 3 {
				for b := range 3 {
					if b != 1 || a != 1 {
						ia := i + a - 1
						jb := j + b - 1
						if ia >= 0 && ia < len(c.bombas) && jb >= 0 && jb < len(c.bombas[0]) {
							if c.bombas[ia][jb].bomba {
								c.bombas[i][j].quanto += 1
							}
						}
					}
				}
			}
		}
	}

	c.inicializado = true
}

func (c *campo) Abrir(cx, cy int) {
	c.bombas[cx][cy].aberto = true
	if c.bombas[cx][cy].quanto == 0 {
		for i := range 3 {
			for j := range 3 {
				if i != 1 || j != 1 {
					xi := cx + i - 1
					yj := cy + j - 1
					if xi >= 0 && xi < len(c.bombas) && yj >= 0 && yj < len(c.bombas[0]) &&
						!c.bombas[xi][yj].aberto && !c.bombas[xi][yj].bomba {
						c.Abrir(xi, yj)
					}
				}
			}
		}
	}
}

func initialModel() campo {
	campo := campo{
		bombas: make([][]cel, 10),
	}
	for i := range campo.bombas {
		campo.bombas[i] = make([]cel, 10)
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

		case "f":
			m.bombas[m.cursorx][m.cursory].bandeira = !m.bombas[m.cursorx][m.cursory].bandeira

		case "enter", " ":
			if !m.inicializado {
				m.Randomizar(10)
			}
			m.Abrir(m.cursorx, m.cursory)

		case "r":
			m.Randomizar(10)
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
			if m.bombas[i][j].bandeira {
				cursor += "F"
			} else if m.bombas[i][j].aberto {
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
