package data

import (
	"context"

	v1 "review-b/api/review/v1"
	"review-b/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type BusinessRepo struct {
	data *Data
	log  *log.Helper
}

// NewBusinessRepo .
func NewBusinessRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &BusinessRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *BusinessRepo) Reply(ctx context.Context, param *biz.ReplyParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data]Reply: %v", param)
	//之前是查询操作数据库，现在需要rpc调用其他服务实现
	ret, err := r.data.rc.ReplyReview(ctx, &v1.ReplyReviewRequest{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	r.log.WithContext(ctx).Infof("[data]Reply: %v", ret)
	if err != nil {
		return 0, err
	}
	return ret.GetReplyID(), nil
}

func (r *BusinessRepo) Appeal(ctx context.Context, param *biz.AppealParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data]Appeal: %v", param)
	//之前是查询操作数据库，现在需要rpc调用其他服务实现
	ret, err := r.data.rc.AppealReview(ctx, &v1.AppealReviewRequest{
		ReviewID: param.ReviewID,
		StoreID:  param.StoreID,
		Content:  param.Content,
		PicInfo:  param.PicInfo,
		Reason:   param.Reason,
	})
	r.log.WithContext(ctx).Infof("[data]Appeal: %v", ret)
	if err != nil {
		return 0, err
	}
	return ret.GetAppealID(), nil
}
