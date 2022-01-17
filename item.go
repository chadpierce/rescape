/*
 *  item.go
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
	"github.com/gdamore/tcell/v2"
)

type itemCategory int
const (
	Potion itemCategory = iota
	Scroll
	Book
	Weapon
	Armor
	Amulet
	Ring
	Tool
	Wand
	Ammo
	Shield
	Buckler
	NoCategory
	AllItems
)

type itemSlot int
const (
	NoSlot itemSlot = iota
	OneHand
	TwoHand
	OneHandRanged
	TwoHandRanged
	Head
	Chest
	Hands
	Feet
	Back
	Neck
	Finger
)

const (
	gPotion = '!'
	gScroll = '?'
	gBook = '+'
	gWeapon = ')'
	gArmor = '['
	gAmulet = '"'
	gRing = '='
	gTool = '('
	gWand = '/'
	gAmmo = ')'
	gNoGlyph = 'X'
)

type itemName int
const (
	VorpalBlade itemName = iota
	BattleAxe
	WeapSwordShort
	WeapFlailDire
	WeapDart
	WeapShortBow
	WeapCrossbow
	AmmoBolt
	AmmoArrow
	InfiniteJest
	BookRage
	RingStrength
	PotHeal
	PotConf
	PotMega
	ArmorChestLeather
	ArmorBootsLeather
	ArmorBootsIron
	ScrollBlink
	ShieldSmall
	ShieldLarge
	ShieldBuckler
	BlankItem
	WeapDexMod
)

type itemBrand int
const (
	NoBrand itemBrand = iota
	//Chopping
	Flaming
	Freezing
	Electrocution
	Venom
	Vorpal //extra dmg (if head pop it off?)
	Speed
	//Antimagic //? maybe
)

type itemDmgRoll int
const (
	R1d2 itemDmgRoll = iota
	R1d4
	R2d4
	R3d4
	R4d4
	R1d6
	R2d6
	R3d6
	R4d6
	R1d8
	R2d8
	R3d8
	R4d8
	R1d10
	R2d10
	R3d10
	R4d10
	R1d12
	R2d12
	R3d12
	R4d12
	R1d20
	R2d20
	R1d100
	NoDamage
)

type BUCStatus int
const (
	Uncursed BUCStatus = iota
	Blessed
	Cursed
)

func (g *Game) discovery(i *Item) {
	for _, d := range g.disc {
		if d.name == i.pname { i.name = i.pname}
	}
}

func (g *Game) makeItem(item itemName, items *[]Item, x, y int) {
	var i Item

	switch item {
	case VorpalBlade: i = makeVorpalBlade(x, y)
	case BattleAxe: i = makeBattleAxe(x, y)
	case WeapSwordShort: i = makeWeapSwordShort(x, y)
	case WeapFlailDire: i = makeWeapFlailDire(x, y)
	case WeapDart: i = makeWeapDart(x, y)
	case WeapShortBow: i = makeWeapShortBow(x, y)
	case WeapCrossbow: i = makeWeapCrossbow(x, y)
	case AmmoArrow: i = makeAmmoArrow(x, y)
	case InfiniteJest: i = makeInfiniteJest(x, y)
	case BookRage: i = makeBookRage(x, y)
	case RingStrength: i = makeRingStrength(x, y)
	case PotHeal: i = g.makePotHeal(x, y)
	case PotConf: i = g.makePotConf(x, y)
	case PotMega: i = g.makePotMega(x, y)
	case ArmorChestLeather: i = makeArmorChestLeather(x, y)
	case ArmorBootsLeather:  i = makeArmorBootsLeather(x, y)
	case ArmorBootsIron: i = g.makeArmorBootsIron(x, y)
	case ShieldBuckler: i = makeShieldBuckler(x, y)
	case ShieldSmall: i = makeShieldSmall(x, y)
	case ShieldLarge: i = makeShieldLarge(x, y)
	case ScrollBlink: i = g.makeScrollBlink(x, y)
	default: return }  //TODO make error here
	//g.hero.addItemInv(i) 
	*items = append(*items, i)
}

func (a *Actor) useItem(id int, g *Game, s tcell.Screen) bool {
	isU := false
	switch (g.hero.inv)[id].item.category {
	case Weapon: isU = a.useWeapon(id, g)
	case Ring: isU = a.useRing(id, g)
	case Potion: isU = a.usePotion(id, g)
	case Armor: isU = a.useArmor(id, g)
	case Scroll: isU = a.useScroll(id, g)
	case Book: isU = a.useBook(id, g, s)
	case Shield, Buckler: isU = a.useShield(id, g)
	}
	return isU
}

func (a *Actor) throwItem(id int, g *Game, s tcell.Screen) bool {
	isU := false
	var targetID int
	tPos := a.projectileTarget(s, g)
	for id, t := range g.floors[g.cur].actors {
		if t.pos == tPos {
			targetID = id
		}
	}
	switch (a.inv)[id].item.category {
	//case Weapon: isU = a.throwWeapon(id tPos g)
	case Potion: isU = a.throwPotion(id, targetID, g)
	}
	return isU
}

func getInvPositionForSlotID(slotID int, inv []Inventory) int {
	for j, i := range inv {
		if i.slot == slotID { return j }
	}
	return -1
}

func (a *Actor) quiverItem(slotId int, g *Game) bool {
	isU := false
	invID := getInvPositionForSlotID(slotId, a.inv)
	if a.quiver != -1 { 
		a.inv[invID].item.quivered = false
	}
	a.quiver = slotId
	a.inv[invID].item.quivered = true
	msgStr := fmt.Sprintf("%v ready to fire", a.inv[invID].item.name)
	g.addMessage(msgStr, tcell.ColorDefault)
	return isU
}

func (a *Actor) quiverEmpty(id int, g *Game) bool {
	isU := false
	if a.quiver == -1 {
		g.addMessage("Your quiver is already empty", tcell.ColorDefault)
	} else {
		a.quiver = -1
		a.inv[getInvPositionForSlotID(id, a.inv)].item.quivered = false
		g.addMessage("Your quiver has been emptied", tcell.ColorDefault)
	}
	return isU
}

func (i *Item) updateBUC(buc BUCStatus) {
	switch buc {
	case Cursed: i.BUC = Cursed
	case Uncursed: i.BUC = Uncursed
	case Blessed: i.BUC = Blessed
	}
	i.BUC = Cursed
}

func (i *Item) addBrand(b itemBrand) {
	switch b {
	case Flaming: i.brand = Flaming
	case Freezing: i.brand = Freezing
	case Electrocution: i.brand = Electrocution
	case Venom: i.brand = Venom
	case Vorpal: i.brand = Vorpal
	case Speed: i.brand = Speed
	}
}

func (i *Item) addEnchant(n int) {
	i.enchant += n
}

func (i *Item) updatePname(n string) {
	i.pname = n
}

///////////////////////////////////////////////////////////////////////////////
// WEAPONS

func (i *Item) isDexMod() bool {
	dexWeapons := []itemName{VorpalBlade,WeapDart,WeapDexMod}
	for _, d := range dexWeapons {
		if d == i.iname { return true }
	}
	return false
}

func (i *Item) isWeapBow() bool {
	w := []itemName{WeapShortBow}
	for _, d := range w {
		if d == i.iname { return true }
	}
	return false
}

func (i *Item) isWeapCrossbow() bool {
	w := []itemName{WeapCrossbow}
	for _, d := range w {
		if d == i.iname { return true }
	}
	return false
}

func (i *Item) isWeapRanged() bool {
	if i.isWeapBow() || i.isWeapCrossbow() {
		return true
	} else { 
		return false
	}
}

func (a *Actor) fireWeapon(s tcell.Screen, g *Game) bool {
	isTurn := false
	var targetID int
	if a.quiver == -1 {
		g.addMessage("You have nothing at the ready", tcell.ColorDefault)
		return false
	} 
	tPos := a.projectileTarget(s, g)
	cancel := Point { -1, -1 }
	if tPos == cancel { g.dbg("cancell");return isTurn}
	for id, t := range g.floors[g.cur].actors {
		if t.pos == tPos {
			targetID = id
			dmg := a.rangedAttack(&g.floors[g.cur].actors[targetID], tPos, g)
			g.addMessage(fmt.Sprintf("ranged attack for %v", fmt.Sprint(dmg)), tcell.ColorDefault)
			return true
		} 
	}
	weaponInvID := getInvPositionForSlotID(a.weapon, a.inv)
	quiverInvID := getInvPositionForSlotID(a.quiver, a.inv)
	if a.inv[quiverInvID].item.category == Ammo {
		if a.weapon == -1 {
			g.addMessage("You can't use that without the proper weapon", tcell.ColorDefault)
			return isTurn
		} else if a.inv[weaponInvID].item.iname == AmmoArrow {
			if !a.inv[weaponInvID].item.isWeapBow() {
				g.addMessage("You need a bow to fire arrows", tcell.ColorDefault)
				return isTurn
			}
		} else if a.inv[quiverInvID].item.iname == AmmoBolt {
			if !a.inv[quiverInvID].item.isWeapCrossbow() {
				g.addMessage("You need a crossbow to fire bolts", tcell.ColorDefault)
				return isTurn
			} 
		}
	} 
	a.dropProjectile(a.inv[getInvPositionForSlotID(a.quiver, a.inv)].item, tPos, g)
	isTurn = true
	a.energy -= SpeedCost
	return isTurn
}

func (a *Actor) useWeapon(id int, g *Game) bool {
	if a.inv[id].item.equipped {
		a.inv[id].item.equipped = false
		a.weapon = -1
		a.useWhichWeapon(id, g)
		return true
	} else {
		if a.canWield(id) {
			a.inv[id].item.equipped = true
			a.weapon = a.inv[id].slot
			a.useWhichWeapon(id, g)
			g.itemDiscovery(a.inv[id].item)
			return true
		} else {
			g.addMessage("Your cannot wield that", tcell.ColorDefault)
			return false
		}
	}
}

// func (a *Actor) useWhichWeapon(id int, g *Game) {
// 	switch a.inv[id].item.pname {
// 	case "Battle Axe": a.useBattleAxe(id, g)
// 	case "Vorpal Blade": a.useVorpalBlade(id, g)
// 	case "Dire Flail": a.useWeapFlailDire(id, g)
// 	}
// }

func (a *Actor) useWhichWeapon(id int, g *Game) {
	switch a.inv[id].item.iname {
	case BattleAxe: a.useBattleAxe(id, g)
	case VorpalBlade: a.useVorpalBlade(id, g)
	case WeapSwordShort: a.useWeapSwordShort(id, g)
	case WeapFlailDire: a.useWeapFlailDire(id, g)
	}
}

func makeVorpalBlade(x, y int) Item {
	i := makeItem(x, y)
	i.iname = VorpalBlade
	i.name = "Vorpal Blade"
	i.pname = "Vorpal Blade"
	i.glyph = gWeapon
	i.fg = tcell.ColorChartreuse
	i.category = Weapon
	i.slot = TwoHand
	i.equipable = true
	i.dmg = R1d8
	i.weight = 10
	//g.discovery(&i)
	return i
}

func (a *Actor) useVorpalBlade(id int, g *Game) {
	if a.inv[id].item.equipped {
		g.addMessage("You feel the urge to pop off some heads", tcell.ColorDefault)
	} else {
		g.addMessage("You regain your calm", tcell.ColorDefault)
	}
}

func makeBattleAxe(x, y int) Item {
	i := makeItem(x, y)
	i.iname = BattleAxe
	i.name = "Battle Axe"
	i.pname = "Battle Axe"
	i.glyph = gWeapon
	i.fg = tcell.ColorYellow
	i.category = Weapon
	i.slot = TwoHand
	i.equipable = true
	i.dmg = R2d12
	i.weight = 12
	return i
}

func (a *Actor) useBattleAxe(id int, g *Game) {
	if a.inv[id].item.equipped {
		g.addMessage("Chop Chop", tcell.ColorDefault)
	} else {
		g.addMessage("No more chopping", tcell.ColorDefault)
	}
}

func makeWeapSwordShort(x, y int) Item {
	i := makeItem(x, y)
	i.iname = WeapSwordShort
	i.name = "Short Sword"
	i.pname = "Short Sword"
	i.glyph = gWeapon
	i.fg = tcell.ColorPink
	i.category = Weapon
	i.slot = OneHand
	i.equipable = true
	i.dmg = R1d8
	i.weight = 6
	return i
}

func (a *Actor) useWeapSwordShort(id int, g *Game) {
	if a.inv[id].item.equipped {
		g.addMessage("You equip the short sword", tcell.ColorDefault)
	} else {
		g.addMessage("You sheath your short sword", tcell.ColorDefault)
	}
}

func makeWeapFlailDire(x, y int) Item {
	i := makeItem(x, y)
	i.iname = WeapFlailDire
	i.name = "Dire Flail"
	i.pname = "Dire Flail"
	i.glyph = gWeapon
	i.fg = tcell.ColorYellow
	i.category = Weapon
	i.slot = TwoHand
	i.equipable = true
	i.dmg = R4d12
	i.weight = 14
	return i
}

func (a *Actor) useWeapFlailDire(id int, g *Game) {
	if a.inv[id].item.equipped {
		g.addMessage("This is one big flail", tcell.ColorDefault)
	} else {
		g.addMessage("Your arms feel lighter", tcell.ColorDefault)
	}
}

func makeWeapDart(x, y int) Item {
	i := makeItem(x, y)
	i.iname = WeapDart
	i.name = "Dart"
	i.pname = "Dart"
	i.glyph = gWeapon
	i.stackable = true
	i.fg = tcell.ColorYellow
	i.category = Weapon
	i.slot = OneHand
	i.equipable = true
	i.charges = 5
	i.dmg = R1d8
	i.weight = 2
	return i
}

func (a *Actor) useWeapDart(id int, g *Game) {
	if a.inv[id].item.equipped {
		g.addMessage("Bullseye", tcell.ColorDefault)
	} else {
		g.addMessage("You put the darts away", tcell.ColorDefault)
	}
}

func makeWeapShortBow(x, y int) Item {
	i := makeItem(x, y)
	i.iname = WeapShortBow
	i.name = "+5 Short Bow"
	i.pname = "+5 Short Bow"
	i.glyph = gWeapon
	i.fg = tcell.ColorYellow
	i.category = Weapon
	i.slot = TwoHand
	i.equipable = true
	i.dmg = R1d10
	i.weight = 6
	return i
}

func (a *Actor) useWeapShortBow(id int, g *Game) {
	if a.inv[id].item.equipped {
		g.addMessage("Hawkeye", tcell.ColorDefault)
	} else {
		g.addMessage("You put the bow away", tcell.ColorDefault)
	}
}

func makeWeapCrossbow(x, y int) Item {
	i := makeItem(x, y)
	i.iname = WeapCrossbow
	i.name = "Crossbow"
	i.pname = "Crossbow"
	i.glyph = gWeapon
	i.fg = tcell.ColorYellow
	i.category = Weapon
	i.slot = TwoHandRanged
	i.equipable = true
	i.dmg = R1d10
	i.weight = 10
	return i
}

func (a *Actor) useWeapCrossbow(id int, g *Game) {
	if a.inv[id].item.equipped {
		g.addMessage("Crossbow at the ready", tcell.ColorDefault)
	} else {
		g.addMessage("You put the crossbow away", tcell.ColorDefault)
	}
}

///////////////////////////////////////////////////////////////////////////////
// AMMO

func makeAmmoArrow(x, y int) Item {
	i := makeItem(x, y)
	i.iname = AmmoArrow
	i.name = "Arrow"
	i.pname = "Arrow"
	i.glyph = gWeapon
	i.stackable = true
	i.fg = tcell.ColorYellow
	i.category = Ammo
	i.charges = 10
	i.enchant = 1
	i.dmg = R1d2
	return i
}

///////////////////////////////////////////////////////////////////////////////
// BOOKS

func (a *Actor) checkBookClass(bookClass, actorClass ActorClass) (bool, string) {
	var msg string
	if bookClass == actorClass {
		return true, ""
	} else {
		switch bookClass {
		case ClassFighter: msg = "You must be skilled in the arts of fighting to use this"
		case ClassRogue: msg = "You must be skilled in the shadow arts to use this"
		case ClassWizard: msg = "You must be skilled in the arcane to use this"

		}
		return false, msg
	}
}

func (a *Actor) useBook(id int, g *Game, s tcell.Screen) bool {
	isTurn := false
	g.itemDiscovery(a.inv[id].item)
	sids, spells, bname, bdesc := a.useWhichBook(id, g)
	if sids == nil {
		return isTurn
	} else {
		for _, known := range a.spells {
			for i, sid := range sids {
				if known == sid { spells[i] = spells[i] + " [known]"}
			}
		} 
		g.drawSpellBook(s, spells, bname, bdesc)
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			switch ev.Rune() {
			case 'a': if len(spells) > 0 { isTurn = a.memorizeSpell(sids[0], g, s) } 
			case 'b': if len(spells) > 1 { isTurn = a.memorizeSpell(sids[1], g, s) }
			case 'c': if len(spells) > 2 { isTurn = a.memorizeSpell(sids[2], g, s) }
			case 'd': if len(spells) > 3 { isTurn = a.memorizeSpell(sids[3], g, s) }
			case 'e': if len(spells) > 4 { isTurn = a.memorizeSpell(sids[4], g, s) }
			case 'f': if len(spells) > 5 { isTurn = a.memorizeSpell(sids[5], g, s) }
			}
			switch ev.Key() { case tcell.KeyEscape: isTurn = false }
		}
		return isTurn
	}
}

func (a *Actor) useWhichBook(id int, g *Game) ([]SpellID, []string, string, string) {
	var spells []string; var sids []SpellID; var bdesc string
	switch a.inv[id].item.iname {
	case InfiniteJest: sids, spells, bdesc = a.useInfiniteJest(id, g)
	case BookRage: sids, spells, bdesc = a.useBookRage(id, g)
	}
	return sids, spells, a.inv[id].item.pname, bdesc
}

func makeInfiniteJest(x, y int) Item {
	i := makeItem(x, y)
	i.iname = InfiniteJest
	i.name = "Infinite Jest"
	i.pname = "Infinite Jest"
	i.glyph = gBook
	i.fg = tcell.ColorYellow
	i.category = Book
	return i
}

func (a *Actor) useInfiniteJest(id int, g *Game) ([]SpellID, []string, string) {
	canReadThis, msg := g.hero.checkBookClass(ClassWizard, g.class)
	if !canReadThis { g.addMessage(msg, tcell.ColorDefault) ; return nil, nil, ""}
	var spells []string
	var spellID []SpellID
	bookDesc := "Mostly footnotes."
	spells = append(spells, "Blink"); spellID = append(spellID, Blink)
	spells = append(spells, "Zap"); spellID = append(spellID, Zap)
	return spellID, spells, bookDesc
}

func makeBookRage(x, y int) Item {
	i := makeItem(x, y)
	i.iname = BookRage
	i.name = "Book of Rage"
	i.pname = "Book of Rage"
	i.glyph = gBook
	i.fg = tcell.ColorYellow
	i.category = Book
	return i
}

func (a *Actor) useBookRage(id int, g *Game) ([]SpellID, []string, string) {
	canReadThis, msg := g.hero.checkBookClass(ClassFighter, g.class)
	if !canReadThis { g.addMessage(msg, tcell.ColorDefault) ; return nil, nil, ""}
	var spells []string
	var spellID []SpellID
	bookDesc := "The Book of Rage is short and to the point. It enables\n"
	bookDesc += "one who is skilled in the art of fighting to become enraged for a\n"
	bookDesc += "short time increasing strength and speed. When the effect wears off\n"
	bookDesc += "there will be a period of fatigue."
	spells = append(spells, "Berserk"); spellID = append(spellID, FighterBerserk)
	return spellID, spells, bookDesc
}

///////////////////////////////////////////////////////////////////////////////
// RINGS

func (a *Actor) useRing(id int, g *Game) bool {
	if a.inv[id].item.equipped {
		a.energy -= SpeedCost
		a.inv[id].item.equipped = false
		a.useWhichRing(id, g)
		g.addMessage("Your pull the ring off your finger", tcell.ColorDefault)
		return true
	} else {
		if a.canPutOnRing(id) {
			a.energy -= SpeedCost
			a.inv[id].item.equipped = true
			g.itemDiscovery(a.inv[id].item)
			a.useWhichRing(id, g)
			return true
		} else {
			g.addMessage("You feel too gaudy", tcell.ColorDefault)
			return false
		}
	}
}

func (a *Actor) useWhichRing(id int, g *Game) bool {
	isU := false
	switch a.inv[id].item.pname {
	case "Ring of Strength": isU = a.useRingStrength(id, g)
	}
	return isU
}

func makeRingStrength(x, y int) Item {
	i := makeItem(x, y)
	i.iname = RingStrength
	i.name = "Ring of htgnertS"
	i.pname = "Ring of Strength"
	i.glyph = gRing
	i.fg = tcell.ColorChartreuse
	i.category = Ring
	return i
}

func (a *Actor) useRingStrength(id int, g *Game) bool {
	strgBonus := 1 + a.inv[id].item.enchant
	if a.inv[id].item.equipped {
		a.strg += strgBonus
		g.addMessage("You feel stronger...", tcell.ColorDefault)
		// TODO if enchanted change message 'you feel much stronger' etc 
	} else {
		a.strg -= strgBonus
		g.addMessage("You feel weaker...", tcell.ColorDefault)
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////
// POTIONS

func (a *Actor) throwPotion(id, tid int, g *Game) bool {
	if a == &g.hero { g.itemDiscovery(a.inv[id].item) }
	g.floors[g.cur].actors[tid].useWhichPotion(id, a.inv[id].item.iname, a.inv[id].item.BUC, g)
	(a.inv)[id].item.charges--
	if (a.inv)[id].item.charges <= 0 {
		copy((a.inv)[id:], (a.inv)[id+1:])
		a.inv[len(a.inv)-1] = Inventory{}
		a.inv = (a.inv)[:len(a.inv)-1]
	}
	return true
}

func (a *Actor) usePotion(id int, g *Game) bool {
	if a == &g.hero { g.itemDiscovery(a.inv[id].item) }
	a.useWhichPotion(id, a.inv[id].item.iname, a.inv[id].item.BUC, g)
	(a.inv)[id].item.charges--
	if (a.inv)[id].item.charges <= 0 {
		copy((a.inv)[id:], (a.inv)[id+1:])
		a.inv[len(a.inv)-1] = Inventory{}
		a.inv = (a.inv)[:len(a.inv)-1]
	}
	a.energy -= 0.5 * SpeedCost
	return true
}

func (a *Actor) useWhichPotion(id int, potName itemName, buc BUCStatus, g *Game) {
	switch potName {
	case PotHeal: a.usePotHeal(buc, g)
	case PotConf: a.usePotConf(buc, g)
	case PotMega: a.usePotMega(buc, g)
	}
}

func (g *Game) makePotHeal(x, y int) Item {
	i := makeItem(x, y)
	i.iname = PotHeal
	i.name = "Potion of gnilaeH"
	i.pname = "Potion of Healing"
	i.glyph = gPotion
	i.stackable = true
	i.fg = tcell.ColorPink
	i.category = Potion
	i.charges = 5
	g.discovery(&i)
	return i
}

func (a *Actor) usePotHeal(buc BUCStatus, g *Game) {
	switch buc {
	case Uncursed:
		a.hp += 11
		if a == &g.hero { 
			g.addMessage("Drank a Potion of Healing", tcell.ColorDefault)
		} else {
			g.addMessage("Threw a Potion of Healing", tcell.ColorDefault)
		}
	case Blessed:
		a.maxHP += 11
		if a == &g.hero { 
			g.addMessage("Drank a Blessed Potion of Healing", tcell.ColorGreen)
		} else {
			g.addMessage("Threw a Blessed Potion of Healing", tcell.ColorGreen)
		}
	case Cursed:
		a.hp -= 11
		if a == &g.hero { 
			g.addMessage("Drank a Cursed Potion of Healing", tcell.ColorRed)
		} else {
			g.addMessage("Threw a Cursed Potion of Healing", tcell.ColorRed)
		}
	}
}

func (g *Game) makePotConf(x, y int) Item {
	i := makeItem(x, y)
	i.iname = PotConf
	i.name = "Potion of Fuseconsion"
	i.pname = "Potion of Confusion"
	i.glyph = gPotion
	i.stackable = true
	i.fg = tcell.ColorAqua
	i.category = Potion
	i.charges = 5
	g.discovery(&i)
	return i
}

func (a *Actor) usePotConf(buc BUCStatus, g *Game) {
	switch buc {
	case Uncursed:
		a.ai = ConfusedAI
		a.addActorEvent(EventBasicAI, g.tick + 5)
		if a == &g.hero { 
			g.addMessage("Drank a Potion of Confusion", tcell.ColorDefault)
		} else {
			g.addMessage("Threw a Potion of Confusion", tcell.ColorDefault)
		}
	case Blessed:
		a.ai = ConfusedAI
		if a == &g.hero { 
			g.addMessage("Drank a Blessed Potion of Confusion", tcell.ColorGreen)
		} else {
			g.addMessage("Threw a Blessed Potion of Confusion", tcell.ColorGreen)
		}
	case Cursed:
		a.ai = ConfusedAI
		if a == &g.hero { 
			g.addMessage("Drank a Cursed Potion of Confusion", tcell.ColorRed)
		} else {
			g.addMessage("Threw a Cursed Potion of Confusion", tcell.ColorRed)
		}
	}

}

func (g *Game) makePotMega(x, y int) Item {
	i := makeItem(x, y)
	i.iname = PotMega
	i.name = "Potion of MegaGood"
	i.pname = "Potion of MegaGood"
	i.glyph = gPotion
	i.stackable = true
	i.fg = tcell.ColorLightBlue
	i.category = Potion
	i.charges = 5
	g.discovery(&i)
	return i
}

func (a *Actor) usePotMega(buc BUCStatus, g *Game) {
	switch buc {
	case Uncursed:
		a.maxHP += 5
		a.hp = a.maxHP
		a.maxMana += 5
		a.mana = a.maxMana
		a.strg += 1
		a.intel += 1
		a.dex += 1
		//a.addActorEvent(EventBasicAI, g.tick + 5)
		if a == &g.hero { 
			g.addMessage("Drank a Potion of MEGA", tcell.ColorDefault)
		} else {
			g.addMessage("Threw a Potion of MEGA", tcell.ColorDefault)
		}
	case Blessed:
		a.ai = ConfusedAI
		if a == &g.hero { 
			g.addMessage("Drank a Blessed Potion of MEGA", tcell.ColorGreen)
		} else {
			g.addMessage("Threw a Blessed Potion of MEGA", tcell.ColorGreen)
		}
	case Cursed:
		a.ai = ConfusedAI
		if a == &g.hero { 
			g.addMessage("Drank a Cursed Potion of MEGA", tcell.ColorRed)
		} else {
			g.addMessage("Threw a Cursed Potion of MEGA", tcell.ColorRed)
		}
	}

}


///////////////////////////////////////////////////////////////////////////////
// SCROLLS

func (a *Actor) useScroll(id int, g *Game) bool {
	g.itemDiscovery(a.inv[id].item)
	a.useWhichScroll(id, g)
	(g.hero.inv)[id].item.charges--
	if (g.hero.inv)[id].item.charges <= 0 {
		copy((g.hero.inv)[id:], (g.hero.inv)[id+1:])
		g.hero.inv[len(g.hero.inv)-1] = Inventory{}
		g.hero.inv = (g.hero.inv)[:len(g.hero.inv)-1]
	}
	a.energy -= SpeedCost // TODO use int as speed multiplier?
	return true
}

func (a *Actor) useWhichScroll(id int, g *Game) bool {
	isU := false
	switch a.inv[id].item.pname {
	case "Scroll of Blinking": isU = a.useScrollBlink(id, g)
	}
	return isU
}

func (g *Game) makeScrollBlink(x, y int) Item {
	i := makeItem(x, y)
	i.iname = ScrollBlink
	i.name = "Scroll of Aoekghd"
	i.pname = "Scroll of Blinking"
	i.glyph = gScroll
	i.stackable = true
	i.fg = tcell.ColorOrange
	i.category = Scroll
	i.charges = 3
	g.discovery(&i)
	return i
}

func (a *Actor) useScrollBlink(id int, g *Game) bool {
	switch (a.inv)[id].item.BUC {
	case Uncursed:
		a.spellBlink(g)
		g.addMessage("You blink.", tcell.ColorDefault)
	case Blessed:
		a.spellBlink(g)
		g.addMessage("You blink.", tcell.ColorDefault)
	case Cursed:
		a.spellBlink(g)
		g.addMessage("You blink.", tcell.ColorDefault)
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////
// ARMOR

func (a *Actor) useArmor(id int, g *Game) bool {
	if a.inv[id].item.equipped {
		a.energy -= SpeedCost // TODO for armor this should be moved to item and based on 'weight'
		a.inv[id].item.equipped = false
		a.useWhichArmor(id, g)
		g.addMessage("You take it off", tcell.ColorDefault)
		return true
	} else {
		if a.canEquipArmor(id) {
			a.energy -= SpeedCost
			a.inv[id].item.equipped = true
			g.itemDiscovery(a.inv[id].item)
			a.useWhichArmor(id, g)
			return true
		} else {
			g.addMessage("You are already wearing something.", tcell.ColorDefault)
			return false
		}
	}
}

func (a *Actor) useWhichArmor(id int, g *Game) bool {
	isU := false
	switch a.inv[id].item.iname {
	case ArmorChestLeather: isU = a.useArmorChestLeather(id, g)
	case ArmorBootsLeather: isU = a.useArmorBootsLeather(id, g)
	case ArmorBootsIron: isU = a.useArmorBootsIron(id, g)
	//case "SpecialArmor1" - so many special items per game
	}
	return isU
}

func makeArmorChestLeather(x, y int) Item {
	i := makeItem(x, y)
	i.iname = ArmorChestLeather
	i.name = "Leather Armor"
	i.pname = "Leather Armor"
	i.glyph = gArmor
	i.fg = tcell.ColorTan
	i.category = Armor
	i.slot = Chest
	i.equipable = true
	return i
}

func (a *Actor) useArmorChestLeather(id int, g *Game) bool {
	ACBonus := 1 + a.dexMod() + a.inv[id].item.enchant
	if a.inv[id].item.equipped {
		a.ac += ACBonus
		g.addMessage("You put on the leather armor", tcell.ColorDefault)
		return true
	} else {
		a.ac -= ACBonus
		return true
	}
}

func makeArmorBootsLeather(x, y int) Item {
	i := makeItem(x, y)
	i.iname = ArmorBootsLeather
	i.name = "Leather Boots"
	i.pname = "Leather Boots"
	i.glyph = gArmor
	i.fg = tcell.ColorTan
	i.category = Armor
	i.slot = Feet
	i.equipable = true
	return i
}

func (a *Actor) useArmorBootsLeather(id int, g *Game) bool {
	ACBonus := 1 + a.dexMod() + a.inv[id].item.enchant
	if a.inv[id].item.equipped {
		a.ac += ACBonus
		g.addMessage("Your feel feel snug...", tcell.ColorDefault)
		return true
	} else {
		a.ac -= ACBonus
		return true
	}
}

func (g *Game) makeArmorBootsIron(x, y int) Item {
	i := makeItem(x, y)
	i.iname = ArmorBootsIron
	i.name = "Iron Boots"
	i.pname = "Iron Boots"
	i.glyph = gArmor
	i.fg = tcell.ColorPurple
	i.category = Armor
	i.slot = Feet
	i.equipable = true
	return i
}

func (a *Actor) useArmorBootsIron(id int, g *Game) bool {
	ACBonus := 3 + a.inv[id].item.enchant
	if a.inv[id].item.equipped {
		a.ac += ACBonus
		g.addMessage("You wish you had socks...", tcell.ColorDefault)
		return true
	} else {
		a.ac -= ACBonus
		return true
	}
}

///////////////////////////////////////////////////////////////////////////////
// Shields

func (a *Actor) useShield(id int, g *Game) bool {
	if a.inv[id].item.equipped {
		a.energy -= SpeedCost // TODO for armor this should be moved to item and based on 'weight'
		a.inv[id].item.equipped = false
		a.useWhichShield(id, g)
		g.addMessage("You put away your shield", tcell.ColorDefault)
		return true
	} else {
		if a.canEquipShield(id) {
			a.energy -= SpeedCost
			a.inv[id].item.equipped = true
			g.itemDiscovery(a.inv[id].item)
			a.useWhichShield(id, g)
			return true
		} else {
			g.addMessage("You can't equip that", tcell.ColorDefault)
			return false
		}
	}
}

func (a *Actor) useWhichShield(id int, g *Game) bool {
	isU := false
	switch a.inv[id].item.iname {
	case ShieldBuckler: isU = a.useShieldBuckler(id, g)
	case ShieldSmall: isU = a.useShieldSmall(id, g)
	case ShieldLarge: isU = a.useShieldLarge(id, g)
	}
	return isU
}

func makeShieldBuckler(x, y int) Item {
	i := makeItem(x, y)
	i.iname = ShieldBuckler
	i.name = "Buckler"
	i.pname = "Buckler"
	i.glyph = gArmor
	i.fg = tcell.ColorBrown
	i.category = Buckler
	i.slot = NoSlot
	i.equipable = true
	return i
}

func (a *Actor) useShieldBuckler(id int, g *Game) bool {
	ACBonus := 1 + a.inv[id].item.enchant
	if a.inv[id].item.equipped {
		a.ac += ACBonus
		g.addMessage("You strap the buckler on your arm", tcell.ColorDefault)
		return true
	} else {
		a.ac -= ACBonus
		return true
	}
}

func makeShieldSmall(x, y int) Item {
	i := makeItem(x, y)
	i.iname = ShieldSmall
	i.name = "Shield"
	i.pname = "Shield"
	i.glyph = gArmor
	i.fg = tcell.ColorBrown
	i.category = Shield
	i.slot = OneHand
	i.equipable = true
	return i
}

func (a *Actor) useShieldSmall(id int, g *Game) bool {
	ACBonus := 3 + a.inv[id].item.enchant
	if a.inv[id].item.equipped {
		a.ac += ACBonus
		g.addMessage("You hide behind the shield", tcell.ColorDefault)
		return true
	} else {
		a.ac -= ACBonus
		return true
	}
}

func makeShieldLarge(x, y int) Item {
	i := makeItem(x, y)
	i.iname = ShieldLarge
	i.name = "Large Shield"
	i.pname = "Large Shield"
	i.glyph = gArmor
	i.fg = tcell.ColorBrown
	i.category = Shield
	i.slot = OneHand
	i.equipable = true
	return i
}

func (a *Actor) useShieldLarge(id int, g *Game) bool {
	ACBonus := 5 + a.inv[id].item.enchant
	if a.inv[id].item.equipped {
		a.ac += ACBonus
		g.addMessage("You hide behind the large shield", tcell.ColorDefault)
		return true
	} else {
		a.ac -= ACBonus
		return true
	}
}


///////////////////////////////////////////////////////////////////////////////

func (a *Actor) canEquipShield(id int) bool {
	canEquip := true
	if a.inv[id].item.category == Shield && a.isTwoHandEquipped() {
		canEquip = false; return canEquip
	}
	for _, item := range a.inv {
		if item.item.equipped == true {
			if item.item.category == Shield || item.item.category == Buckler {
				canEquip = false; break
			}
		} 
	}
	return canEquip
}

func (a *Actor) isShieldEquipped() bool {
	isShield := false
	for _, inv := range a.inv {
		if inv.item.equipped && inv.item.category == Shield {
			isShield = true; break
		}
	}
	return isShield
}

func (a *Actor) isTwoHandEquipped() bool {
	isTwo := false
	for _, inv := range a.inv {
		if inv.item.equipped && inv.item.slot == TwoHand {
			isTwo = true; break
		}
	}
	return isTwo
}

func (a *Actor) canEquipArmor(id int) bool {
	canEquip := true
	for _, item := range a.inv {
		if item.item.equipped == true {
			if item.item.category == a.inv[id].item.category && 
			item.slot == a.inv[id].slot {
				canEquip = false
			}
		}
	}
	return canEquip
}

func (a *Actor) canWield(id int) bool {
	canWield := true
	if a.inv[id].item.weight > a.strg { return false }
	if a.inv[id].item.slot == TwoHand && a.isShieldEquipped() {
		canWield = false; return canWield
	}
	for _, i := range a.inv {
		if i.item.equipped && i.item.category == Weapon {
			canWield = false; break
		} 
	}
	return canWield
}

func (a *Actor) canPutOnRing(id int) bool {
	// you can wear 2 rings at once it doesnt matter which hand
	// as long as there is an open slot
	// TODO can you wear 2 of the same ring?
	canEquip := true
	ringsWorn := 0
	for _, i := range a.inv {
		if i.item.equipped == true && i.item.category == Ring {
			ringsWorn++
		}
		if ringsWorn == 2 { canEquip = false; break }
	}
	return canEquip
}

func getItemDescription(name string) string {
	vorpalBlade := "A blade that pops the head of anything it stabs (if it has one)."
	InfiniteJest := "Mostly footnotes"
	RingStrength := "Enhances the strength of the wearer."
	
	switch name {
	case "Vorpal Blade":
		return vorpalBlade
	case "Infinite Jest":
		return InfiniteJest
	case "Ring of Strength":
		return RingStrength
	default:
		return "ERROR: no description found!"
	}
}

func makeItem(x, y int) Item {
	var i = Item {
		iname: BlankItem,
		name: "Item",
		pname: "Item",
		glyph: gNoGlyph,
		pos: Point { x, y },
		stackable: false,
		fg: tcell.ColorDefault,
		bg: tcell.ColorDefault,
		visible: false,
		visited: false,
		blocks: false,
		blockSight: false,
		category: NoCategory,
		slot: NoSlot,
		equipable: true,
		equipped: false,
		quivered: false,
		BUC: Uncursed,
		identified: false,
		charges: 0,
		enchant: 0,
		brand: NoBrand,
		dmg: NoDamage,
		weight: 0,
	}
	return i
}

//{22 Dart Dart 41 {-1 -1} true 4294967307 0 false false false false 3 1 true false false 0 false 5 0 0 9 2} 