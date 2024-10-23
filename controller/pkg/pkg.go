package activity

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	v1 "stash-go/api"
	"stash-go/model/entity"
	"stash-go/service/activity"
	"stash-go/service/pkg"
	"stash-go/util/httpx"
	"strconv"
)

type Pkg struct {
	activityService *activity.Service
	pkgService      *pkg.Service
}

func New(activityService *activity.Service, pkgService *pkg.Service) *Pkg {
	return &Pkg{
		activityService: activityService,
		pkgService:      pkgService,
	}
}

// Query 获取包
// @Summary 获取包
// @Description 根据活动ID或code获取包
// @Tags 包
// @Param type query int true "包类型"
// @Param activityId query int false "活动ID"
// @Param code query string false "活动code"
// @Router /pkg [get]
func (c *Pkg) Query(ctx *fiber.Ctx) error {
	typ, err := strconv.Atoi(ctx.Query("type", "1"))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "type格式错误"})
	}

	activityId, err := strconv.Atoi(ctx.Query("activityId", "0"))
	code := ctx.Query("code", "")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "activityId格式错误"})
	}
	if activityId == 0 && code == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "activityId或code必传"})
	}

	var act *entity.Activity
	if activityId > 0 {
		act = c.activityService.QueryById(uint(activityId))
	} else {
		act = c.activityService.QueryByCode(code)
	}

	if act == nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "活动不存在"})
	}

	switch typ {
	case 1:
		return v1.HandleSuccess(ctx, c.pkgService.Query(act.ID))
	case 2:
		return v1.HandleSuccess(ctx, c.pkgService.QueryAndGroup(act.ID))
	default:
		return ctx.Status(400).JSON(fiber.Map{"error": "type格式错误"})
	}
}

// Clear 清空
// @Summary 清空包
// @Description 根据活动ID或code清空包
// @Tags 包
// @Param id query int false "活动ID"
// @Param code query string false "活动code"
// @Router /pkg/clear [get]
func (c *Pkg) Clear(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Query("id", "0"))
	code := ctx.Query("code", "")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "id格式错误"})
	}
	if id == 0 && code == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "id或code必传"})
	}

	var act *entity.Activity
	if id > 0 {
		act = c.activityService.QueryById(uint(id))
	} else {
		act = c.activityService.QueryByCode(code)
	}
	if act == nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "活动不存在"})
	}

	c.pkgService.Delete(act.ID)
	return v1.HandleSuccess(ctx, nil)
}

func (c *Pkg) StashPackage(ctx *fiber.Ctx) error {
	uri := ctx.Request().URI()
	method := ctx.Method()
	headers := getHeaders(ctx.GetReqHeaders())
	queries := ctx.Queries()
	bodyStr := string(ctx.Body())
	host := string(uri.Host())
	path := string(uri.Path())
	url := "https://" + host + path

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

		response := ctx.Response()
		response.SetStatusCode(r.StatusCode())
		response.SetBody(r.Body())
		for k, v := range r.Header() {
			response.Header.Set(k, v[0])
		}
		return nil
	}

	p := &entity.Package{
		Url:     url,
		Method:  method,
		Host:    host,
		Path:    path,
		Headers: headers,
		Queries: queries,
		Body:    bodyStr,
		IP:      ctx.IP(),
	}

	c.pkgService.Add(p)
	return nil
}

func getHeaders(headers map[string][]string) map[string]string {
	m := make(map[string]string)
	for k, v := range headers {
		m[k] = v[0]
	}
	return m
}
