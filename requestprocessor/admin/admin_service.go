package admin

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"go_logger_reference/requestprocessor/db"
	"go_logger_reference/requestprocessor/user"
)

func NewAdminService(source *db.DBSource) *AdminService {
	return &AdminService{db: source}
}

type AdminService struct {
	db *db.DBSource
}

func (s *AdminService) Handle(logger *logrus.Logger, w http.ResponseWriter, _ *http.Request, info *user.UserInfo) {
	if info.Role != "admin" {
		logger.Warnf("Admin actions not permitted")
		http.Error(w, "Admin actions not permitted", http.StatusForbidden)
		return
	}

	someList := s.db.SelectSomething(logger, "system")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(strings.Join(someList, " | ")))
	logger.Infof("Admin action performed")
}
