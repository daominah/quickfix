package quickfix

import (
	"io"

	zlog "github.com/daominah/gomicrokit/log"
)

func writeLoop(connection io.Writer, messageOut chan []byte, log Log) {
	for {
		msg, ok := <-messageOut
		if !ok {
			return
		}

		_, err := connection.Write(msg)
		if err != nil {
			log.OnEvent(err.Error())
		}
		zlog.Debugf("conn sent: err %v, msg: %s", err, msg)
	}
}

func readLoop(parser *parser, msgIn chan fixIn) {
	defer close(msgIn)

	for {
		msg, err := parser.ReadMessage()
		if err != nil {
			return
		}
		msgIn <- fixIn{msg, parser.lastRead}
	}
}
