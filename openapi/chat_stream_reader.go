package openapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StreamReader struct {
	isFinished bool

	reader   *bufio.Reader
	response *http.Response
	buffer   *bytes.Buffer
}

func (stream *StreamReader) check() (errResp *ErrorResponse) {
	if stream.buffer.Len() == 0 {
		return
	}
	err := json.Unmarshal(stream.buffer.Bytes(), &errResp)
	if err != nil {
		errResp = nil
	}
	return
}

func (stream *StreamReader) Recv() (response ChatCompletionStreamResponse, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

waitReader:
	line, err := stream.reader.ReadBytes('\n')
	if err != nil {
		resErr := stream.check()
		if resErr != nil {
			err = fmt.Errorf("err,%w", resErr.Error)
		}
		return
	}

	var headerData = []byte("data: ")
	line = bytes.TrimSpace(line)
	if !bytes.HasPrefix(line, headerData) {
		_, Werr := stream.buffer.Write(line)
		if Werr != nil {
			err = Werr
			return
		}
		goto waitReader
	}

	line = bytes.TrimPrefix(line, headerData)
	if string(line) == "[DONE]" {
		stream.isFinished = true
		err = io.EOF
		return
	}

	err = json.Unmarshal(line, &response)
	return
}

func (stream *StreamReader) Close() {
	stream.response.Body.Close()
}
