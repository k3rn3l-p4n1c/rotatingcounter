# Rotating Counter

Counter with rotation based on time. For example you want to count an event in past few seconds.

# Example

```go
package main

import (
	"github.com/k3rn3l-p4n1c/rotatingcounter"
	"time"
	)

func main() {
	counter := rotating.NewCounter(60 * time.Second, time.Second, 0)
	defer counter.Stop()

	counter.Add(10)

	println(counter.Total())

```  

# Benchmark
```
BenchmarkCounter_Add_NonBlocking-8 	 10000000	       219 ns/op
BenchmarkCounter_Add_Blocking-8   	 2000000	       773 ns/op
```