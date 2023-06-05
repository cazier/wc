package exceptions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoResultsFoundError struct {
	Line int
	Col  int
}

func (e *NoResultsFoundError) Error() string {
	return Message(e)
}

type InvalidValueError struct {
	Line int
	Col  int
}

func (e *InvalidValueError) Error() string {
	return Message(e)
}

func JsonResponse(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err)

		switch err.(type) {
		case *strconv.NumError:
			c.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				gin.H{"error": Message(err)},
			)
		case *NoResultsFoundError:
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": Message(err)},
			)
		case *InvalidValueError:
			c.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				gin.H{"error": Message(err)},
			)
		default:
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": Message(err)},
			)
		}
	}
	return err != nil
}

func Message(err error) string {
	switch err.(type) {
	case *strconv.NumError, *InvalidValueError:
		return "the URI parameter was invalid, and could not be parsed"
	case *NoResultsFoundError:
		return "no matching items could be found"
	default:
		return "an unknown error occurred; please try again"
	}
}
