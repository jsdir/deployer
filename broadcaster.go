package deployer

type Broadcaster struct {
	listeners []chan interface{}
}

func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		listeners: []chan interface{}{},
	}
	return b
}

func (b *Broadcaster) Write(message interface{}) {
	for _, channel := range b.listeners {
		channel <- message
	}
}

func (b *Broadcaster) Listen() chan interface{} {
	channel := make(chan interface{})
	b.listeners = append(b.listeners, channel)
	return channel
}
