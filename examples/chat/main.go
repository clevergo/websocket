// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"github.com/valyala/fasthttp"
	"log"
	"text/template"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTemplate = template.Must(template.ParseFiles("home.html"))

func handler(ctx *fasthttp.RequestCtx) {
	log.Println(ctx.URI())
	if bytes.Equal(ctx.Path(), []byte{'/'}) {
		if !ctx.IsGet() {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			return
		}
		ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
		homeTemplate.Execute(ctx, string(ctx.Host()))
		return
	} else if bytes.Equal(ctx.Path(), []byte("/ws")) {
		wsHandler(ctx)
		return
	}
	ctx.NotFound()
}

var wsHandler fasthttp.RequestHandler

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	wsHandler = func(ctx *fasthttp.RequestCtx) {
		serveWs(hub, ctx)
	}
	err := fasthttp.ListenAndServe(*addr, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
