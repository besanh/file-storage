package biz

import (
	"context"
	db "file/internal/data/db/generated"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type DashboardUsecase struct {
	fileRepo FileRepo
	subRepo  SubscriptionRepo
	planRepo PlanRepo
	userRepo UserRepo
	authRepo AuthRepo
	log      *log.Helper
}

func NewDashboardUsecase(fileRepo FileRepo, subRepo SubscriptionRepo, planRepo PlanRepo, userRepo UserRepo, authRepo AuthRepo, logger log.Logger) *DashboardUsecase {
	return &DashboardUsecase{
		fileRepo: fileRepo,
		subRepo:  subRepo,
		planRepo: planRepo,
		userRepo: userRepo,
		authRepo: authRepo,
		log:      log.NewHelper(logger),
	}
}

func (uc *DashboardUsecase) GetSummary(ctx context.Context) (*DashboardSummary, error) {
	_, actorID, err := GetActorInfo(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, fmt.Errorf("invalid user identity string: %v", err)
	}

	// 1. Fetch Storage Stats
	user, err := uc.userRepo.GetUserStorage(ctx, userID)
	if err != nil {
		uc.log.Errorf("failed to get user storage: %v", err)
		user = &User{ID: userID} // Fallback to empty usage
	}

	// 2. Fetch Subscription & Plan
	var activePlanName = "Free"
	var quota int64 = 100 * 1024 * 1024 // Default 100MB
	sub, err := uc.subRepo.GetUserSubscription(ctx, userID)
	if err == nil && sub != nil {
		plan, err := uc.planRepo.GetPlan(ctx, sub.PlanID)
		if err == nil && plan != nil {
			activePlanName = plan.Name
			quota = plan.StorageQuota
		}
	}

	// 3. Fetch Recent Activity
	recentNodes, err := uc.fileRepo.GetRecentFiles(ctx, userID, 10)
	if err != nil {
		uc.log.Errorf("failed to get recent files: %v", err)
	}

	// Use the user's email from the Auth Service for identification.
	userName := "User"
	if profile, err := uc.authRepo.GetUserProfile(ctx, userID.String()); err == nil && profile != nil {
		userName = profile.Email
	}

	return &DashboardSummary{
		WelcomeMessage: fmt.Sprintf("Welcome back, %s", userName),
		Storage: &UserStorageStats{
			Photos:    user.StoragePhotosUsed,
			Videos:    user.StorageVideoUsed,
			Documents: user.StorageDocumentUsed,
			Audio:     user.StorageAudioUsed,
			Compress:  user.StorageCompressUsed,
			Other:     user.StorageOtherUsed,
			Total:     user.TotalStorageUsed(),
			Quota:     quota,
		},
		RecentFiles:    mapFileNodesToBiz(recentNodes),
		ActivePlanName: activePlanName,
	}, nil
}

func mapFileNodesToBiz(nodes []*db.FileNode) []*FileNode {
	res := make([]*FileNode, 0, len(nodes))
	for _, n := range nodes {
		res = append(res, &FileNode{
			ID:           n.ID,
			Name:         n.Name,
			IsFolder:     n.IsFolder,
			FileSize:     n.FileSize.Int64,
			FileType:     n.FileType.String,
			FileExt:      n.FileExt.String,
			MimeType:     n.FileMimeType.String,
			LastAccessed: n.RecentAccessedAt.Time,
		})
	}
	return res
}
