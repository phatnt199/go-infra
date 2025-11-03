package handlers

import (
	"github.com/phatnt199/go-infra/pkg/logger"

	"github.com/phatnt199/go-infra/pkg/adapter/http/httperrors/problemdetails"

	"emperror.dev/errors"
	"github.com/gofiber/fiber/v2"
)

func ProblemDetailErrorHandlerFunc(
	err error,
	c *fiber.Ctx,
	logger logger.Logger,
) error {
	var problem problemdetails.ProblemDetailErr
	logger.Debug("Go here..........: ", problem, errors.As(err, &problem), problem != nil)
	// if error was not problem detail we will convert the error to a problem detail
	if ok := errors.As(err, &problem); !ok {
		logger.Debug("Before: ", err)
		problem = problemdetails.ParseError(err)
		logger.Debug("After: ", problem)
	}

	if problem != nil {
		// Log detailed context for easier tracing
		logger.Error("problem detail error encountered",
			"error", err,
			"problem", problem,
			"status", problem.GetStatus(),
			"method", c.Method(),
			"path", c.Path(),
			"url", c.OriginalURL(),
		)

		// Write problem detail to response
		c.Set(fiber.HeaderContentType, "application/problem+json")
		c.Status(problem.GetStatus())
		if jerr := c.JSON(problem); jerr != nil {
			logger.Error("failed to write problem detail JSON response",
				"json_error", jerr,
				"original_error", err,
				"status", problem.GetStatus(),
				"method", c.Method(),
				"path", c.Path(),
				"url", c.OriginalURL(),
			)
			return jerr
		}
		return nil
	}

	return err
}
