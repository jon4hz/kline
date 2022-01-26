package kline

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	candleHalfShadowUp   = "╷"
	candleShadow         = "│"
	candleCross          = "┼"
	candleHalfShadowDown = "╵"
	candleHalfBodyUp     = "▄"
	candleBody           = "█"
	candleHalfBodyDown   = "▀"
)

type Kline struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	CloseTime time.Time
}

func New(data []Kline, width, height int) Model {
	return Model{
		Klines: data,
		Width:  width,
		Height: height,
	}
}

type Model struct {
	Width         int
	Height        int
	Klines        []Kline
	maxPrice      float64
	minPrice      float64
	pricePerBlock float64
	visibleIndex  int
	visibleOffset int
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.calculate()
}

func (m *Model) calculate() {
	var maxPrice float64
	var minPrice float64
	m.visibleIndex = len(m.Klines) - m.Width
	m.visibleIndex -= m.visibleOffset
	if m.visibleIndex < 0 {
		m.visibleIndex = 0
	}
	for i, c := range m.Klines[m.visibleIndex : len(m.Klines)-m.visibleOffset] {
		if i == 0 {
			minPrice = c.Low
		}
		if c.High > maxPrice {
			maxPrice = c.High
		}
		if c.Low < minPrice {
			minPrice = c.Low
		}
	}
	m.maxPrice = maxPrice
	m.minPrice = minPrice

	priceDelta := maxPrice - minPrice
	m.pricePerBlock = priceDelta / float64(m.Height)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelDown:
			m.visibleOffset++
			if m.visibleOffset > len(m.Klines)-m.Width {
				m.visibleOffset = len(m.Klines) - m.Width
			}
			m.calculate()
			return m, nil

		case tea.MouseWheelUp:
			m.visibleOffset--
			if m.visibleOffset < 0 {
				m.visibleOffset = 0
			}
			m.calculate()
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s []string
	for _, c := range m.Klines[m.visibleIndex : len(m.Klines)-m.visibleOffset] {
		if c.Open <= c.Close {
			s = append(s, m.renderCandle(c, true))
		} else {
			c.Close, c.Open = c.Open, c.Close
			s = append(s, m.renderCandle(c, false))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (m Model) renderCandle(c Kline, green bool) string {
	var (
		candle     string
		paddingTop int
	)
	for i := m.Height; i > 0; i-- {
		j := (float64(i) * m.pricePerBlock) + m.minPrice
		if j > c.High {
			paddingTop++
			continue
		}
		if (j > c.Close && ((float64(i-1))/m.pricePerBlock) > c.Open) || (j > c.Close && c.Open != c.Close) {
			candle += candleShadow
			continue
		}
		if j > c.Open {
			if c.Open == c.Close {
				candle += candleCross
			} else {
				candle += candleBody
			}
			continue
		}
		if j > c.Low {
			candle += candleShadow
			continue
		}

	}
	style := lipgloss.NewStyle().
		Width(1).
		PaddingTop(paddingTop).
		Foreground(lipgloss.Color("#008000"))

	if !green {
		style = style.Foreground(lipgloss.Color("#FF0000"))
	}

	return style.Render(candle)
}
