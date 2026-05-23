package sisr

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/Alia5/SISR/api"
	"github.com/Alia5/SISR/api/handler"
	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/logging"
	"github.com/Alia5/SISR/meta"
	"github.com/Alia5/SISR/middleware"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/cors"
)

func (s *SISR) runAPIServer(
	window *sdl.Window,
	wv webview.WebView,
	deviceStore input.DeviceStore,
	stopFn context.CancelFunc,
) (*http.Server, string) {
	l, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		slog.Error("Failed to listen", "addr", s.ListenAddress, "err", err)
		os.Exit(1)
	}

	resolvedAddr := l.Addr().String()
	serverURL := "http://" + resolvedAddr
	host, port, err := net.SplitHostPort(resolvedAddr)
	if err == nil && (host == "" || host == "0.0.0.0" || host == "::") {
		serverURL = "http://localhost:" + port
	}

	schemaPrefix := "#/components/schemas/"

	apiMux := http.NewServeMux()
	hAPI := humago.New(apiMux, huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   "SISR API",
				Version: meta.Version,
			},
			Components: &huma.Components{
				Schemas: huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer),
			},
			Servers: []*huma.Server{{
				URL:         serverURL,
				Description: "SISR API",
			}},
		},
		OpenAPIPath:   "/openapi",
		SchemasPath:   "/schemas",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		CreateHooks: []func(huma.Config) huma.Config{
			func(c huma.Config) huma.Config {
				linkTransformer := huma.NewSchemaLinkTransformer(schemaPrefix, c.SchemasPath)
				c.OnAddOperation = append(c.OnAddOperation, linkTransformer.OnAddOperation)
				c.Transformers = append(c.Transformers, linkTransformer.Transform)
				return c
			},
		},
		Transformers: []huma.Transformer{
			func(c huma.Context, _ string, v any) (any, error) {
				if err, is := v.(error); is {
					if sw, ok := c.BodyWriter().(*api.StatusWriter); ok {
						sw.Error = err
					}
				}
				return v, nil
			},
		},
	})

	hAPI.Adapter().Handle(&huma.Operation{
		Method: http.MethodGet,
		Path:   "/docs",
	}, func(ctx huma.Context) {
		ctx.SetHeader("Content-Type", "text/html")
		_, _ = ctx.BodyWriter().Write([]byte(`<!doctype html>
			<html>
			<head>
				<title>SISR API</title>
				<meta name="referrer" content="same-origin" />
				<meta charset="utf-8" />
				<meta
				name="viewport"
				content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<script
				id="api-reference"
				data-url="/openapi.yaml"></script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
			</html>`,
		))
	})

	api.RegisterAPI(hAPI, &handler.Env{
		Window:      window,
		WebView:     wv,
		DeviceStore: deviceStore,
		QuitFn:      stopFn,
	})

	allowedOrigins := slices.Concat(
		[]string{serverURL},
		strings.Split(s.CORSOrigins, ","),
	)
	if s.FrontendAddress != "" {
		allowedOrigins = append(allowedOrigins, s.FrontendAddress)
	}

	apiSrv := http.Server{
		Addr: resolvedAddr,
		Handler: middleware.With(
			apiMux,
			logging.Middleware,
			cors.New(cors.Options{
				AllowedOrigins:   allowedOrigins,
				AllowedMethods:   []string{"HEAD", "GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
			}).Handler,
			middleware.UnregisteredRoute(hAPI),
		),
	}
	if os.Getenv("DEV") == "1" {
		yml, err := hAPI.OpenAPI().YAML()
		if err != nil {
			slog.Error("failed to generate OpenAPI YAML", "err", err)
		}

		err = os.WriteFile("../openapi.yaml", yml, 0644)
		if err != nil {
			slog.Error("failed to write OpenAPI YAML to file", "err", err)
		} else {
			slog.Info("wrote OpenAPI YAML to ../openapi.yaml")
		}
	}

	go func() {
		if err := apiSrv.Serve(l); err != nil && err != http.ErrServerClosed {
			slog.Error("API server failed", "err", err)
			os.Exit(1)
		}
	}()

	slog.Info("Started Server", "addr", apiSrv.Addr, "url", serverURL)
	slog.Info("Docs on", "addr", apiSrv.Addr, "url", serverURL+"/docs")

	return &apiSrv, serverURL

}
