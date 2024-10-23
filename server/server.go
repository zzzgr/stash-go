package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/swagger"
	activity2 "stash-go/controller/activity"
	pkgController "stash-go/controller/pkg"
	_ "stash-go/docs" // 导入你的文档
	"stash-go/middleware"
	"stash-go/service/activity"
	"stash-go/service/pkg"
)

type Server struct {
	fiber           *fiber.App
	activityService *activity.Service
	pkgService      *pkg.Service

	activityController *activity2.Controller
	pkgController      *pkgController.Pkg
}

func New(
	activityController *activity2.Controller,
	pkgController *pkgController.Pkg,
) *Server {
	return &Server{
		fiber: fiber.New(fiber.Config{
			ReadBufferSize: 512 * 1024,
		}),
		activityController: activityController,
		pkgController:      pkgController,
	}
}

func (s *Server) Run(port int) {
	s.fiber.Use(middleware.RecoverMiddleware())

	// 注册 Swagger 路由
	s.fiber.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 路径

	s.fiber.All("/", func(c *fiber.Ctx) error {
		return c.SendString("hello, stash-go!")
	})

	s.fiber.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return nil
	})

	// activity_dto
	s.fiber.Get("/activity/query", s.activityController.QueryAll)
	s.fiber.Get("/activity", s.activityController.Get)
	s.fiber.Post("/activity", s.activityController.Save)
	s.fiber.Delete("/activity/:id", s.activityController.Delete)

	// pkg
	s.fiber.Get("/pkg", s.pkgController.Query)
	s.fiber.Get("/pkg/clear", s.pkgController.Clear)
	s.fiber.All("/*", s.pkgController.StashPackage) // 存包

	err := s.fiber.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Http Server 启动失败, 请检查 (%s)", err.Error())
	}
}
