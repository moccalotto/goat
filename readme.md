GOAT Graphics
=============


Graphics lib with stupid names.


## Shed
The namespace where all the mischelaneous utilities are located

## Tractor 
The namespace where the "engine" is located
### Engine
* Singleton. 
* No methods. 
* All "methods" are functions that refer to the singleton.
* The singleton itself is hidden, but it can be replaced, pushed and popped.
* It is lazy-created.

### Window (Wind Shield)

### Controls (keyboard, mouse)

### Vroom (Audio)

## Pilot
The scripting library. Lua.

### Ghost
Available inside lua.
Sprites.

### Turtle
Available inside lua.
The turtle is a state machine. 
It can move around, turn around, and draw lines as it moves. #TurtleTalk

### Goat
Available inside lua.
Draw all the primitives.

### "Background" (better name needed)
Handle background images, scrolling, swapping out, etc.
