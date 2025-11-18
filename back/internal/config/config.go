// internal/config/config.go
package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
	Loc      string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// Config 应用配置
type Config struct {
	Database *DatabaseConfig
	Server   *ServerConfig
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Database: &DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "fjt911008"),
			DBName:   getEnv("DB_NAME", "merchant_admin"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
			Loc:      getEnv("DB_LOC", "Local"),
		},
		Server: &ServerConfig{
			Port: getEnv("SERVER_PORT", "8088"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
	}
}

// buildDSN 构建数据库连接字符串
func buildDSN(config *DatabaseConfig) string {
	return config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.DBName + "?charset=" + config.Charset + "&parseTime=True&loc=" + config.Loc
}

// InitDatabase 初始化数据库连接
func (c *Config) InitDatabase() (*gorm.DB, error) {
	// 构建连接字符串
	dsn := buildDSN(c.Database)

	log.Printf("Connecting to database: %s:%s/%s", c.Database.Host, c.Database.Port, c.Database.DBName)

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启日志
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// 获取底层的 *sql.DB 对象以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully!")
	return db, nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
