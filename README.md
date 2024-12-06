# API REST con Go, GORM y JWT

Esta es una API REST desarrollada en Go que utiliza GORM como ORM y JWT para autenticación. El proyecto implementa un sistema de usuarios con roles y manejo seguro de variables de entorno.

## Características

- 🔐 Autenticación mediante JWT
- 👥 Sistema de usuarios con roles
- 🛡️ Manejo seguro de variables de entorno
- 📝 Respuestas JSON estructuradas
- 🔄 CRUD completo para usuarios
- 🎯 Middleware de autenticación
- 📊 Paginación de resultados
- 🧪 Pruebas automatizadas
- 📚 Documentación completa

## Estructura del Proyecto

```
go-api-orm/
├── config/             # Configuraciones de la aplicación
├── controllers/        # Controladores de la API
├── middlewares/       # Middlewares personalizados
├── models/            # Modelos de la base de datos
├── routes/            # Definición de rutas
├── services/          # Lógica de negocio
├── tools/             # Herramientas útiles
├── tests/             # Pruebas automatizadas
├── .env.example       # Plantilla de variables de entorno
├── .gitignore         # Archivos ignorados por git
├── go.mod            # Dependencias del proyecto
├── go.sum            # Checksums de dependencias
├── main.go           # Punto de entrada de la aplicación
└── README.md         # Esta documentación
```

## Configuración Inicial

### 1. Requisitos Previos

- Go 1.21 o superior
- MySQL, PostgreSQL, SQLite o MongoDB
- Git

### 2. Instalación

```bash
# Clonar el repositorio
git clone https://github.com/tu-usuario/go-api-orm.git
cd go-api-orm

# Instalar dependencias
go mod download

# Configurar variables de entorno
cp .env.example .env

# Generar clave JWT
tools/generate_jwt_key.exe

# Ejecutar la aplicación
go run main.go
```

### 3. Variables de Entorno

El proyecto utiliza un archivo `.env` para la configuración. Para configurarlo:

1. Copia el archivo `.env.example` y renómbralo a `.env`:
   ```bash
   cp .env.example .env
   ```

2. Genera una clave JWT segura usando la herramienta incluida:
   ```bash
   tools/generate_jwt_key.exe
   ```

3. Copia la clave generada y pégala en tu archivo `.env`

4. Ajusta las demás variables según tu entorno:
   ```env
   # JWT Configuration
   JWT_SECRET_KEY=tu_clave_generada
   JWT_EXPIRATION_HOURS=24

   # Server Configuration
   PORT=8080
   GIN_MODE=debug

   # Database Configuration
   DB_DRIVER=mysql
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_NAME=cms-builder
   DB_PASSWORD=
   DB_SQLITE_PATH=./api.db

   # Response Configuration
   SHOW_METADATA=false
   SHOW_PAGINATION=true
   ```

## Uso de la API

### Ejemplos con cURL

#### 1. Crear Usuario
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "contraseña123",
    "username": "usuario1",
    "role_id": 1
  }'
```

#### 2. Iniciar Sesión
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "contraseña123"
  }'
```

#### 3. Obtener Usuario (Autenticado)
```bash
curl -X GET http://localhost:8080/users/1 \
  -H "Authorization: Bearer tu_token_jwt"
```

### Ejemplos con Postman

1. Importa la colección de Postman incluida en `docs/postman/`
2. Configura las variables de entorno en Postman:
   - `base_url`: URL base de tu API
   - `token`: Token JWT después de iniciar sesión

## Desarrollo

### Configuración del Entorno de Desarrollo

1. Instala las herramientas de desarrollo:
   ```bash
   go install github.com/cosmtrek/air@latest  # Hot reload
   go install github.com/golang/mock/mockgen@latest  # Generador de mocks
   ```

2. Ejecuta la aplicación en modo desarrollo:
   ```bash
   air
   ```

### Convenciones de Código

1. **Estructura de Archivos**:
   - Un paquete por directorio
   - Nombres de archivo en snake_case
   - Pruebas en archivos `_test.go`

2. **Nombrado**:
   - Interfaces: nombres terminados en 'er' (ej: `UserService`)
   - Implementaciones: prefijo con el tipo (ej: `SQLUserRepository`)
   - Pruebas: sufijo `Test` (ej: `TestCreateUser`)

3. **Documentación**:
   - Todos los paquetes exportados deben tener documentación
   - Usar ejemplos en la documentación cuando sea posible

### Flujo de Trabajo Git

1. Crear rama para nueva característica:
   ```bash
   git checkout -b feature/nombre-caracteristica
   ```

2. Commits semánticos:
   ```bash
   feat: añadir autenticación JWT
   fix: corregir validación de email
   docs: actualizar README
   test: añadir pruebas para UserService
   ```

## Pruebas

### Ejecutar Pruebas

```bash
# Todas las pruebas
go test ./...

# Con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Pruebas de un paquete específico
go test ./services -v
```

### Tipos de Pruebas

1. **Unitarias**: En archivos `_test.go` junto al código
2. **Integración**: En el directorio `tests/integration`
3. **End-to-End**: En el directorio `tests/e2e`

### Mocks

Usamos `mockgen` para generar mocks:

```bash
# Generar mock para una interfaz
mockgen -source=services/user_service.go -destination=mocks/mock_user_service.go
```

## Despliegue

### Preparación

1. Compilar para producción:
   ```bash
   go build -o api-server
   ```

2. Variables de entorno para producción:
   ```env
   GIN_MODE=release
   PORT=80
   ```

### Docker

```bash
# Construir imagen
docker build -t go-api-orm .

# Ejecutar contenedor
docker run -p 8080:8080 go-api-orm
```

## Monitoreo y Logs

- Los logs se escriben en formato JSON
- Niveles de log: DEBUG, INFO, WARN, ERROR
- Métricas disponibles en `/metrics` (Prometheus)

## Seguridad

### JWT

- Las claves JWT se generan de forma segura usando `crypto/rand`
- La clave nunca se sube al repositorio
- Se proporciona una herramienta dedicada para generar claves seguras

### Variables de Entorno

- Se usa `.env.example` como plantilla
- El archivo `.env` está en `.gitignore`
- Las variables sensibles nunca se exponen en el código

### Mejores Prácticas

1. Todas las contraseñas se hashean antes de almacenarse
2. Implementación de rate limiting
3. Validación de entrada en todas las rutas
4. Headers de seguridad configurados
5. CORS configurado apropiadamente

## Contribuir

1. Haz un fork del repositorio
2. Crea una rama para tu característica (`git checkout -b feature/AmazingFeature`)
3. Haz commit de tus cambios (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

### Guía de Contribución

1. Asegúrate de que las pruebas pasen
2. Actualiza la documentación si es necesario
3. Sigue las convenciones de código
4. Añade pruebas para nuevas características

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## Soporte

Si encuentras un bug o tienes una sugerencia:

1. Revisa los issues existentes
2. Abre un nuevo issue con:
   - Descripción clara del problema
   - Pasos para reproducirlo
   - Comportamiento esperado
   - Capturas de pantalla si aplica
