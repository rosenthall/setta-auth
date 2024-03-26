package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

// InterceptorLogger adapts zap logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func UnaryLoggingInterceptor(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Logging before processing the request
		logger.Infof("New Request - Method:%s; Request:%v", info.FullMethod, hashObject(req))

		// Processing the request
		resp, err := handler(ctx, req)

		// Logging after processing the request
		if err != nil {
			logger.Warnf("Error - Method:%s; Error:%v", info.FullMethod, err)
		} else {
			// Logging the hash of the response
			logger.Debugf("Response - Method:%s; Hash:%s", info.FullMethod, hashObject(resp))
		}

		return resp, err
	}
}

// hashObject generates a hash for the request or response object.
// It converts the response to a JSON string and then computes its SHA-256 hash.
func hashObject(resp interface{}) string {
	// Convert the response to a JSON string
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return "error generating hash"
	}

	// Compute the hash
	hash := sha256.Sum256(jsonResp)
	return hex.EncodeToString(hash[:])
}
