package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"vibe-music-server/internal/config"
	"vibe-music-server/internal/pkg/cache"
	"vibe-music-server/internal/pkg/db"
	_ "vibe-music-server/internal/pkg/validate"
	"vibe-music-server/internal/router"
)

func main() {
	db.Init()
	cache.Init()
	r := router.NewEngine()
	appCfg := config.Get().App
	port := appCfg.Port
	if port == "" {
		if appCfg.SSL {
			port = "8443"
		} else {
			port = "8080"
		}
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		// 可选：统一 TLS 安全配置（即使 SSL=false 也不影响）
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	var err error
	go func() {
		if appCfg.SSL {
			// 启用 HTTPS
			if appCfg.Cert == "" || appCfg.Key == "" {
				log.Fatal("SSL enabled but cert or key path is empty")
			}
			err = server.ListenAndServeTLS(appCfg.Cert, appCfg.Key)
		} else {
			// 启用 HTTP
			err = server.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed to start: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	fmt.Println("Server exiting gracefully")
}
