// package etc read config file defined by XXXROOT env.
package etc

import (
	"errors"
	"os"

	"common/etc/config"
)

var (
	rootEnv  = "YF_PAAS_ROOT"
	confName = "app.cfg"
)

var (
	ErrNoFile = errors.New("not found cfg file.")
	ErrNoCfg  = errors.New("no such config value in cfg file.")
)

var (
	gCfg *config.Config
)

func init() {
	if err := loadCfg(); err != nil {
		panic("no local cfg and no global cfg.")
	}
}

func loadCfg() error {
	// product
	if cfg, err := config.ReadDefault(cfgPath(confName)); err == nil {
		gCfg = cfg
		return err
	}
	// develop
	if cfg, err := config.ReadDefault(cfgPath(confName + ".dev")); err == nil {
		gCfg = cfg
		return err
	}
	// err
	return ErrNoFile
}

func cfgPath(name string) string {
	return RootDir() + "/etc/" + name
}

// -----------------------------------------------------------------------------
// Dirs

func RootDir() string {
	dir := os.Getenv(rootEnv)
	if len(dir) == 0 {
		panic("you should define ROOT environment.")
	}
	return dir
}

func EtcDir() string {
	return RootDir() + "/etc/"
}

// -----------------------------------------------------------------------------
// APIs

func String(session, key string) string {
	s, err := gCfg.String(session, key)
	if err != nil {
		panic(ErrNoCfg.Error() + session + key)
	}
	return s
}

func Int(session, key string) int {
	i, err := gCfg.Int(session, key)
	if err != nil {
		panic(ErrNoCfg.Error() + session + key)
	}
	return i
}

func Bool(session, key string) bool {
	b, err := gCfg.Bool(session, key)
	if err != nil {
		panic(ErrNoCfg.Error() + session + key)

	}
	return b
}

func Float(session, key string) float64 {
	f, err := gCfg.Float(session, key)
	if err != nil {
		panic(ErrNoCfg.Error() + session + key)
	}
	return f
}
