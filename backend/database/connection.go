// database/connection.go
// Provee la inicialización de la conexión a la base de datos y asegura
// que las tablas necesarias existan.
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB es la conexión global usada por los handlers.
var DB *sql.DB

// InitDB abre la conexión a la base de datos, configura el pool y verifica
// la conexión. También llama a createTables para asegurar que la tabla
// `items` exista.
func InitDB() error {
	// String de conexión para la base de datos CockroachDB local
	connStr := "postgresql://root@localhost:26257/defaultdb?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error conectando a la base de datos: %v", err)
	}

	// Configurar conexión
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * 60) // 5 minutos

	// Verificar la conexión
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error haciendo ping: %v", err)
	}

	log.Println("Conectado a CockroachDB exitosamente")
	return createTables()
}

// createTables crea la tabla `items` si no existe. La estructura actual
// refleja los campos esperados por el modelo: ticker, target_from/target_to,
// company, action, brokerage, rating_from, rating_to, time y created_at.
func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		ticker STRING NOT NULL,
		target_from FLOAT NOT NULL,
		target_to FLOAT NOT NULL,
		company STRING NOT NULL,
		action STRING NOT NULL,
		brokerage STRING ,
		rating_from STRING NOT NULL,
		rating_to STRING NOT NULL,
		time TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT now()
	)
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creando tabla: %v", err)
	}

	log.Println("Tabla 'items' creada/verificada exitosamente")
	return nil
}
