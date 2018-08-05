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

type PathCircArc struct {
	A, B            geom.Coord
	R               float64
	LargeArc, Sweep bool
}

var _ PathSegment = (*PathCircArc)(nil)

func NewPathCircArc(a, b geom.Coord, r float64, largeArc, sweep bool) *PathCircArc {
	return &PathCircArc{A: a, B: b, R: r, LargeArc: largeArc, Sweep: sweep}
}

func AlmostEqualsPathCircArc(a, b *PathCircArc) bool {
	if AlmostEqualsCoord(a.A, b.A) &&
		AlmostEqualsCoord(a.B, b.B) &&
		FloatAlmostEqual(a.R, b.R) &&
		a.LargeArc == b.LargeArc &&
		a.Sweep == b.Sweep {
		return true
	}

	if AlmostEqualsCoord(a.A, b.B) &&
		AlmostEqualsCoord(a.B, b.A) &&
		FloatAlmostEqual(a.R, b.R) &&
		a.LargeArc == b.LargeArc &&
		a.Sweep != b.Sweep {
		return true
	}

	return false
}

func (a *PathCircArc) Equals(oi interface{}) bool {
	oa, ok := oi.(*PathCircArc)
	return ok && AlmostEqualsPathCircArc(a, oa)
}

func (a *PathCircArc) Bounds() geom.Rect {
	r := geom.Rect{a.A, a.A}
	r.ExpandToContainCoord(a.B)
	return r
}

func (a *PathCircArc) P1() *geom.Coord { return &a.A }
func (a *PathCircArc) P2() *geom.Coord { return &a.B }
func (a *PathCircArc) PathDraw(svg *SVGWriter) {
	svg.PathCircularArcTo(a.B, a.R, a.LargeArc, a.Sweep)
}
func (a *PathCircArc) Reverse() {
	a.A, a.B = a.B, a.A
	a.Sweep = !a.Sweep
}
