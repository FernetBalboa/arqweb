package apierror

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

// CatchAPIErrors - Middleware to handle API errors.
func CatchAPIErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var errMsg string
			var err error = nil
			var clientError ApiError

			//If recovering from panic
			if panicErr := recover(); panicErr != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				panicInfo := fmt.Sprintf("[Recovery] panic recovered:\nRequest:\n%sError:%s\n", string(httpRequest), panicErr)

				err = NewInternalServerApiError(panicInfo, nil)
				errMsg = err.(ApiError).Message()

				log.Errorf("Recovering from panic:\n %s", errMsg)

				clientError = ToApiError(err)
				c.AbortWithStatusJSON(clientError.Status(), clientError)

			} else if len(c.Errors) > 0 {
				//If an error was registered in the context
				err = c.Errors[len(c.Errors)-1].Err //Get the last recorded error
				errMsg = err.Error()

				log.WithFields(log.Fields{
					"request": parseRequest(c.Request),
				}).Error(err)

				clientError := ToApiError(err)
				c.JSON(clientError.Status(), clientError)
			}

		}()
		c.Next()
	}

}

func ToApiError(err error) ApiError {
	var apiErr ApiError
	if werr, ok := err.(wrappingError); ok {
		switch GetType(werr) {
		case BadRequest:
			apiErr = NewApiError("There was a problem with the request",
				"bad_request", http.StatusBadRequest, CauseList{werr.Error()})
		case Forbidden:
			apiErr = NewApiError("There was a problem with the request",
				"forbidden", http.StatusForbidden, CauseList{werr.Error()})
		case NotFound:
			apiErr = NewApiError("There was a problem with the request",
				"not_found", http.StatusNotFound, CauseList{werr.Error()})
		default:
			//Do not send error to client unless it's a public error
			apiErr = NewInternalServerApiError("There was a problem processing the request", nil)
		}
	} else {
		//Do not send error to client unless it's a public error
		apiErr = NewInternalServerApiError("There was a problem processing the request", nil)
	}
	return apiErr
}

func parseRequest(requestHTTP *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("[Method: %s] [URI: %v] [Protocol: %s]", requestHTTP.Method, requestHTTP.URL, requestHTTP.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("	[Host: %v]", requestHTTP.Host))
	// Loop through headers
	for name, headers := range requestHTTP.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("	[%v: %v]", name, h))
		}
	}
	request = append(request, "\n")

	if requestHTTP.Method == "POST" || requestHTTP.Method == "PUT" {
		// Buffer the body
		bodyBuffer, _ := ioutil.ReadAll(requestHTTP.Body)
		//Usually, this buffer is empty because it was already read
		if len(bodyBuffer) > 0 {
			request = append(request, "	[Body: "+string(bodyBuffer)+"]")
		}
	}

	return strings.Join(request, "\n")
}
