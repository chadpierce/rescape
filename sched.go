/*
 *  sched.go
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
	//"fmt"
	"github.com/gdamore/tcell/v2"
)

type eventType int
const (
	EventBasicAI eventType = iota
	EventConfAI
	EventDOTPois
	EventHealthRegen
	EventManaRegen
	EventEndBerserk
	EventFatigue
	EventEndFatigue
)

func (g *Game) schedule() {
	if g.hero.alive == false { g.heroDeath() }
	g.hero.doActorEvents(g)
	g.hero.schedRegen(g)
	for g.hero.energy < SpeedCost {
		for id, a := range g.floors[g.cur].actors {
			for g.floors[g.cur].actors[id].energy >= SpeedCost {
				if a.alive && a.ai != NoAI { //&& a.target.x != -1 {
					g.floors[g.cur].actors[id].doActorEvents(g)
					g.floors[g.cur].actors[id].schedRegen(g)
					switch a.ai {
					case BasicAI: dumbMove(id, a.pos, g.hero.pos, g)
					case ConfusedAI: confusedMove(id, a.pos, g.hero.pos, g)
					case WanderAI: wanderMove(id, g)
					}
				}
				// TODO how will initiative work for mobs?
				g.floors[g.cur].actors[id].energy -= SpeedCost // TODO this should move to action funcs
			}
			g.floors[g.cur].actors[id].energy += a.spd
		}
		g.hero.energy += g.hero.spd
	}
}

func (a *Actor) doActorEvents(g *Game) {
	// events
	for id, e := range a.events {
		if e.tick == g.tick {
			a.executeActorEvent(id, g)  //; return
		}
	}
	
}

func (a *Actor) addActorEvent(e eventType, t uint) {
	a.events = append(a.events, Event { e, t })
}

func (a *Actor) executeActorEvent(id int, g *Game) {
	msg := ""
	switch a.events[id].eventType {
	case EventBasicAI: a.ai = BasicAI; msg = "You feel normal"
	case EventManaRegen: a.addMana()
	case EventHealthRegen: a.addHealth()
	case EventDOTPois: a.takeDamage(1) //; msg = "mob takes damange"
	case EventEndBerserk: a.endBerserk(g); msg = "You feel slow"
	case EventFatigue: a.fatigue(g)
	case EventEndFatigue: a.endFatigue(g); msg = "You feel normal"
	}

	if a == &g.hero && msg != "" { 
		g.addMessage(msg, tcell.ColorDefault)
	}
}

func (a *Actor) removeActorEvent(t eventType) {
	//TODO this hs not been tested
	for id, e := range a.events {
		if e.eventType == t {
			copy(a.events[id:], a.events[id+1:])
			a.events[len(a.events)-1] = Event{}
			a.events = a.events[:len(a.events)-1]
		}
	}
}

func (a *Actor) doEvent(e eventType) {
	
}

func (a *Actor) schedRegen(g *Game) {

	if g.tick % 3 == 0 {
		a.schedManaRegen(g)
	}
	if g.tick % 4 == 0 {
		a.schedHealthRegen(g)
	}
}

func (a *Actor) schedHealthRegen(g *Game) {
	if a.hp < a.maxHP {
		combinedStats := int(((a.strg + a.intel + a.dex)/3) / 2)
		a.addActorEvent(EventHealthRegen, g.tick + uint(combinedStats))
	}
}

func (a *Actor) schedManaRegen(g *Game) {
	if a.mana < a.maxMana {
		combinedStats := int(((a.strg + a.intel + a.dex)/3) / 2)
		a.addActorEvent(EventManaRegen, g.tick + uint(combinedStats))
	}
}

func (a *Actor) addMana() {
	if a.mana < a.maxMana { a.mana++ }
}

func (a *Actor) addHealth() {
	if a.hp < a.maxHP { a.hp++ }
}

func (a *Actor) endBerserk(g *Game) {
	a.fg = a.pFg
	a.strg = a.pStrg
	a.intel = a.pIntel
	a.canEquip = true
	a.canRead = true
	a.canCast = true
	a.canQuaff = true
	a.fatigue(g)
	a.addActorEvent(EventEndFatigue, g.tick + uint(roll(R2d4)))
}

func (a *Actor) fatigue(g *Game) {
	a.pSpd = a.spd
	a.spd = a.spd/2
}

func (a *Actor) endFatigue(g *Game) {
	a.spd = a.pSpd
}