/*
 *  floors.go
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
    "fmt"

    "github.com/gdamore/tcell/v2"
)

const (
    RoomMaxSize = 12
    RoomMinSize = 6
    MaxRooms = 30
    MaxMobPerRoom = 3  // TODO this should not be const, should change w/ floor/room?
)

type Rectangle struct {
    x1 int
    y1 int
    x2 int
    y2 int
}

func (r *Rectangle) center() (int, int) {
    center_x := (r.x1 + r.x2) / 2;
    center_y := (r.y1 + r.y2) / 2;
    return center_x, center_y
}

func (r *Rectangle) intersects_with(other []Rectangle) bool {
    // returns true if rectangle intersects with existing one
    for i := range other {
        if (r.x1 <= other[i].x2) && 
            (r.x2 >= other[i].x1) && 
            (r.y1 <= other[i].y2) && 
            (r.y2 >= other[i].y1) {
                return true
            }
    }
    return false
}

func create_room(room Rectangle, floor *[][]bool) {
    // make rect tiles empty
    for x := room.x1; x < room.x2; x++ {
        for y := room.y1; y < room.y2; y++ {
            (*floor)[x][y] = false
        }
    }
}

//TODO pass only *tiles instead of game
func (g *Game) genFloor() {
    var f = Floor {
        name: "test floor " + fmt.Sprint(g.cur),
        fovDistance: FOVDistance,
    }
    g.floors = append(g.floors, f)

    // TODO make this expand to existing object vec without returning
    floor := make([][]bool, MaxWidth)
    for i := range floor {
        floor[i] = make([]bool, MaxHeight)
    }
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            floor[w][h] = true
        }
    }

    var rooms []Rectangle
    var upstair Point
    var downstair Point

    for i := 0; i < MaxRooms; i++ {
        // random width and height
        w := getRandNum(RoomMaxSize+10, RoomMinSize)
        h := getRandNum(RoomMaxSize, RoomMinSize)
        x := getRandNum(MaxWidth - w - 1, 1)
        y := getRandNum(MaxHeight - h - 1, 1)
        new_room := Rectangle { x, y, w+x, h+y }
        // make sure room doesnt intersect with others
        // no intersections
        if new_room.intersects_with(rooms) == false {
            new_x, new_y := new_room.center()
            create_room(new_room, &floor)
            if len(rooms) == 0 {
                // first room. hero & upstairs goes here
                upstair = Point { new_x, new_y }
                g.hero.setPos(new_x, new_y)
                g.placeLevelGear(new_x, new_y)
            } else {
                g.placeMobs(new_room)
                g.placeItems(new_room)
                //downstairs in last room generated
                downstair = Point { new_x, new_y }
                //generate mobs - placed here in loop so none gen in starting room
                prev_x, prev_y := rooms[len(rooms) - 1].center()
                // random tunnel directions
                rand.Seed(time.Now().UnixNano())
                if rand.Intn(2) == 0 {
                    make_hor_tunnel(prev_x, new_x, prev_y, floor);
                    make_vert_tunnel(prev_y, new_y, new_x, floor);
                } else {
                    make_vert_tunnel(prev_y, new_y, prev_x, floor);
                    make_hor_tunnel(prev_x, new_x, new_y, floor);
                }
            }
            // append the new room to the list
            rooms = append(rooms, new_room)
        } else {
            continue 
        }
    }  
    //   ▓  ▒  ░  █  ∏  ∆  ∑   ≈  ◊  µ  π  ¿  █
    // make wall object for each 'true' tile, floor for false
    for h := 0; h < MaxHeight; h++ {
    //for h in 0..MAP_HEIGHT {
        for w := 0; w < MaxWidth; w++ {
            if floor[w][h] == true {
                var t = Tile {
                        name: "wall",
                        pos: Point { w, h },
                        glyph: '#',
                        fg: tcell.ColorWhite,
                        bg: tcell.ColorGrey,
                        blocks: true,
                        blockSight: true,
                        visible: false,
                        visited: false,
                    }
                
                g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, t)
            } else {
                var t = Tile {
                        name: "floor",
                        pos: Point { w, h },
                        glyph: '.',
                        fg: tcell.ColorWhite,
                        bg: tcell.ColorDefault,
                        blocks: false,
                        blockSight: false,
                        visible: false,
                        visited: false,
                }
                g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, t)
            }
        }
    }

    var up = Tile {
        name: "upstair",
        pos: upstair,
        glyph: '<',
        fg: tcell.ColorWhite,
        bg: tcell.ColorDefault,
        blocks: false,
        blockSight: false,
        visible: false,
        visited: false,
    }
    g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, up)
    var down = Tile {
        name: "downstair",
        pos: downstair,
        glyph: '>',
        fg: tcell.ColorWhite,
        bg: tcell.ColorDefault,
        blocks: false,
        blockSight: false,
        visible: false,
        visited: false,
    }
    g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, down)
    
    //g.makeItem(PotHeal, &g.floors[g.cur].items, g.hero.pos.x, g.hero.pos.y)
}

func make_hor_tunnel(x1, x2, y int, floor [][]bool) {
    // horizontal tunnel. min/max are in case x1 is bigger than x2
    if x1 <= x2 {
        for x := x1; x <= x2; x++ {
            floor[x][y] = false
        }  
    } else {
        for x := x2; x <= x1; x++ {
            floor[x][y] = false
        }
    }
}

func make_vert_tunnel(y1, y2, x int, floor [][]bool) {
    // vertical tunnel
    if y1 <= y2 {
        for y := y1; y <= y2; y++ {
            floor[x][y] = false
        }  
    } else {
        for y := y2; y <= y1; y++ {
            floor[x][y] = false
        }
    }
}

//TODO pass only *tiles instead of game
func (g *Game) makeTestMap() {

    floor := make([][]bool, MaxWidth)
    for i := range floor {
        floor[i] = make([]bool, MaxHeight)
    }
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            floor[w][h] = true
        }
    }
    var rooms []Rectangle
    
    for i := 0; i < 1; i++ {
        w := 70
        h := 23
        x := 1
        y := 1
        new_room := Rectangle { x, y, w+x, h+y }
        // make sure room doesnt intersect with others
        // no intersections
        if new_room.intersects_with(rooms) == false {
            new_x, new_y := new_room.center()
            create_room(new_room, &floor)
            if len(rooms) == 0 {
                g.hero.setPos(new_x, new_y)
            } else {
                g.placeMobs(new_room)
                prev_x, prev_y := rooms[len(rooms) - 1].center()
                // random tunnel directions
                rand.Seed(time.Now().UnixNano())
                if rand.Intn(2) == 0 {
                    make_hor_tunnel(prev_x, new_x, prev_y, floor);
                    make_vert_tunnel(prev_y, new_y, new_x, floor);
                } else {
                    make_vert_tunnel(prev_y, new_y, prev_x, floor);
                    make_hor_tunnel(prev_x, new_x, new_y, floor);
                }
            }
            // append the new room to the list
            rooms = append(rooms, new_room)
        } else {
            continue 
        }
    }  
    //   ▓  ▒  ░  █  ∏  ∆  ∑   ≈  ◊  µ  π  ¿  █
    // make wall object for each 'true' tile, floor for false
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            if floor[w][h] == true {
                var t = Tile {
                        name: "wall",
                        pos: Point { w, h },
                        glyph: '#',
                        fg: tcell.ColorWhite,
                        bg: tcell.ColorGrey,
                        blocks: true,
                        blockSight: true,
                        visible: false,
                        visited: false,
                    }
                
                g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, t)
            } else {
                var t = Tile {
                        name: "floor",
                        pos: Point { w, h },
                        glyph: '.',
                        fg: tcell.ColorWhite,
                        bg: tcell.ColorDefault,
                        blocks: false,
                        blockSight: false,
                        visible: false,
                        visited: false,
                }
                g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, t)
            }
        }
    }
}


func (g *Game) placeItems(room Rectangle) {
    // choose random number of items
    numItems := getRandNum(5, 0)
    if numItems <= 0 { return }
    for i := 0; i < numItems; i++ {
        // choose random spot for item
        // TODO check for previous mobs at this loc first
        x := getRandNum(room.x2-1, room.x1)
        y := getRandNum(room.y2-1, room.y1)
        //itemType := getRandomItem()
        itemType := getRandomItem(AllItems)
        g.makeItem(itemType, &g.floors[g.cur].items, x, y)
    }
}


// 2d array 
// floor { mob, percetage }
// mob pool based on this array

func (g *Game) placeMobsTest(room Rectangle) {
    // choose random number of monsters
    numMobs := getRandNum(MaxMobPerRoom, 0)
    if numMobs <= 0 { return }
    for i := 0; i < numMobs; i++ {
        // choose random spot for this mob
        x := getRandNum(room.x2-1, room.x1)
        y := getRandNum(room.y2-1, room.y1)
        mobRoll := roll(R1d100)
        if mobRoll <= 50 {  // 80% chance of getting an orc
            makeMob(Rat, &g.floors[g.cur].actors, x, y)
        } else if mobRoll > 50 && mobRoll <= 90 {
            mobID := makeMob(Orc, &g.floors[g.cur].actors, x, y)
            g.floors[g.cur].actors[mobID].inv = append(g.floors[g.cur].actors[mobID].inv, Inventory { 0, makeWeapSwordShort(-1, -1) } )
            g.floors[g.cur].actors[mobID].weapon = 0
            g.floors[g.cur].actors[mobID].inv = append(g.floors[g.cur].actors[mobID].inv, Inventory { 1, makeArmorBootsLeather(-1, -1) } )
            g.floors[g.cur].actors[mobID].inv[1].item.equipped = true
        } else {
            makeMob(Troll, &g.floors[g.cur].actors, x, y)
        }
    }
}


func (g *Game) placeMobs(room Rectangle) {
    // choose random number of monsters
    numMobs := getRandNum(MaxMobPerRoom, 0)
    if numMobs <= 0 { return }
    for i := 0; i < numMobs; i++ {
        // choose random spot for this mob
        // TODO check for previous mobs at this loc first
        x := getRandNum(room.x2-1, room.x1)
        y := getRandNum(room.y2-1, room.y1)
        mobtype := getMob(roll(R1d100), g.cur)
        makeMob(mobtype, &g.floors[g.cur].actors, x, y)
    }
}

func (g *Game) placeLevelGear(x, y int) {

    g.makeItem(PotMega, &g.floors[g.cur].items, x + 1, y + 1)
    // g.makeItem(RingStrength, &g.floors[g.cur].items, g.hero.pos.x + 1, g.hero.pos.y + 1)
    // g.makeItem(RingStrength, &g.floors[g.cur].items, g.hero.pos.x + 0, g.hero.pos.y + 1)
    // g.makeItem(RingStrength, &g.floors[g.cur].items, g.hero.pos.x - 1, g.hero.pos.y + 1)
    // g.makeItem(ScrollBlink, &g.floors[g.cur].items, g.hero.pos.x -1 , g.hero.pos.y - 1)
    // g.makeItem(WeapFlailDire, &g.floors[g.cur].items, g.hero.pos.x + 0 , g.hero.pos.y - 1)

}

func (g *Game) genFloorGroundLevel() {
    var f = Floor {
        name: "Ground Floor" + fmt.Sprint(g.cur),
        fovDistance: -1,
    }
    g.floors = append(g.floors, f)
    // TODO make this expand to existing object vec without returning
    floor := make([][]bool, MaxWidth)
    for i := range floor {
        floor[i] = make([]bool, MaxHeight)
    }
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            floor[w][h] = true
        }
    }
    var rooms []Rectangle
    var upstair Point
    var downstair Point
        new_room := Rectangle { 1, 1, MaxWidth-1, MaxHeight-1 }
        // make sure room doesnt intersect with others
        // no intersections
        if new_room.intersects_with(rooms) == false {
            new_x, new_y := new_room.center()
            create_room(new_room, &floor)
            if len(rooms) == 0 {
                // first room. hero & upstairs goes here
                upstair = Point { new_x, new_y }
                g.hero.setPos(new_x, new_y)
                downstair = Point { new_x, 5 }
            }

            rooms = append(rooms, new_room)
        }
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            if floor[w][h] == true {
                var t = Tile {
                        name: "wall",
                        pos: Point { w, h },
                        glyph: '#',
                        fg: tcell.ColorWhite,
                        bg: tcell.ColorGrey,
                        blocks: true,
                        blockSight: true,
                        visible: false,
                        visited: false,
                    }
                if h == MaxHeight-1 && w >= MaxWidth/2-5 && w <= MaxWidth/2+5 {
                    t.bg = tcell.ColorBrown
                    t.name = "gate"
                }
                g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, t)
            } else {
                var t = Tile {
                        name: "floor",
                        pos: Point { w, h },
                        glyph: '"',
                        fg: tcell.ColorGreen,
                        bg: tcell.ColorDefault,
                        blocks: false,
                        blockSight: false,
                        visible: false,
                        visited: false,
                }
                if (h >= 8 && h <= 9 && w >= 9 && w <= 70) ||
                (h >= 0 && h <= 9 && w >= 9 && w <= 11) ||
                (h >= 0 && h <= 9 && w >= 68 && w <= 70) {
                    t.glyph = '█'
                    t.bg = tcell.ColorBlue
                    t.fg = tcell.ColorBlue
                    t.name = "moat"
                    t.blocks = true
                    t.blockSight = false
                    t.visible = true
                }
                if (h >= 7 && h <= 10 && w >= MaxWidth/2-1 && w <= MaxWidth/2+1) {
                    t.glyph = '▒'
                    t.bg = tcell.ColorBlack
                    t.fg = tcell.ColorBrown
                    t.name = "drawbridge"
                    t.blocks = false
                    t.blockSight = false
                    t.visible = true
                }
                g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, t)
            }
        }
    }
    var up = Tile {
        name: "upstair",
        pos: upstair,
        glyph: '<',
        fg: tcell.ColorWhite,
        bg: tcell.ColorDefault,
        blocks: false,
        blockSight: false,
        visible: false,
        visited: false,
    }
    g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, up)
    var down = Tile {
        name: "downstair",
        pos: downstair,
        glyph: '>',
        fg: tcell.ColorWhite,
        bg: tcell.ColorDefault,
        blocks: false,
        blockSight: false,
        visible: false,
        visited: false,
    }
    g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, down)
    
    var t = Tile {
        name: "torch",
        pos: Point { 2, 2 },
        glyph: '∆',
        fg: tcell.ColorOrange,
        bg: tcell.ColorDefault,
        blocks: true,
        blockSight: false,
        visible: false,
        visited: false,
    }
    g.floors[g.cur].tiles = append(g.floors[g.cur].tiles, t)

    makeMob(DragonRed, &g.floors[g.cur].actors, MaxWidth/2, 5)

    //g.makeItem(PotHeal, &g.floors[g.cur].items, g.hero.pos.x, g.hero.pos.y)
}