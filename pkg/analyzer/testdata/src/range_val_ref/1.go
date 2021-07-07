package main

import "fmt"

type myStruct struct {
	f int
}

func main() {
	lst := []myStruct{{1}, {2}, {3}}
	for _, l := range lst {
		fmt.Println("addr: %v", &l)   // want "taking reference of range variable l or its field"
		fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
		for _, m := range lst {
			fmt.Println("addr: %v", &m)   // want "taking reference of range variable m or its field"
			fmt.Println("addr: %v", &m.f) // want "taking reference of range variable m or its field"
			for _, l := range lst {
				fmt.Println("addr: %v", &l)   // want "taking reference of range variable l or its field"
				fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
			}
			for i := range lst {
				fmt.Println("addr: %v", &i) // want "taking reference of range variable i or its field"
			}
			for range lst {
				fmt.Println("addr: %v", &l)   // want "taking reference of range variable l or its field"
				fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
			}
		}
		fmt.Println("addr: %v", &l)   // want "taking reference of range variable l or its field"
		fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
	}

	l := lst[0]
	fmt.Println("addr: %v", &l)

	for i, l := range lst {
		fmt.Println("addr: %v", &i)   // want "taking reference of range variable i or its field"
		fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
		for j, m := range lst {
			fmt.Println("addr: %v", &j)   // want "taking reference of range variable j or its field"
			fmt.Println("addr: %v", &m.f) // want "taking reference of range variable m or its field"
			for i, l := range lst {
				fmt.Println("addr: %v", &i)   // want "taking reference of range variable i or its field"
				fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
			}
			for i := range lst {
				fmt.Println("addr: %v", &i) // want "taking reference of range variable i or its field"
			}
			for range lst {
				fmt.Println("addr: %v", &i)   // want "taking reference of range variable i or its field"
				fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
			}
		}
		fmt.Println("addr: %v", &i)   // want "taking reference of range variable i or its field"
		fmt.Println("addr: %v", &l.f) // want "taking reference of range variable l or its field"
	}

	mp := map[int]int{}
	for k1, v1 := range mp {
		fmt.Println("addr: %v", &k1) // want "taking reference of range variable k1 or its field"
		fmt.Println("addr: %v", &v1) // want "taking reference of range variable v1 or its field"
		for k2, v2 := range mp {
			fmt.Println("addr: %v", &k1) // want "taking reference of range variable k1 or its field"
			fmt.Println("addr: %v", &v1) // want "taking reference of range variable v1 or its field"
			fmt.Println("addr: %v", &k2) // want "taking reference of range variable k2 or its field"
			fmt.Println("addr: %v", &v2) // want "taking reference of range variable v2 or its field"
		}
	}
}
