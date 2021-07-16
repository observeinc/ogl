package defer_in_loop

func inlineFuncs() {
	defer func() { // this defer is in a different function, not in a loop per-se
	}()

	for i := 0; i < 10; i++ {
		// inline function declaration and call
		_ = func() error {
			defer func() { // this defer is in a different function, not in a loop per-se
			}()
			for {
				defer func() { // want "found defer statement in loop"
				}()
			}
			return nil
		}()

		defer func() { // want "found defer statement in loop"
		}()

		// inline function declaration and call with "go"
		go func() {
			defer func() { // this defer is in a different function, not in the loop per-se
			}()
		}()

		defer func() { // want "found defer statement in loop"
		}()

		// inline function declaration without calling it
		_ = func() {
			defer func() { // this defer is in a different function, not in the loop per-se
			}()
		}

		defer func() { // want "found defer statement in loop"
		}()
	}

	defer func() { // this defer is in a different function, not in a loop per-se
	}()
}
