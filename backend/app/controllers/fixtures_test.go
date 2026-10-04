package controllers_test

import "time"

// canned answers each upstream path a service requests. Every task that adds an
// upstream adds its path here; ok is false for a path no service should call.
func canned(path string, now time.Time) (body []byte, ok bool) {
	switch path {
	}
	return nil, false
}
