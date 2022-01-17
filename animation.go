/*
 *  animation.go
 *  Rescape
 *
 *  Created by Chad Pierce on 1/16/2022.
 *  Copyright 2022. All rights reserved.
 *
 *  This file is part of Rescape.
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
    "time"

    "github.com/gdamore/tcell/v2"
)

func testBeam(s tcell.Screen, g *Game) {
	x := 5
	y := 5
	for i := x; i <= 40; i++ {
		p := Point { i, y }
		for j, t := range g.floors[g.cur].tiles {
			if t.pos == p {
				g.floors[g.cur].tiles[j].bg = tcell.ColorPink
			}
			s.Show()
		}
	time.Sleep(200 * time.Millisecond)
	s.Show()
	}
}

func testBeam2() {
	s := initScreen()
	x := 5
	y := 5

	for i := x; i <= 10; i++ {

		text := "hello world"
        style := tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorPink)
        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
		time.Sleep(200 * time.Millisecond)
		s.Show()
	}

}