/*
 *  mob.go
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

// TODO put definitions in yaml or json file

import (
	"github.com/gdamore/tcell/v2"
)

type MobType int
const (
	Rat MobType = iota
	Troll
	Orc
	Dragon
	DragonRed
)

type MobLevel struct {
	level int
	mobtype MobType
	dieHit int
}

func getMobsForLevel(level int) []MobLevel {
	// this is hacky
	// the numbers are the dice roll range that will generate each mob
	// ex: level 1, 1d100 80 or below will be rat, 80-100 orc
	var mobs []MobLevel
	switch level {
	case 1:
		mobs = append(mobs, MobLevel { 1, Rat, 80 })
		mobs = append(mobs, MobLevel { 1, Orc, 100 })
	case 2:
		mobs = append(mobs, MobLevel { 2, Rat, 50 })
		mobs = append(mobs, MobLevel { 2, Orc, 100 })
	case 3:
		mobs = append(mobs, MobLevel { 3, Rat, 20 })
		mobs = append(mobs, MobLevel { 3, Orc, 80 })
		mobs = append(mobs, MobLevel { 3, Troll, 100 })
	case 4:
		mobs = append(mobs, MobLevel { 4, Orc, 60 })
		mobs = append(mobs, MobLevel { 4, Troll, 100 })
	case 5:
		mobs = append(mobs, MobLevel { 5, Orc, 20 })
		mobs = append(mobs, MobLevel { 5, Troll, 100 })
	case 6:
		mobs = append(mobs, MobLevel { 5, DragonRed, 20 })
		mobs = append(mobs, MobLevel { 5, DragonRed, 100 })
	default:
		mobs = append(mobs, MobLevel { -1, Rat, 100 })
	}
	return mobs
}

func getMob(roll, level int) MobType {
	mobs := getMobsForLevel(level)
	newmob := Rat
	for _, mob := range mobs {
        if roll < mob.dieHit {
			newmob = mob.mobtype
			//getMobSubType(newmob) // TODO create func that gets a dragon from type dragon, etc
			break
		}
    }
	return newmob
}

func makeMob(mob MobType, actors *[]Actor, x, y int) int {
	var m Actor
	switch mob {
	case DragonRed: m = makeDragonRed(x, y)
	case Troll: m = makeTroll(x, y)
	case Orc: m = makeOrc(x, y)
	case Rat: m = makeRat(x, y)
	default: return -1 }  // TODO make error here
	m.randAI()
	*actors = append(*actors, m)
	return len(*actors) - 1
}

func (m *Actor) randAI() {
	if roll(R1d100) <= 50 {
		m.ai = SleepAI
	} else {
		m.ai = WanderAI
	}
}

func makeDragonRed(x, y int) Actor {
	a := makeActorObject(x, y)
	a.name = "a red dragon"
	a.glyph = 'D'
	a.fg = tcell.ColorRed
	a.maxHP = 25
	a.hp = 25
	a.ac = 18
	a.strg = 18
	a.dex = 4
	a.intel = 12
	a.stealth = 0
	a.spd = 6
	i := makeItem(-1, -1)
	i.name = ""
	i.category = Weapon
	i.slot = OneHand
	i.equipable = true
	i.dmg = R2d12
	a.inv = append(a.inv, Inventory { 0, i})
	a.inv[0].item.equipped = true
	a.weapon = 0
	return a
}


func makeTroll(x, y int) Actor {
	a := makeActorObject(x, y)
	a.name = "a gross troll"
	a.glyph = 'T'
	a.fg = tcell.ColorOrange
	a.maxHP = 16
	a.hp = 16
	a.ac = 14
	a.strg = 16
	a.dex = 5
	a.intel = 5
	a.stealth = 0
	a.spd = 7
	i := makeItem(-1, -1)
	i.name = ""
	i.category = Weapon
	i.slot = OneHand
	i.equipable = true
	i.dmg = R1d12
	a.inv = append(a.inv, Inventory { 0, i})
	a.inv[0].item.equipped = true
	a.weapon = 0
	return a
}

func makeOrc(x, y int) Actor {
	a := makeActorObject(x, y)
	a.name = "an orc"
	a.glyph = 'o'
	a.fg = tcell.ColorGreen
	a.maxHP = 5
	a.hp = 5
	a.ac = 5
	a.strg = 5
	a.dex = 5
	a.intel = 5
	a.stealth = 5
	a.rFire = 1
	return a
}

func makeRat(x, y int) Actor {
	a := makeActorObject(x, y)
	a.name = "a rat"
	a.glyph = 'r'
	a.fg = tcell.ColorDefault
	a.maxHP = 5
	a.hp = 5
	a.ac = 2
	a.strg = 2
	a.dex = 20
	a.intel = 2
	a.stealth = 0
	a.rFire = 0

	i := makeItem(-1, -1)
	i.iname = WeapDexMod
	i.name = ""
	i.category = Weapon
	i.slot = OneHand
	i.equipable = true
	i.dmg = R1d6
	a.inv = append(a.inv, Inventory { 0, i})
	a.inv[0].item.equipped = true
	a.weapon = 0
	return a
}


func makeActorObject(x, y int) Actor {
	var actor = Actor {
		name: "mob",
		glyph: 'X',
		pos: Point { x, y },
		//lastKnownPos Point
		fg: tcell.ColorDefault,
		//pFg tcell.Color  //used when status changes fg color
		bg: tcell.ColorDefault,
		visible: false,
		visited: false,
		blocks: true,
		blockSight: false,
		maxHP: 10,
		hp: 10,
		maxMana: 10,
		mana: 10,
		ac: 10,
		pAc: 10,
		strg: 10,
		pStrg: 10,
		dex: 10,
		pDex: 10,
		intel: 10,
		pIntel: 10,
		spd: 10,
		pSpd: 10,
		energy: 0,
		stealth: 10,
		rCold: 0, 
		rFire: 0,
		rElec: 0,
		rPois: 0,
		rConf: 0,
		rSleep: 0,
		rMagic: 0,
		luck: 0,
		//spells
		target: Point { -1, -1 },
		alive: true,
		quiver: -1,
		weapon: -1,
		//ai: WanderAI,
		defaultAI: BasicAI,
		//inv
		//events
		canEquip: true,
		canRead: true,
		canCast: true,
	}
	return actor
}