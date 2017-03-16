package pool

import "net"

type JobPair struct {
	name string
	job  *Job
}

type Job struct {
	conn     net.Conn
	OpCode   int
	Data     []byte
	RespChan chan ResultType
}

type ResultType struct {
	code int
	msg  string
	date []byte
}

func (r *ResultType) GetMsg() string {
	return r.msg
}
type HandlerFunc func(data []byte, resCh chan ResultType)

var code2op = map[int]HandlerFunc{
	0: Hello,
}

func (j *Job) Do() {
	handler := code2op[j.OpCode]
	handler(j.Data, j.RespChan)
}
