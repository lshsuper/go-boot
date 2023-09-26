package mq

type Adapter interface {
}

type MqType string

const (
	Pulsar MqType = "pulsar"
	Rabbit MqType = "rabbit"
)

func (e MqType) String() string {
	return string(e)
}
