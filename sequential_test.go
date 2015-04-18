// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package main

// Omega: Alt+937

import (
	"io/ioutil"
	"net"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func serverBasics(server *ghttp.Server) {
}

var _ = Describe("Testing sequential requests", func() {

	var server *ghttp.Server
	var cliStop, srvStop chan struct{}
	var wstunUrl string
	var wstunToken = "test567890123456"

	BeforeEach(func() {
		server = ghttp.NewServer()
		serverBasics(server)

		l, _ := net.Listen("tcp", "127.0.0.1:0")
		srvStop = wstunsrv([]string{}, l)
		cliStop = wstuncli([]string{
			"-token", wstunToken,
			"-tunnel", "ws://" + l.Addr().String(),
			"-server", server.URL(),
		})
		wstunUrl = "http://" + l.Addr().String()
	})
	AfterEach(func() {
		cliStop <- struct{}{}
		srvStop <- struct{}{}
		server.Close()
	})

	// Perform the test by running main() with the command line args set
	It("Responds to hello requests", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/hello"),
				ghttp.RespondWith(200, `WORLD`, http.Header{"Content-Type": []string{"text/world"}}),
			),
		)

		resp, err := http.Get(wstunUrl + "/_token/" + wstunToken + "/hello")
		Ω(err).ShouldNot(HaveOccurred())
		respBody, err := ioutil.ReadAll(resp.Body)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(string(respBody)).Should(Equal("WORLD"))
		Ω(resp.Header.Get("Content-Type")).Should(Equal("text/world"))
	})

})
