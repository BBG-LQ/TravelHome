// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestayComment

import (
	"context"

	"TravelHome/app/travel/cmd/api/internal/svc"
	"TravelHome/app/travel/cmd/api/internal/types"
	"TravelHome/app/travel/model"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay comment list
func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	// todo: add your logic here and delete this line
	var homestayComment []types.HomestayComment

	list, err := l.svcCtx.HomestayCommentModel.FindPageListByIdDESC(l.ctx, l.svcCtx.HomestayCommentModel.SelectBuilder(), req.LastId, req.PageSize)

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), " homestayComment db err req : %+v , err : %v", req, err)
	}

	if len(list) > 0 {

		for _, item := range list {

			var tempHomestayComment types.HomestayComment
			_ = copier.Copy(&tempHomestayComment, item)

			homestayComment = append(homestayComment, tempHomestayComment)

		}

	}

	return &types.CommentListResp{List: homestayComment}, nil
}
