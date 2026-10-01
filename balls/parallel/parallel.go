package parallel

import (
	"runtime"
	"sync"
)

// ForEach executes the function f for each index from 0 to n-1 in parallel.
func ForEach(n int, f func(i int)) {
	if n <= 0 {
		return
	}

	numWorkers := runtime.NumCPU()
	if n < numWorkers {
		numWorkers = n
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	chunkSize := (n + numWorkers - 1) / numWorkers

	for w := 0; w < numWorkers; w++ {
		go func(w int) {
			defer wg.Done()
			start := w * chunkSize
			end := start + chunkSize
			if end > n {
				end = n
			}

			for i := start; i < end; i++ {
				f(i)
			}
		}(w)
	}

	wg.Wait()
}
