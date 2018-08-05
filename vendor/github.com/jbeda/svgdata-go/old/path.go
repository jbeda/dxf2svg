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
	"container/list"

	"github.com/jbeda/geom"
)

type Path struct {
	segs   *list.List
	Closed bool
}

func (me *Path) PushFront(seg PathSegment) {
	if me.segs == nil {
		me.segs = new(list.List)
	}
	me.segs.PushFront(seg)
}

func (me *Path) PushPathFront(path *Path) {
	if me.segs == nil {
		me.segs = new(list.List)
	}
	me.segs.PushFrontList(path.segs)
}

func (me *Path) PushBack(seg PathSegment) {
	if me.segs == nil {
		me.segs = new(list.List)
	}
	me.segs.PushBack(seg)
}

func (me *Path) PushPathBack(path *Path) {
	if me.segs == nil {
		me.segs = new(list.List)
	}
	me.segs.PushBackList(path.segs)
}

func (me *Path) Reverse() {
	newSegs := new(list.List)
	for e := me.segs.Front(); e != nil; e = e.Next() {
		e.Value.(PathSegment).Reverse()
		newSegs.PushFront(e.Value)
	}
	me.segs = newSegs
}

func (me *Path) Front() PathSegment {
	if me.segs == nil || me.segs.Len() == 0 {
		return nil
	}
	return me.segs.Front().Value.(PathSegment)
}

func (me *Path) FrontPoint() *geom.Coord {
	s := me.Front()
	if s != nil {
		return s.P1()
	}
	return nil
}

func (me *Path) Back() PathSegment {
	if me.segs == nil || me.segs.Len() == 0 {
		return nil
	}
	return me.segs.Back().Value.(PathSegment)
}

func (me *Path) BackPoint() *geom.Coord {
	s := me.Back()
	if s != nil {
		return s.P2()
	}
	return nil
}

func (me *Path) Draw(svg *SVGWriter, s ...string) {
	startP := me.segs.Front().Value.(PathSegment).P1()
	svg.StartPath(*startP, s...)
	for e := me.segs.Front(); e != nil; e = e.Next() {
		e.Value.(PathSegment).PathDraw(svg)
	}

	if me.Closed {
		svg.PathClose()
	}
	svg.EndPath()
}
