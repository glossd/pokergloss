package main

import (
	"fmt"
	"math"
)

func main() {
	var a1 = 0.01
	var a2 = 0.09
	var a3 = 0.1
	var a4 = 0.12
	var a5 = 0.18
	var a6 = 0.3 // 1x
	var a7 = 0.05
	var a8 = 0.025
	var a9 = 0.015
	var a10 = 0.009
	var a11 = 0.001
	sum := a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8 + a9 + a10 + a11
	if math.Round(sum*1000) != 900.0 {
		panic(fmt.Sprintf("sum not 0.9, but %.2f", sum))
	}
	res := a1*0.25 + a2*0.5 + a3*0.75 + a4*0.8 + a5*0.9 + a6*1.0 + a7*1.5 + a8*2.0 + a9*3.0 + a10*5.0 + a11*10.0
	fmt.Println(res)
	if res != 0.9 {
		panic(fmt.Sprintf("res not 0.9, but %.2f", res))
	}
}


