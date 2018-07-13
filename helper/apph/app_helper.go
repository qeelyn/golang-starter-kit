package apph

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qeelyn/go-common/logger"
	"github.com/spf13/viper"
	_ "gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
func LoadConfig(configFile string) (*viper.Viper, error) {
	//var filename, ext string = "app", "yaml"
	realPath, _ := filepath.Abs(configFile)
	file, err := os.Stat(realPath)
	if err != nil {
		return nil, err
	}
	configPath := path.Dir(realPath)
	fn := strings.Split(file.Name(), ".")
	filename := fn[0]
	ext := fn[1]
	cnf := viper.New()
	//cnf.WatchConfig()
	cnf.SetConfigName(filename)
	cnf.SetConfigType(ext)
	cnf.AutomaticEnv()

	cnf.AddConfigPath(configPath)
	cnf.SetDefault("debug", false)

	if err := cnf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read the configuration file: %s", err)
	}
	// local
	localConfig := path.Join(configPath, filename+"-local."+ext)
	if _, err := os.Stat(localConfig); err != nil {
		return nil, err
	}

	cnf.SetConfigName(filename + "-local")
	if err := cnf.MergeInConfig(); err != nil {
		return nil, err
	}

	switch cnf.GetString("appmode") {
	case "debug":
		cnf.Set("debug", true)
	case "release":
	}

	return cnf, nil
}

func NewDb(config map[string]interface{}) *gorm.DB {
	orm, err := gorm.Open(config["dialect"].(string), config["dsn"].(string))
	if err != nil {
		panic(err)
	}
	if _, ok := config["maxidleconns"]; ok {
		orm.DB().SetMaxIdleConns(config["maxidleconns"].(int))
	}
	if _, ok := config["maxopenconns"]; ok {
		orm.DB().SetMaxOpenConns(config["maxopenconns"].(int))
	}
	if _, ok := config["connmaxlifetime"]; ok {
		orm.DB().SetConnMaxLifetime(time.Duration(config["connmaxlifetime"].(int)) * time.Second)
	}
	return orm
}

func NewLogger(config map[string]interface{}) *logger.Logger {
	file := logger.NewFileLogger(config)
	return logger.NewLogger(file)
}
