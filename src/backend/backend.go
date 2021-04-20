package backend

import (
	"fmt"
	"strings"

	"github.com/vj396/milton/src/backend/mysql"
	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"
)

func Connect(logger *zap.Logger, dbType string, conf *types.DatabaseMetadata, modelsDir string) (types.Backend, error) {
	switch strings.ToLower(dbType) {
	case "mysql":
		return mysql.New(logger, conf, modelsDir)
	default:
		return nil, fmt.Errorf("backend type not supported")
	}
}
