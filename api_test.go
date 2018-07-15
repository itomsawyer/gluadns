package gluadns

import (
	"testing"

	glua "github.com/yuin/gopher-lua"
)

func TestDnsExchange(t *testing.T) {

	l := glua.NewState(glua.Options{IncludeGoStackTrace: true})
	l.PreloadModule("dns", Loader)

	err := l.DoString(`
		local dns = require("dns");

		local r, err = dns.exchange("kss.ksyun.com.", "A", "ns1.kscdns.com", "111.11.11.1")
		if not (err == nil) then
			error(err)
		end
		for k,v in pairs(r["answer"]) do
			print(v["value"])
		end

		local r, err = dns.exchange("kss.ksyun.com.", "A", "ns1.kscdns.com", "123.127.0.0")
		if not (err == nil) then
			error(err)
		end
		for k,v in pairs(r["answer"]) do
			print(v["value"])
		end

		local r, err = dns.exchange("kss.ksyun.com.", "A", "ns1.kscdns.com", "1.92.0.0")
		if not (err == nil) then
			error(err)
		end
		for k,v in pairs(r["answer"]) do
			print(v["value"])
		end
	`)

	if err != nil {
		t.Fatal(err)
	}
}
