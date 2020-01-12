package hw05

import (
	"errors"
	"fmt"
	"math"
)

// Run N concurrency tasks
func Run(tasks []func() error, N int, M int) error {
	err := checkInputs(len(tasks), N, M)
	if err != nil {
		return err
	}

	minParallelTaskCount := int(math.Min(float64(N), float64(len(tasks))))

	results := make(chan error, minParallelTaskCount)

	index := 0

	for ; index < minParallelTaskCount; index++ {
		runOne(tasks[index], results)
	}

	errCounter := 0
	readyCounter := 0

	for readyCounter < len(tasks) {
		err := <-results
		readyCounter++
		if err != nil {
			errCounter++
			if errCounter == M {
				//wait remaining tasks
				for index != readyCounter {
					<-results
					readyCounter++
				}

				return fmt.Errorf("%d errors reterned", M)
			}
		}
		if index < len(tasks) {
			runOne(tasks[index], results)
			index++
		}
	}

	return nil
}

// check input "Run" function
func checkInputs(Len, N, M int) error {
	if Len == 0 {
		return errors.New("Empty task list")
	}

	if N <= 0 {
		return errors.New("N must be positive")
	}

	if M <= 0 {
		return errors.New("M must be positive")
	}

	if M > Len {
		return errors.New("M over than length of task")
	}

	return nil
}

//run one task
func runOne(task func() error, results chan<- error) {
	go func() {
		results <- task()
	}()
}
