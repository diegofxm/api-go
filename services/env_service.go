package services

import (
	"fmt"
	"os"
)

// EnsureEnvFile verifica que exista el archivo .env
func EnsureEnvFile() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("\n=== Configuración Inicial Requerida ===")
		fmt.Println("El archivo .env no existe. Por favor, sigue estos pasos:")
		fmt.Println("1. Copia el archivo .env.example y renómbralo a .env")
		fmt.Println("2. Ejecuta el comando: tools/generate_jwt_key.exe")
		fmt.Println("3. Copia la clave JWT generada y pégala en el archivo .env")
		fmt.Println("4. Ajusta las demás configuraciones según tu entorno")
		return fmt.Errorf("archivo .env no encontrado")
	}
	return nil
}
