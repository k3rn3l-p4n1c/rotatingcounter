package rotating_test

import (
	"github.com/k3rn3l-p4n1c/rotatingcounter"
	"testing"
	"time"
)

func TestNonBlockingRotatingCounter(t *testing.T) {
	resolution := 10 * time.Millisecond
	counter := rotating.NewCounter(3*resolution, resolution, 1)
	defer counter.Stop()

	time.Sleep(resolution / 2)

	counter.Add(20)
	time.Sleep(1 * time.Millisecond)
	if counter.Total() != 20 {
		t.Fatal("wrong total number")
	}

	counter.Add(10)
	time.Sleep(1 * time.Millisecond)
	if counter.Total() != 30 {
		t.Fatal("wrong total number when accumulation happen")
	}

	time.Sleep(resolution)
	counter.Add(5)
	time.Sleep(time.Millisecond)

	if counter.Total() != 35 {
		t.Fatal("wrong total number when shift happen")
	}
	time.Sleep(resolution)
	counter.Add(2)
	time.Sleep(resolution)

	if counter.Total() != 7 {
		println(counter.Total())
		t.Fatal("wrong total number when rotation happen")
	}
	time.Sleep(resolution)
	if counter.Total() != 2 {
		t.Fatal("wrong total number when rotation happen")
	}
	time.Sleep(resolution)
	if counter.Total() != 0 {
		println(counter.Total())
		t.Fatal("wrong total number when rotation happen")
	}
}

func TestBlockingRotatingCounter(t *testing.T) {
	counter := rotating.NewCounter(60 * time.Second, time.Second, 0)
	defer counter.Stop()

	counter.Add(10)

	if counter.Total() != 10 {
		t.Fatal("wrong total number in blocking add")
	}
}

func BenchmarkCounter_Add_NonBlocking(b *testing.B) {
	resolution := 100 * time.Millisecond
	counter := rotating.NewCounter(100 * resolution, resolution, 10)

	for i := 0 ; i < b.N ; i++ {
		counter.Add(1)
	}

	println(counter.Total())
}


func BenchmarkCounter_Add_Blocking(b *testing.B) {
	resolution := 100 * time.Millisecond
	counter := rotating.NewCounter(100 * resolution, resolution, 0)

	for i := 0 ; i < b.N ; i++ {
		counter.Add(1)
	}

	println(counter.Total())
}

func TestCounter_Flush(t *testing.T) {
	resolution := 10 * time.Millisecond
	counter := rotating.NewCounter(3*resolution, resolution, 0)
	counter.Add(100)
	time.Sleep(resolution)
	counter.Add(100)

	if counter.Total() != 200 {
		println(counter.Total())
		t.Fatal("wrong total number when rotation happen")
	}

	counter.Flush()
	counter.Add(1)
	if counter.Total() != 1 {
		t.Fatal("wrong total number after flush")
	}
}