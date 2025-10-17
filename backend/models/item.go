// models/item.go
// Contiene los tipos de datos usados por la API: Item y auxiliares.
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Money representa un monto monetario almacenado internamente como float64.
// En la API de entrada acepta tanto números (42.0) como strings con símbolo
// de moneda ("$42.00"). Al serializar JSON devuelve un número (por defecto).
// Implementa interfaces para JSON y para la base de datos (Scanner/Valuer).
type Money float64

func (m *Money) UnmarshalJSON(data []byte) error {
	// Acepta null
	if string(data) == "null" {
		*m = 0
		return nil
	}

	// Intentar decode como número
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*m = Money(num)
		return nil
	}

	// Intentar decode como string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("money: no se pudo parsear json: %w", err)
	}

	// Limpiar simbolo '$' y comas
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")

	if s == "" {
		*m = 0
		return nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("money: formato invalido '%s': %w", s, err)
	}
	*m = Money(f)
	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(m))
}

// driver.Valuer
func (m Money) Value() (driver.Value, error) {
	return float64(m), nil
}

// sql.Scanner
func (m *Money) Scan(src interface{}) error {
	if src == nil {
		*m = 0
		return nil
	}
	switch v := src.(type) {
	case float64:
		*m = Money(v)
		return nil
	case []byte:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return err
		}
		*m = Money(f)
		return nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*m = Money(f)
		return nil
	default:
		return errors.New("money: tipo de dato inesperado")
	}
}

type Item struct {
	ID         uuid.UUID `json:"id"`
	Ticker     string    `json:"ticker"`
	TargetFrom Money     `json:"target_from"`
	TargetTo   Money     `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateItemRequest define la estructura esperada al crear un item vía API.
// Los tags `binding:"required"` son usados por Gin para validar la petición.
type CreateItemRequest struct {
	Ticker     string    `json:"ticker" binding:"required"`
	TargetFrom Money     `json:"target_from" binding:"required"`
	TargetTo   Money     `json:"target_to" binding:"required"`
	Company    string    `json:"company" binding:"required"`
	Action     string    `json:"action" binding:"required"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from" binding:"required"`
	RatingTo   string    `json:"rating_to" binding:"required"`
	Time       time.Time `json:"time" binding:"required"`
}
