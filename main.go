package main

import (
	"fmt"
	"time"
)

func main() {
	beforeQuadTime := time.Now()
	quadraticMain()
	afterQuadTime := time.Since(beforeQuadTime)
	fmt.Printf("%d\n", afterQuadTime)

	before3dQuadTime := time.Now()
	quadratic3dMain()
	after3dQuadTime := time.Since(before3dQuadTime)
	fmt.Printf("%d\n", after3dQuadTime)

	beforeFantasyTime := time.Now()
	fanduelMain()
	afterFantasyTime := time.Since(beforeFantasyTime)
	fmt.Printf("%d\n", afterFantasyTime)
}
