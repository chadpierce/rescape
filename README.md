# Rescape

- Rescape is a text-based roguelike written in Go. 
- At this time it is not a playable game. 
- Development is in progress and being tracked and blogged here: https://rescape.dev

## Design

- The initial idea for Rescape was "*Nethack, but without the kitchen sink*" - meaning a turn based dungeon crawler without all of the spoilers and secrets, and that doesn't take several days and a wiki to finish
- As of this writing Rescape is far from finished (or fully designed)
- I am not a professional software developer and probably do not write good and efficient code
- Use of external libraries and other people's code is kept to a minimum. Currently the only external package is used for i/o ([tcell](https://github.com/gdamore/tcell)). This is similar to the curses library.

## Gameplay

- It uses typical glyphs to represent the hero, monsters, items, etc
- Unlike most other games, you start at the bottom of a dungeon and work your way up
- There are 2 ways to win:
	- SHORT: escape the dungeon and run away
	- LONG: escape the dungeon and then climb the Wizard's tower to defeat them and take their amulet
- Available classes are fighter, rogue, and wizard
- No player race or gender
- No experience points or character levels
- There are 3 primary skills: strength, intelligence, and dexterity
- There are several item categories and types, including weapons, armor, rings, potions, scrolls, books, and amulets
- Skills, weapons, and armor are enhanced with consumable items, either temporarily or permanently (blessed scrolls and potions can permanently alter the hero or some items)
- There are two gods that can optionally be worshiped by the player and some monsters
- Gods may bless your consumable items
- In place of a "hunger clock" that is integral to most roguelikes, there is a chaotic miasma leaking into the dungeon that slowly envelops each level. This generates monsters (or will mutate existing monsters) and can cause harm to the player. This game dynamic keeps the game moving along

*note: many of these components are still in the design phase*

## Build

The game can currently built by downloading the repo: 

	git clone https://github.com/chadpierce/rescape.git

Change directory into the repo:

	cd rescape

Run the application:

	go run . 

Or build the application:

	go build .

- Your terminal may need to be resized to fit the entire game.
- Default colors in a windows shell may be bad.
- Only vi keys are currently supported because that is what I use. 

## Versions

I will create tags and binaries at certain points during development. 

The current versions is playable, but not really much of a game. It is in very early development stages and there is no actual structure. Basically run around the dungeon and try to kill to stuff until you die.

Version 0.1: LINK TBD

## Build

To compile and run the game:

Install Go

	https://go.dev/doc/install

Clone this repo

	git clone git@github.com:chadpierce/rescape.git

Enter the repo directory and run the following command to run the app

	go run .

To build an executable

	go build . 
