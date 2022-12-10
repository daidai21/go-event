package bus

import "sync"

type (
	Event interface{}

	channel chan Event

	channelSlice []channel
)

type EventBus struct {
	subscribers map[string]channelSlice
	lock        sync.Mutex
	isClosed    bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string]channelSlice),
		lock:        sync.Mutex{},
		isClosed:    false,
	}
}

var (
	once sync.Once
	inst *EventBus
)

func GetDefaultEventBus() *EventBus {
	once.Do(func() {
		inst = NewEventBus()
	})
	return inst
}

func (eb *EventBus) Subscribe(topic string, ch channel) {
	defer eb.lock.Unlock()
	eb.lock.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]channel{}, ch)
	}
}

func (eb *EventBus) Publish(topic string, data interface{}) bool {
	defer eb.lock.Unlock()
	eb.lock.Lock()
	if eb.isClosed {
		return false
	}
	if chans, found := eb.subscribers[topic]; found {
		// 这样做是因为切片引用相同的数组，即使它们是按值传递的
		// 因此我们正在使用我们的元素创建一个新切片，从而能正确地保持锁定
		channels := append(channelSlice{}, chans...)
		go func(data Event, dataChannelSlices channelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(data, channels)
	}
	return true
}

func (eb *EventBus) CloseAll() {
	defer eb.lock.Unlock()
	eb.lock.Lock()
	eb.isClosed = true
	for _, chans := range eb.subscribers {
		for _, ch := range chans {
			close(ch)
		}
	}
}
