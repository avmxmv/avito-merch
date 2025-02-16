package handler

import (
	"avito-merch/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	authService service.AuthService
	buyService  service.BuyService
	infoService service.InfoService
	sendService service.SendService
	logger      *zap.Logger
}

func NewHandler(
	authService service.AuthService,
	buyService service.BuyService,
	infoService service.InfoService,
	sendService service.SendService,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		authService: authService,
		buyService:  buyService,
		infoService: infoService,
		sendService: sendService,
		logger:      logger,
	}
}

func (h *Handler) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)
		c.Next()
	}
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		userID, err := h.authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.Use(h.LoggingMiddleware(), gin.Recovery())

	router.POST("/api/auth", h.Login)

	authGroup := router.Group("/").Use(h.AuthMiddleware())
	{
		authGroup.GET("/api/info", h.GetUserInfo)
		authGroup.POST("/api/sendCoin", h.SendCoins)
		authGroup.GET("/api/buy/:item", h.BuyItem)
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	info, err := h.infoService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *Handler) SendCoins(c *gin.Context) {
	var req struct {
		ToUser string `json:"toUser" binding:"required"`
		Amount int    `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromUserID := c.MustGet("userID").(int)
	err := h.sendService.SendCoins(c.Request.Context(), fromUserID, req.ToUser, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) BuyItem(c *gin.Context) {
	item := c.Param("item")
	userID := c.MustGet("userID").(int)

	err := h.buyService.BuyItem(c.Request.Context(), userID, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
