//recover is originated from echo/v4/middleware/recover. Copyright (c) 2017 LabStack with MIT license.

package echologrus

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		Skipper           middleware.Skipper
		StackSize         int  `yaml:"stack_size"`
		DisableStackAll   bool `yaml:"disable_stack_all"`
		DisablePrintStack bool `yaml:"disable_print_stack"`
		//above options are same with original recover's
		Logger EchoLogger
	}
)

// DefaultRecoverConfig returns default recover config with giving echologrus struct.
func DefaultRecoverConfig(Logger EchoLogger) RecoverConfig {
	return RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		Logger:            Logger,
	}
}

// Recover returns a middleware which recovers from panics anywhere in the chain and writes log with given logger.
// Others are same with original's.
func Recover(Logger EchoLogger) echo.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig(Logger))
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig(EchoLogger{nil}).Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig(EchoLogger{nil}).StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						config.Logger.Logger.WithFields(map[string]interface{}{
							"error": err,
							"stack": stack[:length],
						}).Info("Panic recovered")
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
