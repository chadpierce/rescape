/*
 *  rand.go
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
	"math/rand"
    "time"
	//"fmt"
	//"strings"
	//"github.com/gdamore/tcell/v2"
)

func getRandNum(max, min int) int {
    rand.Seed(time.Now().UnixNano()) 
    r := min + rand.Intn(max - min + 1) //+ min
    return r
}

func getItems() []itemName {
	itemArray := []itemName {
		// "VorpalBlade",
		// "BattleAxe",
		// "WeapSwordShort",
		// "WeapFlailDire",
		// "WeapDart",
		// "WeapShortBow",
		// "WeapCrossbow",
		// "AmmoBolt",
		// "AmmoArrow",
		// "InfiniteJest",
		// "BookRage",
		// "RingStrength",
		// "PotHeal",
		// "PotConf",
		// "PotMega",
		// "ArmorChestLeather",
		// "ArmorBootsLeather",
		// "ArmorBootsIron",
		// "ScrollBlink",
		// "ShieldSmall",
		// "ShieldLarge",
		// "ShieldBuckler",
		VorpalBlade,
		BattleAxe,
		WeapSwordShort,
		WeapFlailDire,
		WeapDart,
		WeapShortBow,
		WeapCrossbow,
		AmmoBolt,
		AmmoArrow,
		InfiniteJest,
		BookRage,
		RingStrength,
		PotHeal,
		PotConf,
		PotMega,
		ArmorChestLeather,
		ArmorBootsLeather,
		ArmorBootsIron,
		ScrollBlink,
		ShieldSmall,
		ShieldLarge,
		ShieldBuckler,
	}
	return itemArray
}

func getRandomItemAll() itemName {
	itemArray := getItemArray(AllItems)
	rand.Seed(time.Now().UnixNano())
	theItem := itemArray[rand.Intn(len(itemArray))]
	return theItem
}

func getRandomItemCategory(itemsCat []itemName) itemName {
	//itemArray := getItemArray(AllItems)
	//fmt.Println(itemsCat)
	// TODO this if is probably not needed, but rings are broken -- troubleshooting
	if len(itemsCat) > 1 {
		rand.Seed(time.Now().UnixNano())
		theItem := itemsCat[rand.Intn(len(itemsCat))]
		return theItem
	} else {
		return itemsCat[0]
	}

}

// func 2getItemArray(category itemCategory) []itemName {
// 	// TODO this should get passed item Category Type 
// 	itemArray := getItems()
// 	if category == AllItems {
// 		//fmt.Println(itemArray)
// 		return itemArray
// 	}
// 	cat := string(category)[:3]
// 	var subArray []itemName 
// 	for _, item := range itemArray {
// 		if strings.HasPrefix(string(item), cat) {
// 		fmt.Println("AAAAA")
// 		//if item[:3] == string(category)[:3] {
// 			subArray = append(subArray, item)
// 		}
// 	}
// 	//fmt.Println(subArray)
// 	return subArray
	
// }

func getItemNamesInCategory(category itemCategory, allItemNames []itemName) []itemName {
	var items []Item
	var itemNames []itemName
	for _, item := range allItemNames {
		var gtmp Game
		gtmp.makeItem(item, &items, -1, -1)
	}
	for _, item := range items {
		if item.category == category {
			itemNames = append(itemNames, item.iname)
		}
		
	}
	return itemNames
}
// TODO combine the func above and below?
func getItemArray(category itemCategory) []itemName {
	itemArray := getItems()
	if category == AllItems {
		return itemArray
	}
	items := getItemNamesInCategory(category, itemArray)
	//fmt.Println(items)
	return items
	
}

func getRandomItem(category itemCategory) itemName {
	var item itemName
	if category == AllItems {
		item = getRandomItemAll()
		return item
	} else {
		itemsCat := getItemArray(category)
		item = getRandomItemCategory(itemsCat)
	}
	return item
}

// func printWeapons() {
// 	itemArray := getItems()
// 	for _, item := range itemArray {
// 		if strings.HasPrefix(item, "Weap") {
// 			fmt.Println(item)
// 		}
// 	}
// }