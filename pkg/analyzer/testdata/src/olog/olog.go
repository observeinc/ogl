// In real Observe code, we have a package called olog that requires balanced
// key/value pairs of string,interface. To be able to write test code without
// referencing observe code, dummy up the interface for that module here.
package olog

type Verbose bool

func V(i int) Verbose {
	return Verbose(i < 6)
}

//	This never gets called, but illustrates the requirement.
func assertPairs(args []interface{}) {
	if (len(args) & 1) != 0 {
		panic("uneven pairs")
	}
	for i := 0; i != len(args); i++ {
		if _, is := args[i].(string); !is {
			panic("argument half should be string")
		}
	}
}

func (v Verbose) Info(msg string, args ...interface{}) {
	assertPairs(args)
}

func (v Verbose) Error(msg string, args ...interface{}) {
	assertPairs(args)
}

func Info(msg string, args ...interface{}) {
	assertPairs(args)
}

func Error(msg string, args ...interface{}) {
	assertPairs(args)
}
