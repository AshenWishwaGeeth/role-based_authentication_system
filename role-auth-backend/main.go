package main

import (
	"time"
	"os"
	"fmt"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

var db *sql.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("SUPABASE_URL")+"?sslmode=require")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/me", authMiddleware(), meHandler)
	r.GET("/admin/users", authMiddleware("admin"), adminUsersHandler)

	r.Run(":8080")
}

// Registration
func registerHandler(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // optional, default to "user"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if req.Role == "" {
		req.Role = "user"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hash error"})
		return
	}
	_, err = db.Exec(
		`INSERT INTO users (name, email, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5)`,
		req.Name, req.Email, string(hash), req.Role, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registered"})
}

// Login
func loginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var user User
	err := db.QueryRow(`SELECT id, name, email, password_hash, role, created_at FROM users WHERE email=$1`, req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Middleware
func authMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))
		if len(roles) > 0 {
			role := claims["role"].(string)
			allowed := false
			for _, r := range roles {
				if r == role {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return
			}
		}
		c.Next()
	}
}

// Get logged-in user
func meHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	var user User
	err := db.QueryRow(`SELECT id, name, email, role, created_at FROM users WHERE id=$1`, userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Admin: list all users
func adminUsersHandler(c *gin.Context) {
	rows, err := db.Query(`SELECT id, name, email, role, created_at FROM users`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)
		users = append(users, u)
	}
	c.JSON(http.StatusOK, users)
}
