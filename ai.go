/*
 *  ai.go
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
    //"log"
    //"os"
    //"strconv"

    //"github.com/gdamore/tcell/v2"
)

type actorAI int

const (
    NoAI actorAI = iota
    BasicAI
    RangedAI
    ConfusedAI
    ScaredAI
    WizardAI
    SleepAI
    WanderAI
)

func (f Floor) basic(id, dx, dy int) {
    f.actors[id].pos.x += dx
    f.actors[id].pos.y += dy
}

func (f Floor) aiTakeTurn(id, dx, dy int) {

    ai := f.actors[id].ai
    switch ai {
    case NoAI:
        return
    case BasicAI:
        f.basic(id, dx, dy)
    }
}

func (f *Floor) aiTakeTurnPathTest(id, dx, dy int) {
        f.actors[id].pos.x += dx
	    f.actors[id].pos.y += dy
}

func (g *Game) heroAiTakeTurnPathTest(dx, dy int) {
        g.hero.pos.x += dx
	    g.hero.pos.y += dy
}

func (f *Floor) aiChooseFinalMove(id int, moves []Point) {
	var unblockedMoves []Point
    dx := 0
    dy := 0
    for i, m := range moves {
        isB, _ := isBlocked(f.actors[i].pos.x + m.x, f.actors[i].pos.y + m.y, f.actors, f.tiles)
        if !isB {
            unblockedMoves = append(unblockedMoves, m)
        }
    }
    if len(unblockedMoves) <= 0 {
        fmt.Println("no unblocked moves!!!!")
        dx = 0
        dy = 0
    } else if len(unblockedMoves) == 1 {
        dx = unblockedMoves[0].x
        dy = unblockedMoves[0].y
    } else {
        rand.Seed(time.Now().UnixNano())
        theMove := unblockedMoves[rand.Intn(len(unblockedMoves))]
        dx = theMove.x
        dy = theMove.y 
    }
    dx = dx - f.actors[id].pos.x
    dy = dy - f.actors[id].pos.y
    f.actors[id].pos.x += dx
    f.actors[id].pos.y += dy
}

func distanceNoiseMultiplier(d int) int {
    switch d {
    case 1: return 20
    case 2,3: return 18
    case 4,5: return 16
    case 6,7: return 14
    case 8,9: return 12
    default: return 5
    }
}

func noiseCheck(hero, mob Actor) int {
    // TODO figure out how to calculate this
    // distance, weight of equipped objects, dexterity
    distance := getDistanceBetween(hero.pos, mob.pos)
    distNoise := distanceNoiseMultiplier(distance)
    noise := distNoise - hero.dex
    noise = noise + roll(R2d4)
    return noise
}

func sleepCheck(hero, mob Actor) bool {
    is := false
    noise := noiseCheck(hero, mob)
    if noise >= 10 {
        is = true
    }
    return is
}

func noticeCheck(hero, mob Actor) bool {
    is := false
    noise := noiseCheck(hero, mob)
    if noise >= 5 {
        is = true
    }
    return is
}