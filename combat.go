/*
 *  combat.go
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
	"math"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func (a *Actor) meleeDamage(t *Actor) int {
	return 0
}

func (a *Actor) dexMod() int {
	return (a.dex - 10) / 2 
}

func (a *Actor) strMod() int {
	return (a.strg - 10) / 2 
}

func resistMod(rLevel int) float64 {
	switch rLevel {
	case 0: return 1
	case 1: return .80  // TODO these should be randomized ranges
	case 2: return .60
	case 3: return .40
	case 4: return .20
	default: return 1
	}
}

func (a *Actor) meleeAttack(t *Actor) int {
	var weapon Item
	var dmgValue int
	isCrit := false
	attackRoll := roll(R1d20)
	if attackRoll == 20 || t.ai == SleepAI { 
		isCrit = true
	} else if attackRoll == 1 {
		return 0 // botch
	}
	if a.weapon == -1 {
		weapon = a.getFists()
	} else {
		weapon = a.inv[getInvPositionForSlotID(a.weapon, a.inv)].item
	}
	isRanged := weapon.isWeapRanged()
	if isRanged {
		dmgValue = roll(weapon.dmg) / 2
	} else if weapon.isDexMod() {
		dmgValue = roll(weapon.dmg) + a.dexMod() + weapon.enchant
	} else {
		dmgValue = roll(weapon.dmg) + a.strMod() + weapon.enchant
	}
	if isCrit { dmgValue = dmgValue * 2 }
	dmg := dmgValue - t.ac
	if dmg < 0 { dmg = 0 }
	//brand
	var brandDmg float64
	// TODO make seperate func for other attack types
	// TODO account for weakness to dmg types
	if weapon.brand != NoBrand && !isRanged {
		switch weapon.brand {
		case Flaming: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rFire)
		case Freezing: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rCold)
		case Electrocution: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rElec)
		//case Venom: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rPois)
		//case Vorpal:
		//case Speed:
		default: brandDmg = 0
		}
	}
	dmg = dmg + int(math.Round(brandDmg))
	t.takeDamage(dmg)
	if t.ai == SleepAI || t.ai == ConfusedAI {
		t.ai = t.defaultAI
	}
	return dmg
}


func (a *Actor) dropProjectile(item Item, pos Point, g *Game) {
	for j, i := range a.inv {
		if i.item == item {
			// if stack process that
			if i.item.stackable && i.item.charges > 1 {
				a.inv[j].item.charges -= 1
				item.setPos(pos.x, pos.y)
				item.charges = 1
				g.floors[g.cur].items = append(g.floors[g.cur].items, item)
			} else {
				a.inv[j].item.quivered = false
				g.emptyAltWeapon(a.inv[j].slot)
				g.dbg(fmt.Sprintf("dropping %v", a.inv[j].item.pname))
				copy(a.inv[j:], a.inv[j+1:])
				a.inv = a.inv[:len(a.inv)-1]
				//drop by adding to floor items
				item.setPos(pos.x, pos.y)
				g.floors[g.cur].items = append(g.floors[g.cur].items, item)
				a.quiver = -1
			}
		}
	}
}


func (a *Actor) rangedAttack(t *Actor, tPos Point, g *Game) int {
	var weapon Item
	//var ammo Item
	isCrit := false
	var dmgValue int
	var attackRoll int
	//attacker
	quiverInvID := getInvPositionForSlotID(a.quiver, a.inv)
	weaponInvID := getInvPositionForSlotID(a.weapon, a.inv)
	if a.inv[quiverInvID].item.category == Ammo {
		//you need the proper weapon to shoot this
		if a.weapon == -1 {
			g.addMessage("You can't use that without the proper weapon", tcell.ColorDefault)
			return 0
		}
		weapon = a.inv[weaponInvID].item
		if a.inv[quiverInvID].item.iname == AmmoArrow {
			if !weapon.isWeapBow() {
				g.addMessage("You need a bow to fire arrows", tcell.ColorDefault)
				return 0
			}
		} else if a.inv[quiverInvID].item.iname == AmmoBolt {
			if !weapon.isWeapCrossbow() {
				g.addMessage("You need a crossbow to fire bolts", tcell.ColorDefault)
				return 0
			}
		}
		a.dropProjectile(a.inv[quiverInvID].item, tPos, g) //throw/fire
		dmgValue = roll(weapon.dmg) + a.dexMod() + weapon.enchant
	} else if a.inv[quiverInvID].item.category == Weapon {
		//you dont need another weapon to shoot this
		weapon = a.inv[quiverInvID].item
		a.dropProjectile(weapon, tPos, g) //throw/fire
		dmgValue = roll(weapon.dmg) + a.dexMod() + weapon.enchant
	}
	if attackRoll == 20 || t.ai == SleepAI {
		isCrit = true
	} else if attackRoll == 1 {
		a.energy -= SpeedCost  //TODO change to dex modifier??
		return 0 // botch
	}
	if isCrit { dmgValue = dmgValue * 2 }
	dmg := dmgValue - t.ac
	if dmg < 0 { dmg = 0 }
	//brand
	var brandDmg float64
	if weapon.brand != NoBrand {
		switch weapon.brand {
		case Flaming: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rFire)
		case Freezing: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rCold)
		case Electrocution: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rElec)
		//case Venom: brandDmg = float64(((roll(weapon.dmg)/2) + weapon.enchant)) * resistMod(t.rPois)
		//case Vorpal:
		//case Speed:
		default: brandDmg = 0
		}
	}
	dmg = dmg + int(math.Round(brandDmg))
	t.takeDamage(dmg)
	a.energy -= SpeedCost //TODO change to dex mod??
	if t.ai == SleepAI || t.ai == ConfusedAI {
		t.ai = t.defaultAI
	}
	return dmg
}

func (t *Actor) takeDamage(dmg int) {
	if t.ai == SleepAI {
		t.ai = t.defaultAI
	}
	if dmg > 0 { t.hp -= dmg } 
	if t.hp <= 0 { t.death() }
}

func (s *Actor) getFists() Item {
	var i = Item {
		name: "fists",
		pname: "fists",
		slot: OneHand,
		category: Weapon,
		equipable: true,
		equipped: true,
		BUC: Uncursed,
		brand: NoBrand,
		dmg: R1d4,
	}
	return i
}

func (a *Actor) death() {
	a.alive = false
	a.ai = NoAI
	a.blocks = false
	a.glyph = '%'
	a.fg = tcell.ColorRed
	a.target = Point { -1, -1 }
}

// dice rolls
func rollDice(n, d int) int {
	roll := 0
	for i := 0; i < n; i++ {
		roll += getRandNum(d, 1)
	}
	return roll
}

func roll(d itemDmgRoll) int {
	switch d {
	case R1d2: return rollDice(1, 2)
	case R1d4: return rollDice(1, 4)
	case R2d4: return rollDice(2, 4)
	case R3d4: return rollDice(3, 4)
	case R4d4: return rollDice(4, 4)
	case R1d6: return rollDice(1, 6)
	case R2d6: return rollDice(2, 6)
	case R3d6: return rollDice(3, 6)
	case R4d6: return rollDice(4, 6)
	case R1d8: return rollDice(1, 8)
	case R2d8: return rollDice(2, 8)
	case R3d8: return rollDice(3, 8)
	case R4d8: return rollDice(4, 8)
	case R1d10: return rollDice(1, 10)
	case R2d10: return rollDice(2, 10)
	case R3d10: return rollDice(3, 10)
	case R4d10: return rollDice(4, 10)
	case R1d12: return rollDice(1, 12)
	case R2d12: return rollDice(2, 12)
	case R3d12: return rollDice(3, 12)
	case R4d12: return rollDice(4, 12)
	case R1d20: return rollDice(1, 20)
	case R2d20: return rollDice(2, 20)
	case R1d100: return rollDice(1, 100)
	default: return 0
	}
}