package simpleamqp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageToOuputWithNilMessagesChannelDoesNotBlockForever(t *testing.T) {
	for _, queueTimeout := range []time.Duration{0, 30 * time.Second} {
		closed := make(chan bool, 1)
		go func() {
			closed <- messageToOuput(nil, make(chan AmqpMessage), queueTimeout)
		}()

		select {
		case result := <-closed:
			assert.True(t, result, "a nil messages channel must be reported as closed so the consumer reconnects")
		case <-time.After(2 * time.Second):
			t.Fatalf("messageToOuput blocked forever with a nil messages channel (queueTimeout %v)", queueTimeout)
		}
	}
}
