package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Morphhed/sakila-go-api/auth"
	dbsqlc "github.com/Morphhed/sakila-go-api/db/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// ===============================
// DB CONNECTION
// ===============================
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

// ===============================
// MAIN APP ENTRY
// ===============================
func main() {
	conn := setupDB()
	defer conn.Close()
	q := dbsqlc.New(conn)
	ctx := context.Background()

	r := gin.Default()

	// ==========================================
	//           AUTH ROUTES
	// ==========================================
	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		_, err := q.CreateUser(ctx, dbsqlc.CreateUserParams{
			Username:     req.Username,
			PasswordHash: string(hash),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user created"})
	})

	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		user, err := q.GetUserByUsername(ctx, req.Username)
		if err != nil {
			fmt.Println("Error getting user:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		fmt.Println("User found:", user.Username)

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			return
		}

		token, _ := auth.GenerateToken(user.Username)
		c.JSON(http.StatusOK, gin.H{"token": token})

		fmt.Println("Login attempt:", req.Username)
		fmt.Println("DB hash:", user.PasswordHash)
		fmt.Println("Password input:", req.Password)

	})

	// ==========================================
	//           PROTECTED API ROUTES
	// ==========================================
	api := r.Group("/api", auth.AuthMiddleware())

	// --- ACTOR CRUD ---
	api.GET("/actors", func(c *gin.Context) {
		actors, err := q.ListActors(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, actors)
	})

	api.GET("/actors/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actor id"})
			return
		}
		actor, err := q.GetActor(ctx, uint16(id))
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "actor not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, actor)
	})

	// You can continue with your country/city routes here...
	// (same logic as before, just put them under `api := r.Group("/api", auth.AuthMiddleware())`)

	// ==========================================
	// START SERVER
	// ==========================================
	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	r.Run(":8080")
}
