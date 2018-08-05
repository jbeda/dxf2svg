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
	"fmt"
	"io"
	"strings"

	"github.com/jbeda/geom"
)

type SVGWriter struct {
	writer io.Writer
}

func NewSVG(w io.Writer) *SVGWriter {
	return &SVGWriter{w}
}

func (svg *SVGWriter) printf(format string, a ...interface{}) (n int, errno error) {
	return fmt.Fprintf(svg.writer, format, a...)
}

// BUGBUG: not quoting aware
func extraparams(s []string) string {
	ep := ""
	for i := 0; i < len(s); i++ {
		if strings.Index(s[i], "=") > 0 {
			ep += (s[i]) + " "
		} else if len(s[i]) > 0 {
			ep += fmt.Sprintf("style='%s' ", s[i])
		}
	}
	return ep
}

func onezero(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (svg *SVGWriter) Start(viewBox geom.Rect, s ...string) {
	svg.printf(`<?xml version="1.0"?>
<svg version="1.1"
     viewBox="%f %f %f %f"
     xmlns="http://www.w3.org/2000/svg" %s>
`, viewBox.Min.X, viewBox.Min.Y, viewBox.Width(), viewBox.Height(), extraparams(s))
}

func (svg *SVGWriter) End() {
	svg.printf("</svg>\n")
}

func (svg *SVGWriter) Line(p1 geom.Coord, p2 geom.Coord, s ...string) {
	svg.printf("<line x1='%f' y1='%f' x2='%f' y2='%f' %s/>\n", p1.X, p1.Y, p2.X, p2.Y, extraparams(s))
}

func (svg *SVGWriter) Circle(c geom.Coord, r float64, s ...string) {
	svg.printf("<circle cx='%f' cy='%f' r='%f' %s/>\n", c.X, c.Y, r, extraparams(s))
}

func (svg *SVGWriter) CircularArc(p1, p2 geom.Coord, r float64, largeArc, sweep bool, s ...string) {
	svg.printf("<path d='M%f,%f A%f,%f 0 %s,%s %f,%f' %s/>\n",
		p1.X, p1.Y, r, r, onezero(largeArc), onezero(sweep), p2.X, p2.Y, extraparams(s))
}

func (svg *SVGWriter) QuadBezier(p1 geom.Coord, ctrl1 geom.Coord, p2 geom.Coord, s ...string) {
	svg.printf("<path d='M%f,%f Q%f,%f %f,%f' %s/>\n",
		p1.X, p1.Y, ctrl1.X, ctrl1.Y, p2.X, p2.Y, extraparams(s))
}

func (svg *SVGWriter) CubicBezier(p1, ctrl1, ctrl2, p2 geom.Coord, s ...string) {
	svg.printf("<path d='M%f,%f C%f,%f %f,%f %f,%f' %s/>\n",
		p1.X, p1.Y, ctrl1.X, ctrl1.Y, ctrl2.X, ctrl2.Y, p2.X, p2.Y, extraparams(s))
}

func (svg *SVGWriter) StartPath(p1 geom.Coord, s ...string) {
	svg.printf("<path %sd='M%f,%f", extraparams(s), p1.X, p1.Y)
}

func (svg *SVGWriter) EndPath() {
	svg.printf("'/>\n")
}

func (svg *SVGWriter) PathLineTo(p geom.Coord) {
	svg.printf("\n  L%f,%f", p.X, p.Y)
}

func (svg *SVGWriter) PathCircularArcTo(p geom.Coord, r float64, largeArc, sweep bool) {
	svg.printf("\n  A%f,%f 0 %s,%s %f,%f", r, r, onezero(largeArc), onezero(sweep), p.X, p.Y)
}

func (svg *SVGWriter) PathQuadBezierTo(p, ctrl1 geom.Coord) {
	svg.printf("\n  Q%f,%f, %f,%f", ctrl1.X, ctrl1.Y, p.X, p.Y)
}

func (svg *SVGWriter) PathCubicBezierTo(p, ctrl1, ctrl2 geom.Coord) {
	svg.printf("\n  C%f,%f, %f,%f %f,%f", ctrl1.X, ctrl1.Y, ctrl2.X, ctrl2.Y, p.X, p.Y)
}

func (svg *SVGWriter) PathClose() {
	svg.printf("\n  Z")
}
