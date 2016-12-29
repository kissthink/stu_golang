package que

import (
	"sync"
	"time"
	"fmt"
	"testing"
)

func TestQue(t *testing.T) {
	var wg sync.WaitGroup
	var lr LimitRate
	lr.SetRate(3)

	b:=time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			if lr.Limit() {
				fmt.Println("Got it!", time.Now().String())
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(b))
}