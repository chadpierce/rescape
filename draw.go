/*
 *  draw.go
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
    //"strconv"
    // "log"
    // "os"
    //"strconv"
	//"sync"
    "time"

    "github.com/gdamore/tcell/v2"
)

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, styleNum int) {
    var fg tcell.Color
    var bg tcell.Color
    var fillGlyph rune
    if styleNum == 0 {  //game board
        fg = tcell.ColorReset
        bg = tcell.ColorReset
        fillGlyph = ' '
    } else if styleNum == 1 {  // stats 
        fg = tcell.ColorYellow
        bg = tcell.ColorReset
        fillGlyph = ' '
    } else if styleNum == 2 {  // messages
        fg = tcell.ColorDefault
        bg = tcell.ColorDefault
        fillGlyph = ' '
    } else if styleNum == 3 {  // message list
        fg = tcell.ColorBlue
        bg = tcell.ColorGrey
        fillGlyph = ' '
    } else if styleNum == 4 {  // examine
        fg = tcell.ColorChartreuse
        bg = tcell.ColorGrey
        fillGlyph = '+'
    } else {
        fg = tcell.ColorRed  // error
        bg = tcell.ColorWhite
        fillGlyph = 'X'
    }
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	
    if y2 < y1 {
        y1, y2 = y2, y1
    }
    if x2 < x1 {
        x1, x2 = x2, x1
    }
    // Fill background
    for row := y1; row <= y2; row++ {
        for col := x1; col <= x2; col++ {
            s.SetContent(col, row, fillGlyph, nil, style)
        }
    }
    // draw borders
    if styleNum == 0 {
        style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
        for col := x1; col <= x2; col++ {
            s.SetContent(col, y1, '+', nil, style)
            s.SetContent(col, y2, '+', nil, style)
        }
        for row := y1 + 1; row < y2; row++ {
            s.SetContent(x1, row, '+', nil, style)
            s.SetContent(x2, row, '+', nil, style)
        }
    } else if styleNum == 1 || styleNum == 3 {
        for col := x1; col <= x2; col++ {
            s.SetContent(col, y1, tcell.RuneHLine, nil, style)
            s.SetContent(col, y2, tcell.RuneHLine, nil, style)
        }
        for row := y1 + 1; row < y2; row++ {
            s.SetContent(x1, row, tcell.RuneVLine, nil, style)
            s.SetContent(x2, row, tcell.RuneVLine, nil, style)
        }
        // Only draw corners if necessary
        if y1 != y2 && x1 != x2 {
            s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
            s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
            s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
            s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
        }
    }
}

func drawActors(s tcell.Screen, actors []Actor, layer int) {    
    vfg := tcell.ColorDarkViolet
    lastKnownfg := tcell.ColorViolet
    vbg := tcell.ColorDefault
    // 0 dead, 1 living, 2 both
    switch layer {
    case 0: 
        for _, a := range actors {
            if a.visible && !a.alive {
                s.SetContent(a.pos.x, a.pos.y, a.glyph, nil, 
                    tcell.StyleDefault.Foreground(a.fg).Background(a.bg))
            } else if a.visited && !a.alive {
                s.SetContent(a.pos.x, a.pos.y, a.glyph, nil, 
                    tcell.StyleDefault.Foreground(vfg).Background(vbg))
            }
        }
    case 1:
        for _, a := range actors {
            if a.visible && a.alive {
                s.SetContent(a.pos.x, a.pos.y, a.glyph, nil, 
                    tcell.StyleDefault.Foreground(a.fg).Background(a.bg))
                    //a.lastKnownPos = a.pos
            } else if a.visited && a.alive {
                s.SetContent(a.lastKnownPos.x, a.lastKnownPos.y, a.glyph, nil, 
                    tcell.StyleDefault.Foreground(lastKnownfg).Background(vbg))
            }
        }
    case 2:
        for _, a := range actors {
            if a.visible {
                s.SetContent(a.pos.x, a.pos.y, a.glyph, nil, 
                    tcell.StyleDefault.Foreground(a.fg).Background(a.bg))
            } else if a.visited {
                s.SetContent(a.pos.x, a.pos.y, a.glyph, nil, 
                    tcell.StyleDefault.Foreground(vfg).Background(vbg))
            }
        }
    }
}

func drawActor(s tcell.Screen, a Actor) {    
        s.SetContent(a.pos.x, a.pos.y, a.glyph, nil, 
			tcell.StyleDefault.Foreground(a.fg).Background(a.bg))
}

func drawTerrain(s tcell.Screen, tiles []Tile) {
    //vfg := tcell.ColorDarkViolet
    vfg := tcell.ColorDarkBlue
    vbg := tcell.ColorDefault
    for _, t := range tiles {
        if t.visible {
            s.SetContent(t.pos.x, t.pos.y, t.glyph, nil,
                tcell.StyleDefault.Foreground(t.fg).Background(t.bg))
        } else if t.visited {
            s.SetContent(t.pos.x, t.pos.y, t.glyph, nil,
                tcell.StyleDefault.Foreground(vfg).Background(vbg))
        }
    }
}

func drawItems(s tcell.Screen, items []Item) {    
    vfg := tcell.ColorDarkViolet
    vbg := tcell.ColorDefault
    var fg tcell.Color
    var bg tcell.Color
    var occupiedPos []Point 
    for _, i := range items {
        isInvert := false
        if i.visible {
            for _, o := range occupiedPos {
                if o.x == i.pos.x && o.y == i.pos.y {
                    isInvert = true
                }
            }
            occupiedPos = append(occupiedPos, i.pos)
            if isInvert {
                    fg = tcell.ColorBlack
                    bg = i.fg
                } else {
                    fg = i.fg
                    bg = i.bg
                }
            s.SetContent(i.pos.x, i.pos.y, i.glyph, nil,
                tcell.StyleDefault.Foreground(fg).Background(bg))
        } else if i.visited {
            s.SetContent(i.pos.x, i.pos.y, i.glyph, nil,
                tcell.StyleDefault.Foreground(vfg).Background(vbg))
        }
    }
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, style)
        x++
    }
}

func drawFovActorStats(s tcell.Screen, mobs []string) {
    // get distance in fov func and order by
    StatsStyle := tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault)
    y := 15
    for _, mob := range mobs {
        x := MaxWidth + 1; y++
        for _, r := range []rune(mob) {
            s.SetContent(x, y, r, nil, StatsStyle)
            x++
        }
    }

}

func drawStats(s tcell.Screen, hero Actor, floor int) {
	StatsStyle := tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault)
	quiveredItem := "empty"
    wieldedItem := "empty handed"
    shieldItem := "empty"
    quiverInvID := getInvPositionForSlotID(hero.quiver, hero.inv)
    if hero.quiver != -1 { 
        if hero.inv[quiverInvID].item.stackable {
            quiveredItem = hero.inv[quiverInvID].item.name + " (" + 
            fmt.Sprint(hero.inv[quiverInvID].item.charges) + ")" 
        } else {
            quiveredItem = hero.inv[quiverInvID].item.name
        }
    }
    if hero.weapon != -1 { wieldedItem = hero.inv[getInvPositionForSlotID(hero.weapon, hero.inv)].item.name }
    for _, i := range hero.inv {
        if (i.item.category == Shield && i.item.equipped) || 
        (i.item.category == Buckler && i.item.equipped) { shieldItem = i.item.name; break}
    }
    x := MaxWidth + 1; y := 0
	text := fmt.Sprintf("hero\nx: %v\ny: %v", hero.pos.x, hero.pos.y)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("floor: %v", floor)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y += 2
    text = fmt.Sprintf("HP: %v/%v", hero.hp, hero.maxHP)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y += 1
    text = fmt.Sprintf("Mana: %v/%v", hero.mana, hero.maxMana)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("Str: %v", hero.strg)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("Dex: %v", hero.dex)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("Int: %v", hero.intel)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("AC: %v", hero.ac)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("Shield: %v", shieldItem)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("Weapon: %v", wieldedItem)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
    x = MaxWidth + 1; y++
    text = fmt.Sprintf("Quiver: %v", quiveredItem)
	for _, r := range []rune(text) {
        s.SetContent(x, y, r, nil, StatsStyle)
        x++
    }
}

func drawMessages(s tcell.Screen, m []Message) {
    x := 1
    y := MaxHeight
    for i := 0; i <= 4; i++ {
        text := m[i].text
        style := tcell.StyleDefault.Foreground(m[i].fg).Background(tcell.ColorDefault)
        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
        x = 1
        y++
    }
}

func drawAllMessages(s tcell.Screen, m []Message) {
    //TODO space to get older messages
    drawBox(s, 1, 1, MaxWidth + 5, MaxHeight +5, 3)  // message list
    x := 3
    y := 2
    var start, end int
    fmt.Printf("lenm: %v", len(m))
    if len(m) >= 29 {
        start = len(m) - 28
        end = len(m) -1
    } else {
        start = 4
        end = len(m) - 1
    }

    for i := start; i <= end ; i++ {
        text := m[i].text
        style := tcell.StyleDefault.Foreground(m[i].fg).Background(tcell.ColorDefault)
        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
        x = 3
        y++
    }
    s.Show()
    ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
			return
        }
    }
}

func drawInv(titleStr string, s tcell.Screen, inv []Inventory) {
    decDiff := 97
    drawBox(s, 10, 3, MaxWidth + 10, MaxHeight + 2, 1)
    x := 15; y := 5
    style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
    for _, r := range []rune(titleStr) {
        s.SetContent(x, y, r, nil, style)
        x++
    }
    x = 14; y = 7
    for cnt, i := range inv { 
        equipStr := ""
        selAscii := rune(i.slot + decDiff)
        if i.item.stackable && i.item.charges > 1 {
            equipStr = fmt.Sprintf("(%v)", i.item.charges)
        }
        if i.item.equipped {
            equipStr = equipStr + " [equipped]"
        } else if i.item.quivered {
            equipStr = equipStr + " [quivered]"
        } 
        text := fmt.Sprintf("%v) %v %v", string(selAscii), i.item.name, equipStr)
        style := tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault)
        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
        if cnt < 15 {
            x = 14
        } else if cnt == 15 {
            x = 56
            y = 6
        } else {
            x = 56
        }
        y++
    }
    s.Show()
}

func drawItemList(titleStr string, s tcell.Screen, items []Item) {
    selDecimal := 97
    equipStr := ""
    drawBox(s, 10, 3, MaxWidth + 10, MaxHeight + 2, 1)
    x := 15
    y := 5
    style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
    for _, r := range []rune(titleStr) {
        s.SetContent(x, y, r, nil, style)
        x++
    }
    x = 14
    y = 7
    for i := 0; i <= len(items) - 1 ; i++ {
        selAscii := rune(i + selDecimal)
        if items[i].equipped {
            equipStr = "[equipped]"
        } else if items[i].stackable && items[i].charges > 1 {
            equipStr = fmt.Sprintf("(%v)", items[i].charges)
        } else {
            equipStr = ""
        }
        text := fmt.Sprintf("%v) %v %v", string(selAscii), items[i].name, equipStr)
        style := tcell.StyleDefault.Foreground(items[i].fg).Background(tcell.ColorDefault)

        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
        if i < 15 {
            x = 14
        } else if i == 15 {
            x = 56
            y = 6
        } else {
            x = 56
        }
        y++
    }
    s.Show()
}


func drawKnownSpells(s tcell.Screen, list []SpellID, titleStr string) {
    selDecimal := 97
    drawBox(s, 10, 3, MaxWidth + 10, MaxHeight + 2, 1)
    x := 15; y := 5
    style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
    for _, r := range []rune(titleStr) {
        s.SetContent(x, y, r, nil, style)
        x++
    }
    x = 14; y = 7
    for i := 0; i <= len(list) - 1 ; i++ {
        selAscii := rune(i + selDecimal)
        text := fmt.Sprintf("%v) %v", string(selAscii), getSpellName(list[i]))
        style := tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault)
        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
        x = 14; y++
    }
    s.Show()
}

func (g *Game) drawSpellBook(s tcell.Screen, spells []string, titleStr, descStr string) {
    selDecimal := 97
    equipStr := ""
    drawBox(s, 10, 3, MaxWidth + 10, MaxHeight + 2, 1)
    x := 15; y := 5
    style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorPurple)
    for _, r := range []rune(titleStr) {
        s.SetContent(x, y, r, nil, style)
        x++
    }
    x = 15; y = 7
    style = tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault)
    for _, r := range []rune(descStr) {
        s.SetContent(x, y, r, nil, style)
        x++
    }
    x = 14; y = 12
    for i := 0; i <= len(spells) - 1 ; i++ {
        selAscii := rune(i + selDecimal)
        text := fmt.Sprintf("%v) %v %v", string(selAscii), spells[i], equipStr)
        style := tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault)

        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
        x = 14
        y++
    }
    s.Show()
}

func drawDisc(titleStr string, s tcell.Screen, list []DiscoveredItem) {
    selDecimal := 97
    drawBox(s, 10, 3, MaxWidth + 10, MaxHeight + 2, 1)
    x := 15
    y := 5
    style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
    for _, r := range []rune(titleStr) {
        s.SetContent(x, y, r, nil, style)
        x++
    }
    x = 14
    y = 7
    for i := 1; i <= len(list) - 1 ; i++ {
        selAscii := rune(i + selDecimal-1)
        text := fmt.Sprintf("%v) %v", string(selAscii), list[i].name)
        style := tcell.StyleDefault.Foreground(tcell.ColorDefault)

        for _, r := range []rune(text) {
            s.SetContent(x, y, r, nil, style)
            x++
        }
        x = 14
        y++
    }
    s.Show()
}

func drawTarget(t Point, s tcell.Screen, g Game, items string, targetType int) {
    var style tcell.Style; var style2 tcell.Style
    if targetType == 0 {  // examine
        style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorChartreuse)
        style2 = tcell.StyleDefault.Foreground(tcell.ColorChartreuse).Background(tcell.ColorDefault)

    } else if targetType == 1 {  // magic
        style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorPurple)
        style2 = tcell.StyleDefault.Foreground(tcell.ColorPurple).Background(tcell.ColorDefault)
    } else if targetType == 2 {  // projectiles
        style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorOrange)
        style2 = tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorDefault)

    } else { g.dbg("ERROR: BAD TARGET STYLE") }
    s.Clear()
    glyph, _, _, _ := s.GetContent(t.x, t.y)
    x1 := 0; y1 := 0
    x2 := MaxWidth - 1; y2 := MaxHeight - 1
    drawTerrain(s, g.floors[g.cur].tiles)
    for col := x1; col <= x2; col++ {
        s.SetContent(col, y1, tcell.RuneHLine, nil, style2)
        s.SetContent(col, y2, tcell.RuneHLine, nil, style2)
    }
    for row := y1 + 1; row < y2; row++ {
        s.SetContent(x1, row, tcell.RuneVLine, nil, style2)
        s.SetContent(x2, row, tcell.RuneVLine, nil, style2)
    }
    // Only draw corners if necessary
    if y1 != y2 && x1 != x2 {
        s.SetContent(x1, y1, tcell.RuneULCorner, nil, style2)
        s.SetContent(x2, y1, tcell.RuneURCorner, nil, style2)
        s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style2)
        s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style2)
    }
    drawStats(s, g.hero, g.cur)
    drawMessages(s, g.msg[len(g.msg)-5:len(g.msg)])
	drawTerrain(s, g.floors[g.cur].tiles)
    drawActors(s, g.floors[g.cur].actors, 0)
	drawItems(s, g.floors[g.cur].items)
    drawActors(s, g.floors[g.cur].actors, 1)
    drawActor(s, g.hero)
    s.SetContent(t.x, t.y, glyph, nil, style)
    x := 10; y := MaxHeight - 1
    for _, r := range []rune(items) {
        if x == 90 {
            x = 10
            y++
        }
        s.SetContent(x, y, r, nil, style)
        x++
    }
    s.Show()
}

func drawZapLine(s tcell.Screen, sx, sy, tx, ty int) {
    zDuration := 10
    style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
    line := getLine(sx, sy, tx, ty)
    if len(line) >= 5 && len(line) <= 10 {
        zDuration = 8
    } else if len(line) > 10 && len(line) <= 20  {
        zDuration = 5
    } else if len(line) > 20 {
        zDuration = 2
    }
    for _, l := range line {
        s.SetContent(l.x, l.y, ' ', nil, style)
        s.Show()
        time.Sleep(time.Duration(zDuration) * time.Millisecond) 
    }
}

func draw(s tcell.Screen, g Game) {
	var actorsInView []string
    if g.floors[g.cur].fovDistance != -1 {
        actorsInView = g.floors[g.cur].getFOV(g.hero)
    } else {
        actorsInView = g.floors[g.cur].getAllView(g.hero)
    }
    
	drawStats(s, g.hero, g.cur)
    drawFovActorStats(s, actorsInView)
    drawMessages(s, g.msg[len(g.msg)-5:len(g.msg)])
	drawTerrain(s, g.floors[g.cur].tiles)
    drawActors(s, g.floors[g.cur].actors, 0)
	drawItems(s, g.floors[g.cur].items)
    drawActors(s, g.floors[g.cur].actors, 1)
    drawActor(s, g.hero)
}

// func dmsg (s tcell.Screen, m string) {
// 	x := 60
//     y := 26
//     style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
//         for _, r := range []rune(m) {
//             s.SetContent(x, y, r, nil, style)
//             x++
//         }
//         s.Show()
// }

// func dmsg2 (s tcell.Screen, m string) {
// 	x := 60
//     y := 27
//     style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorLightBlue)
//         for _, r := range []rune(m) {
//             s.SetContent(x, y, r, nil, style)
//             x++
//         }
//         s.Show()
// }

// func dmsg3 (s tcell.Screen, m string) {
// 	x := 60
//     y := 28
//     style := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorLightBlue)
//         for _, r := range []rune(m) {
//             s.SetContent(x, y, r, nil, style)
//             x++
//         }
//         s.Show()
// }