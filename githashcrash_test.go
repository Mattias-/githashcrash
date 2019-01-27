package main

import (
	"os/exec"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	obj, _ := exec.Command("git", "cat-file", "-p", "HEAD").Output()

	for i := 0; i < b.N; i++ {
		run("^0000.*", obj, "11", 1)
	}
}
