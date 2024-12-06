package services

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type ConfigService struct {
	mu       sync.RWMutex
	configs  map[string]string
	required []string
}

func NewConfigService(envFile string, required []string) (*ConfigService, error) {
	service := &ConfigService{
		configs:  make(map[string]string),
		required: required,
	}

	// Cargar variables de entorno desde el archivo
	if err := godotenv.Load(envFile); err != nil {
		return nil, err
	}

	// Cargar todas las variables de entorno en el mapa
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		service.configs[pair[0]] = pair[1]
	}

	// Validar variables requeridas
	if err := service.ValidateRequired(); err != nil {
		return nil, err
	}

	return service, nil
}

func (c *ConfigService) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.configs[key]
}

func (c *ConfigService) GetInt(key string) (int, error) {
	value := c.Get(key)
	return strconv.Atoi(value)
}

func (c *ConfigService) GetBool(key string) bool {
	value := strings.ToLower(c.Get(key))
	return value == "true" || value == "1" || value == "yes"
}

func (c *ConfigService) GetFloat(key string) (float64, error) {
	value := c.Get(key)
	return strconv.ParseFloat(value, 64)
}

func (c *ConfigService) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configs[key] = value
	os.Setenv(key, value)
}

func (c *ConfigService) ValidateRequired() error {
	missing := []string{}
	
	for _, key := range c.required {
		if _, exists := c.configs[key]; !exists {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return errors.New("Missing required environment variables: " + strings.Join(missing, ", "))
	}

	return nil
}

func (c *ConfigService) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// Crear una copia del mapa para evitar modificaciones externas
	configs := make(map[string]string)
	for k, v := range c.configs {
		configs[k] = v
	}
	
	return configs
}

func (c *ConfigService) LoadFromFile(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Actualizar el mapa con las nuevas variables
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		c.configs[pair[0]] = pair[1]
	}

	return nil
}

func (c *ConfigService) GetWithDefault(key, defaultValue string) string {
	if value := c.Get(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *ConfigService) GetIntWithDefault(key string, defaultValue int) int {
	if value, err := c.GetInt(key); err == nil {
		return value
	}
	return defaultValue
}

func (c *ConfigService) GetFloatWithDefault(key string, defaultValue float64) float64 {
	if value, err := c.GetFloat(key); err == nil {
		return value
	}
	return defaultValue
}

func (c *ConfigService) GetBoolWithDefault(key string, defaultValue bool) bool {
	if value := c.Get(key); value != "" {
		return c.GetBool(key)
	}
	return defaultValue
}
