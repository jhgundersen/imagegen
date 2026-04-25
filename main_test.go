package main

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestExtractOutputArgSupportsFlagAfterPrompt(t *testing.T) {
	args, output := extractOutputArg([]string{"a prompt", "--output", "out.png"})

	if output != "out.png" {
		t.Fatalf("output = %q, want %q", output, "out.png")
	}
	wantArgs := []string{"a prompt"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("args = %#v, want %#v", args, wantArgs)
	}
}

func TestExtractOutputArgSupportsEqualsForm(t *testing.T) {
	args, output := extractOutputArg([]string{"-o=out.png", "a prompt"})

	if output != "out.png" {
		t.Fatalf("output = %q, want %q", output, "out.png")
	}
	wantArgs := []string{"a prompt"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("args = %#v, want %#v", args, wantArgs)
	}
}

func TestExtractOutputArgUsesLastValue(t *testing.T) {
	args, output := extractOutputArg([]string{"-o", "first.png", "a prompt", "--output", "second.png"})

	if output != "second.png" {
		t.Fatalf("output = %q, want %q", output, "second.png")
	}
	wantArgs := []string{"a prompt"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("args = %#v, want %#v", args, wantArgs)
	}
}

func TestReadResponseDataClassifiesTaskNotFound(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"code":1,"message":"task not found"}`)),
	}

	_, err := readResponseData(resp)
	if err == nil {
		t.Fatal("err = nil, want task not found error")
	}
	if !isTaskNotFoundError(err) {
		t.Fatalf("isTaskNotFoundError() = false, want true for %v", err)
	}
}

func TestReadResponseDataDoesNotClassifyOther404sAsTaskNotFound(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"code":2,"message":"not found"}`)),
	}

	_, err := readResponseData(resp)
	if err == nil {
		t.Fatal("err = nil, want response error")
	}
	if isTaskNotFoundError(err) {
		t.Fatalf("isTaskNotFoundError() = true, want false for %v", err)
	}
}
