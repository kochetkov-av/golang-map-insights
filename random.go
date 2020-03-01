package main

import "fmt"

func main() {
	keys := []string{"Cleese", "Gilliam", "Idle", "Palin", "Chapman", "Jones",}
	values := []string{"John", "Terry", "Eric", "Michael", "Graham", "Terry",}

	for i := 0; i < 10; i++ {
		m := make(map[string]string)
		for j := 0; j < 6; j++ {
			m[keys[j]] = values[j]
		}

		for k, v := range m {
			fmt.Printf("%10s %10s", k, v)
		}
		fmt.Println()
	}

	//for i := 0; i < 10; i++ {
	m := make(map[int]int)
	for j := 0; j < 10; j++ {
		m[j] = j * 2
	}

	for k, v := range m {
		fmt.Printf("%d*2=%d ", k, v)
	}

	fmt.Printf("\n - - -\n")

	for k, v := range m {
		fmt.Printf("%d*2=%d ", k, v)
	}
	//}
}
