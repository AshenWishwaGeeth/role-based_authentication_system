package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtSecret []byte

type User struct {
	ID        string    `json:"id"` // changed from int to string
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PasswordHash string `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// Load .env file automatically
	_ = godotenv.Load()

	supabaseDBURL := os.Getenv("SUPABASE_DB_URL")
	jwtSecretEnv := os.Getenv("JWT_SECRET")

	if supabaseDBURL == "" || jwtSecretEnv == "" {
		panic("SUPABASE_DB_URL or JWT_SECRET is not set in environment")
	}
	jwtSecret = []byte(jwtSecretEnv)

	var err error
	// Connect using the proper PostgreSQL connection string
	db, err = sql.Open("postgres", supabaseDBURL)
	if err != nil {
		panic(fmt.Sprintf("Cannot open DB: %v", err))
	}
	defer db.Close()

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to DB: %v", err))
	}
	fmt.Println("✓ Connected to Supabase Postgres successfully!")

	// Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/me", authMiddleware(), meHandler)
	r.GET("/admin/users", authMiddleware("admin"), adminUsersHandler)

	// Start server
	fmt.Println("✓ Server starting on :8080")
	r.Run(":8080")
}

// Registration handler
func registerHandler(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and password are required"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hash error"})
		return
	}

	_, err = db.Exec(
		`INSERT INTO users (name, email, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5)`,
		req.Name, req.Email, string(hash), req.Role, time.Now(),
	)
	if err != nil {
		// Check for duplicate email
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

// Login handler
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
		fmt.Println("Login DB error:", err) // Debug log
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials", "details": err.Error()})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		fmt.Println("Password mismatch for user:", req.Email) // Debug log
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials", "details": "Password mismatch"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID, // string
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Auth middleware
func authMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"].(string)) // changed from int(claims["user_id"].(float64))
		c.Set("role", claims["role"].(string))

		if len(roles) > 0 {
			userRole := claims["role"].(string)
			allowed := false
			for _, r := range roles {
				if r == userRole {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
				return
			}
		}

		c.Next()
	}
}

// Get current user
func meHandler(c *gin.Context) {
	userID := c.GetString("user_id") // changed from GetInt to GetString
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
	rows, err := db.Query(`SELECT id, name, email, role, created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error", "details": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, users)
}