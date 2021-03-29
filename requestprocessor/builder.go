package requestprocessor

import (
	"go_logger_reference/requestprocessor/admin"
	"go_logger_reference/requestprocessor/db"
	"go_logger_reference/requestprocessor/report"
)

func BuildService(config string) *Service {
	dbSource := db.NewDBSource()

	reportService := report.NewReportService(dbSource)
	adminService := admin.NewAdminService(dbSource)

	service := NewService(config, reportService, adminService)

	return service
}
