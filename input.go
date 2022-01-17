/*
 *  input.go
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
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func (g *Game) getInput(s tcell.Screen) bool {
	isTurn := false
	quit := func() {
        s.Fini()
        os.Exit(0)
    }
	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		
		switch ev.Rune() {
		case 'C', 'c': s.Clear()
		case 'R': s.Sync()
		case 'h':
			g.hero.pos = heroMove(g, -1, 0)
			isTurn = true
		case 'H':
			g.floors[g.cur].setRunTarget(&g.hero, 7, g); g.state = Autopilot
			isTurn = false
		case 'j':
			g.hero.pos = heroMove(g, 0, 1)
			isTurn = true
		case 'J':
			g.floors[g.cur].setRunTarget(&g.hero, 5, g); g.state = Autopilot
			isTurn = false
		case 'k':
			g.hero.pos = heroMove(g, 0, -1)
			isTurn = true
		case 'K':
			g.floors[g.cur].setRunTarget(&g.hero, 1, g); g.state = Autopilot
			isTurn = false
		case 'l':
			g.hero.pos = heroMove(g, 1, 0)
			isTurn = true
		case 'L':
			g.floors[g.cur].setRunTarget(&g.hero, 3, g); g.state = Autopilot
			isTurn = false
		case 'y':
			g.hero.pos = heroMove(g, -1, -1)
			isTurn = true
		case 'Y':
			g.floors[g.cur].setRunTarget(&g.hero, 0, g); g.state = Autopilot
			isTurn = false
		case 'u':
			g.hero.pos = heroMove(g, +1, -1)
			isTurn = true
		case 'U':
			g.floors[g.cur].setRunTarget(&g.hero, 2, g); g.state = Autopilot
			isTurn = false
		case 'b':
			g.hero.pos = heroMove(g, -1, +1)
			isTurn = true
		case 'B':
			g.floors[g.cur].setRunTarget(&g.hero, 6, g); g.state = Autopilot
			isTurn = false
		case 'n':
			g.hero.pos = heroMove(g, +1, +1)
			isTurn = true
		case 'N':
			g.floors[g.cur].setRunTarget(&g.hero, 4, g); g.state = Autopilot
			isTurn = false
		case '.':
			g.hero.pos = heroMove(g, 0, 0)
			isTurn = true
		case ',':
			isTurn = menuPickup(s, g)
		case '\'':
			g.swapWeapons()
		case 'd':
			isTurn = menuDrop(s, g)
		case 'i':
			menuInv(s, g)
		case '\\':
			menuDisc(s, g)
		case 'x':
			target(s, g)
		case 'z':
			isTurn = g.hero.menuCastSpell(s, g)
		case 'q':
			isTurn = g.menuQuaff(s)
		case 'Q':
			isTurn = g.hero.menuQuiver(s, g)
		case 'f':
			isTurn = g.hero.fireWeapon(s, g)
		case 't':
			isTurn = g.hero.menuThrow(s, g)
		case 'p':
			isTurn = g.menuJewelry(s)
		case 'e':
			isTurn = g.menuEquipArmor(s)
		case 'w':
			isTurn = g.menuWield(s)
		case 'r':
			isTurn = g.menuRead(s)
		case 'M':
			g.hero.castSpell(FighterBerserk, s, g)
		case '%':
			g.dbg(fmt.Sprintf("game tick:: %v", g.tick))
		case '#':
			g.dbg(fmt.Sprintf("1d20:: %v", roll(R1d100)))
		case '[':
			getRandomItem(Scroll)
		case '&':
			g.hero.reorderInv()
		case '>':
			// for _, t := range g.floors[g.cur].tiles {
			// 	if t.pos == g.hero.pos && t.glyph == '>' { g.goDownstairs() }
			// }
			g.goDownstairs()
		case 'D':
			// for _, t := range g.floors[g.cur].tiles {
			// 	if t.pos == g.hero.pos && t.glyph == '>' { g.goDownstairs2() }
			// }
			log.Println("DEBUG: DRAGON LEVEL, hero.pos:", g.hero.pos)
			g.dbg(fmt.Sprintf("DEBUG: DRAGON LEVEL: %v", g.hero.pos))
			g.goDownstairs2()
		case '<':
			// for _, t := range g.floors[g.cur].tiles {
			// 	if t.pos == g.hero.pos  && t.glyph == '<'  { g.goUpstairs() }
			// }
			g.goUpstairs()

		case 'F':
			g.hero.spellBlink(g); isTurn = true
		case '=':
			g.menuReorder(s)
		case '1':
			if len(g.hero.spells) > 0 && g.hero.canCast {g.hero.castSpell(g.hero.spells[0], s, g)}
		case '2':
			if len(g.hero.spells) > 1 && g.hero.canCast {g.hero.castSpell(g.hero.spells[1], s, g)}
		case '3':
			if len(g.hero.spells) > 2 && g.hero.canCast {g.hero.castSpell(g.hero.spells[2], s, g)}
		case '4':
			if len(g.hero.spells) > 3 && g.hero.canCast {g.hero.castSpell(g.hero.spells[3], s, g)}
		case '5':
			if len(g.hero.spells) > 4 && g.hero.canCast {g.hero.castSpell(g.hero.spells[4], s, g)}
		case '6':
			if len(g.hero.spells) > 5 && g.hero.canCast {g.hero.castSpell(g.hero.spells[5], s, g)}
		case '7':
			if len(g.hero.spells) > 6 && g.hero.canCast {g.hero.castSpell(g.hero.spells[6], s, g)}
		case '8':
			if len(g.hero.spells) > 7 && g.hero.canCast {g.hero.castSpell(g.hero.spells[7], s, g)}
		case '9':
			if len(g.hero.spells) > 8 && g.hero.canCast {g.hero.castSpell(g.hero.spells[8], s, g)}
		case '0':
			if len(g.hero.spells) > 9 && g.hero.canCast {g.hero.castSpell(g.hero.spells[9], s, g)}
		}
		switch ev.Key() {
		case tcell.KeyLeft:
			g.hero.pos = heroMove(g, -1, 0); isTurn = true
		case tcell.KeyUp:
			g.hero.pos = heroMove(g, 0, -1); isTurn = true
		case tcell.KeyDown:
			g.hero.pos = heroMove(g, 0, 1); isTurn = true
		case tcell.KeyRight:
			g.hero.pos = heroMove(g, 1, 0); isTurn = true
		case tcell.KeyEscape, tcell.KeyCtrlC: quit()
		case tcell.KeyCtrlL: s.Sync()
		case tcell.KeyCtrlP: drawAllMessages(s, g.msg)
		default: return isTurn
		}
	}
	return isTurn
}