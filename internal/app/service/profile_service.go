package service

import (
	"errors"
	"fmt"
	"student-portal/internal/app/dto"
	"student-portal/internal/app/repository"
	"student-portal/pkg/cache"
	"time"
)

type ProfileService struct {
	studentRepo *repository.StudentRepository
	subjectRepo *repository.SubjectRepository
	cache       *cache.MemoryCache
}

func NewProfileService(studentRepo *repository.StudentRepository,
	subjectRepo *repository.SubjectRepository, cache *cache.MemoryCache) *ProfileService {
	return &ProfileService{
		studentRepo: studentRepo,
		subjectRepo: subjectRepo,
		cache:       cache,
	}
}

func (s *ProfileService) GetProfile(email string) (*dto.ProfileResponse, error) {
	cacheKey := fmt.Sprintf("profile:%s", email)
	if cached, found := s.cache.Get(cacheKey); found {
		if profile, ok := cached.(*dto.ProfileResponse); ok {
			return profile, nil
		}
	}

	student, err := s.studentRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("пользователь не найден")
	}

	studentResponse := dto.StudentResponse{
		ID:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Group: student.GroupName,
	}

	subjects, err := s.subjectRepo.FindByStudentID(student.ID)
	if err != nil {
		return nil, err
	}

	if subjects == nil {
		return &dto.ProfileResponse{
			Student:  studentResponse,
			Subjects: nil,
		}, nil
	}

	var subjectsResponse []dto.SubjectResponse
	for _, subject := range subjects {
		subjectsResponse = append(subjectsResponse, dto.SubjectResponse{
			ID:      subject.ID,
			Name:    subject.Name,
			Teacher: subject.Teacher,
		})
	}

	s.cache.Set(cacheKey, studentResponse, 5*time.Minute)
	s.cache.Set(cacheKey, subjectsResponse, 10*time.Minute)

	return &dto.ProfileResponse{
		Student:  studentResponse,
		Subjects: subjectsResponse,
	}, nil
}
