package olog_unbalanced_args

import "olog"

func BasicFunc(n int) {
	olog.Info("This is the message", "x", 3, "y", 5, "z")  // want "olog key is missing value"
	olog.Info("This is the message", "x", 3, "y", 5)       // ok
	olog.Error("This is the message", "x", 3, "y", 5, "z") // want "olog key is missing value"
	olog.Error("This is the message", "x", 3, "y", 5)      // ok

	// unpacks are OK; they are not type checked more than this
	olog.Info("foo", []interface{}{"x", 3}...)
	olog.Error("foo", []interface{}{"x", 3}...)

	// note: not unpacked!
	olog.Info("foo", []interface{}{"x", "y"})  // want "olog key is missing value"
	olog.Error("foo", []interface{}{"x", "y"}) // want "olog key is missing value"

	if olog.V(3) {
		olog.Info("The message is here", 3, "x") // want "olog key is not a string"
		olog.Info("The message is here", "x", 3) // ok
	}
	if olog.V(3) {
		olog.Error("The message is here", 3, "x") // want "olog key is not a string"
		olog.Error("The message is here", "z", 3) // ok
	}

	for j := 0; j != 2; j++ {
		olog.V(j).Info("This is the message", "foo", j, "bar") // want "olog key is missing value"
		olog.V(j).Info("This is the message", "foo", j)        // ok
	}
	for j := 0; j != 2; j++ {
		olog.V(j).Error("This is the message", "foo", j, "bar") // want "olog key is missing value"
		olog.V(j).Error("This is the message", "foo", j)        // ok
	}
}
