package internal

import (
	"time"
)

type BrokerClient struct {
	bro *BrokerImpl
}

func NewClient() *BrokerClient {
	return &BrokerClient{
		bro: NewBroker(),
	}
}

func (c *BrokerClient) SetConditions(capacity int) {
	c.bro.setConditions(capacity)
}

func (c *BrokerClient) Publish(topic string, msg interface{}) error {
	return c.bro.publish(topic, msg)
}

func (c *BrokerClient) Subscribe(topic string) (<-chan interface{}, error) {
	return c.bro.subscribe(topic)
}

func (c *BrokerClient) Unsubscribe(topic string, sub <-chan interface{}) error {
	return c.bro.unsubscribe(topic, sub)
}

func (c *BrokerClient) Close() {
	c.bro.close()
}

func (c *BrokerClient) GetPayLoad(sub <-chan interface{}) interface{} {
	for val := range sub {
		if val != nil {
			return val
		}
	}
	return nil
}

// 毫秒为单位的
func (c *BrokerClient) GetPayLoadAsync(sub <-chan interface{}, timeout uint32) interface{} {
	_timeout := time.NewTimer(time.Duration(timeout) * time.Millisecond)
	select {
	case val := <-sub:
		return val
	case <-_timeout.C:
		return nil
	}

}
