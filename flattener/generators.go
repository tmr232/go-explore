package main

/*
To Run:

	go build .\flattener\generate\ ; go generate .\flattener\ ; goimports -w .\flattener\generators_gen.go ; go test .\flattener\
*/

//go:generate go run ./generate

func generate_MyGen() int {
	return 1
	return 2
	return 3
}

func generate_IfStmt(flag bool) int {
	if flag {
		return 1
	} else {
		return 2
	}
	return 3
}

func generate_AnotherIfStmt(flag bool) int {
	return 0
	return 1
	if flag {
		if flag {
			return 2
		}
		return 3
	} else {
		return 5
	}
	return 4
}

func generate_RepeatOne() int {
	for {
		return 1
	}
}

func generate_Fib() int {
	var a int
	var b int
	a = 1
	b = 1
	for {
		return a
		a, b = b, a+b
	}
}

func generate_Count() int {
	var i int
	for {
		return i
		i++
	}
}

func generate_NestedScopes() int {
	i := 1
	return i
	{
		z := 2
		i := 6
		return z
		i = 3 + i
	}
	return i
}
