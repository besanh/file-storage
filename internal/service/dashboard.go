package service

import (
	"context"
	"encoding/json"
	pb "file/api/dashboard/v1"
	"file/internal/biz"
	"fmt"
	"net/http"
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
		TotalFiles:     summary.TotalFiles,
		TotalFolders:   summary.TotalFolders,
		TotalShared:    summary.TotalShared,
		Report: &pb.GetDashboardSummaryReply_DashboardReport{
			NewFilesCount:  summary.Report.NewFilesCount,
			NewStorageUsed: summary.Report.NewStorageUsed,
		},
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

func (s *DashboardService) StreamDashboardUpdates(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	eventCh, cleanup, err := s.uc.Subscribe(ctx)
	if err != nil {
		s.log.Errorf("failed to subscribe to dashboard updates: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	defer cleanup()

	// Send initial summary
	s.sendSummary(ctx, w, flusher)

	ticker := time.NewTicker(30 * time.Second) // Heartbeat
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		case <-eventCh:
			s.sendSummary(ctx, w, flusher)
		}
	}
}

func (s *DashboardService) sendSummary(ctx context.Context, w http.ResponseWriter, flusher http.Flusher) {
	summary, err := s.GetDashboardSummary(ctx, &pb.GetDashboardSummaryRequest{})
	if err != nil {
		s.log.Errorf("failed to get dashboard summary for SSE: %v", err)
		return
	}

	fmt.Fprintf(w, "data: ")
	data, _ := json.Marshal(summary)
	w.Write(data)
	fmt.Fprintf(w, "\n\n")
	flusher.Flush()
}
