package utils

func CreateAll [K any] (n int) []chan K {
	channels := make([]chan K, n)
	for i := range channels {
		channels[i] = make(chan K)
	}
	return channels
}

func CloseAll [K any] (channels ...chan K) {
	for _, output := range channels {
		close(output)
	}
}

func BroadCast[K any] (quit <-chan int, input <-chan K, n int) []chan K {
	outputs := CreateAll[K](n)
	go func() {
		defer CloseAll(outputs...)
		var s K
		ok := true
		for ok {
			select {
			case s, ok = <-input:
				if ok {
					for _, output := range outputs {
						output <- s
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return outputs
}