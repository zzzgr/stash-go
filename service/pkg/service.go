package pkg

import (
	"errors"
	"fmt"
	"stash-go/model/dto/pkg"
	"stash-go/model/entity"
	"stash-go/repository"
	"stash-go/service/activity"
	"stash-go/util"
	"strings"
)

type Service struct {
	repo            *repository.Repository
	activityService *activity.Service
}

func NewService(
	repo *repository.Repository,
	activityService *activity.Service) *Service {
	return &Service{
		repo:            repo,
		activityService: activityService,
	}
}
func (s *Service) Add(requestDTO *entity.Package) {

	activities := s.activityService.Query()
	activityId := uint(0)
	accountName := ""
	for _, a := range activities {
		if util.IsMatch(requestDTO.Url, a.UrlPattern) {
			if len(a.QueryPattern) > 0 {
				// 包里的kv串去和库里面的匹配  来获取活动id
				kvMap := make(map[string]string)
				for k, v := range requestDTO.Queries {
					kvMap[fmt.Sprintf("%s=%s", k, v)] = ""
				}
				for _, q := range strings.Split(a.QueryPattern, "&") {
					_, ok := kvMap[q]
					if ok {
						activityId = a.ID
					}
				}
				if activityId == 0 {
					continue
				}
			} else {
				activityId = a.ID
			}

			// 提取field作为accountName
			switch a.Type {
			case 1:
				{
					ckMap := util.GetCookieFieldMap(requestDTO.Headers["Cookie"])
					for _, f := range strings.Split(a.Field, ",") {
						v, ok := ckMap[f]
						if ok {
							accountName = v
							break
						}
					}
					break
				}
			case 2:
				{

					for _, f := range strings.Split(a.Field, ",") {
						v, ok := requestDTO.Headers[f]
						if ok {
							accountName = v
							break
						}
					}
					break
				}
			default:
				// do nothing
			}
			if len(accountName) == 0 {
				accountName = "unknown"
			}
		}
	}

	if activityId > 0 && len(accountName) > 0 {
		requestDTO.ActivityId = activityId
		requestDTO.AccountName = accountName

		s.repo.Db.Save(requestDTO)
	}

}

func (s *Service) Query(activityId uint) []*entity.Package {
	var pkgs = make([]*entity.Package, 0)
	tx := s.repo.Db.Model(&entity.Package{})

	if activityId > 0 {
		tx = tx.Where("activity_id = ?", activityId)
	}

	err := tx.Find(&pkgs).Error
	if err != nil {
		panic(errors.New("查询失败"))
	}
	return pkgs
}

func (s *Service) QueryAndGroup(activityId uint) []*pkg.NamedPackage {
	pkgs := s.Query(activityId)

	namedPackageMap := make(map[string]*pkg.NamedPackage)
	for _, p := range pkgs {

		namedPkg, ok := namedPackageMap[p.AccountName]
		if ok {
			namedPkg.Packages = append(namedPkg.Packages, p)
		} else {
			namedPackageMap[p.AccountName] = &pkg.NamedPackage{
				Name:     p.AccountName,
				Packages: []*entity.Package{p},
			}
		}
	}

	namedPackages := make([]*pkg.NamedPackage, 0)
	for _, p := range namedPackageMap {
		namedPackages = append(namedPackages, p)
	}
	return namedPackages
}

func (s *Service) Delete(activityId uint) {
	s.repo.Db.Unscoped().Delete(&entity.Package{}, "activity_id = ?", activityId)
}
