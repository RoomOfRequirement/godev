package requests

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// code from: https://github.com/syavorsky/reqstrategy

// Validator func type
type Validator = func(resp *http.Response) error

type key string

const validatorKey key = "validators"

// WithValidator introduces a response validator function to the context to be used in Do/Race/All/Some/Retry.
func WithValidator(req *http.Request, validator Validator) *http.Request {
	ctx := req.Context()
	validators, _ := ctx.Value(validatorKey).([]Validator)
	validators = append(validators, validator)
	ctx = context.WithValue(ctx, validatorKey, validators)
	return req.WithContext(ctx)
}

// WithStatusRequired adds the response validator by listing acceptable status codes
func WithStatusRequired(r *http.Request, codes ...int) *http.Request {
	return WithValidator(r, func(r *http.Response) error {
		for _, code := range codes {
			if r.StatusCode == code {
				return nil
			}
		}
		return fmt.Errorf("%s %s: expected response status %v, got %d", r.Request.Method, r.Request.URL, codes, r.StatusCode)
	})
}

// Do is not much different from calling client.Do(request) except it runs the
// response validation. See WithValidator and WithStatusRequired.
func Do(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	validators, _ := req.Context().Value(validatorKey).([]Validator)
	for _, v := range validators {
		if err := v(resp); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// Retry re-attempts request with provided intervals. By manually providing intervals sequence you
// can have different wait strategies like exponential back-off (time.Second, 2 * time.Second, 4 * time.Second)
// or just multiple reties after same interval (time.Second, time.Second, time.Second). If Request had a context
// with timeout cancellation then it will be applied to entire chain.
func Retry(client *http.Client, req *http.Request, intervals ...time.Duration) (*http.Response, error) {
	ctx := req.Context()
	for {
		resp, err := Do(client, req)
		// interval exhausted or request succeeded
		if len(intervals) == 0 || err == nil {
			return resp, err
		}
		select {
		case <-time.After(intervals[0]):
			intervals = intervals[1:]
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// All runs requests simultaneously returning responses in same order or error if at least one request failed.
// Once result is determined all requests are cancelled through the context.
func All(client *http.Client, reqs ...*http.Request) ([]*http.Response, error) {
	results := make(chan result, len(reqs))
	stop := make(chan struct{})
	defer close(stop)
	for i, req := range reqs {
		go do(client, req, i, stop, results)
	}
	cnt := 0
	resps := make([]*http.Response, len(reqs), len(reqs))
	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		cnt++
		resps[res.order] = res.response
		if cnt == len(reqs) {
			break
		}
	}
	return resps, nil
}

// Some runs requests simultaneously returning responses for successful requests and <nil> for failed ones.
// Error is returned only if all requests failed.
func Some(client *http.Client, reqs ...*http.Request) ([]*http.Response, error) {
	results := make(chan result, len(reqs))
	stop := make(chan struct{})
	defer close(stop)
	for i, req := range reqs {
		go do(client, req, i, stop, results)
	}
	cnt, success := 0, 0
	resps := make([]*http.Response, len(reqs), len(reqs))
	for res := range results {
		cnt++
		if res.err == nil {
			success++
			resps[res.order] = res.response
		}
		if cnt == len(reqs) {
			break
		}
	}
	if success == 0 {
		return nil, fmt.Errorf("all requests failed")
	}
	return resps, nil
}

// Race runs requests simultaneously returning first successful result or error if all failed.
// Once result is determined all requests are cancelled through the context.
func Race(client *http.Client, reqs ...*http.Request) (*http.Response, error) {
	results := make(chan result, len(reqs))
	stop := make(chan struct{})
	defer close(stop)
	for i, req := range reqs {
		go do(client, req, i, stop, results)
	}
	cnt := 0
	for res := range results {
		if res.err == nil {
			return res.response, nil
		}
		cnt++
		if cnt == len(reqs) {
			break
		}
	}
	return nil, fmt.Errorf("all requests failed")
}

type result struct {
	order    int // keep resp the same order with req
	response *http.Response
	err      error
}

func do(client *http.Client, req *http.Request, order int, stop <-chan struct{}, results chan<- result) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()
	go func() {
		<-stop
		cancel()
	}()
	resp, err := Do(client, req.WithContext(ctx))
	results <- result{
		order:    order,
		response: resp,
		err:      err,
	}
}
