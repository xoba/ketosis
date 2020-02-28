package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	w := os.Stdout
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
<style>
table {
    border-collapse:collapse;
}

td,th{
    margin:0em;
    padding:0.3em;
    border:1px solid black;
}

.neg {
background-color:#fcc;
}

.low {
background-color:yellow;
}

.med {
background-color:green;
}

.high {
background-color:blue;
}

</style>
  </head>
  <body>
`)

	fmt.Fprintf(w, `<table>`)

	glucose := gen(50, 150, 5)
	ketones := rev(gen(0.1, 4.0, 0.1))

	var rows, cols []float64
	var gf, kf func(float64, float64) float64

	if false {
		rows, cols = glucose, ketones
		gf = func(row, col float64) float64 {
			return row
		}
		kf = func(row, col float64) float64 {
			return col
		}
	} else {
		rows, cols = ketones, glucose
		gf = func(row, col float64) float64 {
			return col
		}
		kf = func(row, col float64) float64 {
			return row
		}
	}

	fmt.Fprintf(w, "<tr><th></th>")
	for _, col := range cols {
		fmt.Fprintf(w, `<th>%.1f</th>`, col)
	}
	fmt.Fprintf(w, "</tr>")
	for _, row := range rows {
		fmt.Fprintf(w, "<tr>")
		fmt.Fprintf(w, "<th>%.1f</th>", row)
		for _, col := range cols {
			g, k := gf(row, col), kf(row, col)
			i := gki(g, k)
			var class string
			switch {
			case i < 1:
				class = "neg"
			case i < 3:
				class = "low"
			case i < 6:
				class = "med"
			case i < 9:
				class = "high"
			}
			if false {
				i = lgki(g, k)
			}
			fmt.Fprintf(w, `<td class=%q>%.2f</td>`, class, i)
		}
		fmt.Fprintf(w, "</tr>")
	}
	fmt.Fprintf(w, `</table></body></html>`)
}

func rev(in []float64) (out []float64) {
	n := len(in)
	for i := range in {
		out = append(out, in[n-i-1])
	}
	return
}

func gen(x0, x1, dx float64) (out []float64) {
	for x := x0; x <= x1; x += dx {
		out = append(out, x)
	}
	return
}

func lgki(glucose, ketones float64) float64 {
	return math.Log(gki(glucose, ketones))
}

func gki(glucose, ketones float64) float64 {
	return glucose / 18 / ketones
}
