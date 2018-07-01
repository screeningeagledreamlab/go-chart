package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apremalal/go-chart"
	"github.com/apremalal/go-chart/drawing"
	"strconv"
)

func drawChart(res http.ResponseWriter, req *http.Request) {
	err := drawChart2("HLD", 2.25, 5.25, []float64{2.3, 4.5, 0.0, 6.2, 7.2, 3.4 }, res)
	if err != nil {
		log.Fatalf("failed: %s", err.Error())
	}
	res.Header().Set("Content-Type", "image/png")

}

func drawChart2(unit string, min, max float64, values []float64, res http.ResponseWriter) error {
	// create bar values
	defaultBarStyle := chart.Style{
		Show: true,
		FillColor: drawing.ColorFromHex("009ca1"),
		FontColor: drawing.ColorWhite,
	}
	warningBarStyle := chart.Style{
		Show: true,
		FillColor: drawing.ColorFromHex("ff7e2b"),
	}

	bars := make([]chart.Value, 0)
	for i, v := range values {
		if v == 0 {
			bars = append(bars, chart.Value{Value: v + 0.1, Label: strconv.Itoa(i), Style: warningBarStyle})
		} else if (v < max) {
			bars = append(bars, chart.Value{Value: v, Label: strconv.Itoa(i), Style: defaultBarStyle})
		} else {
			bars = append(bars, chart.Value{Value: v, Label: strconv.Itoa(i), Style: warningBarStyle})
		}
	}
	n := len(values)
	sbc := chart.BarChart{
		Min: min,
		Max: max,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
				Left: 40,
			},
		},
		Height:   512,
		Width: chart.DefaultChartWidth,
		BarWidth: (chart.DefaultChartWidth - (20 * n)) / n,
		BarSpacing: 20,
		XAxis: chart.Style{
			Show: true,
			StrokeWidth: 1,
			StrokeColor: drawing.ColorBlack,
		},
		YAxis: chart.YAxis{
			AxisType: chart.YAxisSecondary,
			NameStyle: chart.StyleShow(),
			Style: chart.Style{
				Show: true,
				StrokeWidth: 1,
				StrokeColor: drawing.ColorBlack,
			},
			Ticks: []chart.Tick{
				{0.0, "0"},
				{1.0, "200"},
				{2.0, "400"},
				{3.0, "600"},
				{4.0, "800"},
				{5.0, "1000"},
				{6.0, "1200"},
			},
		},
		Bars: bars,
		Series: []chart.Series{
			chart.AnnotationSeries{
				Annotations: []chart.Value2{
					{Style: defaultBarStyle, XValue: 1.0, YValue: min, Label: strconv.FormatFloat(min, 'f', 0, 64) +
					" " + unit},
					{Style: defaultBarStyle,XValue: 1.0, YValue: max, Label:  strconv.FormatFloat(max, 'f', 0, 64) + " " + unit},
				},
			},
		},
	}
	err := sbc.Render(chart.PNG, res)
	if err != nil {
		return err
	}
	return err
}

func port() string {
	if len(os.Getenv("PORT")) > 0 {
		return os.Getenv("PORT")
	}
	return "8080"
}

func main() {
	listenPort := fmt.Sprintf(":%s", port())
	fmt.Printf("Listening on %s\n", listenPort)
	http.HandleFunc("/", drawChart)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
