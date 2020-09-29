# Morfeusz

Go bindings for [Morfeusz 2](http://morfeusz.sgjp.pl/),
a morphological analyser for Polish.

## Installation

1. Install the C/C++ headers for Morfeusz, as described on
[the page about Morfeusz's programming tools](http://morfeusz.sgjp.pl/download/).

2. `go get github.com/go-morfeusz/morfeusz`

## Usage example

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-morfeusz/morfeusz"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	m, err := morfeusz.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	r := m.Analyse(b)
	for r.Next() {
		t := r.TokenInfo()
		fmt.Println(
			t.StartNode(), t.EndNode(), t.Orth(), t.Lemma(),
			t.Tag(m), t.Name(m), t.LabelsAsString(m))
	}
}
```

## Author

Marcin Ciura < mciura at gmail dot com >

## License

The Go package `morfeusz`, like Morfeusz 2 itself,
is licensed under [BSD 2-Clause License](LICENSE).
