package watermill

import (
	"reflect"

	"strings"
	"unicode"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/golang/protobuf/proto"
	"github.com/satori/uuid"
)

type MuPublisher struct {
	publisher message.Publisher
	topicName func(interface{}) string
}

func (p *MuPublisher) Publish(msgs ...proto.Message) error {
	for _, msg := range msgs {
		topic := p.topicName(msg)

		payload, err := proto.Marshal(msg)
		if err != nil {
			return err
		}

		watermillMessage := message.NewMessage(uuid.NewV4().String(), payload)
		if err := p.publisher.Publish(topic, watermillMessage); err != nil {
			return err
		}
	}

	return nil
}

func DefaultTopicName(m interface{}) string {
	t := reflect.TypeOf(m)

	name := ""
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	var sb strings.Builder

	for i, c := range name {
		if unicode.IsUpper(c) {
			if i != 0 {
				sb.WriteRune('.')
			}

			sb.WriteRune(unicode.ToLower(c))
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}
