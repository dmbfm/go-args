package main

import (
	"fmt"

	"github.com/dmbfm/go-args"
)

var options struct {
	All   bool
	Force bool
	Title string
	files []string
}

func main() {
	var err error

	args.AddBool("all", "a", &options.All, "Apply for all")
	args.AddBool("force", "f", &options.Force, "Force action")
	args.AddString("title", "t", &options.Title, "The title")

	options.files, err = args.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", options)

}
