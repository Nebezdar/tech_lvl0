package messaging

import (
	"log"
)

type Nats struct {
	conn *nats.Conn
}

func NewNats(url string) (*Nats, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Nats{conn: conn}, nil
}

func (n *Nats) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	subscription, err := n.conn.Subscribe(subject, handler)
	if err != nil {
		log.Printf("Error subscribing to subject %s: %v", subject, err)
		return nil, err
	}

	return subscription, nil
}
