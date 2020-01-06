package main

import (
	"net"

	"github.com/meyskens/adns/pkg/proxy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

type serveCmdOptions struct {
	BindAddr string
	Port     int
	ProxyTo  string
}

// NewServeCmd generates the `serve` command
func NewServeCmd() *cobra.Command {
	s := serveCmdOptions{}
	c := &cobra.Command{
		Use:   "serve",
		Short: "Serves the DoH proxy",
		Long:  `Serves the DoH proxy on the given bind address and port`,
		RunE:  s.RunE,
	}
	c.Flags().StringVarP(&s.BindAddr, "bind-address", "b", "0.0.0.0", "address to bind port to")
	c.Flags().IntVarP(&s.Port, "port", "p", 53, "Port to listen on")
	c.Flags().StringVarP(&s.ProxyTo, "proxy-to", "t", "https://1.1.1.1/dns-query", "DoH endpoint to proxy requests to")

	viper.BindPFlags(c.Flags())

	return c
}

func (s *serveCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(s.BindAddr), Port: s.Port})
	if err != nil {
		return err
	}

	p := proxy.NewProxyForConn(s.ProxyTo, conn)

	// hard coded on purpose!
	p.AllowRegexMatch(`.*\.licenses\.adobe.com$`)
	p.AllowRegexMatch(`^api-cna01\.adobe-services\.com$`)
	p.AllowRegexMatch(`^supportanyware\.adobe\.io$`)
	p.AllowRegexMatch(`^www\.adobe\.com`)
	p.AllowRegexMatch(`^lcs-cops\.adobe\.io$`)
	p.AllowRegexMatch(`^genuine\.adobe\.com$`)
	p.AllowRegexMatch(`^prod\.adobegenuine\.com$`)
	p.AllowRegexMatch(`^gocart-web-prod-.*\.elb.amazonaws\.com$`)
	p.AllowRegexMatch(`^na1e\.services\.adobe\.com$`)
	p.AllowRegexMatch(`^auth\.services\.adobe\.com$`)
	p.AllowRegexMatch(`^auth-api\.services\.adobe\.com$`)
	p.AllowRegexMatch(`.*\.adobelogin\.com$`)
	p.AllowRegexMatch(`^adobeid-na1\.services\.adobe\.com$`)
	p.AllowRegexMatch(`^na1e-acc\.services\.adobe\.com$`)
	p.AllowRegexMatch(`^na1r\.services\.adobe\.com$`)
	p.AllowRegexMatch(`^ams\.adobe\.com$`)
	p.AllowRegexMatch(`^oobe\.adobe\.com$`)
	p.AllowRegexMatch(`^federatedid-na1\.services\.adobe\.com$`)
	p.AllowRegexMatch(`^adobelogin\.prod\.ims\.adobejanus.com$`)
	p.AllowRegexMatch(`^services\.prod\.ims\.adobejanus\.com$`)
	p.AllowRegexMatch(`^www-prod\.adobesunbreak\.com$`)
	p.AllowRegexMatch(`^.*\.okta\.com$`)
	p.AllowRegexMatch(`^.*\.oktapreview\.com$`)
	p.AllowRegexMatch(`^.*\.adobess\.com$`)
	p.AllowRegexMatch(`.*\.digicert.com`) // OCSP and CRL

	p.ListenAndServe()

	return nil
}
