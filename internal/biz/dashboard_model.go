package biz

type DashboardSummary struct {
	WelcomeMessage string
	Storage        *UserStorageStats
	RecentFiles    []*FileNode
	ActivePlanName string
	TotalFiles     int64
	TotalFolders   int64
	TotalShared    int64
	Report         *DashboardReport
}

type DashboardReport struct {
	NewFilesCount  int64
	NewStorageUsed int64
}

type UserStorageStats struct {
	Photos    int64
	Videos    int64
	Documents int64
	Audio     int64
	Compress  int64
	Other     int64
	Total     int64
	Quota     int64
}
