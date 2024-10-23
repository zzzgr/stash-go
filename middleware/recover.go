package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	v1 "stash-go/api"
)

// RecoverMiddleware 自定义 recover 中间件
func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if rec := recover(); rec != nil {
				// 记录错误信息
				log.Printf("panic: %v", rec)

				// 返回自定义格式的响应
				switch v := rec.(type) {
				case string:
					{
						_ = v1.HandleError(c, 200, errors.New(v), nil)
					}
				default:
					{
						_ = v1.HandleError(c, 200, errors.New("服务器内部错误，请稍后再试。"), nil)
					}
				}

			}
		}()
		return c.Next() // 继续执行后续中间件或处理程序
	}
}
