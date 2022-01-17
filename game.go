/*
 *  game.go
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
    "log"
	"sort"
	
	"github.com/gdamore/tcell/v2"
)

// Playable area max width and height
const (
    MaxWidth int = 80
    MaxHeight int = 25
	MaxInvHero int = 26
	MaxInvMob int = 8
	SpeedCost = 10
)

type GameState int
const (
    Playing GameState = iota
    Menu
	Autopilot
	Confusion
    Target
	Over
)

type Game struct {
	tick uint
	debugMode bool
	state GameState
	cur int  	//current floor
	hero Actor
	altWeapons []int
	floors []Floor
	msg []Message
	disc []DiscoveredItem
	class ActorClass
}

type ActorClass int
const (
	ClassFighter ActorClass = iota
	ClassRogue
	ClassWizard
	ClassBeast
	ClassNone
)

type Floor struct {
	name string
	tiles []Tile
	actors []Actor
	items []Item
	fovDistance int
}

type Actor struct {
	name string
	glyph rune
	pos Point
	lastKnownPos Point
	fg tcell.Color
	pFg tcell.Color  //used when status changes fg color
	bg tcell.Color
	visible bool
	visited bool
	blocks bool
	blockSight bool
	maxHP int
	hp int
	maxMana int
	mana int
	ac int
	pAc int
	strg int
	pStrg int  //p status used for temp statuses
	dex int
	pDex int
	intel int
	pIntel int
	spd int
	pSpd int
	energy int
	stealth int
	rCold int
	rFire int
	rElec int
	rPois int
	rConf int
	rSleep int
	rMagic int
	luck int
	spells []SpellID
	target Point
	alive bool
	ai actorAI
	defaultAI actorAI
	quiver int
	weapon int
	inv []Inventory
	events []Event
	canEquip bool
	canRead bool
	canCast bool
	canQuaff bool
	state []actorState
}

type actorState int
const (
	AStateNone actorState = iota
	AStateBerserk
	AStateWhirlwind
	AStateRiposte
	AStateConcus
	AStateDoubleSwing
	AStateDoubleThrow
	AStateHide
	AStateFrozen
	AStateInvis
)

type Event struct {
	eventType eventType
	tick uint
}

type Item struct {
	iname itemName
	name string
	pname string
	glyph rune
	pos Point
	stackable bool
	fg tcell.Color
	bg tcell.Color
	visible bool
	visited bool
	blocks bool
	blockSight bool
	category itemCategory
	slot itemSlot
	equipable bool
	equipped bool
	quivered bool
	BUC BUCStatus
	identified bool
	charges int
	enchant int
	brand itemBrand
	dmg itemDmgRoll
	weight int
}

type Tile struct {
	name string
	pos Point
	glyph rune
	fg tcell.Color
	bg tcell.Color
	blocks bool
	blockSight bool
	visible bool
	visited bool
}

type Point struct {
	x, y int
}

type Inventory struct {
	slot int
	item Item
}

type Message struct {
	text string
	fg tcell.Color
}

type DiscoveredItem struct {
	name string
	category itemCategory
}

func (a *Actor) setPos(x, y int) {
	a.pos.x = x
	a.pos.y = y
}

func (i *Item) setPos(x, y int) {
	i.pos.x = x
	i.pos.y = y
}

func (a *Actor) getNextAvailInvSlot() int {
	for j := 0; j <= MaxInvHero; j++ {
		isEmpty := true
		for _, i := range a.inv {
			if i.slot == j { isEmpty = false }
		}
		if isEmpty { return j }
	}
	return -1
}

func (a *Actor) addItemInv(item Item) bool {
	if len(a.inv) > MaxInvHero {
		return false
	} else {
		slot := a.getNextAvailInvSlot()
		a.inv = append(a.inv, Inventory { slot, item })
	}
	return true
}

func (a *Actor) reorderInv() {
	sort.SliceStable(a.inv, func(i, j int) bool {
		return a.inv[i].slot < a.inv[j].slot
	})
}

func (g *Game) emptyAltWeapon(slotID int) {
	if g.altWeapons[0] == slotID {
		g.altWeapons[0] = -1
	} else if g.altWeapons[1] == slotID {
		g.altWeapons[1] = -1
	} else if g.altWeapons[2] == slotID {
		g.altWeapons[2] = -1
	}
}

func (g *Game) swapWeapons() {
	origWeap := g.hero.weapon
	origQuiv := g.hero.quiver
	origShield := -1
	for _, i := range g.hero.inv {
		if (i.item.category == Shield && i.item.equipped) ||
		(i.item.category == Buckler && i.item.equipped) {
			origShield = i.slot
		}
	}
	if g.hero.weapon != -1 {
		g.hero.inv[getInvPositionForSlotID(g.hero.weapon, g.hero.inv)].item.equipped = false
	}
	if g.hero.quiver != -1 {
		g.hero.inv[getInvPositionForSlotID(g.hero.quiver, g.hero.inv)].item.quivered = false
	}
	if origShield != -1 {
		g.hero.useShield(getInvPositionForSlotID(origShield, g.hero.inv), g)
	}
	g.hero.weapon = g.altWeapons[0]
	if g.hero.weapon != -1 {
		g.hero.inv[getInvPositionForSlotID(g.hero.weapon, g.hero.inv)].item.equipped = true
	}
	g.hero.quiver = g.altWeapons[1]
	if g.hero.quiver != -1 {
		g.hero.inv[getInvPositionForSlotID(g.hero.quiver, g.hero.inv)].item.quivered = true
	}
	if g.altWeapons[2] != -1 {
		g.hero.useShield(getInvPositionForSlotID(g.altWeapons[2], g.hero.inv), g)
	}
	g.altWeapons[0] = origWeap
	g.altWeapons[1] = origQuiv
	g.altWeapons[2] = origShield
}

func (g *Game) initHeroClass() {
	// TODO init starting gear here
	switch g.class {
	case ClassFighter: g.hero.strg += 2
	case ClassRogue: g.hero.dex += 2
	case ClassWizard: g.hero.intel += 2 
	}
}


func (g *Game) initGame() {
	    // empty messages are a hack to not pass an empty slice to draw func
		g.addMessage("", tcell.ColorLightBlue)
		g.addMessage("", tcell.ColorLightBlue)
		g.addMessage("", tcell.ColorLightBlue)
		g.addMessage("", tcell.ColorLightBlue)
		g.addMessage("You begin your escape...", tcell.ColorLightBlue)
		// this is a hack to make adding items to this slice easier
		g.disc = append(g.disc, DiscoveredItem {"placeholder", Amulet})
		g.hero = Actor {
				name: "Hero",
				glyph: '@',
				pos: Point { -1, -1 },
				fg: tcell.ColorDefault,
				bg: tcell.ColorDefault,
				visible: true,
				blocks: true,
				blockSight: false,
				maxHP: 10,
				hp: 10,
				maxMana: 5,
				mana: 3, 
				ac: 10,
				strg: 10,
				dex: 10,
				intel: 10,
				spd: 10,
				energy: 0,
				stealth: 10,
				target: Point {-1, -1},
				alive: true,
				ai: NoAI,
				quiver: -1,
				weapon: -1,
				canEquip: true,
				canRead: true,
				canCast: true,
				canQuaff: true,
		}
		g.altWeapons = []int{-1, -1, -1}
		g.class = ClassFighter
		g.initHeroClass()

		// g.hero.inv = append(g.hero.inv, Inventory { 0, makeVorpalBlade(-1, -1) } )
		g.hero.inv = append(g.hero.inv, Inventory { 0, makeBattleAxe(-1, -1) } )
		g.hero.inv = append(g.hero.inv, Inventory { 1, makeInfiniteJest(-1, -1) } )
		g.hero.inv = append(g.hero.inv, Inventory { 2, g.makePotHeal(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 4, g.makePotConf(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 5, makeArmorBootsLeather(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 6, makeWeapDart(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 7, makeVorpalBlade(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 8, makeWeapDart(-1, -1) } )
		g.hero.inv = append(g.hero.inv, Inventory { 3, makeWeapShortBow(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 10, makeShieldSmall(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 11, makeShieldBuckler(-1, -1) } )
		// g.hero.inv = append(g.hero.inv, Inventory { 12, makeBookRage(-1, -1) } )
		g.hero.inv = append(g.hero.inv, Inventory { 4, makeWeapSwordShort(-1, -1) } )
		g.hero.inv = append(g.hero.inv, Inventory { 5, makeArmorChestLeather(-1, -1 ) } )
		
		g.genFloor()
		g.makeItem(AmmoArrow, &g.floors[g.cur].items, g.hero.pos.x, g.hero.pos.y - 1)

		// g.makeItem(RingStrength, &g.floors[g.cur].items, g.hero.pos.x + 1, g.hero.pos.y + 1)
		// g.makeItem(RingStrength, &g.floors[g.cur].items, g.hero.pos.x + 0, g.hero.pos.y + 1)
		// g.makeItem(RingStrength, &g.floors[g.cur].items, g.hero.pos.x - 1, g.hero.pos.y + 1)
		// g.makeItem(ScrollBlink, &g.floors[g.cur].items, g.hero.pos.x -1 , g.hero.pos.y - 1)
		// g.makeItem(WeapFlailDire, &g.floors[g.cur].items, g.hero.pos.x + 0 , g.hero.pos.y - 1)

}

func initScreen() tcell.Screen {
	// Initialize tcell screen
	defStyle := 
        tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
    s, err := tcell.NewScreen()
    if err != nil {
        log.Fatalf("%+v", err)
    }
    if err := s.Init(); err != nil {
        log.Fatalf("%+v", err)
    }
    s.SetStyle(defStyle)
    s.Clear()
	return s
}

func (g *Game) addMessage(text string, color tcell.Color) {
	var newMsg = Message { text: text, fg: color }	
	g.msg = append(g.msg, newMsg)
}

func (g *Game) dbg(text string) {
	//var newMsg = Message { text: text, fg: tcell.ColorRed }	
	//g.msg = append(g.msg, newMsg)
	g.addMessage(text, tcell.ColorRed)
}

func (g *Game) itemDiscovery(i Item) {
	isDisc := false
	for _, d := range g.disc {
		if i.pname == d.name { isDisc = true }
	}
	if isDisc == false {
		g.disc = append(g.disc, DiscoveredItem{ i.pname, i.category})
		i.name = i.pname
		//make every existing item show pname
		for j, invItem := range g.hero.inv {
			if invItem.item.pname == i.pname {
				g.hero.inv[j].item.name = g.hero.inv[j].item.pname
			}
		}
		for _, floor := range g.floors {
			for j, _ := range floor.items {
				if floor.items[j].pname == i.pname {
					floor.items[j].name = floor.items[j].pname
				}
			}
		}
	}
}

func (t Point) targetToTile(g rune, tiles []Tile) Point {
	for _, tile := range tiles {
		if tile.glyph == g {
			t = tile.pos
		}
	}
	return t
}

func target(s tcell.Screen, g *Game) {
	t := g.hero.pos
	drawTarget(t, s, *g, "", 0)  // initialize target ui
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
			case '>': t = t.targetToTile('>', g.floors[g.cur].tiles)
			case '<': t = t.targetToTile('<', g.floors[g.cur].tiles)
			}
			switch ev.Key() {
			case tcell.KeyEscape: return
			case tcell.KeyTab: g.hero.pos = t; return
			case tcell.KeyEnter: g.hero.target = t; g.state = Autopilot; return
			case tcell.KeyLeft: t.x += -1; t.y += 0
			case tcell.KeyDown: t.x += 0; t.y += 1
			case tcell.KeyUp: t.x += 0; t.y += -1
			case tcell.KeyRight: t.x += 1; t.y += 0
			}
		}

		var itemStr string
		for _, items := range g.floors[g.cur].items {
			if items.pos == t {
				if itemStr == "" {
					itemStr = "Here lies: " + items.name
				} else {
					itemStr = itemStr + ", " + items.name
				}
			}
		}
		for _, tile := range g.floors[g.cur].tiles {
			if tile.pos == t && tile.name != "floor" {
				if itemStr == "" {
					itemStr = "Here lies: " + tile.name
				} else {
					itemStr = itemStr + ", " + tile.name
				}
			}
		}
		for _, a := range g.floors[g.cur].actors {
			if a.pos == t && a.name != "floor" {
				if itemStr == "" {
					aAI := ""
					if a.ai == SleepAI {
						aAI = "(sleeping)"
					} else if a.ai == ConfusedAI { 
						aAI = "(confused)"
					}
					if a.weapon == -1 || a.inv[a.weapon].item.name == "" {
						itemStr = "Here lies: " + a.name  
					} else {
						itemStr = "Here lies: " + a.name + " [" + a.inv[a.weapon].item.name + "]"
					}
					if aAI != "" {
						itemStr = itemStr + " " + aAI
					}
				} else {
					itemStr = itemStr + ", " + a.name
				}
			}
		}
		drawTarget(t, s, *g, itemStr, 0)
	}
}


func (a *Actor) projectileTarget(s tcell.Screen, g *Game) Point {
	t := a.pos
	drawTarget(t, s, *g, "", 2)  // initialize target ui
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
			case 'f': return t 
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
		drawTarget(t, s, *g, itemStr, 2)
	}
}

func (g *Game) heroDeath() {
	g.dbg("The game is over.")
}
