//go:build prod

package utils

import (
	"go.uber.org/zap"
	"log"
)

func init() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("logger creation error: %v", err)
	}

	Sugar = Logger.Sugar()
}
