package ovoclient

import (
	"strconv"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/maxzerbini/ovoclient/model"
)
// global log flag
var LogEnabled  bool // Log request and response

// Http Session
type Session struct {
	Client *http.Client
	// Optional defaults - can be overridden in a Request
	Header *http.Header
	Params *url.Values
	// Ovo Node 
	node *model.OvoTopologyNode
	port string
}

// create a new Session
func NewSession() *Session {
	return &Session{Client:&http.Client{}}
}

func (s *Session) SetNode(node *model.OvoTopologyNode){
	s.node = node
	s.port = strconv.Itoa(node.Port)
}

// Send constructs and sends an HTTP request.
func (s *Session) Send(r *Request) (response *Response, err error) {
	r.Method = strings.ToUpper(r.Method)
	//
	// Create a URL object from the raw url string.  This will allow us to compose
	// query parameters programmatically and be guaranteed of a well-formed URL.
	//
	u, err := url.Parse(r.Url)
	if err != nil {
		logInfo("URL", r.Url)
		logInfo(err)
		return
	}
	//
	// Default query parameters
	//
	p := url.Values{}
	if s.Params != nil {
		for k, v := range *s.Params {
			p[k] = v
		}
	}
	//
	// Parameters that were present in URL
	//
	if u.Query() != nil {
		for k, v := range u.Query() {
			p[k] = v
		}
	}
	//
	// User-supplied params override default
	//
	if r.Params != nil {
		for k, v := range *r.Params {
			p[k] = v
		}
	}
	//
	// Encode parameters
	//
	u.RawQuery = p.Encode()
	//
	// Attach params to response
	//
	r.Params = &p
	//
	// Create a Request object; if populated, Data field is JSON encoded as
	// request body
	//
	header := http.Header{}
	if s.Header != nil {
		for k, _ := range *s.Header {
			v := s.Header.Get(k)
			header.Set(k, v)
		}
	}
	var req *http.Request
	var buf *bytes.Buffer
	if r.Payload != nil {
		if r.RawPayload {
			var ok bool
			// buf can be nil interface at this point
			// so we'll do extra nil check
			buf, ok = r.Payload.(*bytes.Buffer)
			if !ok {
				err = errors.New("Payload must be of type *bytes.Buffer if RawPayload is set to true")
				return
			}
		} else {
			var b []byte
			b, err = json.Marshal(&r.Payload)
			if err != nil {
				logInfo(err)
				return
			}
			buf = bytes.NewBuffer(b)
		}
		if buf != nil {
			req, err = http.NewRequest(r.Method, u.String(), buf)
		} else {
			req, err = http.NewRequest(r.Method, u.String(), nil)
		}
		if err != nil {
			logInfo(err)
			return
		}
		// Overwrite the content type to json since we're pushing the payload as json
		header.Set("Content-Type", "application/json")
	} else { // no data to encode
		req, err = http.NewRequest(r.Method, u.String(), nil)
		if err != nil {
			logInfo(err)
			return
		}

	}
	//
	// Merge Session and Request options
	//
	if r.Header != nil {
		for k, v := range *r.Header {
			header.Set(k, v[0]) // Is there always guarnateed to be at least one value for a header?
		}
	}
	if header.Get("Accept") == "" {
		header.Add("Accept", "application/json") // Default, can be overridden with Opts
	}
	req.Header = header
	r.timestamp = time.Now()
	var client *http.Client
	if s.Client != nil {
		client = s.Client
	} else {
		client = &http.Client{}
		s.Client = client
	}
	resp, err := client.Do(req)
	if err != nil {
		logInfo(err)
		return
	}
	defer resp.Body.Close()
	r.status = resp.StatusCode
	r.response = resp

	//
	// Unmarshal
	//
	r.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logInfo(err)
		return
	}
	if string(r.body) != "" {
		if resp.StatusCode < 300 && r.Result != nil {
			err = json.Unmarshal(r.body, r.Result)
		}
		if resp.StatusCode >= 400 && r.Error != nil {
			json.Unmarshal(r.body, r.Error) // Should we ignore unmarshall error?
		}
	}
	if r.CaptureResponseBody {
		r.ResponseBody = bytes.NewBuffer(r.body)
	}
	rsp := Response(*r)
	response = &rsp
	return
}

// Get sends a GET request.
func (s *Session) Get(url string, p *url.Values, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method: "GET",
		Url:    url,
		Params: p,
		Result: result,
		Error:  errMsg,
	}
	return s.Send(&r)
}

// Post sends a POST request.
func (s *Session) Post(url string, payload, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method:  "POST",
		Url:     url,
		Payload: payload,
		Result:  result,
		Error:   errMsg,
	}
	return s.Send(&r)
}

// Put sends a PUT request.
func (s *Session) Put(url string, payload, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method:  "PUT",
		Url:     url,
		Payload: payload,
		Result:  result,
		Error:   errMsg,
	}
	return s.Send(&r)
}

// Delete sends a DELETE request.
func (s *Session) Delete(url string, p *url.Values, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method: "DELETE",
		Url:    url,
		Params: p,
		Result: result,
		Error:  errMsg,
	}
	return s.Send(&r)
}

// Debug method for logging
// Centralizing logging in one method
// avoids spreading conditionals everywhere
func logInfo(args ...interface{}) {
	if LogEnabled {
		log.Println(args...)
	}
}

// Debug method for logging
// Centralizing logging in one method
// avoids spreading conditionals everywhere
func logInfof(message string, args ...interface{}) {
	if LogEnabled {
		log.Printf(message, args...)
	}
}