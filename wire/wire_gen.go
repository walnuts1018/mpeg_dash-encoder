// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/walnuts1018/mpeg_dash-encoder/config"
	"github.com/walnuts1018/mpeg_dash-encoder/infra/ffmpeg"
	"github.com/walnuts1018/mpeg_dash-encoder/infra/jwt"
	"github.com/walnuts1018/mpeg_dash-encoder/infra/minio"
	"github.com/walnuts1018/mpeg_dash-encoder/router"
	"github.com/walnuts1018/mpeg_dash-encoder/router/handler"
	"github.com/walnuts1018/mpeg_dash-encoder/router/middleware"
	"github.com/walnuts1018/mpeg_dash-encoder/usecase"
)

// Injectors from wire.go:

func CreateUsecase(ctx context.Context, cfg config.Config) (*usecase.Usecase, error) {
	jwtSigningKey := cfg.JWTSigningKey
	manager := jwt.NewManager(jwtSigningKey)
	ffmpegFFMPEG, err := ffmpeg.NewFFMPEG(cfg)
	if err != nil {
		return nil, err
	}
	sourceClientBucketName := cfg.MinIOSourceUploadBucket
	client, err := minio.NewMinIOClient(cfg)
	if err != nil {
		return nil, err
	}
	sourceClient := minio.NewSourceClient(sourceClientBucketName, client)
	encodedObjectBucketName := cfg.MinIOOutputBucket
	encodedObjectClient := minio.NewEncodedObjectClient(encodedObjectBucketName, client)
	usecaseUsecase, err := usecase.NewUsecase(cfg, manager, ffmpegFFMPEG, sourceClient, encodedObjectClient)
	if err != nil {
		return nil, err
	}
	return usecaseUsecase, nil
}

func CreateRouter(ctx context.Context, cfg config.Config, usecase2 *usecase.Usecase) (*gin.Engine, error) {
	handlerHandler, err := handler.NewHandler(cfg, usecase2)
	if err != nil {
		return nil, err
	}
	adminToken := cfg.AdminToken
	middlewareMiddleware := middleware.NewMiddleware(adminToken, usecase2)
	engine, err := router.NewRouter(cfg, handlerHandler, middlewareMiddleware)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// wire.go:

var jwtSet = wire.NewSet(jwt.NewManager, wire.Bind(new(usecase.TokenIssuer), new(*jwt.Manager)))

var minioSourceClientSet = wire.NewSet(minio.NewSourceClient, wire.Bind(new(usecase.SourceRepository), new(*minio.SourceClient)))

var minioEncodedObjectClientSet = wire.NewSet(minio.NewEncodedObjectClient, wire.Bind(new(usecase.EncodedObjectRepository), new(*minio.EncodedObjectClient)))

var ffmpegSet = wire.NewSet(ffmpeg.NewFFMPEG, wire.Bind(new(usecase.Encoder), new(*ffmpeg.FFMPEG)))

var UsecaseConfigSet = wire.FieldsOf(new(config.Config),
	"JWTSigningKey",
	"MinIOSourceUploadBucket",
	"MinIOOutputBucket",
)

var RouterConfigSet = wire.FieldsOf(new(config.Config),
	"AdminToken",
)
