// Copyright 2018 Joe Beda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"reflect"

	"github.com/jbeda/geom"
	svgdata "github.com/jbeda/svgdata-go/old"
	dxfcore "github.com/rpaloschi/dxf-go/core"
	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
)

//var dlog = log.New(os.Stderr, "DEBUG ", 0)
var dlog = log.New(ioutil.Discard, "", 0)

func polarToCartesian(center dxfcore.Point, radius, angleDeg float64) geom.Coord {
	angleRad := (angleDeg * math.Pi / 180.0)
	center.Y = -center.Y
	return geom.Coord{
		X: center.X + (radius * math.Cos(angleRad)),
		Y: center.Y - (radius * math.Sin(angleRad)),
	}
}

func dxfCoord2GeomCoordExt(p dxfcore.Point, extrusion dxfcore.Point) geom.Coord {
	if extrusion.Z == 1 {
		return geom.Coord{X: p.X, Y: -p.Y}
	}
	return geom.Coord{X: -p.X, Y: -p.Y}
}

func dxfCoord2GeomCoord(p dxfcore.Point) geom.Coord {
	return geom.Coord{X: p.X, Y: -p.Y}
}

func geomCoordExtAdj(c geom.Coord, extrusion dxfcore.Point) geom.Coord {
	if extrusion.Z == -1 {
		c.X = -c.X
	}
	return c
}

func main() {
	dxfcore.Log.SetOutput(ioutil.Discard)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dxf-file>\n", os.Args[0])
		os.Exit(1)
	}

	infn := os.Args[1]
	inext := path.Ext(infn)
	outfn := infn[0:len(infn)-len(inext)] + ".svg"

	file, err := os.Open(infn)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := document.DxfDocumentFromStream(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	var opc svgdata.OptimizedPathCollection
	var els []svgdata.Element

	for _, entity := range doc.Entities.Entities {
		switch e := entity.(type) {
		case *entities.Line:
			dlog.Printf("Processing Line\n")
			opc.AddSegment(
				svgdata.NewPathLine(
					dxfCoord2GeomCoordExt(e.Start, e.ExtrusionDirection),
					dxfCoord2GeomCoordExt(e.End, e.ExtrusionDirection)))
		case *entities.Circle:
			dlog.Printf("Processing Circle\n")
			els = append(els, &svgdata.Circle{
				Center: dxfCoord2GeomCoordExt(e.Center, e.ExtrusionDirection),
				Radius: e.Radius,
			})
		case *entities.Arc:
			dlog.Printf("Processing Arc. Radius: %f, SA: %f, EA: %f\n",
				e.Radius, e.StartAngle, e.EndAngle)
			start := polarToCartesian(e.Center, e.Radius, e.StartAngle)
			end := polarToCartesian(e.Center, e.Radius, e.EndAngle)

			startAngle, endAngle := e.StartAngle, e.EndAngle
			if endAngle < startAngle {
				endAngle += 360
			}
			largeArc := (endAngle - startAngle) > 180

			sweep := e.ExtrusionDirection.Z == -1

			dlog.Printf("  startAngle: %f, endAngle: %f, largeArc: %t, sweep: %t\n",
				startAngle, endAngle, largeArc, sweep)

			opc.AddSegment(
				svgdata.NewPathCircArc(
					geomCoordExtAdj(start, e.ExtrusionDirection),
					geomCoordExtAdj(end, e.ExtrusionDirection),
					e.Radius, largeArc, sweep))
		case *entities.Polyline:
			dlog.Printf("Processing Polyline\n")
			for i := 0; i < len(e.Vertices)-1; i++ {
				opc.AddSegment(
					svgdata.NewPathLine(
						dxfCoord2GeomCoordExt(e.Vertices[i].Location, e.ExtrusionDirection),
						dxfCoord2GeomCoordExt(e.Vertices[i+1].Location, e.ExtrusionDirection)))
			}
			if e.Closed && len(e.Vertices) > 1 {
				opc.AddSegment(
					svgdata.NewPathLine(
						dxfCoord2GeomCoordExt(e.Vertices[len(e.Vertices)].Location, e.ExtrusionDirection),
						dxfCoord2GeomCoordExt(e.Vertices[0].Location, e.ExtrusionDirection)))
			}
		case *entities.LWPolyline:
			dlog.Printf("Processing LWPolyLine\n")
			for i := 0; i < len(e.Points)-1; i++ {
				opc.AddSegment(
					svgdata.NewPathLine(
						dxfCoord2GeomCoordExt(e.Points[i].Point, e.ExtrusionDirection),
						dxfCoord2GeomCoordExt(e.Points[i+1].Point, e.ExtrusionDirection)))
			}
			if e.Closed && len(e.Points) > 1 {
				opc.AddSegment(
					svgdata.NewPathLine(
						dxfCoord2GeomCoordExt(e.Points[len(e.Points)].Point, e.ExtrusionDirection),
						dxfCoord2GeomCoordExt(e.Points[0].Point, e.ExtrusionDirection)))
			}
		default:
			log.Printf("Unknown entity %s\n", reflect.TypeOf(entity))
		}
	}

	opc.Optimize()

	file, err = os.Create(outfn)
	if err != nil {
		log.Fatal(err)
	}
	w := svgdata.NewSVG(file)
	w.Start(geom.Rect{
		Min: geom.Coord{X: 0, Y: -11},
		Max: geom.Coord{X: 19.5, Y: 0}},
		"width=\"19.5in\"", "height=\"11in\"")
	opc.Draw(w, "fill: none; stroke: black; stroke-width: 0.01")
	for _, el := range els {
		el.Draw(w, "fill: none; stroke: black; stroke-width: 0.01")
	}
	w.End()
	file.Close()
}
