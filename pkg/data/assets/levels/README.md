# Level Definition

## T (Title)  **string**

### *The title of the level*

As it says.

## S (TileSet) **string**

### *The tileset to use*

The tileset to use for map cells. The name corresponds to a subdirectory in the images directory.

## W (Waves) **[]Wave**

### *The waves configuration*

These define wave configurations for the spawners. There should be an amount of "W" lines equal to the amount of spawners. A Wave definition line corresponds to the next spawner in the map, reading from top-left to bottom-right.

The syntax is for a single spawn is `<AMOUNT>[@<TICK DELAY>] <ENEMY>[&<ENEMY>...]`, with multiple spawns in a wave separate by a `,`,  and multiple waves by using a `;` delimiter.

For example, 3 waves could be defined as follows: `5@20 walker,2@20 runner;10@20 walker;15@10 walker`. This would result in 3 waves, with the first consisting of 1 walker spawning every 20 ticks 5 times, then 1 runner spawning every 20 ticks 2 times. The second wave would be 1 walker spawning every 20 ticks. The third would be 15 walkers spawning every 10 ticks.
