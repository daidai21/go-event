package bus

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	ch1 := make(chan Event)
	ch2 := make(chan Event)
	ch3 := make(chan Event)

	GetDefaultEventBus().Subscribe("topic1", ch1)
	GetDefaultEventBus().Subscribe("topic2", ch2)
	GetDefaultEventBus().Subscribe("topic2", ch3)

	fmt.Println(GetDefaultEventBus().Publish("topic1", "Hi topic 1"))
	fmt.Println(GetDefaultEventBus().Publish("topic2", "Welcome to topic 2"))

	for {
		select {
		case d := <-ch1:
			go func() {
				fmt.Println("ch1 ", d)
			}()
		case d := <-ch2:
			go func() {
				fmt.Println("ch2 ", d)
			}()
		case d := <-ch3:
			go func() {
				fmt.Println("ch3 ", d)
			}()
		case <-time.After(3 * time.Second):
			goto JumpLoop
		}
	}
JumpLoop:
	GetDefaultEventBus().CloseAll()
	fmt.Println(GetDefaultEventBus().Publish("topic1", "asda"))

}
