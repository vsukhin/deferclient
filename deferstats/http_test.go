package deferstats

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

type TestJSON struct {
	Title string
}

func TestHTTPPost(t *testing.T) {

	dps := NewClient("token")

	mux := http.NewServeMux()
	mux.HandleFunc("/", dps.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var tj TestJSON
		err = json.Unmarshal(body, &tj)
		if tj.Title != "sample title in json" {
			t.Error("not parsing the POST body correctly")
		}

	}))

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func TestHTTPPostHandler(t *testing.T) {

	dps := NewClient("token")

	mux := http.NewServeMux()

	posth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var tj TestJSON
		err = json.Unmarshal(body, &tj)
		if tj.Title != "sample title in json" {
			t.Error("not parsing the POST body correctly")
		}

	})

	mux.Handle("/", dps.HTTPHandler(posth))

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func TestHTTPHeader(t *testing.T) {

	dps := NewClient("token")

	mux := http.NewServeMux()
	mux.HandleFunc("/", dps.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "some header" {
			t.Error("headers not being passed through")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("headers not being passed through")
		}
	}))

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func TestHTTPHeaderHandler(t *testing.T) {

	dps := NewClient("token")

	mux := http.NewServeMux()

	posth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "some header" {
			t.Error("headers not being passed through")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("headers not being passed through")
		}
	})

	mux.Handle("/", dps.HTTPHandler(posth))

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func TestSOA(t *testing.T) {

	dps := NewClient("token")

	LatencyThreshold = -1

	mux := http.NewServeMux()
	mux.HandleFunc("/", dps.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		dpsi := r.Header.Get("X-dpparentspanid")

		okey := r.FormValue("other_key")

		if dpsi != "8103318854963911860" {
			t.Error("span not accessible")
		}

		if okey != "2" {
			t.Error("other_key not accessible")
		}

	}))

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	data := url.Values{}
	data.Set("other_key", "2")

	client := &http.Client{}
	r, _ := http.NewRequest("POST", lurl, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("X-dpparentspanid", "8103318854963911860")

	resp, err := client.Do(r)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if len(curlist.list) == 0 {
		t.Error("should have a http in the list")
	}

	if curlist.list[0].ParentSpanId != 8103318854963911860 {
		t.Error("not tracking our parent_span_id")
	}

}

func TestSOAHandler(t *testing.T) {

	dps := NewClient("token")

	LatencyThreshold = -1

	mux := http.NewServeMux()

	posth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		dpsi := r.Header.Get("X-dpparentspanid")
		okey := r.FormValue("other_key")

		if dpsi != "8103318854963911860" {
			t.Error("span not accessible")
		}

		if okey != "2" {
			t.Error("other_key not accessible")
		}

	})

	mux.Handle("/", dps.HTTPHandler(posth))

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	data := url.Values{}
	data.Set("other_key", "2")

	client := &http.Client{}
	r, _ := http.NewRequest("POST", lurl, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("X-dpparentspanid", "8103318854963911860")

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if len(curlist.list) == 0 {
		t.Error("should have a http in the list")
	}

	if curlist.list[0].ParentSpanId != 8103318854963911860 {
		t.Error("not tracking our parent_span_id")
	}

}
