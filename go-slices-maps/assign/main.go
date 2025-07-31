package main

type address struct {
	city string
}

type person struct {
	name string
	age  int
	addr address
}

func main() {
	pv1 := person{
		name: "Alice",
		age:  30,
		addr: address{
			city: "Wonderland",
		},
	}

	pv2 := pv1

	pv2.name = "Bob"
	pv2.addr.city = "Builderland"

	println(pv1.name, pv1.age, pv1.addr.city) // Alice 30 Wonderland
	println(pv2.name, pv2.age, pv2.addr.city) // Bob 30 Builderland

	v1 := 10
	v2 := v1
	v2 = 20

	println(v1) // 10
	println(v2) // 20

	pp1 := &person{
		name: "Charlie",
		age:  25,
		addr: address{
			city: "Chocoland",
		},
	}

	pp2 := pp1

	pp2.name = "Dave"
	pp2.addr.city = "Daveland"
	println(pp1.name, pp1.age, pp1.addr.city) // Dave 25 Daveland
	println(pp2.name, pp2.age, pp2.addr.city) // Dave 25 Daveland
}
