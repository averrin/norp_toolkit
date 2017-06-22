package main

import "C"

// Name is plugin ID
const Name string = "customSpeak"

// Version is version
const Version string = "0.3.0"

//Description is description
const Description string = "Custom Discord layer for OBS"

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
