package defer_in_loop

func inlineFuncs() {
	for i := 0; i < 10; i++ {
		// inline function declaration and call
		_ = func() error {
			defer func() { // this defer is in a different function, not in a loop per-se
			}()
			return nil
		}()

		// inline function declaration and call with "go"
		go func() {
			defer func() { // this defer is in a different function, not in the loop per-se
			}()
		}()

		// inline function declaration without calling it
		_ = func() {
			defer func() { // this defer is in a different function, not in the loop per-se
			}()
		}
	}
}
