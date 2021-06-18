package generator

import (
	"fmt"
	"os"
	"testing"

	"github.com/longXboy/grpc-gateway/v2/internal/codegenerator"
)

func Test_Gen(t *testing.T) {
	var err error
	f, err := os.Open("./req.bin")
	if err != nil {
		t.Fatal(err)
	}
	req, err := codegenerator.ParseRequest(f)
	if err != nil {
		t.Fatal(err)
	}
	var g Generator

	resp, err := g.Gen(req, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range resp.File {
		fmt.Println(*file.Content)
	}
}
