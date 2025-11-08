package main

import "testing"

func TestParsePatchDirectivesDefaultsToNginx(t *testing.T) {
	directives, err := parsePatchDirectives("./patches/foo.patch")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(directives))
	}
	d := directives[0]
	if d.target != defaultPatchTarget {
		t.Fatalf("expected default target %q, got %q", defaultPatchTarget, d.target)
	}
	if len(d.paths) != 1 || d.paths[0] != "./patches/foo.patch" {
		t.Fatalf("unexpected paths: %#v", d.paths)
	}
}

func TestParsePatchDirectivesMixedTargets(t *testing.T) {
	input := "foo.patch,openssl=bar.patch"
	directives, err := parsePatchDirectives(input)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(directives) != 2 {
		t.Fatalf("expected 2 directives, got %d", len(directives))
	}
	if directives[0].target != defaultPatchTarget {
		t.Fatalf("first directive should target %s, got %s", defaultPatchTarget, directives[0].target)
	}
	if len(directives[0].paths) != 1 || directives[0].paths[0] != "foo.patch" {
		t.Fatalf("unexpected default paths: %#v", directives[0].paths)
	}
	if directives[1].target != "openssl" {
		t.Fatalf("second directive should target openssl, got %s", directives[1].target)
	}
	if len(directives[1].paths) != 1 || directives[1].paths[0] != "bar.patch" {
		t.Fatalf("unexpected openssl paths: %#v", directives[1].paths)
	}
}
