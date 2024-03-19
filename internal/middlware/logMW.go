package middlware

import (
	"net/http"
	"time"

	"films_library/pkg/logger"

	"github.com/sirupsen/logrus"
)

type ResponseWriterWrap struct {
	http.ResponseWriter
	Status int
	Length int
}

func (r *ResponseWriterWrap) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseWriterWrap) Write(bytes []byte) (int, error) {
	r.Length = len(bytes)

	return r.ResponseWriter.Write(bytes)
}

type LoggingMiddleware struct {
	log logger.Interface
}

func NewLoggingMiddleware(log logger.Interface) *LoggingMiddleware {
	return &LoggingMiddleware{
		log: log,
	}
}

func (m *LoggingMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wr := &ResponseWriterWrap{
			ResponseWriter: w,
			Status:         200,
		}

		next.ServeHTTP(wr, r)

		status := wr.Status
		length := wr.Length

		logEntry := m.log.WithFields(logrus.Fields{
			"time":       time.Now(),
			"duration":   time.Since(startTime),
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     status,
			"remote-IP":  r.RemoteAddr,
			"byteLen":    length,
			"user-agent": r.UserAgent(),
		})
		switch {
		case status >= http.StatusInternalServerError:
			logEntry.Error("Server Error")
		case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
			logEntry.Warn("Client Error")
		case status >= http.StatusMultipleChoices && status < http.StatusBadRequest:
			logEntry.Info("Redirect")
		case status >= http.StatusOK && status < http.StatusMultipleChoices:
			logEntry.Info("Success")
		default:
			logEntry.Info("Informational")
		}
	}
	return http.HandlerFunc(fn)
}
