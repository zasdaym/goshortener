package link

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Server is Link HTTP server.
type Server struct {
	svc     service
	timeout time.Duration
	router  *echo.Echo
}

// ServerOpts is options to create new Link HTTP server.
type ServerOpts struct {
	Svc     service
	Timeout time.Duration
}

// NewServer creates a new Link HTTP server.
func NewServer(opts ServerOpts) *Server {
	srv := &Server{
		svc:     opts.Svc,
		timeout: opts.Timeout,
	}
	router := echo.New()
	router.HideBanner = true
	router.GET("/links", srv.getLinks)
	router.POST("/links", srv.createLink)
	srv.router = router
	return srv
}

// Start starts Link HTTP server.
func (srv *Server) Start(addr string) error {
	return srv.router.Start(addr)
}

func (srv *Server) getLinks(c echo.Context) error {
	var req getLinksRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()
	links, err := srv.svc.getLinks(ctx, req)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, links)
}

func (srv *Server) createLink(c echo.Context) error {
	var req createLinkRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := req.validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), srv.timeout)
	defer cancel()
	if err := srv.svc.createLink(ctx, req); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}
