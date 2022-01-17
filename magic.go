/*
 *  magic.go
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
    //"log"
    // "os"
    // "strconv"
	
	"github.com/gdamore/tcell/v2"
)

type SpellID int
const (
	Blink SpellID = iota
	Zap
	FighterBerserk
	RogueSmokeBomb
	RogueShadowStep

)

func getSpellName(sid SpellID) string {
	switch sid {
	case Blink: return "Blink"
	case Zap: return "Zap Something"
	case FighterBerserk: return "Berserk"
	default: return "poof"
	}
}

func (a *Actor) castSpell(spell SpellID, s tcell.Screen, g *Game) bool {
	isTurn := false
	switch spell {
	case Blink: a.spellBlink(g); isTurn = true
	case Zap: isTurn = a.spellZap(s, g)
	case FighterBerserk: isTurn = a.spellBerserk(g)
	}
	return isTurn
}

func (a *Actor) spellBerserk(g *Game) bool {
	isTurn := false
	manaCost := 1
	if a.isCastable(manaCost, g) { 
		a.pFg = a.fg
		a.fg = tcell.ColorRed
		a.pStrg = a.strg
		a.pIntel = a.intel
		a.strg = a.strg + a.strg/2
		a.intel = a.intel - a.intel/2
		a.canRead = false
		a.canEquip = false
		a.canCast = false
		a.canQuaff = false
		a.addActorEvent(EventEndBerserk, g.tick + uint(roll(R3d4)))
		isTurn = true
		a.mana -= manaCost
		a.energy -= SpeedCost
		g.addMessage("You see red", tcell.ColorDefault)
	}
	return isTurn
}

func (a *Actor) spellZap(s tcell.Screen, g *Game) bool {
	isTurn := false
	manaCost := 1
	if !a.isCastable(manaCost, g) { return false }
	tPos := a.spellTarget(s, g)
	cancel := Point { -1, -1 }
	if tPos == cancel { g.dbg("cancell");return isTurn}
	for id, t := range g.floors[g.cur].actors {
		if t.pos == tPos {
			isTurn = true
			a.mana -= manaCost
			a.energy -= SpeedCost
			drawZapLine(s, a.pos.x, a.pos.y, t.pos.x, t.pos.y)
			g.floors[g.cur].actors[id].takeDamage(3)
		}
	}
	return isTurn
}


func (a *Actor) spellBlink(g *Game) {
	// TODO should cursed blink make monsters blink towards you?, or just make them blink randomly?
	// TODO completely inefficient should probably rewrite
	manaCost := 1
	if !a.isCastable(manaCost, g) { return }

	for {
		
		f := getRandNum(2, 1)
		dx := getRandNum(5, 2)
		if f == 1 { dx = a.pos.x + dx
		} else if f == 2 { dx = a.pos.x - dx }
		f = getRandNum(2, 1)
		dy := getRandNum(5, 2)
		if f == 1 { dy = a.pos.y + dy
		} else if f == 2 { dy = a.pos.y - dy }
		isB, _ := isBlocked(dx, dy, g.floors[g.cur].actors, g.floors[g.cur].tiles)
		if !isB { a.pos.x = dx; a.pos.y = dy
			a.mana -= manaCost; return }
	}
}

func (a *Actor) spellRogueShadowStep(g *Game) {
	manaCost := 1
	if !a.isCastable(manaCost, g) { return }

	for {
		
		f := getRandNum(2, 1)
		dx := getRandNum(5, 2)
		if f == 1 { dx = a.pos.x + dx
		} else if f == 2 { dx = a.pos.x - dx }
		f = getRandNum(2, 1)
		dy := getRandNum(5, 2)
		if f == 1 { dy = a.pos.y + dy
		} else if f == 2 { dy = a.pos.y - dy }
		isB, _ := isBlocked(dx, dy, g.floors[g.cur].actors, g.floors[g.cur].tiles)
		if !isB { a.pos.x = dx; a.pos.y = dy
			a.mana -= manaCost; return }
	}
}


func (a *Actor) isCastable(cost int, g *Game) bool {
	if a.mana >= cost { return true 
	} else { g.addMessage("You don't have enough mana", tcell.ColorDefault); return false }
}

func (a *Actor) spellTarget(s tcell.Screen, g *Game) Point {
	t := a.pos
	drawTarget(t, s, *g, "", 1)  // initialize target ui
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			switch ev.Rune() {
			case 'h': t.x += -1; t.y += 0; case 'j': t.x += 0; t.y += 1
			case 'k': t.x += 0; t.y += -1; case 'l': t.x += 1; t.y += 0
			case 'H': t.x += -5; t.y += 0; case 'J': t.x += 0; t.y += 5
			case 'K': t.x += 0; t.y += -5; case 'L': t.x += 5; t.y += 0
			case 'y': t.x += -1; t.y += -1; case 'Y': t.x += -5; t.y += -5
			case 'u': t.x += 1; t.y += -1; case 'U': t.x += 5; t.y += -5
			case 'b': t.x += -1; t.y += 1; case 'B': t.x += -5; t.y += 5
			case 'n': t.x += 1; t.y += 1; case 'N': t.x += 5; t.y += 5 
			}
			switch ev.Key() {
			case tcell.KeyEscape: return Point { -1, -1 }
			case tcell.KeyEnter: return t
			case tcell.KeyLeft: t.x += -1; t.y += 0
			case tcell.KeyDown: t.x += 0; t.y += 1
			case tcell.KeyUp: t.x += 0; t.y += -1
			case tcell.KeyRight: t.x += 1; t.y += 0 
			} 
		}
		var itemStr string
		for _, a := range g.floors[g.cur].actors {
			if a.pos == t && a.name != "floor" {
				aAI := ""
				if a.ai == SleepAI {
					aAI = "(sleeping)"
				} else if a.ai == ConfusedAI { 
					aAI = "(confused)"
				}
				if a.weapon != -1 {
					itemStr = "Here lies: " + a.name + " [" + a.inv[a.weapon].item.name + "]"  
				} else {
					itemStr = "Here lies: " + a.name
				}
				if aAI != "" {
					itemStr = itemStr + " " + aAI
				}
			}
		}
		drawTarget(t, s, *g, itemStr, 1)
	}
}

func (a *Actor) memorizeSpell(sid SpellID, g *Game, s tcell.Screen) bool {
	a.spells = append(a.spells, sid)
	g.addMessage(fmt.Sprintf("You have learned %v", getSpellName(sid)), tcell.ColorDefault)
	return false
}


