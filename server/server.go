package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	v1 "stash-go/api"
	"stash-go/const/cache"
	"stash-go/httpx"
	"stash-go/jsonx"
	"stash-go/model"
	"stash-go/redis"
	"time"
)

func Run(port int) {
	app := fiber.New()

	app.All("/", func(c *fiber.Ctx) error {
		return c.SendString("hello, stash-go!")
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/get", func(c *fiber.Ctx) error {
		list := redis.List("stash:package")
		res := make([]*model.Package, 0)
		for _, item := range list {
			m, err := jsonx.To(item, &model.Package{})
			if err != nil {
				continue
			}
			res = append(res, m)
		}
		return v1.HandleSuccess(c, res)
	})

	app.Get("/clear", func(c *fiber.Ctx) error {
		err := redis.Del(cache.PACKAGE)
		if err != nil {
			return v1.HandleError(c, 200, err, nil)
		}
		return v1.HandleSuccess(c, nil)
	})

	app.All("/*", func(c *fiber.Ctx) error {

		uri := c.Request().URI()
		method := c.Method()
		headers := getHeaders(c.GetReqHeaders())
		queries := c.Queries()
		bodyStr := string(c.Body())
		protocol := c.Protocol()
		host := string(uri.Host())
		path := string(uri.Path())
		url := protocol + "://" + host + path

		// options 真实请求
		if method == "OPTIONS" {
			r, err := httpx.RestyClient.R().
				SetHeaders(headers).
				SetBody(bodyStr).
				SetQueryParams(queries).
				Options(url)
			if err != nil {
				log.Info("OPTIONS forward err: ", err.Error())
				return err
			}

			response := c.Response()
			response.SetStatusCode(r.StatusCode())
			response.SetBody(r.Body())
			for k, v := range r.Header() {
				response.Header.Set(k, v[0])
			}
		}

		// 存包
		p := model.Package{
			Url:     url,
			Method:  method,
			Headers: headers,
			Queries: queries,
			Body:    bodyStr,
		}

		_ = redis.LPush(cache.PACKAGE, jsonx.ToStr(p), time.Hour*24)

		return nil
	})

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Http Server 启动失败, 请检查 (%s)", err.Error())
	}
}

func getHeaders(headers map[string][]string) map[string]string {
	m := make(map[string]string)
	for k, v := range headers {
		m[k] = v[0]
	}
	return m
}
