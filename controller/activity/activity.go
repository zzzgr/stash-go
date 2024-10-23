package activity

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	v1 "stash-go/api"
	"stash-go/model/dto/activity_dto"
	"stash-go/model/entity"
	"stash-go/service/activity"
	"stash-go/service/pkg"
	"strconv"
)

type Controller struct {
	activityService *activity.Service
	pkgService      *pkg.Service
}

func New(activityService *activity.Service, pkgService *pkg.Service) *Controller {
	return &Controller{
		activityService: activityService,
		pkgService:      pkgService,
	}
}

// QueryAll 查询所有活动
// @Summary 查询所有活动
// @Description 获取所有活动的列表
// @Tags 活动
// @Router /activity/query [get]
func (c *Controller) QueryAll(ctx *fiber.Ctx) error {
	return v1.HandleSuccess(ctx, c.activityService.Query())
}

// Save 保存活动
// @Summary 保存活动
// @Description 添加或更新活动信息
// @Tags 活动
// @Param activity body activity_dto.SaveRequestDTO true "活动信息"
// @Router /activity [post]
func (c *Controller) Save(ctx *fiber.Ctx) error {
	var reqDTO activity_dto.SaveRequestDTO
	if err := ctx.BodyParser(&reqDTO); err != nil {
		return v1.HandleError(ctx, fiber.StatusBadRequest, err, nil)
	}

	c.activityService.Save(&reqDTO)
	return v1.HandleSuccess(ctx, nil)
}

// Get 获取活动
// @Summary 获取活动
// @Description 根据 Id 或 code 查询活动
// @Tags 活动
// @Param id query int false "活动ID"
// @Param code query string false "活动code"
// @Router /activity [get]
func (c *Controller) Get(ctx *fiber.Ctx) error {
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
		panic("活动不存在")
	}

	return v1.HandleSuccess(ctx, act)
}

// Delete 根据 ID 删除活动
// @Summary 根据 ID 删除活动
// @Description 删除指定 ID 的活动
// @Tags 活动
// @Param id path int true "活动 ID"
// @Router /activity/{id} [delete]
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return v1.HandleError(ctx, fiber.StatusBadRequest, errors.New("无效的活动 ID"), nil)
	}

	c.activityService.Delete(uint(id))
	return ctx.SendStatus(fiber.StatusNoContent) // 返回 204 No Content
}
