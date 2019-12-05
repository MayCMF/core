package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	global *Config
)

// LoadGlobal - Load global configuration
func LoadGlobal(fpath string) error {
	c, err := Parse(fpath)
	if err != nil {
		return err
	}
	global = c
	return nil
}

// Global - Get global configuration
func Global() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

// Parse configuration file
func Parse(fpath string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(fpath, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Config parameters
type Config struct {
	RunMode     string      `toml:"run_mode"`
	WWW         string      `toml:"www"`
	Swagger     string      `toml:"swagger"`
	Store       string      `toml:"store"`
	HTTP        HTTP        `toml:"http"`
	I18n        i18n        `toml:"i18n"`
	Permission  Permission  `toml:"permission"`
	Casbin      Casbin      `toml:"casbin"`
	Log         Log         `toml:"log"`
	LogGormHook LogGormHook `toml:"log_gorm_hook"`
	Root        Root        `toml:"root"`
	JWTAuth     JWTAuth     `toml:"jwt_auth"`
	Monitor     Monitor     `toml:"monitor"`
	Captcha     Captcha     `toml:"captcha"`
	RateLimiter RateLimiter `toml:"rate_limiter"`
	CORS        CORS        `toml:"cors"`
	Redis       Redis       `toml:"redis"`
	Gorm        Gorm        `toml:"gorm"`
	MySQL       MySQL       `toml:"mysql"`
	Postgres    Postgres    `toml:"postgres"`
	Sqlite3     Sqlite3     `toml:"sqlite3"`
	FileManager FileManager `toml:"filemanager"`
}

// IsDebugMode - Is it debug mode?
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// i18n configuration parameters
type i18n struct {
	Enable          bool     `toml:"enable"`
	Default         string   `toml:"Default"`
	Languages       []string `toml:"languages"`
	Data            string   `toml:"data"`
	CountriesEnable bool     `toml:"countries_enable"`
	CountriesData   string   `toml:"countries_data"`
}

// Permission configuration parameters
type Permission struct {
	Enable bool   `toml:"enable"`
	Data   string `toml:"data"`
}

// Casbin configuration parameters
type Casbin struct {
	Enable           bool   `toml:"enable"`
	Debug            bool   `toml:"debug"`
	Model            string `toml:"model"`
	AutoLoad         bool   `toml:"auto_load"`
	AutoLoadInternal int    `toml:"auto_load_internal"`
}

// Log configuration parameters
type Log struct {
	Level         int    `toml:"level"`
	Format        string `toml:"format"`
	Output        string `toml:"output"`
	OutputFile    string `toml:"output_file"`
	EnableHook    bool   `toml:"enable_hook"`
	Hook          string `toml:"hook"`
	HookMaxThread int    `toml:"hook_max_thread"`
	HookMaxBuffer int    `toml:"hook_max_buffer"`
}

// Log Gorm Hook configuration
type LogGormHook struct {
	DBType       string `toml:"db_type"`
	MaxLifetime  int    `toml:"max_lifetime"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	Table        string `toml:"table"`
}

// Root User
type Root struct {
	UserName string `toml:"user_name"`
	Password string `toml:"password"`
	RealName string `toml:"real_name"`
}

// JWTAuth User Authentication
type JWTAuth struct {
	SigningMethod string `toml:"signing_method"`
	SigningKey    string `toml:"signing_key"`
	Expired       int    `toml:"expired"`
	Store         string `toml:"store"`
	FilePath      string `toml:"file_path"`
	RedisDB       int    `toml:"redis_db"`
	RedisPrefix   string `toml:"redis_prefix"`
}

// HTTP configuration parameters
type HTTP struct {
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	ShutdownTimeout int    `toml:"shutdown_timeout"`
}

// Monitor configuration parameters
type Monitor struct {
	Enable    bool   `toml:"enable"`
	Addr      string `toml:"addr"`
	ConfigDir string `toml:"config_dir"`
}

// Captcha - Graphic verification code configuration parameter
type Captcha struct {
	Store       string `toml:"store"`
	Length      int    `toml:"length"`
	Width       int    `toml:"width"`
	Height      int    `toml:"height"`
	RedisDB     int    `toml:"redis_db"`
	RedisPrefix string `toml:"redis_prefix"`
}

// RateLimiter - Request frequency limit configuration parameter
type RateLimiter struct {
	Enable  bool  `toml:"enable"`
	Count   int64 `toml:"count"`
	RedisDB int   `toml:"redis_db"`
}

// CORS Cross-domain request configuration parameters
type CORS struct {
	Enable           bool     `toml:"enable"`
	AllowOrigins     []string `toml:"allow_origins"`
	AllowMethods     []string `toml:"allow_methods"`
	AllowHeaders     []string `toml:"allow_headers"`
	AllowCredentials bool     `toml:"allow_credentials"`
	MaxAge           int      `toml:"max_age"`
}

// Redis configuration parameters
type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

// Gorm configuration parameters
type Gorm struct {
	Debug             bool   `toml:"debug"`
	DBType            string `toml:"db_type"`
	MaxLifetime       int    `toml:"max_lifetime"`
	MaxOpenConns      int    `toml:"max_open_conns"`
	MaxIdleConns      int    `toml:"max_idle_conns"`
	TablePrefix       string `toml:"table_prefix"`
	EnableAutoMigrate bool   `toml:"enable_auto_migrate"`
}

// Postgres Configuration parameter
type FileManager struct {
	Dir         string   `toml:"dir"`
	MaxSize     int64    `toml:"maxsize"`
	ImagesDir   string   `toml:"images_dir"`
	AllowImages []string `toml:"allow_images"`
	FilesDir    string   `toml:"files_dir"`
	AllowFiles  []string `toml:"allow_files"`
}

// MySQL configuration parameters
type MySQL struct {
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DBName     string `toml:"db_name"`
	Parameters string `toml:"parameters"`
}

// DSN Database connection string
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

// Postgres Configuration parameter
type Postgres struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	SSLMode  string `toml:"ssl_mode"`
}

// DSN Database connection string
func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

// Sqlite3 configuration parameter
type Sqlite3 struct {
	Path string `toml:"path"`
	Name string `toml:"name"`
	Dir  string `toml:"dir"`
}

// DSN Database connection string
func (a Sqlite3) DSN() string {
	return a.Path
}
