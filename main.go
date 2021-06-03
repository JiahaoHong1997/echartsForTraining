package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func cvsToSlice(file string) ([]string,[]opts.LineData){
	epochs := []string{}
	losses := make([]opts.LineData, 0)
	csvfile, err := os.Open(file)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		epochs = append(epochs, record[1])

		loss1, err := strconv.ParseFloat(record[2],32)
		if err == nil {
			loss, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", loss1), 64) // 保留4位小数
			losses = append(losses, opts.LineData{Value:loss})
		}
	}
	return epochs[1:], losses
}

func lineSplitLine(epochs []string, losses []opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Trainint loss",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			SplitLine: &opts.SplitLine{
				Show: true,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.SetXAxis(epochs).AddSeries("Category A", losses,
		charts.WithLabelOpts(
			opts.Label{Show: true},
		))
	return line
}

func lineMarkPoint(epochs []string, losses []opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Training loss",
		}),
		charts.WithYAxisOpts(opts.YAxis{
		Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		)

	line.SetXAxis(epochs).AddSeries("Category A", losses).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(
				opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
				opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
				opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
			),
			charts.WithMarkPointStyleOpts(
				opts.MarkPointStyle{Label: &opts.Label{Show: true}}),
		)
	return line
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	epochs, losses := cvsToSlice("train_loss.csv")
	line := lineMarkPoint(epochs, losses)
	line.Render(w)
}

func main() {

	// 本地生成
	//epochs, losses := cvsToSlice("train_loss.csv")
	//line := lineMarkPoint(epochs, losses)

	//f, err := os.Create("./loss.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//line.Render(f)

	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)

}
