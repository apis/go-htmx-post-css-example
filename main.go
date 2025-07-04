package main

import (
	"htmx-example/internal/pkg/config"
	"htmx-example/internal/pkg/staticAssets"
	"htmx-example/internal/pkg/storage"
	"htmx-example/internal/pkg/viewModels"
	"htmx-example/internal/pkg/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Windows configuration examples
// cmd /V /C "set HTMX_APP_PORT=321&& htmx-example.exe"
// htmx-example.exe --Port 123

// Linux configuration examples
// HTMX_APP_PORT=321 ./htmx-example
// ./htmx-example --Port 123

const applicationName = "htmx-example"
const serverShutdownTimeout = 5 * time.Second
const embedFsRoot = "web/frontend/dist"
const templatesDir = "web/templates"
const uiUrlPrefix = "/ui"
const defaultUiUrl = "index.html"

var (
	//go:embed all:web/templates/*
	templateFS embed.FS

	//go:embed web/frontend/dist
	embedFs embed.FS
)

func main() {
	setupZerolog()

	log.Info().Msg("Parsing configuration")
	appConfig := &applicationConfig{}
	config.Parse(appConfig, applicationName)

	log.Info().Msg("Starting up")

	listener := createNetListener(appConfig)
	server := startHttpServer(listener, appConfig.SimulatedDelay)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	log.Info().Msg("Application stopping")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server.Shutdown failed")
	}

	log.Info().Msg("Application stopped")
}

func startHttpServer(listener net.Listener, simulatedDelay int) *http.Server {
	templates, err := web.TemplateParseFSRecursive(templateFS, templatesDir, ".html", nil)
	if err != nil {
		log.Panic().Err(err).Msg("template parsing failed")
	}

	jsonStorage := storage.NewJsonStorage("./data/companies.json")
	companiesViewModel := viewModels.NewCompaniesViewModel(templates, jsonStorage)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	httpLogger := httplog.NewLogger("backend-api", httplog.Options{
		LogLevel: slog.LevelDebug,
		JSON:     true,
		Concise:  true,
		//RequestHeaders:   true,
		//ResponseHeaders:  true,
	})

	router := chi.NewRouter()
	router.Use(httplog.RequestLogger(httpLogger))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*",
			"http://*",
		},
	}))

	router.Handle(uiUrlPrefix+"*", staticAssets.Handler(embedFs, embedFsRoot, uiUrlPrefix, defaultUiUrl))

	router.Handle("GET /company/add", web.Handler{Request: companiesViewModel.AddCompany,
		SimulatedDelay: simulatedDelay})
	router.Handle("POST /company", web.Handler{Request: companiesViewModel.SaveNewCompany,
		SimulatedDelay: simulatedDelay})
	router.Handle("GET /company", web.Handler{Request: companiesViewModel.CancelSaveNewCompany,
		SimulatedDelay: simulatedDelay})
	router.Handle("GET /company/edit/{id}", web.Handler{Request: companiesViewModel.EditCompany,
		SimulatedDelay: simulatedDelay})
	router.Handle("PUT /company/{id}", web.Handler{Request: companiesViewModel.SaveExistingCompany,
		SimulatedDelay: simulatedDelay})
	router.Handle("GET /company/{id}", web.Handler{Request: companiesViewModel.CancelSaveExistingCompany,
		SimulatedDelay: simulatedDelay})
	router.Handle("DELETE /company/{id}", web.Handler{Request: companiesViewModel.DeleteCompany,
		SimulatedDelay: simulatedDelay})

	router.Handle("GET /companies", web.Handler{Request: companiesViewModel.Index,
		SimulatedDelay: simulatedDelay})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, uiUrlPrefix, http.StatusPermanentRedirect)
	})

	server := &http.Server{
		Handler: router,
	}

	go func() {
		log.Info().Msg("Server is about to start")

		err := server.Serve(listener)
		if err != nil {
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("server.ListenAndServe failed")
			}
		}

		log.Info().Msg("Server stopped")
	}()
	return server
}

func createNetListener(appConfig *applicationConfig) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("net.Listen failed")
	}

	return listener
}

func setupZerolog() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = zerolog.New(os.Stderr).
		With().
		Timestamp().
		//Caller().
		Logger()
}
