package laisky_blog_graphql

import (
	"net/http"

	ginMiddlewares "github.com/Laisky/go-utils/gin-middlewares"

	"github.com/gin-gonic/gin"

	"github.com/99designs/gqlgen/handler"

	utils "github.com/Laisky/go-utils"
	"github.com/Laisky/zap"
)

var (
	server = gin.New()
	auth   *ginMiddlewares.Auth
)

func setupAuth() (err error) {
	auth, err = ginMiddlewares.NewAuth([]byte(utils.Settings.GetString("settings.secret")))
	return
}

func RunServer(addr string) {
	server.Use(gin.Recovery())
	if !utils.Settings.GetBool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := setupAuth(); err != nil {
		utils.Logger.Panic("try to setup auth got error", zap.Error(err))
	}
	server.Use(LoggerMiddleware)
	ginMiddlewares.EnableMetric(server)

	server.Any("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})

	server.Any("/ui/", ginMiddlewares.FromStd(handler.Playground("GraphQL playground", "/graphql/query/")))
	server.Any("/query/", ginMiddlewares.FromStd(handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{}}))))

	utils.Logger.Info("listening on http", zap.String("addr", addr))
	utils.Logger.Panic("httpServer exit", zap.Error(server.Run(addr)))
}

func LoggerMiddleware(ctx *gin.Context) {
	utils.Logger.Debug("request",
		zap.String("path", ctx.Request.RequestURI),
		zap.String("method", ctx.Request.Method),
	)
	ctx.Next()
}
