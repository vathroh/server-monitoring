package service

import (
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
)

type SettingService interface {
	GetAllSettings() (map[string]string, error)
	GetSetting(key string) (string, error)
	UpsertSettings(settings map[string]string) error
}

type settingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{repo: repo}
}

func (s *settingService) GetAllSettings() (map[string]string, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}

func (s *settingService) GetSetting(key string) (string, error) {
	setting, err := s.repo.GetByKey(key)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *settingService) UpsertSettings(settings map[string]string) error {
	for k, v := range settings {
		err := s.repo.Upsert(&domain.Setting{Key: k, Value: v})
		if err != nil {
			return err
		}
	}
	return nil
}
