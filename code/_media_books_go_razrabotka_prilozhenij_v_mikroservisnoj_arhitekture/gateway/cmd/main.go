package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/yuliapopova/book_all/gateway/internal/config"
	accountpb "github.com/yuliapopova/book_all/contracts/account"
	authpb "github.com/yuliapopova/book_all/contracts/auth"
	transactionpb "github.com/yuliapopova/book_all/contracts/transaction"
	"google.golang.org/grpc"
)

func main() {
	logger := config.InitLogger()
	cfg := config.New()

	logger.Info().Msg("gateway service starting")

	accountConn, err := grpc.Dial(cfg.AccountGRPCHost, grpc.WithInsecure())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to account service")
	}
	defer accountConn.Close()

	authConn, err := grpc.Dial(cfg.AuthGRPCHost, grpc.WithInsecure())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to auth service")
	}
	defer authConn.Close()

	transactionConn, err := grpc.Dial(cfg.TransactionGRPCHost, grpc.WithInsecure())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to transaction service")
	}
	defer transactionConn.Close()

	accountClient := accountpb.NewAccountServiceClient(accountConn)
	authClient := authpb.NewAuthServiceClient(authConn)
	transactionClient := transactionpb.NewTransactionServiceClient(transactionConn)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(authMiddleware(cfg.JWTSecret))

	r.POST("/auth/register", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
			Name     string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := authClient.Register(c.Request.Context(), &authpb.RegisterRequest{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := authClient.Login(c.Request.Context(), &authpb.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/auth/refresh", func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := authClient.Refresh(c.Request.Context(), &authpb.RefreshRequest{
			RefreshToken: req.RefreshToken,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/auth/logout", func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := authClient.Logout(c.Request.Context(), &authpb.LogoutRequest{
			RefreshToken: req.RefreshToken,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/auth/me", func(c *gin.Context) {
		accessToken, _ := c.Get("access_token")

		resp, err := authClient.GetCurrentUser(c.Request.Context(), &authpb.GetCurrentUserRequest{
			AccessToken: accessToken.(string),
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/account/deposit", func(c *gin.Context) {
		var req struct {
			UserID uint64  `json:"user_id" binding:"required"`
			Amount float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := accountClient.Deposit(c.Request.Context(), &accountpb.DepositRequest{
			UserId: req.UserID,
			Amount: req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/account/withdraw", func(c *gin.Context) {
		var req struct {
			AccountID uint64  `json:"account_id" binding:"required"`
			Amount    float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := accountClient.Withdraw(c.Request.Context(), &accountpb.WithdrawRequest{
			AccountId: req.AccountID,
			Amount:    req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/account/transfer", func(c *gin.Context) {
		var req struct {
			AccountID   uint64  `json:"account_id" binding:"required"`
			RecipientID uint64  `json:"recipient_id" binding:"required"`
			Amount      float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := accountClient.Transfer(c.Request.Context(), &accountpb.TransferRequest{
			AccountId:   req.AccountID,
			RecipientId: req.RecipientID,
			Amount:      req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/account/balance/:account_id", func(c *gin.Context) {
		var accountID uint64
		fmt.Sscanf(c.Param("account_id"), "%d", &accountID)

		resp, err := accountClient.GetBalance(c.Request.Context(), &accountpb.GetBalanceRequest{
			AccountId: accountID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/transaction/deposit", func(c *gin.Context) {
		var req struct {
			UserID uint64  `json:"user_id" binding:"required"`
			Amount float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := transactionClient.Deposit(c.Request.Context(), &transactionpb.DepositRequest{
			UserId: req.UserID,
			Amount: req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/transaction/withdraw", func(c *gin.Context) {
		var req struct {
			AccountID uint64  `json:"account_id" binding:"required"`
			Amount    float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := transactionClient.Withdraw(c.Request.Context(), &transactionpb.WithdrawRequest{
			AccountId: req.AccountID,
			Amount:    req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/transaction/transfer", func(c *gin.Context) {
		var req struct {
			UserID    uint64  `json:"user_id" binding:"required"`
			Recipient uint64  `json:"recipient" binding:"required"`
			Amount    float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := transactionClient.Transfer(c.Request.Context(), &transactionpb.TransferRequest{
			UserId:    req.UserID,
			Recipient: req.Recipient,
			Amount:    req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/transactions", func(c *gin.Context) {
		var userID uint64
		fmt.Sscanf(c.Query("user_id"), "%d", &userID)

		var limit uint32 = 10
		var offset uint32 = 0
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		if o := c.Query("offset"); o != "" {
			fmt.Sscanf(o, "%d", &offset)
		}

		resp, err := transactionClient.GetTransactions(c.Request.Context(), &transactionpb.GetTransactionsRequest{
			UserId: userID,
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	logger.Info().Str("port", cfg.HTTP_PORT).Msg("HTTP server starting")
	if err := r.Run(fmt.Sprintf(":%s", cfg.HTTP_PORT)); err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}
}

func authMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if c.Request.URL.Path == "/auth/login" || c.Request.URL.Path == "/auth/register" || c.Request.URL.Path == "/auth/refresh" {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString := authHeader[7:]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
			c.Set("access_token", tokenString)
		}

		c.Next()
	}
}

var _ = zerolog.TimeFieldFormat
