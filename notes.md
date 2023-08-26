CUSTOM PROTAGONIST CLASS 
========================
* Implements Draw() and Update()

* Primary Sprite (the ship)
* Current Weapon *Weapon (maybe draw extra stuff on screen)
* Shield (may draw extra (transparent) ssprite on top of primary sprite)
* Mods (upgrades that can alter velocity, acceleration, rotation) 
    * can have time limits

* Orientation (angle)
    * Current
    * Max
    * Min
    * Wanted ( 0 deg if not moving. +5 deg if moving up, -5 deg if moving down )

* Scale (x and y)
* Location
    * X, Y
    * Max X and Y
    * Min X and Y

* Velocity
    * Current (direction and magnitude)
    * Max Magnitude
    * Min Magnitude

* Acceleration
    * Current (direction and magnitude)
    * Max Magnitude
    * Min Magnitude

* Rotation (tumbling - change in orientation)
    * Radians per second

* Collision

WEAPONS
* Beam Weapons Bzzzz
* Turning (changing direction of velocity without affecting magnitude)
    * Rotate velocity vector
    * Radians per second 
