package main

import (
	"fmt"
	"os"

	"github.com/prybintsev/memepool/mempool"
)

func main() {
	m := mempool.NewMemPool()

	input, err := os.Open("transactions.txt")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	err = m.ReadTransactions(input)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	output, err := os.Create("prioritized-transactions.txt")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer func() {
		err = output.Close()
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}()

	err = m.WriteTransactions(output)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
