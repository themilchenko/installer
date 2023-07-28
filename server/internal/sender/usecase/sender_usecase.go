package sender

import (
	"errors"
	"os"
)

const (
	defaultSourceFolder = "/home/milchenko/programming/aktiv/base.agent/bin/"
)

type SenderUsecase struct{}

func NewSenderUsecase() SenderUsecase {
	return SenderUsecase{}
}

func (u SenderUsecase) ReadFile(fileName string) (*os.File, error) {
	stat, err := os.Stat(defaultSourceFolder + fileName)
	if err != nil {
		return nil, err
	}
	if stat.Name() != fileName {
		return nil, errors.New("it's not requested filename")
	}

	f, err := os.Open(defaultSourceFolder + fileName)
	if err != nil {
		return nil, err
	}

	return f, nil
}
