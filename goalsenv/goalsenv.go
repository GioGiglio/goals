// Package goalsenv exports variable related to goals program
package goalsenv

import (
	"go/build"
	"os"
)

// GoPath refers to go env variable GOPATH
var GoPath string

// GoalsPath refers to GOPATH/src/Goals/
var GoalsPath string

func init() {
	GoPath := os.Getenv("GOPATH")
	if GoPath == "" {
		GoPath = build.Default.GOPATH
	}
	GoalsPath = GoPath + "/src/goals/"
}
