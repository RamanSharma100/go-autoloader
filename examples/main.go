package main

import (
	"fmt"

	autoloader "github.com/RamanSharma100/go-autoloader"
)

func main() {
	engine := autoloader.Load("./")

	fmt.Println()
	fmt.Println("=== MODULES ===")
	fmt.Println()

	for _, module := range engine.Modules.Items {

		fmt.Println("Name:", module.Name)
		fmt.Println("Path:", module.Path)
		fmt.Println("Kind:", module.Kind)
		fmt.Println("Extension:", module.Extension)

		module.Get("Login")

		if module.HasFunction("Login") {
			fmt.Println("Found Login")
			fmt.Println("Calling Login")
			module.Call("Login")
		}

		if len(module.Functions) > 0 {
			fmt.Println("Functions:")

			for _, fn := range module.Functions {
				fmt.Println(" -", fn.Name)
			}
		}

		fmt.Println()
	}
}
