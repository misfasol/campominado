package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	estiloBorda    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	estiloTitulo   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	estiloBomba    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#FF4444"))
	estiloFechado  = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	estiloBandeira = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	estiloCursor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
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

// melhorar a funcao Randomizar() e a Abrir() (talvez juntar elas)
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

// melhorar a funcao Randomizar() e a Abrir() (talvez juntar elas)
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

// adicionar detecção de bomba ou completado

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
	s := estiloTitulo.Render("Campo Minado")
	s += "\n"

	campo := ""

	for i := range m.bombas {
		for j := range m.bombas[0] {
			celula := ""
			cnesse := m.cursorx == i && m.cursory == j
			if cnesse {
				celula += estiloCursor.Render("[")
			} else {
				celula += " "
			}
			if m.bombas[i][j].bandeira {
				celula += estiloBandeira.Render("█")
			} else if m.bombas[i][j].aberto {
				if m.bombas[i][j].bomba {
					celula += estiloBomba.Render("O")
				} else {
					// tem que melhorar isso aqui vvv
					if m.bombas[i][j].quanto == 0 {
						numero := strconv.FormatInt(int64(m.bombas[i][j].quanto), 10)
						celula += estiloFechado.Render(numero)
					} else {
						celula += strconv.FormatInt(int64(m.bombas[i][j].quanto), 10)
					}
				}
			} else {
				celula += estiloFechado.Render("#")
			}
			if cnesse {
				celula += estiloCursor.Render("]")
			} else {
				celula += " "
			}
			campo += celula
		}
		campo += "\n"
	}

	s += campo
	s += "Pressione q para sair."

	// adicionar mais controles

	s = estiloBorda.Render(s)
	s += "\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
