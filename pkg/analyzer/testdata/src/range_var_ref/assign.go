package main

func assignLists() {
	lst := []myStruct{}
	for i, l := range lst {
		i = 10         // want "modifying range variable i or its field"
		l = myStruct{} // want "modifying range variable l or its field"
		l.f = 10       // want "modifying range variable l or its field"
		l.o.b = 10     // This is ok. `o` is a pointer, so l.o.b is same as lst[i].o.b
		_, _ = i, l
	}
}

func assignMaps() {
	mp := map[myStruct]myStruct{}
	for k, v := range mp {
		k = myStruct{} // want "modifying range variable k or its field"
		k.f = 10       // want "modifying range variable k or its field"
		k.o.b = 10     // This is ok. `o` is a pointer, so l.o.b is same as lst[i].o.b
		v = myStruct{} // want "modifying range variable v or its field"
		v.f = 10       // want "modifying range variable v or its field"
		v.o.b = 10     // This is ok. `o` is a pointer, so l.o.b is same as lst[i].o.b
		_, _ = k, v
	}
}

func assignPointers() {
	lst := []*myStruct{}
	for i, l := range lst {
		i = 10          // want "modifying range variable i or its field"
		l = &myStruct{} // want "modifying range variable l or its field"
		l.f = 10        // This is ok. `l` is a pointer, so l.o.b is same as lst[i].o.b
		l.o.b = 10      // This is ok. `l` is a pointer, so l.o.b is same as lst[i].o.b
		_, _ = i, l
	}
	mp := map[*myStruct]*myStruct{}
	for k, v := range mp {
		k = &myStruct{} // want "modifying range variable k or its field"
		k.f = 10        // This is ok. `k` is a pointer, so l.o.b is same as lst[i].o.b
		k.o.b = 10      // This is ok. `k` is a pointer, so l.o.b is same as lst[i].o.b
		v = &myStruct{} // want "modifying range variable v or its field"
		v.f = 10        // This is ok. `v` is a pointer, so l.o.b is same as lst[i].o.b
		v.o.b = 10      // This is ok. `v` is a pointer, so l.o.b is same as lst[i].o.b
		_, _ = k, v
	}
}
