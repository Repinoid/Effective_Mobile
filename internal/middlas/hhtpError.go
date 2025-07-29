package middlas

import (
	"emobile/internal/models"
	"net/http"
	"time"
)

// ErrorLoggerMiddleware - middleware для логирования ошибок HTTP
func ErrorLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем обертку для ResponseWriter, чтобы перехватывать ошибки
		lrw := &logErrResponseWriter{
			ResponseWriter: w,
			request:        r,
		}

		// Вызываем следующий обработчик
		next.ServeHTTP(lrw, r)
	})
}

// logErrResponseWriter - обертка для http.ResponseWriter для перехвата ошибок
type logErrResponseWriter struct {
	http.ResponseWriter
	request      *http.Request
	errorMessage string
	statusCode   int
}

// WriteHeader перехватывает статус код ответа
func (lrw *logErrResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write перехватывает запись тела ответа
func (lrw *logErrResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode >= 400 {
		lrw.errorMessage = string(b)
	}
	return lrw.ResponseWriter.Write(b)
}

// logError логирует ошибку
func (lrw *logErrResponseWriter) logError() {
	if lrw.statusCode >= 400 && lrw.errorMessage != "" {

		models.Logger.Debug("Ошибка вышла",
			"Time", time.Now().Format(time.RFC3339),
			"Method", lrw.request.Method,
			"uri", lrw.request.URL.Path,
			"code", lrw.statusCode,
			"message", lrw.errorMessage,
			"RemoteAddr", lrw.request.RemoteAddr,
			"UserAgent", lrw.request.UserAgent(),
		)

	}
}

// Ensure logging happens when the response is finished
func (lrw *logErrResponseWriter) finish() {
	lrw.logError()
}
