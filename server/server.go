package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	v1 "stash-go/api"
	"stash-go/config"
	"stash-go/const/cache"
	"stash-go/httpx"
	"stash-go/jsonx"
	"stash-go/model"
	"stash-go/redis"
	"stash-go/util"
	"strconv"
	"strings"
	"time"
)

func Run(port int) {
	app := fiber.New(fiber.Config{
		ReadBufferSize: 512 * 1024,
	})

	app.All("/", func(c *fiber.Ctx) error {
		return c.SendString("hello, stash-go!")
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/get", func(c *fiber.Ctx) error {
		typ, err := strconv.Atoi(c.Query("type", "1"))
		if err != nil {
			typ = 1
		}

		list := redis.List("stash:package")

		switch typ {
		case 1:
			{
				return v1.HandleSuccess(c, handleType1(list))
			}

		case 2:
			{
				return v1.HandleSuccess(c, handleType2(list))
			}
		default:
			return v1.HandleSuccess(c, nil)
		}

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
		host := string(uri.Host())
		path := string(uri.Path())
		url := "https://" + host + path

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
			return nil
		}

		// 存包
		p := model.Package{
			Url:     url,
			Method:  method,
			Host:    host,
			Path:    path,
			Headers: headers,
			Queries: queries,
			Body:    bodyStr,
			Ip:      c.IP(),
		}

		_ = redis.LPush(cache.PACKAGE, jsonx.ToStr(p), time.Hour*24)

		response := c.Response()
		if strings.Index(url, "watsonsvip.com.cn") >= 0 {
			response.Header.Add("access-control-allow-origin", "*")
		}

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

func handleType1(list []string) (res []*model.Package) {
	for _, item := range list {
		m, err := jsonx.To(item, &model.Package{})
		if err != nil {
			continue
		}
		res = append(res, m)
	}
	return res
}

func handleType2(list []string) []*model.NamedPackage {
	m := make(map[string]*model.NamedPackage)
	for _, item := range list {
		p, err := jsonx.To(item, &model.Package{})
		if err != nil {
			continue
		}
		for _, strategy := range config.Conf.SplitStrategy {
			if !util.IsMatch(p.Host, strategy.HostPattern) {
				continue
			}

			field := strategy.Field
			typ := strategy.Type

			v := ""
			if typ == 1 { // cookie
				cookie := p.Headers["Cookie"]
				v = util.GetCookieFieldMap(cookie)[field]
			} else if typ == 2 { // header
				v = p.Headers[field]
			} else {
				continue
			}

			namedPkg, ok := m[v]
			if ok {
				namedPkg.Packages = append(namedPkg.Packages, p)
			} else {

				m[v] = &model.NamedPackage{
					Name:     v,
					Packages: []*model.Package{p},
				}
			}

		}

	}

	r := make([]*model.NamedPackage, 0)
	for _, p := range m {
		r = append(r, p)
	}

	return r
}
