package main

import (
	"chi_soccer/services"
	"os"
	"testing"
)

var app services.Application

func TestMain(m *testing.M) {
    os.Exit(m.Run())
}
