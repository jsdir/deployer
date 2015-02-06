package deployer

type Broadcaster struct {
	listeners map[int]chan interface{}
	incr      int
}

func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		listeners: make(map[int]chan interface{}),
	}
	return b
}

func (b *Broadcaster) Write(message interface{}) {
	for _, channel := range b.listeners {
		channel <- message
	}
}

func (b *Broadcaster) Listen() (int, chan interface{}) {
	channel := make(chan interface{})
	b.incr++
	b.listeners[b.incr] = channel
	return b.incr, channel
}

func (b *Broadcaster) Unregister(index int) {
	delete(b.listeners, index)
}
