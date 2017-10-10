package scan

import "fmt"

type Scan interface {
	Do()
}

func Do(s Scan) {
	fmt.Println(s)
	s.Do()
}
