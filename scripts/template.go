package main

import (
	"io"
	"os"
)

// Creates a .go file containing index.html as a string literal.
// This allows us not needing to deploy index.html.
//
// Note: run 'go generate' before 'go build'
func main() {
	// open source
	in, err := os.Open("index.html")
	panicIf(err)
	defer in.Close()

	// open destination
	out, err := os.Create("index.html.go")
	panicIf(err)
	defer out.Close()

	// write header
	_, err = out.Write([]byte("package main\n\nconst (\nindexHTML = `"))
	panicIf(err)

	// copy content
	_, err = io.Copy(out, in)
	panicIf(err)

	// write footer
	_, err = out.Write([]byte("`\n)\n"))
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
