package messaging

type Publisher interface {
	Publish(message interface{}) error
}

type Subscriber interface {
	Subscribe(function interface{}) error
}
