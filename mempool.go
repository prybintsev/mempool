package memepool

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const memPoolCapacity = 5000

type MemPool struct {
	queue *PriorityQueue
}

func NewMemPool() MemPool {
	q := NewPriorityQueue(memPoolCapacity)
	return MemPool{queue: &q}
}

func (m MemPool) ReadTransactions(reader io.Reader) error {
	bReader := bufio.NewReader(reader)
	var isEof bool
	for !isEof {
		line, err := bReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				isEof = true
				err = nil
			} else {
				return err
			}
		}
		// Skip whitespace lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		tx, err := ReadTransaction(line)
		if err != nil {
			return err
		}

		m.queue.Push(tx)
	}

	return nil
}

func (m MemPool) WriteTransactions(writer io.Writer) error {
	if m.queue.Len() == 0 {
		return nil
	}

	for  {
		item := m.queue.Pop()
		_, err := writer.Write([]byte(fmt.Sprintf("%s", item)))
		if err != nil {
			return err
		}
		if m.queue.Len() == 0 {
			break
		}
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
