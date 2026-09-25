package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"pet-social-game/server/internal/config"
	"pet-social-game/server/internal/handler"
	"pet-social-game/server/internal/middleware"
	"pet-social-game/server/internal/repo"
	"pet-social-game/server/internal/security"
	"pet-social-game/server/internal/service"
	"pet-social-game/server/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	cfgPath := flag.String("config", "configs/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		// try relative to executable / module root
		alt := filepath.Join("server", *cfgPath)
		cfg, err = config.Load(alt)
		if err != nil {
			log.Fatalf("load config: %v", err)
		}
	}

	gin.SetMode(cfg.Server.Mode)
	db, err := repo.OpenMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	rdb := repo.OpenRedis(cfg.Redis)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := rdb.Ping(ctx).Err(); err != nil {
		cancel()
		log.Fatalf("redis: %v", err)
	}
	cancel()

	store := repo.NewStore(db, rdb)
	detector := security.NewDetector(cfg.Security)
	svc := service.New(cfg, store, detector)
	api := handler.New(svc)
	hub := ws.NewHub(rdb, svc, cfg.JWT.Secret)

	runCtx, stop := context.WithCancel(context.Background())
	defer stop()
	go hub.Run(runCtx)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())
	r.GET("/healthz", api.Healthz)
	r.GET("/ws", hub.Handle)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", api.Login)
		authz := v1.Group("")
		authz.Use(middleware.JWT(cfg.JWT.Secret))
		{
			authz.GET("/me", api.Me)
			authz.POST("/pet", api.CreatePet)
			authz.POST("/pet/feed", api.Feed)
			authz.POST("/pet/outfit", api.Outfit)
			authz.GET("/friends", api.ListFriends)
			authz.POST("/friends/request", api.RequestFriend)
			authz.POST("/friends/accept", api.AcceptFriend)
			authz.GET("/visit/:userId", api.VisitHome)
			authz.POST("/visit/:userId/interact", api.VisitInteract)
			authz.GET("/chat/history", api.ChatHistory)
		}
	}

	// Browser playable demo (no Cocos required)
	playDir := resolvePlayDir()
	if playDir != "" {
		r.Static("/play", playDir)
		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/play/")
		})
		log.Printf("playable demo: http://127.0.0.1%s/play/", cfg.Server.Addr)
	}

	srv := &http.Server{Addr: cfg.Server.Addr, Handler: r}
	go func() {
		log.Printf("pet-social api listening on %s", cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	stop()
	shutdownCtx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	_ = srv.Shutdown(shutdownCtx)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func resolvePlayDir() string {
	candidates := []string{
		filepath.Join("..", "play"),
		"play",
		filepath.Join("..", "..", "play"),
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "play"),
			filepath.Join(wd, "..", "play"),
		)
	}
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	return ""
}
