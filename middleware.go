package echologrus

import (
	"io"
	"strconv"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

// New create new EchoLogger struct with new logrus struct.
func New() EchoLogger {
	return EchoLogger{logrus.New()}
}

// Attach create and attach EchoLogger to echo, use for simple code.
func Attach(e *echo.Echo) EchoLogger {
	el := New()
	e.Logger = el
	e.Use(el.Hook())
	e.Use(Recover(el))
	/*
		e.HTTPErrorHandler = func(err error, c echo.Context) {
			el.HTTPErrorHandler(err, c)
			e.DefaultHTTPErrorHandler(err, c)
		}
	*/
	return el
}

/*
EchoLogger is logrus wrapper for satisfying echo.Logger interface.
Set this to echo.Logger for using logrus logger in echo requset logging.
Original logger implementation : https://pkg.go.dev/github.com/labstack/gommon@v0.3.0/log?tab=doc#Logger
Interface definition : https://pkg.go.dev/github.com/labstack/echo/v4?tab=doc#Logger
*/
type EchoLogger struct {
	*logrus.Logger
}

// Level returns logger level
func (l EchoLogger) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		return log.OFF //original implemetation panics in this case, but is panicking good choice in log function?
	}
}

// HTTPErrorHandler log HTTP error through logrus and pass error to echo's default error handler.
func (l EchoLogger) HTTPErrorHandler(err error, cxt echo.Context) {
	l.Logger.WithFields(map[string]interface{}{
		"time_rfc3339": time.Now().Format(time.RFC3339),
		"context":      cxt,
		"error":        err,
	}).Error("Handled HTTP error")
	cxt.Echo().DefaultHTTPErrorHandler(err, cxt)
}

/*
It's empty because actual header function controlled in logrus.
So, just defined for satisfying interface.
*/
func (l EchoLogger) SetHeader(_ string) {}

/*
It's empty because actual prefix function controlled in logrus.
So, just defined for satisfying interface.
*/
func (l EchoLogger) SetPrefix(s string) {}

/*
It's empty because actual prefix function controlled in logrus.
So, just defined for satisfying interface and return just empty string.
*/
func (l EchoLogger) Prefix() string {
	return ""
}

// SetLevel set level to logger from given log.Lvl
func (l EchoLogger) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		l.Logger.SetLevel(logrus.DebugLevel)
	case log.WARN:
		l.Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		l.Logger.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		l.Logger.SetLevel(logrus.InfoLevel)
	default:
		l.Logger.SetLevel(logrus.TraceLevel) //Same reason with Level()
	}
}

// Output logger output func
func (l EchoLogger) Output() io.Writer {
	return l.Out
}

// SetOutput change output, default os.Stdout
func (l EchoLogger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

// Printj print json log
func (l EchoLogger) Printj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Print()
}

// Debugj debug json log
func (l EchoLogger) Debugj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Debug()
}

// Infoj info json log
func (l EchoLogger) Infoj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Info()
}

// Warnj warning json log
func (l EchoLogger) Warnj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Warn()
}

// Errorj error json log
func (l EchoLogger) Errorj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Error()
}

// Fatalj fatal json log
func (l EchoLogger) Fatalj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Fatal()
}

// Panicj panic json log
func (l EchoLogger) Panicj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Panic()
}

// Print string log
func (l EchoLogger) Print(i ...interface{}) {
	l.Logger.Print(i...)
}

// Debug string log
func (l EchoLogger) Debug(i ...interface{}) {
	l.Logger.Debug(i...)
}

// Info string log
func (l EchoLogger) Info(i ...interface{}) {
	l.Logger.Info(i...)
}

// Warn string log
func (l EchoLogger) Warn(i ...interface{}) {
	l.Logger.Warn(i...)
}

// Error string log
func (l EchoLogger) Error(i ...interface{}) {
	l.Logger.Error(i...)
}

// Fatal string log
func (l EchoLogger) Fatal(i ...interface{}) {
	l.Logger.Fatal(i...)
}

// Panic string log
func (l EchoLogger) Panic(i ...interface{}) {
	l.Logger.Panic(i...)
}

// handler for real part of printing log
func (l EchoLogger) handler(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()

	start := time.Now()
	if err := next(c); err != nil {
		c.Error(err)
	}
	stop := time.Now()

	p := req.URL.Path
	bytesIn := req.Header.Get(echo.HeaderContentLength)
	l.Logger.WithFields(map[string]interface{}{
		"time_rfc3339":  time.Now().Format(time.RFC3339),
		"remote_ip":     c.RealIP(),
		"host":          req.Host,
		"uri":           req.RequestURI,
		"method":        req.Method,
		"path":          p,
		"referer":       req.Referer(),
		"user_agent":    req.UserAgent(),
		"status":        res.Status,
		"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
		"latency_human": stop.Sub(start).String(),
		"bytes_in":      bytesIn,
		"bytes_out":     strconv.FormatInt(res.Size, 10),
	}).Info("Handled request")

	return nil
}

// Hook is a function to process middleware.
func (l EchoLogger) Hook() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return l.handler(c, next)
		}
	}
}
