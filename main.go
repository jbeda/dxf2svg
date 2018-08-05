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
	"os"
	"path"
	"reflect"

	"github.com/jbeda/geom"
	"github.com/jbeda/svgdata-go/old"
	dxfcore "github.com/rpaloschi/dxf-go/core"
	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
)

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
		case entities.Line:
			opc.AddSegment(
				svgdata.NewPathLine(
					geom.Coord{e.Start.X, e.Start.Y},
					geom.Coord{e.End.X, e.End.Y}))
		case entities.Circle:
			els = append(els, &svgdata.Circle{
				Center: geom.Coord{e.Center.X, e.Center.Y},
				Radius: e.Radius,
			})
		default:
			fmt.Printf("Unknown entity %s\n", reflect.TypeOf(entity))
		}
	}

	opc.Optimize()

	file, err = os.Create(outfn)
	if err != nil {
		log.Fatal(err)
	}
	w := svgdata.NewSVG(file)
	w.Start(geom.Rect{geom.Coord{0, 0}, geom.Coord{19, 11}}, "width=\"19in\"", "height=\"11in\"")
	opc.Draw(w, "fill: none; stroke: black; stroke-width: 0.01in")
	for _, el := range els {
		el.Draw(w, "fill: none; stroke: black; stroke-width: 0.01in")
	}
	w.End()
	file.Close()
}
