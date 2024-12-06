package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func generateSecretKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func main() {
	key, err := generateSecretKey()
	if err != nil {
		fmt.Printf("Error generando la clave: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Generador de Clave JWT ===")
	fmt.Println()
	fmt.Println("Se ha generado una nueva clave JWT segura:")
	fmt.Printf("JWT_SECRET_KEY=%s\n", key)
	fmt.Println()
	fmt.Println("Instrucciones:")
	fmt.Println("1. Copia la línea completa 'JWT_SECRET_KEY=...'")
	fmt.Println("2. Pégala en tu archivo .env")
	fmt.Println("3. Reinicia tu aplicación")
	fmt.Println()
	fmt.Println("Nota: Mantén esta clave segura y nunca la compartas")
}
