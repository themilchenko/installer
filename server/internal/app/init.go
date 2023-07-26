package app

import (
	"github.com/labstack/echo/v4"

	"server/internal/domain"
	"server/internal/sender/delivery"
	sender "server/internal/sender/usecase"
	"server/pkg/logger"
)

type Server struct {
	Echo *echo.Echo

	senderHandler *httpSender.Handler

	senderUsecase domain.SenderUsecase
}

func New(echo *echo.Echo) *Server {
	return &Server{
		Echo: echo,
	}
}

func (s *Server) Start() error {
	if err := s.init(); err != nil {
		return err
	}
	return s.Echo.Start("localhost:8080")
}

func (s *Server) init() error {
	s.makeEchoLogger()
	s.makeUsecases()
	s.makeHandlers()
	s.makeRouter()

	return nil
}

func (s *Server) makeHandlers() {
	s.senderHandler = httpSender.NewSenderHandler(s.senderUsecase)
}

func (s *Server) makeUsecases() {
	s.senderUsecase = sender.NewSenderUsecase()
}

func (s *Server) makeRouter() {
	v := s.Echo.Group("/api")
	v.GET("/download", s.senderHandler.GetFile)
}

func (s *Server) makeEchoLogger() {
	s.Echo.Logger = logger.GetInstance()
	s.Echo.Logger.SetLevel(logger.ToLevel("info"))
	s.Echo.Logger.Info("server started")
}
