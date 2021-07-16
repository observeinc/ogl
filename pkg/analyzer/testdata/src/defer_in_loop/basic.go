package defer_in_loop

import "fmt"

func myFunc(n int) {
	defer func() {}()
	for i := 0; i < n; i++ {
		defer fmt.Printf("%d", i) // want "found defer statement in loop"
		if true {
			defer fmt.Printf("%d", i) // want "found defer statement in loop"
		} else if false {
			defer fmt.Printf("%d", i) // want "found defer statement in loop"
		} else {
			defer fmt.Printf("%d", i) // want "found defer statement in loop"
		}
		for _ = range []int{0, 1} {
			switch i {
			case 0:
				defer fmt.Printf("%d", i) // want "found defer statement in loop"
				if true {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				} else if false {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				} else {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				}
			case 1:
				defer fmt.Printf("%d", i) // want "found defer statement in loop"
				if true {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				} else if false {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				} else {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				}
			default:
				defer fmt.Printf("%d", i) // want "found defer statement in loop"
				if true {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				} else if false {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				} else {
					defer fmt.Printf("%d", i) // want "found defer statement in loop"
				}
			}
		}
		{
			defer fmt.Printf("%d", i) // want "found defer statement in loop"
		}
		// TODO: select/comm statements
	}
	defer func() {}()
}
