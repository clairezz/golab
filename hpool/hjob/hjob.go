package hjob

// hjob 由用户交给 worker

type HandlerFunc func()

type Hjob struct {
	Handler HandlerFunc
}
