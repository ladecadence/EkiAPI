package msgbus

type MsgBus struct {
	clients []chan string
}

func New() *MsgBus {
	return &MsgBus{
		clients: []chan string{},
	}
}

func (mb *MsgBus) Subscribe() <-chan string {
	ch := make(chan string)
	mb.clients = append(mb.clients, ch)
	return ch
}

func (mb *MsgBus) Publish(msg string) {
	for _, ch := range mb.clients {
		go func(c chan string) {
			select {
			case c <- msg:
				// nothing
			default:
				// nothing
			}
		}(ch)
	}
}
