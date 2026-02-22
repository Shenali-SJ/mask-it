package main

import (
	"fmt"
	"io"
	"os"

	"maskit/mask"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(mask.Mask(string(input)))
}
