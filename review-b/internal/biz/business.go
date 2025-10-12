package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	Reason    string
	PicInfo   string
	VideoInfo string
}

// BusinessRepo is a Greater repo.
type BusinessRepo interface {
	Reply(context.Context, *ReplyParam) (int64, error)
	Appeal(context.Context, *AppealParam) (int64, error)
}

// BusinessUsecase is a Business usecase.
type BusinessUsecase struct {
	repo BusinessRepo
	log  *log.Helper
}

// NewBusinessUsecase new a Business usecase.
func NewBusinessUsecase(repo BusinessRepo, logger log.Logger) *BusinessUsecase {
	return &BusinessUsecase{repo: repo, log: log.NewHelper(logger)}
}

// service层调用biz层方法
func (uc *BusinessUsecase) CreatReply(ctx context.Context, param *ReplyParam) (int64, error) {
	//调用data层方法
	uc.log.WithContext(ctx).Infof("[biz]CreatReply: %v", param)
	return uc.repo.Reply(ctx, param)
}

func (uc *BusinessUsecase) CreatAppeal(ctx context.Context, param *AppealParam) (int64, error) {
	//调用data层方法
	uc.log.WithContext(ctx).Infof("[biz]CreatAppeal: %v", param)
	return uc.repo.Appeal(ctx, param)
}
