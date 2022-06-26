# Entity Definition

## T (Title)  **string**

### *The name of the entity*

Needs to match the name of the config file.

## C (Cost/Points) **int**

### *How many points something is worth*

For turrets, how many points they cost to build.
For enemies, how many points you get for defeating them.
For players, you get nothing.

## H (Health) **int**

### *The number of hit points*

For players and enemies, this is the amount of damage they can sustain before perishing.

Ignored for turrets as they are not subject to damage (yet?)

## r (Hitbox Range) **float**

### *The distance used to judge collision*

For enemies, this determines the range within which a collision is detected.

Ignored for turrets and players.

## D (Damage) **int**

### *The damage the entity does*

For players and turrets this is the damage per projecticle.

Ignored for enemies. (perhaps this could be the damage on contact with player?)

## R (Attack Range) **float**

### *The range of the entity*

For turrets, this is the range they can acquire and fire at targets within.

Ignored for player and enemies.

## X (Attack Rate) **float**

### *The rate of projecticle firing*

For players and turrets, this is the rate at which they can fire projecticles.

Ignored for enemies (perhaps this could be the rate of contact damage).

## N (Number of Projecticles) **int**

## *How many projecticles are shot each attack*

For turrets and players, this determines how many projecticles they fire in a 45degree angle with one shot

## O (Projecticle Speed) **float**

### *The speed of entity's projecticles*

For players and turrets, this is the speed of their projecticles.

Ignored for enemies.

## S (Speed) **float**

### *The speed of the entity*

For players and enemies, this is their movement speed, how fast they can travel throughout the level.

Ignored for turrets because they are immobile :(

## P (Polarity) **positive/neutral/negative**

### *The polarity of the entity*

For players and turrets this all sets the polarity of them and their projecticles.

Ignored for enemies, as their polarity is determined by the polarity of the spawner they emerge from.

## M (Magnetic) **boolean**

### *Whether or not the entity produces a magnetic field*

For enemies, should be set to true
For turrets and players, this might cause some unintended behavior

## Y (Magnet Strength) **float**

### *The strength of the magnetic field*

For all entities, this determines the strength of the magnetic field that surrounds them.

## Z (Magnetic Radius) **float**

### *The radius of the magnetic field*

For all entities, this sets the radius of their magnetic field

## I (Image Prefix) **string**

### *The prefix used to identify an entity's images*

For all entities, this prefix is used to identify the entity's images located in the `assets/images/` folder.

## i (Head Image Prefix) **string**

### *The prefix used to identify an entity's head images*

For turrets, this is used to set their "barrel" image

## W (Walk Image Prefix) **string**

### *The prefix used to identify an entity's walk images*

For players, this is to used to identify the images for the walk animation
