package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AntonBraer/urlShorter/internal/usecase"
	"github.com/AntonBraer/urlShorter/pkg/logger"
)

type linkRoutes struct {
	t usecase.Link
	l logger.Interface
}

func init() {
	prometheus.MustRegister(LinkVisited)
}

func newLinkRoutes(handler *gin.RouterGroup, t usecase.Link, l logger.Interface) {
	r := &linkRoutes{t, l}

	h := handler.Group("/link")
	{
		h.POST("/add", r.add)
		h.GET("/:hash", r.redirect)
	}
}

type addLinkRequest struct {
	ToLink string `json:"to_link" binding:"required"`
}

func (r *linkRoutes) add(c *gin.Context) {
	var (
		request addLinkRequest
		err     error
	)
	if err = c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - add")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	var hash string
	if hash, err = r.t.Add(c.Request.Context(), request.ToLink); err != nil {
		r.l.Error(err, "http - v1 - add")
		errorResponse(c, http.StatusInternalServerError, "link service problems")
		return
	}

	c.JSON(http.StatusOK, map[string]string{"Hash": hash})
}

type redirectRequest struct {
	Hash string `uri:"hash" binding:"required"`
}

func (r *linkRoutes) redirect(c *gin.Context) {
	var (
		request redirectRequest
		err     error
	)

	if err = c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - redirect")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	LinkVisited.WithLabelValues(request.Hash).Inc()
	var toLink string
	if toLink, err = r.t.GetLink(c.Request.Context(), request.Hash); err != nil {
		r.l.Error(err, "http - v1 - GetLink")
		errorResponse(c, http.StatusInternalServerError, "link service problems")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, toLink)
}

var LinkVisited = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "linkVisited",
		Help: "How many people visited this link",
	}, []string{"link"})

//count by(link) (rate(link_visited[5m]))
