package gin

type Proxy struct {
}

func NewProxy(builder Builder, runner Runner) *Proxy {
	return &Proxy{}
}

func (p *Proxy) Run(config *Config) error {
	return nil
}
