package gin

import (
	"fmt"
	"net"
	"net/http"
)

type Proxy struct {
	listener net.Listener
}

func NewProxy(builder Builder, runner Runner) *Proxy {
	return &Proxy{}
}

func (p *Proxy) Run(config *Config) error {
	var err error
	p.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}

	go http.Serve(p.listener, http.HandlerFunc(p.defaultHandler))
	return nil
}

func (p *Proxy) Close() error {
	return p.listener.Close()
}

func (p *Proxy) defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Serving")
}
