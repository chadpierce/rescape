/*
 *  path.go
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


// NOTES
// basic pathfinding is a bit wonky and needs work
// mobs will align diagonally when chasing to the east

// the following mit course was used to develop the bfs algorithm:
// https://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-006-introduction-to-algorithms-fall-2011/lecture-videos/lecture-13-breadth-first-search-bfs/

// this was ported from rust and needs some additional changes to make more go-like

type PathPoint struct {
    i, x, y int
}

func isAdjacent(sx, sy, tx, ty int) bool {

    if tx == sx - 1 && ty == sy - 1 { return true 
    } else if tx == sx && ty == sy - 1 { return true 
    } else if tx == sx + 1 && ty == sy - 1 { return true
    } else if tx == sx + 1 && ty == sy { return true
    } else if tx == sx + 1 && ty ==  sy + 1 { return true
    } else if tx == sx && ty == sy + 1 { return true
    } else if tx == sx - 1 && ty == sy + 1 { return true
    } else if tx == sx - 1 && ty == sy { return true 
    } else { return false }
}

func getAdjacentBlockers(x, y int, actors []Actor, items []Item, grid *[][]bool) {

    var moves []Point
    moves = append(moves, 
        Point { x-1, y-1 }, 
        Point { x, y-1 },
        Point { x+1, y-1 },
        Point { x+1, y },
        Point { x+1, y+1 },
        Point { x, y+1 },
        Point { x-1, y+1 },
        Point { x-1, y } )
    for _, m := range moves {
        for _, a := range actors {
            if a.pos.x == m.x && a.pos.y == m.y && a.blocks == true {
                (*grid)[m.x][m.y] = true
            } 
        }
        for _, i := range items {
            if i.pos.x == m.x && i.pos.y == m.y && i.blocks == true {
                (*grid)[m.x][m.y] = true
            } 
        }
    }
}


func getNeighbors(source_x, source_y int, grid [][]bool, level [][]int) []Point {
    var moves []Point
    if grid[source_x - 1][source_y - 1] == false && level[source_x - 1][source_y - 1] == -1 {
        moves = append(moves, Point{source_x - 1, source_y - 1})
    }
    if grid[source_x][source_y - 1] == false && level[source_x][source_y - 1] == -1 {
        moves = append(moves, Point{source_x, source_y - 1})
    }
    if grid[source_x + 1][source_y - 1] == false && level[source_x + 1][source_y - 1] == -1 {
        moves = append(moves, Point{source_x + 1, source_y - 1})
    }
    if grid[source_x + 1][source_y] == false && level[source_x + 1][source_y] == -1 {
        moves = append(moves, Point{source_x + 1, source_y})
    }
    if grid[source_x + 1][source_y + 1] == false && level[source_x + 1][source_y + 1] == -1 {
        moves = append(moves, Point{source_x + 1, source_y + 1})
    }
    if grid[source_x][source_y + 1] == false && level[source_x][source_y + 1] == -1 {
        moves = append(moves, Point{source_x, source_y + 1})
    }
    if grid[source_x - 1][source_y + 1] == false && level[source_x - 1][source_y + 1] == -1 {
        moves = append(moves, Point{source_x - 1, source_y + 1})
    }	
    if grid[source_x - 1][source_y] == false && level[source_x - 1][source_y] == -1 {
        moves = append(moves, Point{source_x - 1, source_y}) 
    }
    return moves
}

func getPNeighbors(source_x, source_y int, grid [][]bool) []Point {
    var moves []Point
    if grid[source_x - 1][source_y - 1] == false {
        moves = append(moves, Point{source_x - 1, source_y - 1})
    }
    if grid[source_x][source_y - 1] == false {
        moves = append(moves, Point{source_x, source_y - 1})
    }
    if grid[source_x + 1][source_y - 1] == false {
        moves = append(moves, Point{source_x + 1, source_y - 1})
    }
    if grid[source_x + 1][source_y] == false {
        moves = append(moves, Point{source_x + 1, source_y})
    }
    if grid[source_x + 1][source_y + 1] == false {
        moves = append(moves, Point{source_x + 1, source_y + 1})
    }
    if grid[source_x][source_y + 1] == false {
        moves = append(moves, Point{source_x, source_y + 1})
    }
    if grid[source_x - 1][source_y + 1] == false {
        moves = append(moves, Point{source_x - 1, source_y + 1})
    }	
    if grid[source_x - 1][source_y] == false {
        moves = append(moves, Point{source_x - 1, source_y}) 
    }
    return moves
}

// TODO pass only cur floor
func getBFSGrid(g *Game) [][]bool {
    grid := make([][]bool, MaxWidth)
    for i, _ := range grid {
        grid[i] = make([]bool, MaxHeight)
    }
    for h:= 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            grid[w][h] = false
        }
    }
    
    for h := 0; h < MaxHeight; h++ {
       for w := 0; w < MaxWidth; w++ {
            for _, t := range g.floors[g.cur].tiles {
                if t.blocks == true &&
                  t.pos.x == w && t.pos.y == h {
                    grid[w][h] = true
                }
            } 
        }
    } 
    return grid
}

//this is not being used - slow and buggy - bfsTest2 is better 
func (g *Game) bfsTest(id, sx, sy, tx, ty int) {
    grid := getBFSGrid(g)
    getAdjacentBlockers(sx, sy, g.floors[g.cur].actors, g.floors[g.cur].items, &grid)
    levels := make([][]int, MaxWidth)
    for i, _ := range levels {
        levels[i] = make([]int, MaxHeight)
    }
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            levels[w][h] = -1
        }
    }
    lvl := 1
    levels[sx][sy] = lvl - 1
    frontier := getNeighbors(sx, sy, grid, levels)
    if len(frontier) <= 0 { return } // if no moves
    for _, f := range frontier {
        levels[f.x][f.y] = lvl
    }
    lvl++
    for len(frontier) > 0 {
        var next []Point
        for _, f := range frontier {
            neighbors := getNeighbors(f.x, f.y, grid, levels)
            for _, n := range neighbors {
                    levels[n.x][n.y] = lvl
                    next = append(next, Point { n.x, n.y })
                    
                    // for id, tile := range g.floors[g.cur].tiles { 
                    // 	//fmt.Printf("n: %v, %v", n.y, n.y)
                    // 	if tile.pos.x == n.x && tile.pos.y == n.y {
                    // 		//fmt.Printf("m:%v %v l: %v", n.x, n.y, level[n.x][n.y])
                    // 		switch levels[n.x][n.y] {
                    // 			case -1: g.floors[g.cur].tiles[id].glyph = 'w'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color20
                    // 			case 0: g.floors[g.cur].tiles[id].glyph = '0'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color21
                    // 			case 1: g.floors[g.cur].tiles[id].glyph = '1'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color22
                    // 			case 2: g.floors[g.cur].tiles[id].glyph = '2'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color23
                    // 			case 3: g.floors[g.cur].tiles[id].glyph = '3'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color24
                    // 			case 4: g.floors[g.cur].tiles[id].glyph = '4'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color25
                    // 			case 5: g.floors[g.cur].tiles[id].glyph = '5'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color26
                    // 			case 6: g.floors[g.cur].tiles[id].glyph = '6'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color27
                    // 			case 7: g.floors[g.cur].tiles[id].glyph = '7'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color28
                    // 			case 8: g.floors[g.cur].tiles[id].glyph = '8'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color29
                    // 			case 9: g.floors[g.cur].tiles[id].glyph = '9'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color30
                    // 			case 10: g.floors[g.cur].tiles[id].glyph = 'A'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color31
                    // 			case 11: g.floors[g.cur].tiles[id].glyph = 'B'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color32
                    // 			case 12: g.floors[g.cur].tiles[id].glyph = 'C'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color33
                    // 			case 13: g.floors[g.cur].tiles[id].glyph = 'D'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color34
                    // 			case 14: g.floors[g.cur].tiles[id].glyph = 'E'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color35
                    // 			case 15: g.floors[g.cur].tiles[id].glyph = 'F'
                    // 			default: 
                    // 				g.floors[g.cur].tiles[id].glyph = 'x'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color37
                    // 		}
                    // 	}
                    // }
                    //fmt.Printf("next: %v", next)
                    // fmt.Printf("lnext: %v", len(next))
                    // fmt.Printf("lh: %v", level[tx][ty])
                    
                //}
            }
           
        }
        frontier = next
        next = nil
        lvl++
    }
    //get neighbors of target
    px := tx
    py := ty
    cnt := 0
    for isAdjacent(px, py, sx, sy) == false { 
        var shortPath []Point
        var short int
        pNeighbors := getPNeighbors(px, py, grid)
        var neighbors []Point
        for _, pN := range pNeighbors {
            pNPos := Point { pN.x, pN.y }
            if g.floors[g.cur].actors[id].pos != pNPos {
                neighbors = append(neighbors, pN)
            }
        }
        for i, n := range neighbors {
            if i==0 || levels[n.x][n.y] < short {
                short = levels[n.x][n.y]
            }
        }
        for _, sh := range neighbors {
                if levels[sh.x][sh.y] == short { //&& !isB {
                    shortPath = append(shortPath, sh)
                }
        }
        if len(shortPath) < 1 {
            g.dbg("no moves?!?")    
            return
        } else {
            rand.Seed(time.Now().UnixNano())
            theMove := shortPath[rand.Intn(len(shortPath))]
            px = theMove.x
            py = theMove.y
        }
        cnt++
        if cnt > 1000 {
            g.dbg("break - hack to fix infinite loop!")
            break
        }
    }
    dx := px - sx
    dy := py - sy
    g.floors[g.cur].aiTakeTurnPathTest(id, dx, dy)
}
   
func moveLine (id int, sPos, tPos Point, g *Game) {
    line := getLine(sPos.x, sPos.y, tPos.x, tPos.y)
    dx := line[1].x - sPos.x
    dy := line[1].y - sPos.y
    if isAdjacent(sPos.x, sPos.y, tPos.x, tPos.y) {
        g.dbg("attack! moveLine")
    } else {
        isB, _ := isBlocked(sPos.x + dx, sPos.y + dy, g.floors[g.cur].actors, g.floors[g.cur].tiles)
        if !isB {
            g.floors[g.cur].aiTakeTurnPathTest(id, dx, dy)
        }
    }
}

func confusedMove(id int, sPos, tPos Point, g *Game) {
    
    if isAdjacent(sPos.x, sPos.y, tPos.x, tPos.y) {
        dmg := g.floors[g.cur].actors[id].meleeAttack(&g.hero)
        g.addMessage(fmt.Sprintf("%v attacks for %v damage!", g.floors[g.cur].actors[id].name, dmg), tcell.ColorRed)
    } else {
        grid := getBFSGrid(g)
        n := getPNeighbors(sPos.x, sPos.y, grid)
        n = append(n, Point { sPos.x, sPos.y })
        rand.Seed(time.Now().UnixNano())
        theMove := n[rand.Intn(len(n))]
        dx := theMove.x - sPos.x
        dy := theMove.y - sPos.y
        g.floors[g.cur].aiTakeTurnPathTest(id, dx, dy)
    }
}

func getRandomFloorTile(g [][]bool) Point {
    var emptyTiles []Point
    for h:= 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            if g[w][h] == false {
                emptyTiles = append(emptyTiles, Point { w, h }) 
            }
        }
    }
    rand.Seed(time.Now().UnixNano())
    randTile := emptyTiles[rand.Intn(len(emptyTiles))]
    return randTile
}

func wanderMove(id int, g *Game) {
    //mobs with no target choose a random target,
    //if they reach their target they choose a new one
    a := g.floors[g.cur].actors[id]
    curPos := a.pos 
    if a.target.x == -1 || a.pos == a.target {
        g.floors[g.cur].actors[id].target = getRandomFloorTile(getBFSGrid(g))
    }
    g.bfsTest2(id, a.pos.x, a.pos.y, g.floors[g.cur].actors[id].target.x, g.floors[g.cur].actors[id].target.y)
    //after moving, if they are in the same pos they choose a new target
    if g.floors[g.cur].actors[id].pos == curPos {
        g.floors[g.cur].actors[id].target = getRandomFloorTile(getBFSGrid(g))
    }
}

func dumbMove(id int, sPos, tPos Point, g *Game) {


    // TODO WORKING check for los and move straight at target if not blocked

    if isAdjacent(sPos.x, sPos.y, tPos.x, tPos.y) {
        dmg := g.floors[g.cur].actors[id].meleeAttack(&g.hero)
        g.addMessage(fmt.Sprintf("%v attacks for %v damage!", g.floors[g.cur].actors[id].name, dmg), tcell.ColorRed)

    } else {
        g.bfsTest2(id, sPos.x, sPos.y, tPos.x, tPos.y)
    }
}
 
func (g *Game) heroDumbMove() {
    sPos := Point { g.hero.pos.x, g.hero.pos.y }
    tPos := Point { g.hero.target.x, g.hero.target.y }
    isBlockedPath := false
    if isAdjacent(sPos.x, sPos.y, tPos.x, tPos.y) {
        g.dbg("hero hits it")
    } else {
        line := getLine(sPos.x, sPos.y, tPos.x, tPos.y)
        for _, l := range line {
            // if actor cant make straigt line, do bfs
            isB, _ := isBlocked(l.x, l.y, g.floors[g.cur].actors, g.floors[g.cur].tiles)
            if isB {
                isBlockedPath = true
            } 
            if isBlockedPath {
                g.heroBfsTest() 
            } else {
                g.dbg("move that dir")
                dx := line[1].x - sPos.x
                dy := line[1].y - sPos.y
                heroMove(g, dx, dy)
                sawSomething, sawWhat := g.floors[g.cur].isActorsInFOV(g.hero)
                if sawSomething {
                    g.addMessage(fmt.Sprintf("You see %v.", sawWhat), tcell.ColorDefault)
                    g.state = Playing
                } else if g.hero.target == g.hero.pos {
                    g.state = Playing
                }
            }
        }
    }
}

func (f *Floor) isExitCorridor(x, y, dir int) bool {
    var s1 bool
    var s2 bool
    var f1 bool
    var f2 bool
    switch dir {
    case 1: 
        for _, tile := range f.tiles {
            if tile.pos == (Point { x-1, y }) { s1 = tile.blocks }
            if tile.pos == (Point { x+1, y }) { s2 = tile.blocks }
            if tile.pos == (Point { x-1, y-1 }) { f1 = tile.blocks }
            if tile.pos == (Point { x+1, y-1 }) { f2 = tile.blocks }
        } 
    case 3:
        for _, tile := range f.tiles {
            if tile.pos == (Point { x, y-1 }) { s1 = tile.blocks }
            if tile.pos == (Point { x, y+1 }) { s2 = tile.blocks }
            if tile.pos == (Point { x+1, y-1 }) { f1 = tile.blocks }
            if tile.pos == (Point { x+1, y+1 }) { f2 = tile.blocks }
        }
    case 5:
        for _, tile := range f.tiles {
            if tile.pos == (Point { x-1, y }) { s1 = tile.blocks }
            if tile.pos == (Point { x+1, y }) { s2 = tile.blocks }
            if tile.pos == (Point { x-1, y+1 }) { f1 = tile.blocks }
            if tile.pos == (Point { x+1, y+1 }) { f2 = tile.blocks }
        }
    case 7:
        for _, tile := range f.tiles {
            if tile.pos == (Point { x-1, y-1 }) { s1 = tile.blocks }
            if tile.pos == (Point { x-1, y+1 }) { s2 = tile.blocks }
            if tile.pos == (Point { x-2, y-1 }) { f1 = tile.blocks }
            if tile.pos == (Point { x-2, y+1 }) { f2 = tile.blocks }
        }
    default: return false
    }
    if (s1 == true && f1 == false) || (s2 == true && f2 == false) {
        return true
    } else {
        return false
    }
}

func (f *Floor) setRunTarget(hero *Actor, dir int, g *Game) {
	var tx int
	var ty int
	switch dir {
	case 0: tx = -1; ty = -1
	case 1: tx = 0; ty = -1
	case 2: tx = +1; ty = -1
	case 3: tx = +1; ty = 0
	case 4: tx = 1; ty = 1
	case 5: tx = 0; ty = 1
	case 6: tx = -1; ty = 1
	case 7: tx = -1; ty = 0
	default: tx = hero.pos.x; ty = hero.pos.y // TODO handle this better
	}
	curx := hero.pos.x
	cury := hero.pos.y
	out:
	for {
		curx += tx
		cury += ty
		for i, tile := range f.tiles {
			if tile.pos == (Point{curx, cury}) {
				if tile.blocks == true || f.isExitCorridor(curx, cury, dir) {
					if dir == 1 || dir == 5 {
                        hero.target = Point { curx, f.tiles[i-1].pos.y }
                    } else if dir == 3 {
                        hero.target = Point { f.tiles[i].pos.x, cury }
                        // 3 and 7 are split up as a hack
                        // because running east was sticking
                    } else if dir == 7 {
                        hero.target = Point { f.tiles[i-1].pos.x, cury }
                    } else {
                        hero.target = Point { f.tiles[i-1].pos.x, f.tiles[i-1].pos.y }
                    }
                    break out
				} 
			}
		}
	}
}

func (g *Game) heroRun() {
    sawSomething, sawWhat := g.floors[g.cur].isActorsInFOV(g.hero)
    if sawSomething {
        g.addMessage(fmt.Sprintf("You see %v.", sawWhat), tcell.ColorDefault)
        g.state = Playing
    }
    if g.hero.pos != g.hero.target && !sawSomething {
        line := getLine(g.hero.pos.x, g.hero.pos.y, g.hero.target.x, g.hero.target.y)
        isB, _ := isBlocked(line[1].x, line[1].y, g.floors[g.cur].actors, g.floors[g.cur].tiles)
            if isB {
                g.state = Playing
            } else {
                heroMove(g, line[1].x - g.hero.pos.x, line[1].y - g.hero.pos.y)
                if g.hero.target == g.hero.pos {
                    g.state = Playing
                } 
            }
    } else {
        g.state = Playing
    }
}

func (g *Game) heroConfused() Point {                
    grid := getBFSGrid(g)
    n := getPNeighbors(g.hero.pos.x, g.hero.pos.y, grid)
    n = append(n, Point { g.hero.pos.x, g.hero.pos.y })
    rand.Seed(time.Now().UnixNano())
    theMove := n[rand.Intn(len(n))]
    dx := theMove.x - g.hero.pos.x
    dy := theMove.y - g.hero.pos.y
    heroMove(g, dx, dy)
    g.dbg(fmt.Sprintf("heromove: %v, %v ", dx, dy))
    return g.hero.pos
}

// TODO combine this with normal AI mob bfs move func
func (g *Game) heroBfsTest() {
    tx := g.hero.target.x
    ty := g.hero.target.y
    sx := g.hero.pos.x
    sy := g.hero.pos.y

    grid := getBFSGrid(g)
    getAdjacentBlockers(sx, sy, g.floors[g.cur].actors, g.floors[g.cur].items, &grid)
    levels := make([][]int, MaxWidth)
    for i, _ := range levels {
        levels[i] = make([]int, MaxHeight)
    }
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            levels[w][h] = -1
        }
    }
    lvl := 1
    levels[sx][sy] = lvl - 1
    frontier := getNeighbors(sx, sy, grid, levels)
    if len(frontier) <= 0 { return } // if no moves
    for _, f := range frontier {
        levels[f.x][f.y] = lvl
    }
    lvl++

    for len(frontier) > 0 {
        var next []Point
        for _, f := range frontier {
            neighbors := getNeighbors(f.x, f.y, grid, levels)
            for _, n := range neighbors {
                    levels[n.x][n.y] = lvl
                    next = append(next, Point { n.x, n.y })
            }
           
        }
        frontier = next
        next = nil
        lvl++
    }
    px := tx
    py := ty
    cnt := 0
    for isAdjacent(px, py, sx, sy) == false { 
        var shortPath []Point
        var short int
        pNeighbors := getPNeighbors(px, py, grid)
        var neighbors []Point
        
        for _, pN := range pNeighbors {
            pNPos := Point { pN.x, pN.y }
            if g.hero.pos != pNPos {
                neighbors = append(neighbors, pN)
            }
        }
        for i, n := range neighbors {
            if i==0 || levels[n.x][n.y] < short {
                short = levels[n.x][n.y]
            }
        }
        for _, sh := range neighbors {
                if levels[sh.x][sh.y] == short { //&& !isB {
                    shortPath = append(shortPath, sh)
                }
        }
        if len(shortPath) < 1 {
            g.dbg("no moves?!?")    
            return
        } else {
     
            rand.Seed(time.Now().UnixNano())
            theMove := shortPath[rand.Intn(len(shortPath))]
      
            px = theMove.x
            py = theMove.y
        }
        cnt++
        if cnt > 2000 {
            g.dbg("break - hack to fix infinite loop!")
            break
        }
    }
    dx := px - sx
    dy := py - sy
    heroMove(g, dx, dy)
    sawSomething, sawWhat := g.floors[g.cur].isActorsInFOV(g.hero)
    if sawSomething {
        g.addMessage(fmt.Sprintf("You see %v.", sawWhat), tcell.ColorDefault)
        g.state = Playing
    } else if g.hero.target == g.hero.pos {
        g.state = Playing
    }
}

// this one does bfs from destination to source and is faster
func (g *Game) bfsTest2(id, sx, sy, tx, ty int) {
    grid := getBFSGrid(g)
    getAdjacentBlockers(sx, sy, g.floors[g.cur].actors, g.floors[g.cur].items, &grid)
    levels := make([][]int, MaxWidth)
    for i, _ := range levels {
        levels[i] = make([]int, MaxHeight)
    }
    for h := 0; h < MaxHeight; h++ {
        for w := 0; w < MaxWidth; w++ {
            levels[w][h] = -1
        }
    }
    lvl := 1
    levels[tx][ty] = lvl - 1
    frontier := getNeighbors(tx, ty, grid, levels)
    if len(frontier) <= 0 { return } // if no moves
    for _, f := range frontier {
        levels[f.x][f.y] = lvl
    }
    lvl++
    isSourceReached := false
    for len(frontier) > 0 {
        var next []Point
        for _, f := range frontier {
            neighbors := getNeighbors(f.x, f.y, grid, levels)
            for _, n := range neighbors {
                    levels[n.x][n.y] = lvl
                    next = append(next, Point { n.x, n.y })
                    if n.x == sx && n.y == sy { isSourceReached = true }
                    // for id, tile := range g.floors[g.cur].tiles { 
                    // 	if tile.pos.x == n.x && tile.pos.y == n.y {
                    // 		switch levels[n.x][n.y] {
                    // 			case -1: g.floors[g.cur].tiles[id].glyph = 'w'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color20
                    // 			case 0: g.floors[g.cur].tiles[id].glyph = '0'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color21
                    // 			case 1: g.floors[g.cur].tiles[id].glyph = '1'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color22
                    // 			case 2: g.floors[g.cur].tiles[id].glyph = '2'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color23
                    // 			case 3: g.floors[g.cur].tiles[id].glyph = '3'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color24
                    // 			case 4: g.floors[g.cur].tiles[id].glyph = '4'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color25
                    // 			case 5: g.floors[g.cur].tiles[id].glyph = '5'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color26
                    // 			case 6: g.floors[g.cur].tiles[id].glyph = '6'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color27
                    // 			case 7: g.floors[g.cur].tiles[id].glyph = '7'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color28
                    // 			case 8: g.floors[g.cur].tiles[id].glyph = '8'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color29
                    // 			case 9: g.floors[g.cur].tiles[id].glyph = '9'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color30
                    // 			case 10: g.floors[g.cur].tiles[id].glyph = 'A'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color31
                    // 			case 11: g.floors[g.cur].tiles[id].glyph = 'B'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color32
                    // 			case 12: g.floors[g.cur].tiles[id].glyph = 'C'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color33
                    // 			case 13: g.floors[g.cur].tiles[id].glyph = 'D'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color34
                    // 			case 14: g.floors[g.cur].tiles[id].glyph = 'E'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color35
                    // 			case 15: g.floors[g.cur].tiles[id].glyph = 'F'
                    // 			default: 
                    // 				g.floors[g.cur].tiles[id].glyph = 'x'
                    // 				g.floors[g.cur].tiles[id].bg = tcell.Color37
                    // 		}
                    // 	}
                    // }
            }
        }
        if isSourceReached { break }
        frontier = next
        next = nil
        lvl++
    }
    moveLevel := levels[sx][sy] - 1
    pNeighbors := getPNeighbors(sx, sy, grid)
    if pNeighbors == nil { g.dbg("nil"); return }
    var theMoves []Point
    for _, p := range pNeighbors {
        if levels[p.x][p.y] == moveLevel {
            theMoves = append(theMoves, p)
        } 
    }
    var theMove Point
    rand.Seed(time.Now().UnixNano())
    if theMoves == nil {
        theMove = Point { sx, sy }
    } else if len(theMoves) == 1 {
        theMove = theMoves[0]
    } else {
        theMove = theMoves[1]
    }
    dx := theMove.x - sx
    dy := theMove.y - sy
    g.floors[g.cur].aiTakeTurnPathTest(id, dx, dy)
}