package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type ModuleData struct {
	ModuleName string
	StructName string
}

func CreateModule(moduleName, structName string) error {
	modulePath := filepath.Join("module", moduleName)
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	data := ModuleData{
		ModuleName: moduleName,
		StructName: structName,
	}

	files := map[string]string{
		"init.go.tmpl":       "init.go",
		"model.go.tmpl":      "model.go",
		"dto.go.tmpl":        "dto.go",
		"handler.go.tmpl":    "handler.go",
		"service.go.tmpl":    "service.go",
		"repository.go.tmpl": "repository.go",
	}

	for templateName, outputName := range files {
		templateContent := getModuleTemplate(templateName)
		if templateContent == "" {
			continue
		}

		tmpl, err := template.New(templateName).Parse(templateContent)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", templateName, err)
		}

		outputPath := filepath.Join(modulePath, outputName)
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outputPath, err)
		}

		if err := tmpl.Execute(file, data); err != nil {
			file.Close()
			return fmt.Errorf("failed to execute template %s: %w", templateName, err)
		}
		file.Close()
	}

	return nil
}

func getModuleTemplate(name string) string {
	templates := map[string]string{
		"init.go.tmpl": `// module/{{.ModuleName}}/init.go
package {{.ModuleName}}

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"{{.ModuleName}}/server/core"
)

var _ core.Module = (*Module)(nil)

type Module struct {
	Service *Service
	handler *Handler
}

func Init(db *gorm.DB) *Module {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	return &Module{
		Service: svc,
		handler: h,
	}
}

func (m *Module) Routes(r *gin.RouterGroup) {
	m.handler.register(r)
}
`,
		"model.go.tmpl": `// module/{{.ModuleName}}/model.go
package {{.ModuleName}}

import (
	"time"

	"gorm.io/gorm"
)

type {{.StructName}} struct {
	ID          uint           ` + "`" + `gorm:"primaryKey" json:"id"` + "`" + `
	Title       string         ` + "`" + `gorm:"size:255;not null" json:"title"` + "`" + `
	Content     string         ` + "`" + `gorm:"type:text" json:"content"` + "`" + `
	Status      int8           ` + "`" + `gorm:"default:1" json:"status"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `json:"createdAt"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `json:"updatedAt"` + "`" + `
	DeletedAt   gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (*{{.StructName}}) TableName() string {
	return "{{.ModuleName}}s"
}
`,
		"dto.go.tmpl": `// module/{{.ModuleName}}/dto.go
package {{.ModuleName}}

import "time"

type Create{{.StructName}}Req struct {
	Title   string ` + "`" + `json:"title" binding:"required"` + "`" + `
	Content string ` + "`" + `json:"content"` + "`" + `
}

type Update{{.StructName}}Req struct {
	Title   string ` + "`" + `json:"title"` + "`" + `
	Content string ` + "`" + `json:"content"` + "`" + `
	Status  *int8  ` + "`" + `json:"status"` + "`" + `
}

type {{.StructName}}Resp struct {
	ID        uint      ` + "`" + `json:"id"` + "`" + `
	Title     string    ` + "`" + `json:"title"` + "`" + `
	Content   string    ` + "`" + `json:"content"` + "`" + `
	Status    int8      ` + "`" + `json:"status"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"createdAt"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updatedAt"` + "`" + `
}
`,
		"handler.go.tmpl": `// module/{{.ModuleName}}/handler.go
package {{.ModuleName}}

import (
	"strconv"

	"{{.ModuleName}}/server/helper"
	"{{.ModuleName}}/server/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) register(r *gin.RouterGroup) {
	group := r.Group("/{{.ModuleName}}s")
	{
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.POST("", middleware.JwtAuth(), h.Create)
		group.PUT("/:id", middleware.JwtAuth(), h.Update)
		group.DELETE("/:id", middleware.JwtAuth(), h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	var req Create{{.StructName}}Req
	if !helper.MustBindJSON(c, &req) {
		return
	}

	resp, err := h.svc.Create(c, userID, req)
	helper.Respond(c, err, resp)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	resp, total, err := h.svc.List(c, page, pageSize)
	helper.Respond(c, err, gin.H{
		"list":     resp,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	resp, err := h.svc.Get(c, id)
	helper.Respond(c, err, resp)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	var req Update{{.StructName}}Req
	if !helper.MustBindJSON(c, &req) {
		return
	}

	err := h.svc.Update(c, id, req)
	helper.Respond(c, err, gin.H{"message": "updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	err := h.svc.Delete(c, id)
	helper.Respond(c, err, gin.H{"message": "deleted"})
}
`,
		"service.go.tmpl": `// module/{{.ModuleName}}/service.go
package {{.ModuleName}}

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID uint, req Create{{.StructName}}Req) (*{{.StructName}}Resp, error) {
	item := &{{.StructName}}{
		Title:     req.Title,
		Content:   req.Content,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return toResp(item), nil
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]{{.StructName}}Resp, int64, error) {
	items, total, err := s.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var resp []{{.StructName}}Resp
	for _, item := range items {
		resp = append(resp, *toResp(&item))
	}

	return resp, total, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*{{.StructName}}Resp, error) {
	item, err := s.repo.First(ctx, id)
	if err != nil {
		return nil, err
	}

	return toResp(item), nil
}

func (s *Service) Update(ctx context.Context, id uint, req Update{{.StructName}}Req) error {
	item, err := s.repo.First(ctx, id)
	if err != nil {
		return err
	}

	if req.Title != "" {
		item.Title = req.Title
	}
	if req.Content != "" {
		item.Content = req.Content
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	item.UpdatedAt = time.Now()

	return s.repo.Update(ctx, item)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func toResp(item *{{.StructName}}) *{{.StructName}}Resp {
	var resp {{.StructName}}Resp
	_ = copier.Copy(&resp, item)
	return &resp
}
`,
		"repository.go.tmpl": `// module/{{.ModuleName}}/repository.go
package {{.ModuleName}}

import (
	"context"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type Repository struct {
	*dbx.BaseRepo[{{.StructName}}]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		BaseRepo: dbx.NewBaseRepo[{{.StructName}}](db),
	}
}

func (r *Repository) FindAll(ctx context.Context, page, pageSize int) ([]{{.StructName}}, int64, error) {
	var items []{{.StructName}}
	var total int64

	db := r.DB().WithContext(ctx).Model(&{{.StructName}}{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
`,
	}

	return templates[name]
}
