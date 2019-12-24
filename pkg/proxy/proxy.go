package proxy

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
)

// Proxy is a DoH to DNS proxy
type Proxy struct {
	conn      *net.UDPConn
	dohAddr   string
	allowList []*regexp.Regexp
}

// NewProxyForConn gives a new Proxy for a given socket
func NewProxyForConn(dohAddr string, conn *net.UDPConn) Proxy {
	return Proxy{
		conn:    conn,
		dohAddr: dohAddr,
	}
}

func (p *Proxy) AllowRegexMatch(in string) error {
	xp, err := regexp.Compile(in)
	if err != nil {
		return err
	}

	p.allowList = append(p.allowList, xp)
	return nil
}

func (p *Proxy) ListenAndServe() {
	for {
		var raw [512]byte
		n, addr, err := p.conn.ReadFromUDP(raw[:512])
		if err != nil {
			log.Printf("could not read: %s", err)
			continue
		}
		go p.proxyRequest(addr, raw[:n])
	}
}

func (p *Proxy) proxyRequest(addr *net.UDPAddr, raw []byte) {

	header, reqs, err := parseDNSRequest(raw)

	if err != nil {
		log.Printf("could not read request: %s", err)
		return
	}

	allowedReqs := []dnsResourceRecord{}
	for _, rr := range reqs {
		ok := false
		for _, xp := range p.allowList {
			if xp.MatchString(rr.DomainName) {
				allowedReqs = append(allowedReqs, rr)
				ok = true
				break
			}
		}
		if !ok {
			log.Println(rr.DomainName)
		}
	}

	if len(allowedReqs) != len(reqs) {
		// TODO: add multi domain behaviour
		p.conn.WriteToUDP(writeNoFound(header, reqs), addr)
		return
	}

	enc := base64.RawURLEncoding.EncodeToString(raw)
	url := fmt.Sprintf("%s?dns=%s", p.dohAddr, enc)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("could not create request: %s", err)
		return
	}
	r.Header.Set("Content-Type", "application/dns-message")
	r.Header.Set("Accept", "application/dns-message")

	c := http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		log.Printf("could not perform request: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("wrong response from DOH server got %s", http.StatusText(resp.StatusCode))
		return
	}

	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("could not read message from response: %s", err)
		return
	}

	if _, err := p.conn.WriteToUDP(msg, addr); err != nil {
		log.Printf("could not write to udp connection: %s", err)
		return
	}
}
