package requestprocessor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"go_logger_reference/requestprocessor/admin"
	"go_logger_reference/requestprocessor/report"
	"go_logger_reference/requestprocessor/user"
	"go_logger_reference/utils"
)

func NewService(config string, reportService *report.ReportService, adminService *admin.AdminService) *Service {
	return &Service{
		loggerConfig:  config,
		logger:        utils.NewLoggerFromConfig(config),
		reportService: reportService,
		adminService:  adminService,
	}
}

type Service struct {
	nextRequetID  uint64
	logger        *log.Logger
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

	localLogger := utils.NewLoggerFromConfig(s.loggerConfig)
	localLogger.AddHook(utils.LogDefaultField("who", "service"))
	localLogger.AddHook(utils.LogDefaultField("requestID", requestID))
	localLogger.WithField("url", r.URL.Path).Infof("Processing request")

	userInfo, err := oauthValidate(localLogger, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userInfo == nil {
		s.redirect(localLogger, w, r)
		return
	}

	localLogger.AddHook(utils.LogDefaultField("user", userInfo.Username))
	localLogger.AddHook(utils.LogDefaultField("role", userInfo.Role))

	switch {
	case strings.HasPrefix(r.URL.Path, "/admin/"):
		s.adminService.Handle(localLogger, w, r, userInfo)

	case strings.HasPrefix(r.URL.Path, "/report/"):
		s.reportService.Handle(localLogger, w, r, userInfo)

	default:
		s.notFound(localLogger, w, r)
	}
}

func oauthValidate(logger *log.Logger, r *http.Request) (*user.UserInfo, error) {
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
		logger.WithError(err).Errorf("failed to unmarshal token")
		return nil, err
	}

	if tokenStruct.Username == "" || tokenStruct.Role == "" {
		logger.Errorf("user or role is empty. user: %s, role: %s", tokenStruct.Username, tokenStruct.Role)
		return nil, nil
	}

	userInfo := &user.UserInfo{Username: tokenStruct.Username, Role: tokenStruct.Role}

	logger.Infof("Successfully validated user %s, role %s", userInfo.Username, userInfo.Role)

	return userInfo, nil
}

func (s *Service) getNextRequestID() string {
	return fmt.Sprintf("req-%d", atomic.AddUint64(&s.nextRequetID, 1))
}

func (s *Service) redirect(logger *log.Logger, w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	logger.Infof("Redirect for authorization sent")
}

func (s *Service) notFound(logger *log.Logger, w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	logger.Infof("'Page not found' sent")
}
