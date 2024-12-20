package logus

import (
	"github.com/darklab8/go-typelog/typelog"
	_ "github.com/darklab8/redock/settings" // enverant.json injection to env
)

var Log *typelog.Logger = typelog.NewLogger(
	"redock",
	typelog.WithLogLevel(typelog.LEVEL_INFO),
)
