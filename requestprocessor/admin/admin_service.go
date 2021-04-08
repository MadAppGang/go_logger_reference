package admin

import (
	"net/http"
	"strings"

	"go_logger_reference/requestprocessor/db"
	"go_logger_reference/requestprocessor/user"

	"go.uber.org/zap"
)

func NewAdminService(source *db.DBSource) *AdminService {
	return &AdminService{db: source}
}

type AdminService struct {
	db *db.DBSource
}

func (s *AdminService) Handle(logger *zap.SugaredLogger, w http.ResponseWriter, _ *http.Request, info *user.UserInfo) {
	if info.Role != "admin" {
		logger.Warn("Admin actions not permitted")
		http.Error(w, "Admin actions not permitted", http.StatusForbidden)
		return
	}

	someList := s.db.SelectSomething(logger, "system")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(strings.Join(someList, " | ")))
	logger.Info("Admin action performed")
}
