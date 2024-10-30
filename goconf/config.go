package goconf

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var Props AppProps
var ProjectRoot = getRootProjectDir()

type AppProps struct {
	Port   int
	Domain string
	DB
	Stdout
	GCP
	TableService `mapstructure:"table_service"`
	Bank
	Ws
	Bonus
	Mail
	ProfileService `mapstructure:"profile_service"`
	Survival
}

type DB struct {
	Scheme   string
	Host     string
	Port     *int
	Username string
	Password string
}

type Stdout struct {
	Level string
}

type GCP struct {
	Credentials string
	ProjectID   string `mapstructure:"project_id"`
	Metrics
}

type Metrics struct {
	PushDuration time.Duration `mapstructure:"push_duration"`
}

type Bank struct {
	RankerDuration time.Duration `mapstructure:"ranker_duration"`
}

type Mail struct {
	Host     string
	Port     int
	From     string
	Username string
	Password string
}

type Ws struct {
	Domain  string
	LogoURL string `mapstructure:"logo_url"`
}

type Bonus struct {
	Cron string
}

type ProfileService struct {
	GCS
}

type GCS struct {
	AvatarBucket string `mapstructure:"avatar_bucket"`
}

type Survival struct {
	MaxAnonymousCounter int `mapstructure:"max_anonymous_counter"`
}

func init() {
	Init()
}

func Init() {
	if IsLocal() {
		os.Setenv("PG_JWT_VERIFICATION_DISABLE", "true")
	}
	ParseConfig(ProjectRoot, &Props)
	InitLogs(Props.Stdout.Level)

	log.Debugf("Configuration properties: %+v", Props)
}

func ParseConfig(projectRoot string, props interface{}) {
	viper.SetConfigName("config")
	viper.AddConfigPath(projectRoot)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	loadEnvVariables()

	err := viper.Unmarshal(props)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func loadEnvVariables() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("PG") // IMPORTANT!!!
	viper.AutomaticEnv()
}

func InitLogs(levelStr string) {
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		log.Errorf("Couldn't parse log level, falling back to the INFO level")
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	//log.SetFormatter(NewGCEFormatter(true))
}

func getRootProjectDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(b))
}
