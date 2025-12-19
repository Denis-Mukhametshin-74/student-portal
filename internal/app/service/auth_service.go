package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"student-portal/internal/app/dto"
	"student-portal/internal/app/model"
	"student-portal/internal/app/repository"
	"student-portal/pkg/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	studentRepo       *repository.StudentRepository
	oauthRepo         *repository.OAuthRepository
	googleOAuthConfig *oauth2.Config
}

func NewAuthService(studentRepo *repository.StudentRepository,
	oauthRepo *repository.OAuthRepository) *AuthService {
	googleOAuthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &AuthService{
		studentRepo:       studentRepo,
		oauthRepo:         oauthRepo,
		googleOAuthConfig: googleOAuthConfig,
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	student, err := s.studentRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(student.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("неверный пароль")
	}

	token, err := jwt.GenerateToken(student.ID, student.Email)
	if err != nil {
		return nil, errors.New("ошибка генерации токена")
	}

	return &dto.LoginResponse{
		Token: token,
		Student: dto.StudentResponse{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
			Group: student.GroupName,
		},
	}, nil
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.StudentResponse, error) {
	exists, err := s.studentRepo.Exists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка хеширования пароля")
	}

	student := &model.Student{
		Name:         req.Name,
		Email:        req.Email,
		GroupName:    req.Group,
		PasswordHash: string(hashedPassword),
	}

	err = s.studentRepo.Create(student)
	if err != nil {
		return nil, err
	}

	return &dto.StudentResponse{
		ID:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Group: student.GroupName,
	}, nil
}

func (s *AuthService) GetGoogleOAuthURL(state string) string {
	return s.googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *AuthService) HandleGoogleCallback(code string) (*dto.LoginResponse, error) {
	token, err := s.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("ошибка обмена code: %w", err)
	}

	client := s.googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных пользователя: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	oauthProvider, err := s.oauthRepo.FindByProviderAndExternalID("google", userInfo.ID)
	if err != nil {
		return nil, err
	}

	var student *model.Student

	if oauthProvider != nil {
		student, err = s.studentRepo.FindByID(oauthProvider.StudentID)
		if err != nil {
			return nil, err
		}
	} else {
		student, err = s.studentRepo.FindByEmail(userInfo.Email)
		if err != nil {
			return nil, err
		}

		if student == nil {
			student = &model.Student{
				Name:         userInfo.Name,
				Email:        userInfo.Email,
				GroupName:    "Не указана",
				PasswordHash: "",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			if err := s.studentRepo.Create(student); err != nil {
				if strings.Contains(err.Error(), "duplicate") {
					return nil, errors.New("пользователь с таким email уже существует")
				}
				return nil, err
			}
		}

		oauth := &model.OAuthProvider{
			StudentID:    student.ID,
			ProviderName: "google",
			ExternalID:   userInfo.ID,
			Email:        userInfo.Email,
		}

		if err := s.oauthRepo.Create(oauth); err != nil {
			if !strings.Contains(err.Error(), "duplicate") {
				return nil, err
			}
		}
	}

	jwtToken, err := jwt.GenerateToken(student.ID, student.Email)
	if err != nil {
		return nil, errors.New("ошибка генерации токена")
	}

	return &dto.LoginResponse{
		Token: jwtToken,
		Student: dto.StudentResponse{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
			Group: student.GroupName,
		},
	}, nil
}
