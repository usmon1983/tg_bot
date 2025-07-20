package main

import (
	"image/color"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func CreateStatsChart(data map[string]float64, filename string) error {
	p := plot.New()
	p.Title.Text = "Расходы за последние 7 дней"
	p.X.Label.Text = "Дата"
	p.Y.Label.Text = "Сумма"

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pts := make(plotter.XYs, len(keys))
	for i, k := range keys {
		pts[i].X = float64(i)
		pts[i].Y = data[k]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	line.LineStyle.Width = vg.Points(2)
	line.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(line)
	p.NominalX(keys...)

	return p.Save(10*vg.Centimeter, 5*vg.Centimeter, filename)
}
