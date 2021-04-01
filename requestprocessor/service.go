package requestprocessor

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"go_logger_reference/requestprocessor/admin"
	"go_logger_reference/requestprocessor/report"
	"go_logger_reference/requestprocessor/user"
)

func NewService(config string, reportService *report.ReportService, adminService *admin.AdminService, logger *zap.SugaredLogger) *Service {
	return &Service{
		loggerConfig:  config,
		logger:        logger,
		reportService: reportService,
		adminService:  adminService,
	}
}

type Service struct {
	nextRequetID  uint64
	logger        *zap.SugaredLogger
	loggerConfig  string
	reportService *report.ReportService
	adminService  *admin.AdminService
}

func (s *Service) StartListening() error {
	if s.reportService == nil || s.adminService == nil {
		return fmt.Errorf("s.reportService == nil or s.adminService == nil")
	}

	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(s.entryPoint))

		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return nil
}

func (s *Service) entryPoint(w http.ResponseWriter, r *http.Request) {
	requestID := s.getNextRequestID()

	localLogger := s.logger.With("who", "service", "requestID", requestID, "url", r.URL.Path)
	localLogger.Infow("Processing request")

	userInfo, err := oauthValidate(localLogger, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userInfo == nil {
		s.redirect(localLogger, w, r)
		return
	}

	localLogger = localLogger.With("user", userInfo.Username, "role", userInfo.Role)

	switch {
	case strings.HasPrefix(r.URL.Path, "/admin/"):
		s.adminService.Handle(localLogger, w, r, userInfo)

	case strings.HasPrefix(r.URL.Path, "/report/"):
		s.reportService.Handle(localLogger, w, r, userInfo)

	default:
		s.notFound(localLogger, w, r)
	}
}

func oauthValidate(logger *zap.SugaredLogger, r *http.Request) (*user.UserInfo, error) {
	token := r.Header.Get("token")
	if token == "" {
		logger.Warn("Token is absent")
		return nil, fmt.Errorf("Authentication required")
	}

	tokenStruct := struct {
		Username string
		Role     string
	}{}

	err := json.Unmarshal([]byte(token), &tokenStruct)
	if err != nil {
		logger.Errorw("failed to unmarshal token", "err", err)
		return nil, err
	}

	if tokenStruct.Username == "" || tokenStruct.Role == "" {
		logger.Errorf("user or role is empty. user: %s, role: %s", tokenStruct.Username, tokenStruct.Role)
		return nil, nil
	}

	userInfo := &user.UserInfo{Username: tokenStruct.Username, Role: tokenStruct.Role}

	logger.Infow("Successfully validated user", "username", userInfo.Username, "role", userInfo.Role)

	return userInfo, nil
}

func (s *Service) getNextRequestID() string {
	return fmt.Sprintf("req-%d", atomic.AddUint64(&s.nextRequetID, 1))
}

func (s *Service) redirect(logger *zap.SugaredLogger, w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	logger.Info("Redirect for authorization sent")
}

func (s *Service) notFound(logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	logger.Info("'Page not found' sent")
}
