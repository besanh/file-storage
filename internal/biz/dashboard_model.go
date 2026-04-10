package biz

type DashboardSummary struct {
	WelcomeMessage string
	Storage        *UserStorageStats
	RecentFiles    []*FileNode
	ActivePlanName string
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
