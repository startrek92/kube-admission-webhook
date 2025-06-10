package main

import ("fmt")

// go seperates each line either by trailing \n or ; semi colon
// a line cannot start with {


var a int = 1;
var b int = 1;

var x, z int = 4,5;

var h, j = 5, "Hello"; // type will be inferred by compiler when compiling the code

// := "hello" -> Error here as := can only be used inside function


// Constants -> cannot update, once declared, can be declared inside & outside the funciton
// generally constants are names in uppercase for differentiating from variables

const PI float32 = 3.14; // type inferred at compile time
// multiple declaration in blocks
const (
	c1 = 3
	c3 int = 4
	c5 string = "cruel world"
);

func printVars() {

	t, y := 88, "Cruel World";
	fmt.Println(a, b);
	fmt.Println(x, z);
	fmt.Println(t, y, h, j);
	fmt.Println(c1, c5, c3);
	
	// print functions
	/*
	
	Print()
		- prints in single line, it non string then func adds a space, does not add 
		a new line in the end
	
	Println()
		- adds space between args passed in the function and adds a new line in the end
	
	Printf()
		- prints arguments based on the flags passed in the string
	
		*/
	fmt.Printf("c5 value: %v and type %T \n", c5, t);
	// %v for value, %T for type
}


func loops() {
	/*
	go has only for loop and no bracktes are required for loops
	*/

	for i := 0; i<10; i++ {
		fmt.Println("hell", i);
	}


	/*
	range, go uses range for iterating over slices or maps
	*/
	var arr = []int {1, 3, 4, 5, 6};
	for index, value := range arr {
		fmt.Println(index, value)
	}
}


func ifElse() {
	/*
	we can define variable itself in if statements
		but it will be available only in subsequent if else statements
		not after else statements 
	*/

	// var x int = 10;
	if x1:= 100; x1 >= 10 {
		fmt.Println("gt 10");
	} else {
		fmt.Println("n gt 10");
	}

	// fmt.Println(x1) --> will raise error
}

func main() {
	// printVars()
	// loops()
	ifElse()
}

