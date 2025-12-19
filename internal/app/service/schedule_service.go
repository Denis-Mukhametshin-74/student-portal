package service

import (
	"errors"
	"fmt"
	"student-portal/internal/app/dto"
	"student-portal/internal/app/repository"
	"student-portal/pkg/cache"
	"time"
)

type ScheduleService struct {
	scheduleRepo *repository.ScheduleRepository
	studentRepo  *repository.StudentRepository
	cache        *cache.MemoryCache
}

func NewScheduleService(scheduleRepo *repository.ScheduleRepository,
	studentRepo *repository.StudentRepository, cache *cache.MemoryCache) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		studentRepo:  studentRepo,
		cache:        cache,
	}
}

func (s *ScheduleService) GetSchedule(email string) (*dto.ScheduleResponse, error) {
	student, err := s.studentRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("пользователь не найден")
	}

	groupName := student.GroupName
	cacheKey := fmt.Sprintf("schedule:%s", groupName)

	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(*dto.ScheduleResponse), nil
	}

	schedules, err := s.scheduleRepo.FindByGroup(groupName)
	if err != nil {
		return nil, err
	}

	if schedules == nil {
		return &dto.ScheduleResponse{
			Group: groupName,
			Days:  nil,
		}, nil
	}

	scheduleMap := make(map[string][]dto.ScheduleItem)
	for _, schedule := range schedules {
		timeStr := fmt.Sprintf(
			"%s-%s",
			schedule.StartTime.Format("15:04"),
			schedule.EndTime.Format("15:04"),
		)

		item := dto.ScheduleItem{
			Time:    timeStr,
			Subject: schedule.Subject,
			Room:    schedule.Room,
			Teacher: schedule.Teacher,
		}

		scheduleMap[schedule.DayOfWeek] = append(
			scheduleMap[schedule.DayOfWeek], item,
		)
	}

	var days []dto.ScheduleDay
	dayOrder := []string{
		"Понедельник", "Вторник", "Среда",
		"Четверг", "Пятница", "Суббота",
	}

	for _, day := range dayOrder {
		if lessons, ok := scheduleMap[day]; ok && len(lessons) > 0 {
			days = append(days, dto.ScheduleDay{
				Day:     day,
				Lessons: lessons,
			})
		}
	}

	scheduleResponse := &dto.ScheduleResponse{
		Group: groupName,
		Days:  days,
	}

	s.cache.Set(cacheKey, scheduleResponse, 30*time.Minute)

	return scheduleResponse, nil
}
