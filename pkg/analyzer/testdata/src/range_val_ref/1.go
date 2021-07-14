package main

import (
	"fmt"
	"os"
)

type otherStruct struct {
	b int
}

type myStruct struct {
	f int
	o *otherStruct
}

func lists() {
	lst := []myStruct{}
	for _, l := range lst {
		fmt.Println("addr: %v", &l)   // want "taking reference of range variable l or its field"
		fmt.Println("addr: %v", &l.o) // want "taking reference of range variable l or its field"
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

	for i, l := range os.Environ() {
		fmt.Printf("addr: %v", &i) // want "taking reference of range variable i or its field"
		fmt.Printf("addr: %v", &l) // want "taking reference of range variable l or its field"
	}

	for i, l := range os.Args {
		fmt.Printf("addr: %v", &i) // want "taking reference of range variable i or its field"
		fmt.Printf("addr: %v", &l) // want "taking reference of range variable l or its field"
	}

	for i, l := range []myStruct{} {
		fmt.Printf("addr: %v", &i) // want "taking reference of range variable i or its field"
		fmt.Printf("addr: %v", &l) // want "taking reference of range variable l or its field"
	}

	for i, l := range []*myStruct{} {
		fmt.Printf("addr: %v", &i) // want "taking reference of range variable i or its field"
		fmt.Printf("addr: %v", &l) // want "taking reference of range variable l or its field"
	}
}

func maps() {
	mp := map[myStruct]myStruct{}
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

func retListPtr() []*myStruct {
	return nil
}

func pointers() {
	lst := []myStruct{}
	for i, l := range lst {
		fmt.Printf("addr: %v", &i)     // want "taking reference of range variable i or its field"
		fmt.Printf("addr: %v", &l)     // want "taking reference of range variable l or its field"
		fmt.Printf("addr: %v", &l.o)   // want "taking reference of range variable l or its field"
		fmt.Printf("addr: %v", &l.o.b) // This is ok. `o` is a pointer, so &l.o.f is same as &lst[i].o.f
	}

	plst := []*myStruct{}
	for i, l := range plst {
		fmt.Printf("addr: %v", &i)     // want "taking reference of range variable i or its field"
		fmt.Printf("addr: %v", &l)     // want "taking reference of range variable l or its field"
		fmt.Printf("addr: %v", &l.o)   // This is ok. `l` is a pointer, so &l.o is same as &lst[i].o
		fmt.Printf("addr: %v", &l.o.b) // This is ok. `l` is a pointer, so &l.o.f is same as &lst[i].o.f
	}

	for i, l := range retListPtr() {
		fmt.Printf("addr: %v", &i)     // want "taking reference of range variable i or its field"
		fmt.Printf("addr: %v", &l)     // want "taking reference of range variable l or its field"
		fmt.Printf("addr: %v", &l.o)   // This is ok. `l` is a pointer, so &l.o is same as &lst[i].o
		fmt.Printf("addr: %v", &l.o.b) // This is ok. `l` is a pointer, so &l.o.f is same as &lst[i].o.f
	}

	mp1 := map[myStruct]myStruct{}
	for k1, v1 := range mp1 {
		fmt.Println("addr: %v", &k1)     // want "taking reference of range variable k1 or its field"
		fmt.Println("addr: %v", &k1.o)   // want "taking reference of range variable k1 or its field"
		fmt.Println("addr: %v", &k1.o.b) // This is ok. `o` is a pointer
		fmt.Println("addr: %v", &v1)     // want "taking reference of range variable v1 or its field"
		fmt.Println("addr: %v", &v1.o)   // want "taking reference of range variable v1 or its field"
		fmt.Println("addr: %v", &v1.o.b) // This is ok. `o` is a pointer
	}

	mp2 := map[myStruct]*myStruct{}
	for k1, v1 := range mp2 {
		fmt.Println("addr: %v", &k1)     // want "taking reference of range variable k1 or its field"
		fmt.Println("addr: %v", &k1.o)   // want "taking reference of range variable k1 or its field"
		fmt.Println("addr: %v", &k1.o.b) // This is ok. `o` is a pointer
		fmt.Println("addr: %v", &v1)     // want "taking reference of range variable v1 or its field"
		fmt.Println("addr: %v", &v1.o)   // This is ok. `v1` is a pointer
		fmt.Println("addr: %v", &v1.o.b) // This is ok. `v1` is a pointer
	}

	mp3 := map[*myStruct]myStruct{}
	for k1, v1 := range mp3 {
		fmt.Println("addr: %v", &k1)     // want "taking reference of range variable k1 or its field"
		fmt.Println("addr: %v", &k1.o)   // This is ok. `k1` is a pointer
		fmt.Println("addr: %v", &k1.o.b) // This is ok. `k1` is a pointer
		fmt.Println("addr: %v", &v1)     // want "taking reference of range variable v1 or its field"
		fmt.Println("addr: %v", &v1.o)   // want "taking reference of range variable v1 or its field"
		fmt.Println("addr: %v", &v1.o.b) // This is ok. `o` is a pointer
	}

	mp4 := map[*myStruct]*myStruct{}
	for k1, v1 := range mp4 {
		fmt.Println("addr: %v", &k1)     // want "taking reference of range variable k1 or its field"
		fmt.Println("addr: %v", &k1.o)   // This is ok. `k1` is a pointer
		fmt.Println("addr: %v", &k1.o.b) // This is ok. `k1` is a pointer
		fmt.Println("addr: %v", &v1)     // want "taking reference of range variable v1 or its field"
		fmt.Println("addr: %v", &v1.o)   // This is ok. `v1` is a pointer
		fmt.Println("addr: %v", &v1.o.b) // This is ok. `v1` is a pointer
	}
}
