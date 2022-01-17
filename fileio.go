/*
 *  fileio.go
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
	"encoding/gob"
	//"fmt"
	//"compress/gzip"
	"os"
	//"bytes"
)

func (g *Game) loadGame() error {

    fi, err := os.Open("save.gob")
    if err !=nil {
        return err
    }
    defer fi.Close()

    // fz, err := gzip.NewReader(fi)
    // if err !=nil {
    //     return err
    // }
    // defer fz.Close()

	//decoder := gob.NewDecoder(fz)
    decoder := gob.NewDecoder(fi)
    err = decoder.Decode(&g.hero)
    if err !=nil {
        return err
    }

    return nil
}

func (g *Game) saveGame() error {

    fi, err := os.Create("save.gob")
    if err !=nil {
        return err
    }
    defer fi.Close()

    // fz := gzip.NewWriter(fi)
    // defer fz.Close()

	//encoder := gob.NewEncoder(fz)
    encoder := gob.NewEncoder(fi)
    err = encoder.Encode(g.hero)
    if err !=nil {
        return err
    }

    return nil
}




// func (g *Game) saveGame() {
// 	//data := []int{101, 102, 103}
// 	var data Game
// 	// data = *g
// 	// buf := new(bytes.Buffer)
// 	// //glob encoding
// 	// enc := gob.NewEncoder(buf)
// 	// enc.Encode(data)
// 	//fmt.Println("Encoded:", data)  //Encoded: ABC


// 	// create a file
// 	dataFile, err := os.Create("savegame.gob")

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	 // serialize the data
// 	dataEncoder := gob.NewEncoder(dataFile)
// 	dataEncoder.Encode(data)

// 	dataFile.Close()
// }

// func (g *Game) loadGame(){
// 	var data Game

// 	// //glob decoding
// 	// d := gob.NewDecoder(buf)
// 	// d.Decode(data)


// 	// open data file
// 	dataFile, err := os.Open("savegame.gob")
// 	if err != nil {
// 		fmt.Println("333")
// 		os.Exit(1)
// 	}

// 	// d := gob.NewDecoder(dataFile)
// 	// d.Decode(&data)

// 	dataDecoder := gob.NewDecoder(dataFile)
// 	err = dataDecoder.Decode(&data)

// 	// if err != nil {
// 	// 	fmt.Println("222")
// 	// 	os.Exit(1)
// 	// }

// 	dataFile.Close()
// 	//g = nil
// 	//g = &data
// 	//fmt.Println(data)
// 	g = &data
// 	//return data
// }