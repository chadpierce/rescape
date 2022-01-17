/*
 *  action.go
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
	"math/rand"
    "time"

	"github.com/gdamore/tcell/v2"
)

func (g *Game) msgUnderfoot(heroPos Point) {
	var itemStr string
	for _, items := range g.floors[g.cur].items {
		if items.pos == heroPos {
			stackStr := ""
			if items.stackable && items.charges > 1 {
				stackStr = " (" + fmt.Sprint(items.charges) + ")"
			}
			if itemStr == "" {
				itemStr = "Here lies: " + items.name + stackStr
			} else {
				itemStr = itemStr + ", " + items.name + stackStr
			}
			
		}
	}
	for _, tile := range g.floors[g.cur].tiles {
		if tile.pos == heroPos && tile.name != "floor" {
			if itemStr == "" {
				itemStr = "Here lies: " + tile.name
			} else {
				itemStr = itemStr + ", " + tile.name
			}
		}
	}
	if itemStr != "" {
		g.addMessage(itemStr, tcell.ColorDefault)
	}
}

func (g *Game) invStackAdd(item Item) {
	for id, inv := range g.hero.inv {
		if inv.item.pname == item.pname {
			g.hero.inv[id].item.charges += item.charges
		}
	} 
}

func (g *Game) isHoldingStack(item Item) bool {
	isH := false
	for _, i := range g.hero.inv { if item.pname == i.item.pname { isH = true } }
	return isH
}

func (g *Game) getUnderfoot(item Item) bool {
	isPickup := false
	itemCnt := 0
	if item.stackable && g.isHoldingStack(item) {
		g.invStackAdd(item)
		isPickup = true
	} else {
		isAdded := g.hero.addItemInv(item)
		if isAdded == false {
			g.addMessage("You cannot carry anymore", tcell.ColorDefault)
			return false
		}
	}
	itemCnt += item.charges
	for j, i := range g.floors[g.cur].items {
		if i == item {
			// remove element wihout changing order
			copy(g.floors[g.cur].items[j:], g.floors[g.cur].items[j+1:])
			g.floors[g.cur].items = g.floors[g.cur].items[:len(g.floors[g.cur].items)-1]
			break
		}
	}
	cntStr := ""
	if itemCnt > 1 { cntStr = " (" + fmt.Sprint(itemCnt) + ")" }
	itemStr := "You picked up " + item.name + cntStr
	isPickup = true
	if itemStr != "" {
		g.addMessage(itemStr, tcell.ColorDefault)
	}
	// TODO use item weight as multiplier?
	if isPickup { g.hero.energy -= SpeedCost * .8 }
	return isPickup
}

func (g *Game) dropStack(item Item) bool {
	return false
}

func (g *Game) groundStackAdd(item Item) {
	for id, i := range g.floors[g.cur].items {
		if g.hero.pos == i.pos && item.pname == i.pname {
			g.floors[g.cur].items[id].charges += item.charges
		}
	}
}

func (g *Game) isStackOnGround(item Item) bool {
	isG := false
	for _, i := range g.floors[g.cur].items {
		if i.pname == item.pname && i.pos == item.pos { isG = true } 
	}
	return isG
}

//TODO make actor agnostic 
func (g *Game) dropInv(slotID int) bool {
	isDropped := false
	cntStr := ""
	itemCnt := 0
	var item Item
	var id int
	for j, i := range g.hero.inv {
		if i.slot == slotID {
			id = j
		}
	}
	if g.hero.inv[id].item.equipped || g.hero.inv[id].item.quivered { //|| (i.equipped && j == g.hero.quiver) {
		g.addMessage("You are using that.", tcell.ColorDefault)
	} else {
		item = g.hero.inv[id].item
		g.hero.energy -= SpeedCost * .6 // TODO use item weight as multiplier?
		
		copy(g.hero.inv[id:], g.hero.inv[id+1:])
		g.hero.inv[len(g.hero.inv)-1] = Inventory{}
		g.hero.inv = g.hero.inv[:len(g.hero.inv)-1]
		if item.stackable && g.isStackOnGround(item) {
			g.groundStackAdd(item)
			isDropped = true
			itemCnt += item.charges
		} else {
			item.setPos(g.hero.pos.x, g.hero.pos.y)
			g.floors[g.cur].items = append(g.floors[g.cur].items, item)
			itemCnt += item.charges
			isDropped = true
		}
	}
	if isDropped {
		g.emptyAltWeapon(slotID)
		if itemCnt > 1 {
			cntStr = " (" + fmt.Sprint(itemCnt) + ")" }
		itemStr := "You dropped " + item.name + cntStr
		if itemStr != "" {
				g.addMessage(itemStr, tcell.ColorDefault)
		}
	}
	return isDropped
}

func isBlocked(x, y int, actors []Actor, tiles []Tile) (bool, string) {
	for _, t := range tiles {
		if t.pos.x == x && t.pos.y == y && t.blocks == true {
			return true, string(t.name)
		} 
	}
	for _, a := range actors {
		if a.pos.x == x && a.pos.y == y && a.blocks == true {
			return true, string(a.name)
		} 
	}
	return false, ""
}

func isBlockedByID(x, y int, actors []Actor, tiles []Tile) (bool, int, int) {
	for i, t := range tiles {
		if t.pos.x == x && t.pos.y == y && t.blocks == true {
			return true, i, 0
		} 
	}
	for i, a := range actors {
		if a.pos.x == x && a.pos.y == y && a.blocks == true {
			return true, i, 1
		} 
	}
	return false, -1, -1
}

func heroAttack(g *Game, dx, dy int) {	
	for i, a := range g.floors[g.cur].actors {
		if a.pos.x == g.hero.pos.x + dx && a.pos.y == g.hero.pos.y + dy {
			dmg := g.hero.meleeAttack(&g.floors[g.cur].actors[i])
			g.dbg(fmt.Sprintf("attack for: %v -  mob hp: %v", dmg, a.hp))
		}
	}
}

// TODO make this a method
func heroMove(g *Game, dx, dy int) Point {
	g.hero.energy -= SpeedCost // TODO move this - only use energy if you move/attack
	isB, id, t := isBlockedByID(g.hero.pos.x + dx, g.hero.pos.y + dy, g.floors[g.cur].actors, g.floors[g.cur].tiles)
	if isB && t == 0 {
		message := fmt.Sprintf("You bump into the %v", g.floors[g.cur].tiles[id].name) 
		g.addMessage(message, tcell.ColorGreen)
	} else if isB && t == 1 {
		heroAttack(g, dx, dy)
	} else {
		if g.hero.pos.x + dx >= MaxWidth || 
		g.hero.pos.x + dx <= 0 || g.hero.pos.y + dy >= MaxHeight || g.hero.pos.y + dy <= 0 {
			g.addMessage("You hit a border wall", tcell.ColorGreen)
		} else {
			if g.hero.ai == ConfusedAI {
				grid := getBFSGrid(g)
				n := getPNeighbors(g.hero.pos.x, g.hero.pos.y, grid)
				n = append(n, Point { g.hero.pos.x, g.hero.pos.y })
				rand.Seed(time.Now().UnixNano())
				theMove := n[rand.Intn(len(n))]
				dx := theMove.x - g.hero.pos.x
				dy := theMove.y - g.hero.pos.y
				g.hero.pos.x += dx
				g.hero.pos.y += dy
			} else {
				g.hero.pos.x += dx
				g.hero.pos.y += dy
			}
		}
	}
	g.msgUnderfoot(g.hero.pos)
	
    return g.hero.pos
}

func (g *Game) goDownstairs() {
	g.cur++
	if g.cur >= len(g.floors) {
		g.genFloor()
	} else {
		for _, u := range g.floors[g.cur].tiles {
			if u.name == "upstair" { g.hero.pos = u.pos }
		}
	}
}

func (g *Game) goDownstairs2() {
	g.cur++
	if g.cur >= len(g.floors) {
		g.genFloorGroundLevel()
	} else {
		for _, u := range g.floors[g.cur].tiles {
			if u.name == "upstair" { g.hero.pos = u.pos }
		}
	}
}

func (g *Game) goUpstairs() {
	if g.cur == 0 {
		g.dbg("You can't leave...")
		return
	}
	g.cur--
	for _, d := range g.floors[g.cur].tiles {
		if d.name == "downstair" {
			g.hero.pos = d.pos
		}
	}
}