package quickfix

import (
	"io"
	"strings"

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
		zlog.Debugf("writeLoop: err %v, msg: %s", err, msg)
	}
}

func readLoop(parser *parser, msgIn chan fixIn) {
	defer close(msgIn)

	for {
		msg, err := parser.ReadMessage()
		if err == nil {
			msgIn <- fixIn{msg, parser.lastRead}
		}
		if err != nil {
			if strings.Contains(err.Error(), errFromIOReader) {
				zlog.Infof("readLoop returned, err: %v", err)
				return
			}
		}
	}
}
