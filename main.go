/*
 *  main.go
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
    "os"
    "log"
)

var enableDebug bool = true
var errorLog string = "./error.log"

func main() {
    var g Game
    var isTurn bool

    // error log initialization
    logFile, err := os.OpenFile(errorLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    if err != nil {
            log.Fatal(err)
    }
    log.SetOutput(logFile)
    defer logFile.Close()  // is this needed here?

    g.initGame()
    log.Println("GAME INITIALIZED")
	s := initScreen()
    for {
        g.tick++
        s.Clear()
        draw(s, g)
        s.Show()
        switch g.state {
        case Playing, Confusion:
            isTurn = g.getInput(s)
            if isTurn {
                g.schedule()
            }
        case Autopilot:
            g.heroRun()
            g.schedule()
        default:
            log.Println("ERROR: A PROBLEM HAS OCCURED")
        }
    }
}