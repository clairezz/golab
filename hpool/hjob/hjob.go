package hjob

// hjob 由用户交给 worker

type HandlerFunc func([]byte, chan string)

type Hjob struct {
	Handler HandlerFunc
	Data []byte
	RespCh chan string
}
