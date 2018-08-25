package rendering

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/francoishill/gomponents/clienterror"
)

type Service interface {
	RenderError(w http.ResponseWriter, r *http.Request, err error, logFields map[string]interface{}, defaultStatus int)
}

func ChiService() *chiService { return &chiService{} }

type chiService struct{}

func (*chiService) RenderError(w http.ResponseWriter, r *http.Request, err error, logFields map[string]interface{}, defaultStatus int) {
	logger := logrus.NewEntry(logrus.StandardLogger())
	if logFields != nil {
		logger = logger.WithFields(logFields)
	}

	userMsg := err.Error()
	logger.WithError(err).Error(userMsg)

	if errWithStatus, ok := err.(clienterror.Error); ok {
		w.WriteHeader(errWithStatus.Status())
	} else {
		w.WriteHeader(defaultStatus)
	}

	render.JSON(w, r, map[string]interface{}{
		"Error": userMsg,
	})
}
