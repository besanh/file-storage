package service

import (
	"context"
	pb "file/api/dashboard/v1"
	"file/internal/biz"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type DashboardService struct {
	pb.UnimplementedDashboardServiceServer
	uc  *biz.DashboardUsecase
	log *log.Helper
}

func NewDashboardService(uc *biz.DashboardUsecase, logger log.Logger) *DashboardService {
	return &DashboardService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *DashboardService) GetDashboardSummary(ctx context.Context, req *pb.GetDashboardSummaryRequest) (*pb.GetDashboardSummaryReply, error) {
	summary, err := s.uc.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetDashboardSummaryReply{
		WelcomeMessage: summary.WelcomeMessage,
		Storage: &pb.GetDashboardSummaryReply_StorageBreakdown{
			Photos:    summary.Storage.Photos,
			Videos:    summary.Storage.Videos,
			Documents: summary.Storage.Documents,
			Audio:     summary.Storage.Audio,
			Compress:  summary.Storage.Compress,
			Other:     summary.Storage.Other,
			Total:     summary.Storage.Total,
			Quota:     summary.Storage.Quota,
		},
		RecentActivity: mapRecentActivityToPb(summary.RecentFiles),
		ActivePlanName: summary.ActivePlanName,
	}, nil
}

func mapRecentActivityToPb(nodes []*biz.FileNode) []*pb.GetDashboardSummaryReply_RecentActivity {
	res := make([]*pb.GetDashboardSummaryReply_RecentActivity, 0, len(nodes))
	for _, n := range nodes {
		res = append(res, &pb.GetDashboardSummaryReply_RecentActivity{
			Id:             n.ID.String(),
			Name:           n.Name,
			Type:           n.FileType,
			Icon:           "", // Can be determined by type in FE
			Size:           formatByteSize(n.FileSize),
			LastAccessedAt: n.LastAccessed.Format(time.RFC3339),
			IsFolder:       n.IsFolder,
		})
	}
	return res
}

func formatByteSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
