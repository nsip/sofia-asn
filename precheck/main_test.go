package main

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

func TestMisc(t *testing.T) {
	fmt.Printf("%05s", "a")
}
