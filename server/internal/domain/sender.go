package domain

import "os"

type SenderUsecase interface {
	ReadFile(fileName string) (*os.File, error)
}
