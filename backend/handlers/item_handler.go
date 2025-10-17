// handlers/user_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"mi-api-rest/models"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// ItemHandler agrupa la dependencia a la base de datos y expone los métodos
// HTTP para manipular los items.
type ItemHandler struct {
	DB *sql.DB
}

// NewItemHandler crea un nuevo handler con la conexión a la base de datos.
func NewItemHandler(db *sql.DB) *ItemHandler {
	return &ItemHandler{DB: db}
}

// CreateItem crea un nuevo item a partir del JSON recibido.
// Espera un payload que cumpla con models.CreateItemRequest (o models.Item) y
// devuelve el item creado con su `id` y `created_at`.
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req models.Item

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.Item
	query := `INSERT INTO items (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at`

	// Ejecutar INSERT y obtener id y created_at
	var createdAt time.Time
	err := h.DB.QueryRow(query, req.Ticker, float64(req.TargetFrom), float64(req.TargetTo), req.Company, req.Action, req.Brokerage, req.RatingFrom, req.RatingTo, req.Time).Scan(
		&item.ID, &createdAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Completar campos para la respuesta
	item.Ticker = req.Ticker
	item.TargetFrom = req.TargetFrom
	item.TargetTo = req.TargetTo
	item.Company = req.Company
	item.Action = req.Action
	item.Brokerage = req.Brokerage
	item.RatingFrom = req.RatingFrom
	item.RatingTo = req.RatingTo
	item.Time = req.Time
	item.CreatedAt = createdAt

	c.JSON(http.StatusCreated, item)
}

// GetItems lista todos los items ordenados por fecha de creación (descendente).
// Responde con un array JSON de models.Item.
func (h *ItemHandler) GetItems(c *gin.Context) {
	query := `SELECT id, ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time, created_at FROM items ORDER BY created_at DESC`

	rows, err := h.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Ticker, &item.TargetFrom, &item.TargetTo, &item.Company, &item.Action, &item.Brokerage, &item.RatingFrom, &item.RatingTo, &item.Time, &item.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

// GetItem devuelve un item por su ID. Si no existe retorna 404.
func (h *ItemHandler) GetItem(c *gin.Context) {
	id := c.Param("id")

	var item models.Item
	query := `SELECT id, ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time, created_at FROM items WHERE id = $1`

	err := h.DB.QueryRow(query, id).Scan(
		&item.ID, &item.Ticker, &item.TargetFrom, &item.TargetTo, &item.Company, &item.Action, &item.Brokerage, &item.RatingFrom, &item.RatingTo, &item.Time, &item.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItems sincroniza los items con la API externa.
// Comportamiento:
//  1. Lee la API key desde la variable de entorno SWE_API_KEY.
//  2. Borra todos los registros de la tabla `items`.
//  3. Consume la API externa paginada (cabecera Authorization: Bearer <key>),
//     siguiendo `next_page` hasta que sea nulo/ vacío.
//  4. Inserta los registros recibidos en la tabla `items`.
//
// Nota: este método borra todos los datos locales antes de insertar los nuevos.
func (h *ItemHandler) UpdateItems(c *gin.Context) {

	apiKey := os.Getenv("SWE_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SWE_API_KEY no definida en variables de entorno"})
		return
	}
	log.Println("api key : ", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	baseURL := "https://api.karenai.click/swechallenge/list"
	nextURL := baseURL

	// Eliminar todos los registros antes de insertar los nuevos. Esto evita
	// la necesidad de hacer upserts y garantiza que la tabla refleje exactamente
	// la información de la API después de la sincronización.
	if _, err := h.DB.Exec("DELETE FROM items"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error borrando registros: " + err.Error()})
		return
	}

	// Inserción de datos paginados
	for {
		// Construir request
		reqHTTP, err := http.NewRequest("GET", nextURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reqHTTP.Header.Set("Authorization", "Bearer "+apiKey)
		reqHTTP.Header.Set("Accept", "application/json")

		resp, err := client.Do(reqHTTP)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		log.Println("Fetched URL: ", nextURL, " Status: ", resp.Status)

		// Estructura esperada de la respuesta externa. `items` contiene los
		// registros y `next_page` un token o URL para la siguiente página.
		var body struct {
			Data     []models.Item `json:"items"`
			NextPage *string       `json:"next_page"`
		}

		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&body); err != nil {
			resp.Body.Close()
			c.JSON(http.StatusBadGateway, gin.H{"error": "error decodificando respuesta externa: " + err.Error()})
			return
		}
		resp.Body.Close()

		// Insertar todos los registros de la página actual.
		for _, it := range body.Data {
			insertQ := `INSERT INTO items (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
			_, err := h.DB.Exec(insertQ, it.Ticker, float64(it.TargetFrom), float64(it.TargetTo), it.Company, it.Action, it.Brokerage, it.RatingFrom, it.RatingTo, it.Time)
			if err != nil {
				// Loguear y continuar con otros registros. No abortamos toda la
				// sincronización por un registro fallido; se registran los errores
				// en el log para análisis posterior.
				log.Println("Error insertando la data: ", err)
				continue
			}
		}

		// Manejar paginación
		if body.NextPage == nil || *body.NextPage == "" {
			break
		}
		// Si next_page es una URL absoluta, usarla; si es solo un token, anexarlo como ?next_page=token
		np := *body.NextPage
		if u, err := url.Parse(np); err == nil && u.Scheme != "" && u.Host != "" {
			nextURL = np
		} else {
			// construir URL con query
			q := url.Values{}
			q.Set("next_page", np)
			nextURL = baseURL + "?" + q.Encode()
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sincronización completada"})
}
