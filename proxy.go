package gin

type Proxy struct {
}

func NewProxy(builder Builder, runner Runner) *Proxy {
	return &Proxy{}
}
