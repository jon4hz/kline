package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/kline"
)

type model struct {
	kline kline.Model
}

func main() {
	data, err := getData()
	if err != nil {
		log.Fatalln(err)
	}
	/* data = []ohlc{
		{Open: 1, High: 7, Low: 0, Close: 2},
		{Open: 2, High: 5, Low: 1, Close: 4},
		{Open: 1, High: 2, Low: 0, Close: 1},
		{Open: 4, High: 10, Low: 4, Close: 6},
		{Open: 2, High: 5, Low: 1, Close: 4},
		{Open: 4.1, High: 6.3, Low: 3.3, Close: 5.6},
	} */
	/* data = []ohlc{
		{Open: 1, High: 7, Low: 0, Close: 2},
		{Open: 2, High: 5, Low: 1, Close: 4},
		{Open: 4, High: 7, Low: 3, Close: 6},
		{Open: 6, High: 7, Low: 2, Close: 4},
		{Open: 4, High: 4, Low: 1, Close: 3},
		{Open: 3, High: 5, Low: 1, Close: 4},
	} */
	/* data = []ohlc{
		{Open: 41852.89, High: 41947.08, Low: 41699.46, Close: 41915.98},
		{Open: 41933.07, High: 42030.15, Low: 41862.73, Close: 41955.43},
	} */

	m := model{
		kline: kline.Model{
			Klines: data,
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalln(err)
	}
}

func getData() ([]kline.Kline, error) {
	c := binance.NewClient("", "")
	klines, err := c.NewKlinesService().Symbol("BTCUSDC").Interval("5m").Do(context.Background())
	if err != nil {
		return nil, err
	}
	const amount = 500
	data := make([]kline.Kline, amount)
	for i := 0; i < amount; i++ {
		open, err := strconv.ParseFloat(klines[i].Open, 64)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(klines[i].High, 64)
		if err != nil {
			return nil, err
		}
		low, err := strconv.ParseFloat(klines[i].Low, 64)
		if err != nil {
			return nil, err
		}
		close, err := strconv.ParseFloat(klines[i].Close, 64)
		if err != nil {
			return nil, err
		}
		data[i] = kline.Kline{
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			CloseTime: time.UnixMilli(klines[i].CloseTime),
		}
	}
	return data, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.kline.SetSize(msg.Width, msg.Height)

	}
	return m, nil
}

func (m model) View() string {
	return m.kline.View()
}
