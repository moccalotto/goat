FEATURES
=========================
* Text
    * Level 0 Sprites: used for text. A default LVL sprite is included (but lazy loaded)
        * TextShader
    * Can load different atlas.
    * Can "map" text characters to text sprites.
        * Such info might also be provided in the XML
        * Can map multiple runes to same image (to be case insensitive or to mark many runes as "?")
* Sprites
    * Easy to create
    * ALWAYS stored in atlasses
    * Only two atlasses (slots) at a time 
        * Level 2 Sprites: used for foreground objects.
            * FgSpriteShader
        * Level 3 Sprites: Used for movable background objects.
            * BgSpriteShader
    * Can load new atlasses for new scenes
    * Always `gl_clamp`
* Backgrounds (Maybe called scenes)
    * Maybe called scenes to not clash with bg sprites.
    * Easy to create
    * Can scroll
    * Can load new backgrounds for new scenes,
    * Always `gl_wrap`
* Geometries
    * Lines, rects, circles, etc. do not count as sprites
    * Easy to create

* Machine/Engine
    * The object is a singleton
    * All its *methods* become public *functions* that refer to the (current) global singleton object

* Behavior
    * All sprites, backgrounds, primitives have two functions: *draw* and *update*
    the only way to increase the number of behaviors, is to wrap the object in a similar object, 
    but with another update() func, that refers to the original update() func

* Rendering order
    1. Background
    1. bgSprites
    1. fgSprites
    1. Geometries
    1. TxtSprites

Rectangles 
==========================
* Fully filled rectangle
* Any scale, position, and orientation
* Any color
* Runded corners

Fancy Rectangles
==========
* Corner radius
* Stroke width
* Infill color
* In any orientation, scale, and location


Lines
=====
(are basically just rectangles)
* And line segments


Circles & ellipses
==================
* and slices thereof
* with and without infill
* line thickness
* segment count


Images
======
* Load an image
* Put it anywhere on screen
* in any rotation
* in any scale
* possibly other matrix operations


Polygon, Polyline, Path
=======================
* https://developer.mozilla.org/en-US/docs/Web/SVG/Tutorial/Basic_Shapes#path
