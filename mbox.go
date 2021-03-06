// The MIT License (MIT)
//
// Copyright (c) 2016 aerth
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// Package mbox saves a form to a local .mbox file (opengpg option)
/*

Usage of mbox library is as follows:

Define mbox.Destination variable in your program

Accept an email, populate the mbox.Form struct like this:
	mbox.From = "joe"
	mbox.Email = "joe@blowtorches.info
	mbox.Message = "hello world"
	mbox.Subject = "re: hello joe"
	mbox.Save()


*/
package mbox

import (

	//	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	// email validation
	"github.com/microcosm-cc/bluemonday" // input sanitizaation
)

// Form is a single email. No Attachments yet.
type Form struct {
	From, Subject, Message string
	Sent, Received                      time.Time
	Body []byte
}

var (
	// ValidationLevel should be set to something other than 1 to resolve hostnames and validate emails
	ValidationLevel = 1
	// Destination is the address where mail is "sent", its useful to change this to the address you will be replying to.
	Destination = "mbox@localhost"

	// Mail is the local mbox, implemented as a logger
	Mail *os.File
)

// ParseQuery returns a mbox.Form from url.Values
func ParseQuery(query url.Values) *Form {
	p := bluemonday.StrictPolicy()
	form := new(Form)
	additionalFields := ""
	for k, v := range query {
		k = strings.ToLower(k)
		if k == "email" || k == "name" {
			form.From = v[0]
			form.From = p.Sanitize(form.From)
		} else if k == "subject" {
			form.Subject = v[0]
			form.Subject = p.Sanitize(form.Subject)
		} else if k == "message" {
			form.Message = k + ": " + v[0] + "<br>\n"
			form.Message = p.Sanitize(form.Message)
		} else if k != "cosgo" && k != "captchaid" && k != "captchasolution" {
			additionalFields = additionalFields + k + ": " + v[0] + "<br>\n"
		}
	}
	if form.Subject == "" || form.Subject == " " {
		form.Subject = "[New Message]"
	}
	if additionalFields != "" {
		if form.Message == "" {
			form.Message = form.Message + "Message:\n<br>" + p.Sanitize(additionalFields)
		} else {
			form.Message = form.Message + "\n<br>Additional:\n<br>" + p.Sanitize(additionalFields)
		}
	}

	return form
}

// rel2real Relative to Real path name
func rel2real(file string) (realpath string) {
	pathdir, _ := path.Split(file)
	if pathdir == "" {
		realpath, _ = filepath.Abs(file)
	} else {
		realpath = file
	}
	return realpath
}
