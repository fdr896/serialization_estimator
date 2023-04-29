package support

import (
	"os"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func InitLogger() {
    zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}