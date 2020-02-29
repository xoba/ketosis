package main

import (
	"fmt"
	"math"
	"os"

	"github.com/lucasb-eyer/go-colorful"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// potentially lightens a color
func lighten(s string) string {
	x, err := colorful.Hex(s)
	check(err)
	h, c, l := x.Hcl()
	return colorful.Hcl(h, c, l).Hex()
}

func neg() string {
	c, err := colorful.Hex("#ffcccc")
	check(err)
	return c.Hex()
}

func low() string {
	c, err := colorful.Hex("#ffff00")
	check(err)
	return c.Hex()
}

func med() string {
	c, err := colorful.Hex("#ccffcc")
	check(err)
	return c.Hex()
}

func high() string {
	c, err := colorful.Hex("#ccccff")
	check(err)
	return c.Hex()
}

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
background-color:%s;
}

.low {
background-color:%s;
}

.med {
background-color:%s;
}

.high {
background-color:%s;
}

</style>
  </head>
  <body>
`, lighten(neg()), lighten(low()), lighten(med()), lighten(high()))

	fmt.Fprintf(w, `<table>`)

	glucose := gen(50, 150, 21)
	ketones := reverse(gen(0.1, 4.0, 40))

	var rows, cols []float64
	var gf, kf func(float64, float64) float64

	const glucoseRows = false
	if glucoseRows {
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
			// per https://keto-mojo.com/glucose-ketone-index-gki/
			switch {
			case i <= 1: // ≤1 highest therapeutic level of ketosis
				class = "neg"
			case i < 3: // 1-3: high therapeutic level of ketosis
				class = "low"
			case i < 6: // 3-6: moderate level of ketosis
				class = "med"
			case i < 9: // 6-9: low level of ketosis ("ideal")
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

func reverse(in []float64) (out []float64) {
	n := len(in)
	for i := range in {
		out = append(out, in[n-i-1])
	}
	return
}

func gen(x0, x1 float64, n int) (out []float64) {
	dx := (x1 - x0) / float64(n-1)
	for i := 0; i < n; i++ {
		out = append(out, x0+float64(i)*dx)
	}
	return
}

func lgki(glucose, ketones float64) float64 {
	return math.Log(gki(glucose, ketones))
}

func gki(glucose, ketones float64) float64 {
	return glucose / 18 / ketones
}
