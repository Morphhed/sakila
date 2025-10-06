package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "github.com/Morphhed/sakila-go-api/db/sqlc"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// setupDB membuka koneksi ke MySQL Docker
func setupDB() *sql.DB {
	dsn := "root:123@tcp(127.0.0.1:3306)/sakila?parseTime=true&multiStatements=true"
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Gagal konek ke database:", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("❌ Tidak bisa ping ke database:", err)
	}
	fmt.Println("✅ Terhubung ke database Sakila")
	return conn
}

func main() {
	conn := setupDB()
	defer conn.Close()
	queries := db.New(conn)

	r := gin.Default()

	// ==========================================
	//               ACTOR CRUD
	// ==========================================

	// GET /actors
	r.GET("/actors", func(c *gin.Context) {
		actors, err := queries.ListActors(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, actors)
	})

	// GET /actors/:id
	r.GET("/actors/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actor id"})
			return
		}
		actor, err := queries.GetActor(context.Background(), uint16(id)) // ✅ pakai uint16
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "actor not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, actor)
	})

	// PUT /actors/:id
	r.PUT("/actors/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actor id"})
			return
		}

		var req struct {
			FirstName string `json:"first_name" binding:"required"`
			LastName  string `json:"last_name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = queries.UpdateActor(context.Background(), db.UpdateActorParams{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			ActorID:   uint16(id), // ✅ pakai uint16
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		actor, err := queries.GetActor(context.Background(), uint16(id)) // ✅ pakai uint16
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, actor)
	})

	// DELETE /actors/:id
	r.DELETE("/actors/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actor id"})
			return
		}
		err = queries.DeleteActor(context.Background(), uint16(id)) // ✅ pakai uint16
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "actor deleted"})
	})

	// ==========================================
	//               COUNTRY CRUD
	// ==========================================

	r.GET("/countries", func(c *gin.Context) {
		countries, err := queries.ListCountries(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, countries)
	})

	r.GET("/countries/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
			return
		}
		country, err := queries.GetCountry(context.Background(), uint16(id))
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "country not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, country)
	})

	r.POST("/countries", func(c *gin.Context) {
		var req struct {
			Country string `json:"country" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := queries.CreateCountry(context.Background(), req.Country)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		country, err := queries.GetCountry(context.Background(), uint16(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, country)
	})

	r.PUT("/countries/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
			return
		}
		var req struct {
			Country string `json:"country" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = queries.UpdateCountry(context.Background(), db.UpdateCountryParams{
			Country:   req.Country,
			CountryID: uint16(id),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		country, err := queries.GetCountry(context.Background(), uint16(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, country)
	})

	r.DELETE("/countries/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
			return
		}
		err = queries.DeleteCountry(context.Background(), uint16(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "country deleted"})
	})

	// ==========================================
	//                 CITY CRUD
	// ==========================================

	r.GET("/cities", func(c *gin.Context) {
		cities, err := queries.ListCities(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cities)
	})

	r.GET("/cities/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid city id"})
			return
		}
		city, err := queries.GetCity(context.Background(), uint16(id))
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "city not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, city)
	})

	r.POST("/cities", func(c *gin.Context) {
		var req struct {
			City      string `json:"city" binding:"required"`
			CountryID uint16 `json:"country_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := queries.CreateCity(context.Background(), db.CreateCityParams{
			City:      req.City,
			CountryID: req.CountryID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		city, err := queries.GetCity(context.Background(), uint16(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, city)
	})

	r.PUT("/cities/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid city id"})
			return
		}
		var req struct {
			City      string `json:"city" binding:"required"`
			CountryID uint16 `json:"country_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = queries.UpdateCity(context.Background(), db.UpdateCityParams{
			City:      req.City,
			CountryID: req.CountryID,
			CityID:    uint16(id),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		city, err := queries.GetCity(context.Background(), uint16(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, city)
	})

	r.DELETE("/cities/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid city id"})
			return
		}
		err = queries.DeleteCity(context.Background(), uint16(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "city deleted"})
	})

	// Filter city by country
	r.GET("/countries/:id/cities", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
			return
		}
		cities, err := queries.ListCitiesByCountry(context.Background(), uint16(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cities)
	})

	// ==========================================
	// Jalankan server
	// ==========================================
	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	r.Run(":8080")
}
