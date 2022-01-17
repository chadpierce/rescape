# rescape

## story

The name Rescape is a portmanteau of Rogue Escape.

The game’s story is that you begin trapped at the bottom of a dungeon. Your goal is to escape, and if you choose to do so, climb the wizards tower and take the Amulet of Yendor.

Perhaps your roguish adventurer was captured on their first attempt and this is what happens once they break free.

In the depths of the dungeon one of the Wizard’s experiments has unleashed chaotic magic into the world. Your bonds have been broken and the miasma forces you up to the surface.

## general design

A simple ascii roguelike style using typical glyphs to represent the hero, monsters, items, etc
Unlike most other games, you start at the bottom of a dungeon and work your way up
There are 2 ways to win:
SHORT: escape the dungeon and run away
LONG: escape the dungeon and then climb the Wizard’s tower to defeat them and take their amulet
Available classes are fighter, rogue, and wizard
No player race or gender
No experience points or character levels
There are 3 primary skills: strength, intelligence, and dexterity
There are several item categories and types, including weapons, armor, rings, potions, scrolls, books, and amulets
Skills, weapons, and armor are enhanced with consumable items, either temporarily or permanently (blessed scrolls and potions can permanently alter the hero or some items)
There are two gods that can optionally be worshipped by the player and some monsters
Gods may bless your consumable items
In place of a “hunger clock” that is integral to most roguelikes, there is a chaotic miasma leaking into the dungeon that slowly envelops each level. This generates monsters (or will mutate existing monsters) and can cause harm to the player. This game dynamic keeps the game moving along
note: many of these components are still in the design phase

## hero classes, starting items, stat bonus

- fighter: long sword/axe, armor (+str)
- rogue:  dagger, cloak, darts (+dex)
- wizard: staff, robe, hat (+int) 

## religion

- there are 2 gods, light and dark (I'd like to come up with better names for these)
- gods can bless items at alters
- blessed consumables often have permanent effects
- blessings are limited in some way - need to figure this out
possible idea: you can enter opposing temples to kill light/dark beasts to please your god
this grants blessings from your god's temple. clearing and destroying a temple grants a god gift. 

## dungeon features

- dungeon levels
- ground level
- tower/castle levels
- temples (dark and light)
- treasure rooms? (randomly generated)

## mobs/monsters

- typical roguelike monsters
- prison guards
- wizard of yendor
- others tbd

## items

- potions - normal temporarily raise stats, blessed permanently raise stats, cursed temporarily lower stats
- scrolls
- books
- rings - can be enchanted, blessed
- amulets - can be enchanted, blessed
- weapons - can be enchanted, blessed, have brands
- armor - can be enchanted, blessed, have brands

## machines

- traps? maybe
- doors 
