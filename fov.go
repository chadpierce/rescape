/*
 *  fov.go
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
	"sort"
    //"fmt"
    // "log"
    // "os"
    //"strconv"
	//"sync"

    //"github.com/gdamore/tcell/v2"
)

const FOVDistance = 15

type Points []Point

func pushPoint(p []Point, px, py int) []Point {
	var newPoint = Point {
		x: px,
		y: py,
	}
	p = append(p, newPoint)
	return p
}

func getFOVPerimeterVector(hero Point) []Point {
    var p []Point
    //build square around hero - 
    fx := hero.x - FOVDistance;
    fy := hero.y - FOVDistance;
    fw := fx + FOVDistance*2;
    fh := fy + FOVDistance*2;

	for i := fx; i <= fw; i++ {
        p = pushPoint(p, i, fy)
		p = pushPoint(p, i, fh)
    }
	for i := fy; i <= fh; i++ {
        p = pushPoint(p, fx, i)
		p = pushPoint(p, fw, i)
    }
    return p
}

func (f *Floor) getAllView(hero Actor) []string {
	var actorsInView []string
	actorsInView = append(actorsInView, "fix this!")
	for i, _ := range f.tiles {
		f.tiles[i].visible = true
	}
	for i, _ := range f.actors {
		f.actors[i].visible = true
	}
	for i, _ := range f.items {
		f.items[i].visible = true
	}
	return actorsInView
}

func (f *Floor) isActorsInFOV(hero Actor) (bool, string) {
	isA := false
	seen := ""
    perimeterVec := getFOVPerimeterVector(hero.pos);
    for _, p := range perimeterVec {
        points := getLine(hero.pos.x, hero.pos.y, p.x, p.y);
        cnt := 0
		out:
		for _, line := range points {
			for _, t := range f.tiles {
				if line.x == t.pos.x && line.y == t.pos.y {
					if t.blockSight == true {
						break out
					}
				}
			}
			for _, a := range f.items {
				if line.x == a.pos.x && line.y == a.pos.y {
					if a.blockSight == true {
						break out
					}
				}
			}
			for _, a := range f.actors {
				if line.x == a.pos.x && line.y == a.pos.y && a.alive {
					isA = true
					seen = a.name
					break out
				}
			}
		}
		cnt++
		if cnt == FOVDistance { break }
    }
	return isA, seen
}

 type ActorDistance struct {
	 actor string
	 distance int
 }

func (f *Floor) getFOV(hero Actor) []string {
    perimeterVec := getFOVPerimeterVector(hero.pos);
	var actorDistance []ActorDistance
	var actorsInView []string
	var processedActors []int
	// normal gameplay
	for i, _ := range f.tiles {
		f.tiles[i].visible = false
	}
	for i, _ := range f.actors {
		f.actors[i].visible = false
	}
	for i, _ := range f.items {
		f.items[i].visible = false
	}
	//see everything all the time
	// for i, _ := range f.tiles {
	// 	f.tiles[i].visible = true
	// }
	// for i, _ := range f.actors {
	// 	f.actors[i].visible = true
	// }
	// for i, _ := range f.items {
	// 	f.items[i].visible = true
	// }
	// see map but not visible 
	// for i, _ := range f.tiles {
	// 	f.tiles[i].visible = false
	// 	f.tiles[i].visited = true
	// }
	// for i, _ := range f.actors {
	// 	f.actors[i].visible = false
	// 	f.tiles[i].visited = true
	// }
	// for i, _ := range f.items {
	// 	f.items[i].visible = false
	// 	f.tiles[i].visited = true
	// }
	//outout:
    for _, p := range perimeterVec {
        points := getLine(hero.pos.x, hero.pos.y, p.x, p.y);
        cnt := 0
		//out:
		out:  // label for break
        for lineID, line := range points {
			for i, a := range f.actors {
				isProcessed := false
				for _, p := range processedActors {
					if i == p { isProcessed = true ; break }
				}
				if !isProcessed {
					if a.alive {
						itemStr := ""
						if line.x == a.pos.x && line.y == a.pos.y {
							aAI := ""
							if a.ai == SleepAI {
								aAI = "(sleeping)"
							} else if a.ai == ConfusedAI { 
								aAI = "(confused)"
							}
							if a.weapon == -1 || a.inv[a.weapon].item.name == "" {
								itemStr = a.name
							} else {
								itemStr = a.name + " [" + a.inv[a.weapon].item.name + "]"
							}
							if aAI != "" {
								itemStr = itemStr + " " + aAI
							}
							actorDistance = append(actorDistance, ActorDistance { itemStr, lineID})
							processedActors = append(processedActors, i)
						
							f.actors[i].visible = true
							f.actors[i].visited = true
							if f.actors[i].ai == SleepAI {
								if sleepCheck(hero, a) {
									//wake up
									f.actors[i].ai = f.actors[i].defaultAI
									f.actors[i].target = hero.pos
								}
							} else if f.actors[i].ai == WanderAI {
								if noticeCheck(hero, a) {
									f.actors[i].ai = f.actors[i].defaultAI
									f.actors[i].target = hero.pos
								}
							} else if f.actors[i].ai == ConfusedAI {
							} else {
								f.actors[i].target = hero.pos
							}
							f.actors[i].lastKnownPos = a.pos
						} else if line.x == a.lastKnownPos.x && line.y == a.lastKnownPos.y {
							f.actors[i].lastKnownPos = Point { -1, -1 }
							f.actors[i].visible = false
						} else {  // alive
							if line.x == a.pos.x && line.y == a.pos.y {  
								f.actors[i].visible = true
								f.actors[i].visited = true
							}
						}
					} else { // dead
						if line.x == a.pos.x && line.y == a.pos.y {  
							f.actors[i].visible = true
							f.actors[i].visited = true
						}
				}
				if a.blockSight == true {
					break out
				}
			}
			}
			for i, o := range f.items {
				if line.x == o.pos.x && line.y == o.pos.y {
					f.items[i].visible = true
					f.items[i].visited = true
					if o.blockSight == true {
						break out
					}
				}
			}
			for i, t := range f.tiles {
				if line.x == t.pos.x && line.y == t.pos.y {
					f.tiles[i].visible = true
					f.tiles[i].visited = true
					if t.blockSight == true {
						break out
					}
				}
			}
        }
		cnt++
		if cnt == FOVDistance { break }
    }
	sort.SliceStable(actorDistance, func(i, j int) bool {
		return actorDistance[i].distance < actorDistance[j].distance
	})
	for _, a := range actorDistance {
		actorsInView = append(actorsInView, a.actor)
	}
	return actorsInView
}

func getDistanceBetween(a, b Point) int {
	line := getLine(a.x, a.y, b.x, b.y)
	return len(line)
}
// call this one
func getLine(ax, ay, bx, by int) []Point {
    p1 := Point { x: ax, y: ay }
    p2 := Point { x: bx, y: by }
    points := calcLine(p1, p2);
    return points
}

//Bresenham's line algorithm from roguebasin
func calcLine(pos1, pos2 Point) (points []Point) {
	x1, y1 := pos1.x, pos1.y
	x2, y2 := pos2.x, pos2.y
 
	isSteep := abs(y2-y1) > abs(x2-x1)
	if isSteep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}
 
	reversed := false
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		reversed = true
	}
 
	deltaX := x2 - x1
	deltaY := abs(y2 - y1)
	err := deltaX / 2
	y := y1
	var ystep int
 
	if y1 < y2 {
		ystep = 1
	} else {
		ystep = -1
	}
 
	for x := x1; x < x2+1; x++ {
		if isSteep {
			points = append(points, Point{y, x})
		} else {
			points = append(points, Point{x, y})
		}
		err -= deltaY
		if err < 0 {
			y += ystep
			err += deltaX
		}
	}
 
	if reversed {
		//Reverse the slice
		for i, j := 0, len(points)-1; i < j; i, j = i+1, j-1 {
			points[i], points[j] = points[j], points[i]
		}
	}
 
	return
}
 
func abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0
	}
	return x
}