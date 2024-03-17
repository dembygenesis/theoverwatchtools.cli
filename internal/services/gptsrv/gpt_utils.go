package gptsrv

import (
	"fmt"
	"github.com/atotto/clipboard"
)

type GptUtils interface {
	ClipCodingStandardsPreface() error
}

func New() GptUtils {
	return &gptUtils{}
}

type gptUtils struct {
}

func (g *gptUtils) ClipCodingStandardsPreface() error {
	if err := clipboard.WriteAll(preface); err != nil {
		return fmt.Errorf("clip preface: %v", err)
	}
	return nil
}
