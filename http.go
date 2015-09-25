package stats

import (
	"net/http"
	"time"
)

// LogResponseWritter wraps the standard http.ResponseWritter allowing for more
// verbose logging
type LogResponseWritter struct {
	status int
	size   int
	http.ResponseWriter
}

// Status provides an easy way to retrieve the status code
func (w *LogResponseWritter) Status() int {
	return w.status
}

// Size provides an easy way to retrieve the response size in bytes
func (w *LogResponseWritter) Size() int {
	return w.size
}

// Header returns & satisfies the http.ResponseWriter interface
func (w *LogResponseWritter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write satisfies the http.ResponseWriter interface and
// captures data written, in bytes
func (w *LogResponseWritter) Write(data []byte) (int, error) {

	written, err := w.ResponseWriter.Write(data)
	w.size += written

	return written, err
}

// WriteHeader satisfies the http.ResponseWriter interface and
// allows us to cach the status code
func (w *LogResponseWritter) WriteHeader(statusCode int) {

	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// HTTPRequest contains information about the life of an http request
type HTTPRequest struct {
	URL                   string        `json:"url"`
	Method                string        `json:"method"`
	RequestContentLength  int64         `json:"reqContent"`
	Headers               http.Header   `json:"headers"`
	Start                 time.Time     `json:"start"`
	End                   time.Time     `json:"end"`
	Duration              time.Duration `json:"duration"`
	ResponseContentLength int64         `json:"resContent"`
	StatusCode            int           `json:"status"`
	HasErrors             bool          `json:"hasErrs"`
	Error                 string        `json:"err"`
	writer                *LogResponseWritter
	clientStats           *ClientStats
}

// NewHTTPRequest creates a new HTTPRequest for monitoring which wraps the ResponseWriter in order
// to collect stats so you need to call the Writer() function from the HTTPRequest created by this call
func (s *ClientStats) NewHTTPRequest(w http.ResponseWriter, r *http.Request) *HTTPRequest {

	return &HTTPRequest{
		Start:                time.Now().UTC(),
		URL:                  r.URL.String(),
		Method:               r.Method,
		RequestContentLength: r.ContentLength,
		Headers:              r.Header,
		writer:               &LogResponseWritter{status: 200, ResponseWriter: w},
		clientStats:          s,
	}
}

// Writer returns a wrapped http.ResponseWriter for logging purposes
func (r *HTTPRequest) Writer() http.ResponseWriter {
	return r.writer
}

// Failure records an HTTP failure and automatically completes the request
func (r *HTTPRequest) Failure(err string) {
	r.HasErrors = true
	r.Error = err
	r.Complete()
}

// Complete finalizes an HTTPRequest and logs it.
func (r *HTTPRequest) Complete() {

	r.End = time.Now().UTC()
	r.Duration = r.End.Sub(r.Start)
	r.ResponseContentLength = int64(r.writer.Size())
	r.StatusCode = r.writer.Status()
	r.clientStats.httpStats.Add(r)
}
