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

package svgdata

import (
	"github.com/jbeda/geom"
)

type PathLine struct {
	A, B geom.Coord
}

var _ PathSegment = (*PathLine)(nil)

func NewPathLine(a, b geom.Coord) *PathLine {
	return &PathLine{A: a, B: b}
}

func AlmostEqualsPathLines(a, b *PathLine) bool {
	return (AlmostEqualsCoord(a.A, b.A) && AlmostEqualsCoord(a.B, b.B)) ||
		(AlmostEqualsCoord(a.A, b.B) && AlmostEqualsCoord(a.B, b.A))
}

func (cl *PathLine) Equals(oi interface{}) bool {
	ocl, ok := oi.(*PathLine)
	return ok && AlmostEqualsPathLines(cl, ocl)
}

func (cl *PathLine) Bounds() geom.Rect {
	r := geom.Rect{cl.A, cl.A}
	r.ExpandToContainCoord(cl.B)
	return r
}

func (cl *PathLine) P1() *geom.Coord { return &cl.A }
func (cl *PathLine) P2() *geom.Coord { return &cl.B }
func (cl *PathLine) PathDraw(svg *SVGWriter) {
	svg.PathLineTo(cl.B)
}
func (cl *PathLine) Reverse() {
	cl.A, cl.B = cl.B, cl.A
}
