package main

import "github.com/php-any/origami/internal/pseudocode"

func main() {
	if err := pseudocode.Generate("docs/std"); err != nil {
		panic(err)
	}
}
