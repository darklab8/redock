package settings

import (
	"log"
	"path/filepath"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type RedockCfg struct {
	utils_settings.UtilsEnvs

	DockerHost string
}

var Env RedockCfg

var Environ *enverant.Enverant

var Workdir string

func init() {
	log.Println("attempt to load settings")
	Environ = enverant.NewEnverant()
	LoadEnv(Environ)
}

func LoadEnv(envs *enverant.Enverant) {
	Env = RedockCfg{
		UtilsEnvs:  utils_settings.GetEnvs(envs),
		DockerHost: envs.GetStrOr("DOCKER_HOST", "unix:///var/run/docker.sock"),
	}
	Workdir = filepath.Dir(utils_os.GetCurrentFolder().ToString())
}
