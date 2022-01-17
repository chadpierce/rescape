/*
 *  menu.go
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
    //"strconv"
	
	"github.com/gdamore/tcell/v2"
)

const (
	CancelRune = '%'
)

func menuInv(s tcell.Screen, g *Game) {
	drawInv("Inventory", s, g.hero.inv)
	quit := func() { return }
	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
			quit()
		} else if ev.Key() == tcell.KeyCtrlL {
			s.Sync()
		} else if ev.Rune() == 'a' {
			fmt.Println("you selected a")
		} else {
			fmt.Println("bad cmd")
		}
	}
}

func (g *Game) menuWield(s tcell.Screen) bool {
	isTurn := false
	var weap []Inventory
	var itemID []int
	var invID []int
	if !g.hero.canEquip { g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isTurn }
	for id, i := range g.hero.inv {
		if i.item.category == Weapon {
			weap = append(weap, i)
			itemID = append(itemID, i.slot)
			invID = append(invID, id)
		}
	}
	drawInv("Wield something", s, weap)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		isTurn = g.menuUse(s, itemID, invID, sel)
	}
	return isTurn
}

func (g *Game) menuEquipArmor(s tcell.Screen) bool {
	isTurn := false
	var armor []Inventory
	var itemID []int
	var invID []int
	if !g.hero.canEquip { 
		g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isTurn }
	for id, i := range g.hero.inv {
		if i.item.category == Armor || i.item.category == Shield || 
		i.item.category == Buckler {
			armor = append(armor, i)
			itemID = append(itemID, i.slot)
			invID = append(invID, id)
		}
	}
	drawInv("Equip Armor", s, armor)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		isTurn = g.menuUse(s, itemID, invID, sel)
	}
	return isTurn
}

func (g *Game) menuQuaff(s tcell.Screen) bool {
	isTurn := false
	var pot []Inventory
	var itemID []int
	var invID []int
	if !g.hero.canQuaff { g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isTurn }
	for id, i := range g.hero.inv {
		if i.item.category == Potion {
			pot = append(pot, i)
			itemID = append(itemID, i.slot)
			invID = append(invID, id)
		}
	}
	drawInv("Quaff Potion", s, pot)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		isTurn = g.menuUse(s, itemID, invID, sel)
	}
	return isTurn
}

func (a *Actor) menuThrow(s tcell.Screen, g *Game) bool {
	isTurn := false
	var itemID []int
	var invID []int
	var items []Inventory
	for id, i := range a.inv {
		if i.item.category == Potion || i.item.category == Weapon {
			items = append(items, i)
			itemID = append(itemID, i.slot)
			invID = append(invID, id)
		}
	}
	drawInv("Throw something", s, items)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		noItemStr := "There is nothing there."
		isInList, _invID := isItemInList(itemID, invID, sel)
		if isInList {
			isTurn = a.throwItem(_invID, g, s)
		} else {
			g.addMessage(noItemStr, tcell.ColorDefault)
		}
	}
	return isTurn
}

func (a *Actor) menuQuiver(s tcell.Screen, g *Game) bool {
	isTurn := false
	var itemID []int
	var invID []int
	var items []Inventory
	if !g.hero.canEquip { g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isTurn }
	for id, i := range a.inv {
		if i.item.category == Ammo || (i.item.category == Weapon && !i.item.isWeapRanged()) {
			items = append(items, i)
			itemID = append(itemID, i.slot)
			invID = append(invID, id)
		}
	}
	drawInv("Ready ranged weapon", s, items)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		if sel == '-' { 
			isTurn = a.quiverEmpty(a.quiver, g) 
		} else {
			qid := getItemForRune(sel, items)
			if qid != -1 {
				isTurn = a.quiverItem(qid, g)
			} else {
				g.addMessage("That won't work", tcell.ColorDefault)
			}
		}
	}
	return isTurn
}

func (g *Game) menuJewelry(s tcell.Screen) bool {
	isTurn := false
	var jewelry []Inventory
	var itemID []int
	var invID []int
	if !g.hero.canEquip {
		g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isTurn }
	for id, i := range g.hero.inv {
		if i.item.category == Ring || i.item.category == Amulet {
			jewelry = append(jewelry, i)
			itemID = append(itemID, i.slot)
			invID = append(invID, id)
		}
	}
	drawInv("Put on Jewelry", s, jewelry)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		isTurn = g.menuUse(s, itemID, invID, sel)
	}
	return isTurn
}

func (g *Game) menuRead(s tcell.Screen) bool {
	isTurn := false
	var items []Inventory
	var itemID []int
	var invID []int
	if !g.hero.canRead { g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isTurn }
	for id, i := range g.hero.inv {
		if i.item.category == Scroll || i.item.category == Book {
			items = append(items, i)
			itemID = append(itemID, i.slot)
			invID = append(invID, id)
		}
	}
	drawInv("Read", s, items)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		isTurn = g.menuUse(s, itemID, invID, sel)
	}
	return isTurn
}

func (a *Actor) menuCastSpell(s tcell.Screen, g *Game) bool {
	isCast := false
	noSpellStr := "You don't know any spells"
	if !g.hero.canCast { g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isCast }
	if len(a.spells) < 1 {
		g.addMessage(noSpellStr, tcell.ColorDefault); return isCast
	} 
	drawKnownSpells(s, a.spells, "Cast a spell")
	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		switch ev.Rune() {
			case 'a': isCast = a.castSpell(a.spells[0], s, g)
			case 'b': if len(a.spells) > 1 { isCast = a.castSpell(a.spells[1], s, g)
				} else { g.addMessage(noSpellStr, tcell.ColorDefault) }
			case 'c': if len(a.spells) > 2 { isCast = a.castSpell(a.spells[2], s, g)
				} else { g.addMessage(noSpellStr, tcell.ColorDefault) }
		}
	}
	return isCast
}

func menuDrop(s tcell.Screen, g *Game) bool {
	isTurn := false
	if len(g.hero.inv) < 1 {
		g.addMessage("You have nothing to drop.", tcell.ColorDefault)
		return isTurn
	} 
	var itemID []int
	var invID []int
	if !g.hero.canEquip { g.addMessage("You can't do that right now", tcell.ColorDefault)
		return isTurn }
	for id, i := range g.hero.inv {
		itemID = append(itemID, i.slot)
		invID = append(invID, id)
	}
	drawInv("Drop", s, g.hero.inv)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		slotID := getItemForRune(sel, g.hero.inv)
		if slotID != -1 {
			isTurn = g.dropInv(slotID)
			g.hero.reorderInv()
		} else {
			g.addMessage("That won't work", tcell.ColorDefault)
		}
	}
	return isTurn
}

func menuPickup(s tcell.Screen, g *Game) bool {
	var items []Item
	if !g.hero.canEquip { g.addMessage("You can't do that right now", tcell.ColorDefault)
		return false }
	for _, i := range g.floors[g.cur].items {
		if i.pos == g.hero.pos {
			items = append(items, i)
		}
	}
	if len(items) < 1 {
		g.addMessage("There is nothing here.", tcell.ColorDefault)
		return false
	} else if len(items) == 1 {
		g.getUnderfoot(items[0])
		g.hero.reorderInv()
		return true
	} else if len(g.hero.inv) >= MaxInvHero {  // TODO make this MaxInvHero
		g.addMessage("You are holding too much.", tcell.ColorDefault)
		return false
	}
	drawItemList("Pickup", s, items)
	sel := g.menuSelect(s)
	if sel != CancelRune {
		noItemStr := "There is nothing there."
		if int(sel-97) <= len(items)-1 {
			g.getUnderfoot(items[int(sel) - 97])
			g.hero.reorderInv()
			return true
		} else {
			g.addMessage(noItemStr, tcell.ColorDefault)
		}
	}
	return false
}

func menuDisc(s tcell.Screen, g *Game) {
	drawDisc("Discovered Items", s, g.disc)
	quit := func() { return }
	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
			quit()
		} else if ev.Key() == tcell.KeyCtrlL {
			s.Sync()
		} else if ev.Rune() == 'a' {
			fmt.Println("you selected a")
		} else {
			fmt.Println("bad cmd")
		}
	}
}

func (g *Game) menuReorder(s tcell.Screen) {
	//TODO reorder inv, skills/spells
	var firstSpell SpellID
	var secondSpell SpellID
	var firstSlot int
	var secondSlot int
	noSpellStr := "You don't know any spells"
	if len(g.hero.spells) < 1 {
		g.addMessage(noSpellStr, tcell.ColorDefault); return
	} 
	drawKnownSpells(s, g.hero.spells, "Swap which spell?")
	sel := g.menuSelect(s)
	if sel != CancelRune {
		noItemStr := "There is nothing there."
		if int(sel-97) <= len(g.hero.spells)-1 {
			firstSpell = g.hero.spells[int(sel)-97]
			firstSlot = int(sel)-97
		} else {
			g.addMessage(noItemStr, tcell.ColorDefault)
		}
	} else { return }
	drawKnownSpells(s, g.hero.spells, "Swap it with what?")
	sel = g.menuSelect(s)
	if sel != CancelRune {
		noItemStr := "There is nothing there."
		if int(sel-97) <= len(g.hero.spells)-1 {
			secondSpell = g.hero.spells[int(sel)-97]
			secondSlot = int(sel)-97
		} else {
			g.addMessage(noItemStr, tcell.ColorDefault)
		}
	} else { return }
	//gotta make this swap work!
	g.hero.spells[firstSlot] = secondSpell
	g.hero.spells[secondSlot] = firstSpell
}

func getItemForRune(sel rune, items []Inventory) int {
	for _, i := range items {
		if rune(i.slot + 97) == sel {
			return i.slot
		}
	}
	return -1
}

func isItemInList(itemID, invID []int, r rune) (bool, int) {
	isInList := false
	_invID := -1
	for id, i := range itemID {
		if rune(i + 97) == r {
			isInList = true
			_invID = invID[id]
		}
	}
	return isInList, _invID
}

func (g *Game) menuUse(s tcell.Screen, itemID, invID []int, r rune) bool {
	noItemStr := "There is nothing there."
	isTurn := false
	isInList, _invID := isItemInList(itemID, invID, r)
	if isInList {
		isTurn = g.hero.useItem(_invID, g, s)
	} else {
		g.addMessage(noItemStr, tcell.ColorDefault)
	}
	return isTurn
}

func (g *Game) menuSelect(s tcell.Screen) rune {
	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		switch ev.Rune() {
		case 'a': return 'a'
		case 'b': return 'b'
		case 'c': return 'c'
		case 'd': return 'd'
		case 'e': return 'e'
		case 'f': return 'f'
		case 'g': return 'g'
		case 'h': return 'h'
		case 'i': return 'i'
		case 'j': return 'j'
		case 'k': return 'k'
		case 'l': return 'l'
		case 'm': return 'm'
		case 'n': return 'n'
		case 'o': return 'o'
		case 'p': return 'p'
		case 'q': return 'q'
		case 'r': return 'r'
		case 's': return 's'
		case 't': return 't'
		case 'u': return 'u'
		case 'v': return 'v'
		case 'w': return 'w'
		case 'x': return 'x'
		case 'y': return 'y'
		case 'z': return 'z'
		case '-': return '-'
		default: return CancelRune
		}
	}
	return CancelRune
}