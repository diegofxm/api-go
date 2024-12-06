package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go-api-orm/migrations"
	"go-api-orm/models"
	"go-api-orm/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func createDatabaseIfNotExists(driver string) error {
	dbName := os.Getenv("DB_NAME")
	
	switch driver {
	case "mysql":
		// Crear conexión sin especificar la base de datos
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"))
		
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		// Crear la base de datos si no existe usando backticks para escapar el nombre
		createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_spanish_ci", dbName)
		_, err = db.Exec(createDBQuery)
		if err != nil {
			return fmt.Errorf("error creating database: %v", err)
		}

	case "postgres":
		// Conectar a la base de datos 'postgres' por defecto
		dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_PORT"))

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		// Verificar si la base de datos existe
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
		err = db.QueryRow(query).Scan(&exists)
		if err != nil {
			return err
		}

		// Crear la base de datos si no existe
		if !exists {
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InitDB() {
	var err error
	driver := os.Getenv("DB_DRIVER")

	// Crear la base de datos si no existe (para MySQL y PostgreSQL)
	if driver == "mysql" || driver == "postgres" {
		if err := createDatabaseIfNotExists(driver); err != nil {
			log.Printf("Error creating database: %v", err)
		}
	}

	// Configurar el logger JSON personalizado
	logger := utils.NewJSONLogger()

	config := &gorm.Config{
		Logger: logger,
	}

	switch driver {
	case "mysql":
		DB, err = connectMySQL(config)
	case "postgres":
		DB, err = connectPostgres(config)
	default: // sqlite como default
		DB, err = connectSQLite(config)
	}

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Auto-migrar los modelos
	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Role{},
	)
	if err != nil {
		log.Fatalf("Error auto-migrating database: %v", err)
	}

	// Crear roles por defecto
	if err := migrations.SeedDefaultRoles(DB); err != nil {
		log.Printf("Error seeding default roles: %v", err)
	}
}

func connectMySQL(config *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	
	return gorm.Open(mysql.Open(dsn), config)
}

func connectPostgres(config *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	
	return gorm.Open(postgres.Open(dsn), config)
}

func connectSQLite(config *gorm.Config) (*gorm.DB, error) {
	dbPath := os.Getenv("DB_SQLITE_PATH")
	if dbPath == "" {
		dbPath = "api.db"
	}
	return gorm.Open(sqlite.Open(dbPath), config)
}
