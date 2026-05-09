package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func TracingMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer("wallet-service")

	return func(c *gin.Context) {
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		ctx, span := tracer.Start(ctx, c.FullPath())
		defer func() {
			span.SetAttributes(attribute.String("http.method", c.Request.Method))
			span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
			span.End()
		}()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("error", c.Errors.String()))
		}
	}
}