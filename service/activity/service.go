package activity

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"stash-go/model/dto/activity_dto"
	"stash-go/model/entity"
	"stash-go/repository"
)

type Service struct {
	repo  *repository.Repository
	cache []*entity.Activity
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo:  repo,
		cache: make([]*entity.Activity, 0),
	}
}
func (s *Service) Save(requestDTO *activity_dto.SaveRequestDTO) {
	id := requestDTO.Id

	var activityDb entity.Activity

	if id == 0 {
		// 新增逻辑
		if requestDTO.Name == "" {
			panic("Name不能为空")
		}
		if requestDTO.Code == "" {
			panic("Code不能为空")
		}
		if requestDTO.UrlPattern == "" {
			panic("UrlPattern不能为空")
		}

		// 检查是否存在相同活动
		err := s.repo.Db.Where("name = ? OR code = ? OR (url_pattern = ? AND query_pattern = ?)",
			requestDTO.Name, requestDTO.Code, requestDTO.UrlPattern, requestDTO.QueryPattern).
			First(&activityDb).Error

		if activityDb.ID != 0 {
			panic(fmt.Sprintf("相同的活动已存在, code为: %s", activityDb.Code))
		}

		// 创建新记录
		err = s.repo.Db.Transaction(func(tx *gorm.DB) error {
			newActivity := entity.Activity{
				Code:         requestDTO.Code,
				Name:         requestDTO.Name,
				Cron:         requestDTO.Cron,
				Advance:      requestDTO.Advance,
				Interval:     requestDTO.Interval,
				UrlPattern:   requestDTO.UrlPattern,
				QueryPattern: requestDTO.QueryPattern,
				Type:         requestDTO.Type,
				Field:        requestDTO.Field,
				AlertAhead:   requestDTO.AlertAhead,
				Status:       requestDTO.Status,
			}
			if err := tx.Create(&newActivity).Error; err != nil {
				return err
			}
			// 更新缓存
			s.cache = append(s.cache, &newActivity)
			return nil
		})

		if err != nil {
			log.Warnf("新建活动发生错误：%s", err.Error())
		}
	} else {
		// 更新逻辑
		err := s.repo.Db.First(&activityDb, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			panic("活动不存在")
		}

		err = s.repo.Db.Transaction(func(tx *gorm.DB) error {
			activityDb.Code = requestDTO.Code
			activityDb.Name = requestDTO.Name
			activityDb.Cron = requestDTO.Cron
			activityDb.Advance = requestDTO.Advance
			activityDb.Interval = requestDTO.Interval
			activityDb.UrlPattern = requestDTO.UrlPattern
			activityDb.QueryPattern = requestDTO.QueryPattern
			activityDb.Type = requestDTO.Type
			activityDb.Field = requestDTO.Field
			activityDb.AlertAhead = requestDTO.AlertAhead
			activityDb.Status = requestDTO.Status

			if err := tx.Save(&activityDb).Error; err != nil {
				return err
			}
			// 更新缓存
			for i, activityCache := range s.cache {
				if activityCache.ID == id {
					s.cache[i] = &activityDb // 更新缓存中的记录
					break
				}
			}
			return nil
		})

		if err != nil {
			log.Warnf("更新活动发生错误：%s", err.Error())
		}
	}
}

func (s *Service) Query() []*entity.Activity {
	if len(s.cache) > 0 {
		return s.cache // 返回缓存数据，无错误
	}

	// 从数据库查询活动
	err := s.repo.Db.Model(&entity.Activity{}).Find(&s.cache).Error
	if err != nil {
		panic("查询失败")
	}
	return s.cache
}

func (s *Service) QueryById(activityId uint) *entity.Activity {
	var activity *entity.Activity
	err := s.repo.Db.Where("id=?", activityId).First(&activity).Error
	if err != nil {
		return nil
	}
	return activity
}

func (s *Service) QueryByCode(code string) *entity.Activity {

	var activity *entity.Activity
	err := s.repo.Db.Where("code=?", code).First(&activity).Error
	if err != nil {
		return nil
	}
	return activity
}

func (s *Service) Delete(id uint) {
	var activity entity.Activity
	err := s.repo.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Delete(&activity, id).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&entity.Package{}, "activity_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("删除发生异常：%s", err.Error()))
	} else {
		// 更新缓存
		for i, activityCache := range s.cache {
			if activityCache.ID == id {
				s.cache = append(s.cache[:i], s.cache[i+1:]...) // 删除缓存中的记录
				break
			}
		}
	}
}
