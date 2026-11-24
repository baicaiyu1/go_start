package upload

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	limits "hertz-tutorial/demo/biz/consts"
	"hertz-tutorial/demo/biz/errs"
	"hertz-tutorial/demo/biz/response"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	// uploadDir 上传文件落盘目录。
	uploadDir = "./uploads"
	// maxUploadSize 单文件大小上限，不能超过 main.go 的 MaxRequestBodySize。
	maxUploadSize = limits.MaxUploadSize
)

// allowedExt 允许上传的扩展名白名单。
var allowedExt = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".webp": {},
	".pdf":  {},
	".doc":  {},
	".docx": {},
	".xls":  {},
	".xlsx": {},
	".ppt":  {},
	".pptx": {},
	".txt":  {},
	".md":   {},
	".zip":  {},
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Upload 上传单个文件
//
// @Summary 上传单个文件
// @Description form-data 字段 file；单文件 ≤5MB，扩展名走白名单校验，落盘到 ./uploads
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} UploadResp
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /api/v1/upload/one [post]
func (h *Handler) Upload(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "file is required")
		return
	}
	if file.Size <= 0 {
		response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "file is empty")
		return
	}
	if file.Size > maxUploadSize {
		response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "file too large")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedExt[ext]; !ok {
		response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "unsupported file type")
		return
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		hlog.CtxErrorf(ctx, "create upload dir: %v", err)
		response.Fail(c, consts.StatusInternalServerError, errs.CodeInternal, "save failed")
		return
	}

	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		hlog.CtxErrorf(ctx, "save upload: %v", err)
		response.Fail(c, consts.StatusInternalServerError, errs.CodeInternal, "save failed")
		return
	}

	response.OK(c, utils.H{
		"name": name,
		"size": file.Size,
	})
}

// BatchUpload 批量上传文件
//
// @Summary 批量上传文件
// @Description form-data 字段 files；一次 1-5 个文件，单文件 ≤5MB，扩展名走白名单校验
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "文件（1-5 个）"
// @Success 200 {object} BatchResp
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /api/v1/upload/batch [post]
func (h *Handler) BatchUpload(ctx context.Context, c *app.RequestContext) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "invalid multipart form")
		return
	}
	files := form.File["files"]
	if len(files) == 0 || len(files) > 5 {
		response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "1-5 files required")
		return
	}

	names := make([]string, 0, len(files))
	for _, f := range files {
		if f.Size > 5<<20 {
			response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "file too large")
			return
		}
		ext := strings.ToLower(filepath.Ext(f.Filename))
		if _, ok := allowedExt[ext]; !ok {
			response.Fail(c, consts.StatusBadRequest, errs.CodeBadRequest, "unsupported file type")
			return
		}
		name := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), len(names), ext)
		if err := c.SaveUploadedFile(f, filepath.Join("./uploads", name)); err != nil {
			response.Fail(c, consts.StatusInternalServerError, errs.CodeInternal, "save failed")
			return
		}
		names = append(names, name)
	}
	response.OK(c, utils.H{"names": names})
}
