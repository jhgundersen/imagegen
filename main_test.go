package main

import (
	"reflect"
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
