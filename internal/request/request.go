package request

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"httpproto/utils"
)

type (
	Request struct {
		RequestLine RequestLine
		status      int
		buf         []byte
		read        int
		parsed      int
	}

	RequestLine struct {
		HttpVersion   string
		RequestTarget string
		Method        string
	}
)

const (
	RequestDelimiter  = "\r\n"
	SupportedVersion  = "1.1"
	PARSE_IN_PROGRESS = 0
	PARSE_DONE        = 1
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	var r Request
	for {
		b := make([]byte, 8)
		rd, err := reader.Read(b)
		if rd == 0 || err == io.EOF {
			break
		}
		_, err = r.parse(b[:rd])
		if err != nil {
			return nil, err
		}
		if r.done() {
			break
		}
	}

	return &r, nil
}

func parseRequestLine(data string) (*RequestLine, int, error) {
	var rl RequestLine
	pieces := strings.Split(data, RequestDelimiter)
	if len(pieces) == 0 {
		return nil, 0, nil
	}
	rlString := strings.Trim(pieces[0], " ")
	rlPieces := strings.Split(rlString, " ")
	if len(rlPieces) != 3 {
		return nil, 0, errors.New("invalid request line: missing parts")
	}
	method, target, versionLine := rlPieces[0], rlPieces[1], rlPieces[2]
	if !isValidMethod(method) {
		return nil, 0, errors.New("invalid method")
	}
	version, err := getVersion(versionLine)
	if err != nil {
		return nil, 0, err
	}
	if !isValidVersion(version) {
		return nil, 0, errors.New("invalid version")
	}
	rl.HttpVersion = version
	rl.Method = method
	rl.RequestTarget = target

	return &rl, len([]byte(data)), nil
}

func isValidMethod(method string) bool {
	if !utils.IsUppercase(method) {
		return false
	}

	return true
}

func getVersion(version string) (string, error) {
	vPieces := strings.Split(version, "/")
	if len(vPieces) != 2 {
		return "", errors.New("version should be divided by /")
	}

	return vPieces[1], nil
}

func isValidVersion(version string) bool {
	return version == SupportedVersion
}

func (r *Request) done() bool {
	return r.status == PARSE_DONE
}

func (r *Request) parse(data []byte) (int, error) {
	if r.buf == nil {
		r.buf = make([]byte, 0, 8)
	}
	r.buf = append(r.buf, data...)
	r.read += len(data)
	idx := bytes.Index(r.buf, []byte(RequestDelimiter))
	if idx != -1 {
		idx += len([]byte(RequestDelimiter))
		rl, p, err := parseRequestLine(string(r.buf[:idx]))
		if err != nil {
			return 0, err
		}
		defer func() {
			r.buf = nil
		}()
		r.parsed += p
		r.RequestLine = *rl
		r.status = PARSE_DONE

		return p, nil
	}

	return 0, nil
}
