package main

import "C"

// Name is plugin ID
const Name string = "diceSpy"

// Version is version
const Version string = "0.3.1"

//Description is description
const Description string = "Roll20 dice rolls adapter for OBS"

// GetName is name getter
func GetName() string {
	return Name
}

// GetVersion is version getter
func GetVersion() string {
	return Version
}

// GetDescription is desc getter
func GetDescription() string {
	return Description
}
