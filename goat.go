package main // rename to goat, eventually

import "goat/shed" // Goat helps you draw primitives on screen

type Goat struct {
	options []*GOptions
}

// Goat options. Used to save and retrieve drawing settings
// And turn the Goat into a state machine.
type GOptions struct {
	fgColor       shed.V4
	bgColor       shed.V4
	lineThickness float32
	// mask, zoom, stuff
}
