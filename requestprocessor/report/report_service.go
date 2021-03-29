package report

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"go_logger_reference/requestprocessor/db"
	"go_logger_reference/requestprocessor/user"
)

func NewReportService(source *db.DBSource) *ReportService {
	return &ReportService{dbSource: source}
}

type ReportService struct {
	dbSource *db.DBSource
}

func (s *ReportService) Handle(logger *logrus.Logger, w http.ResponseWriter, _ *http.Request, info *user.UserInfo) {
	if info.Role != "reporter" {
		logger.Warnf("User tried to access report service which he has no access to")
		http.Error(w, "Report service is for reporters ONLY", http.StatusForbidden)
		return
	}

	reportItems := s.dbSource.SelectSomething(logger, "months")
	w.WriteHeader(200)
	report := fmt.Sprintf("Affected months:\n%s\n---------------", strings.Join(reportItems, "\n"))
	_, _ = w.Write([]byte(report))
	logger.Infof("Report generated, %d items written", len(reportItems))
}
