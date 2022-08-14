package populate

import (
	"log"
	"os"
	"path"

	"go.uber.org/zap"
)

const (
	logpathEnvVarName = "AH_LOGPATH"
	logpathFileName   = "populate.log"
)

// The are identifiers for functions that use logging.
const (
	populateCoinbaseProCandlesFn = "fn-e540"
)

var logger *zap.Logger

func logpath() string { return os.Getenv(logpathEnvVarName) }
func logfile() string { return path.Join(logpath(), logpathFileName) }

func init() {
	cfg := zap.NewDevelopmentConfig()
	// cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logfile()}

	var err error
	logger, err = cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger.Sugar().Info("Log started")
}
