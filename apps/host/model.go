package host

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
)

var (
	validate = validator.New()
)

type HostSet struct {
	Total int    `json:"total"`
	Items []*Host `json:"items"`
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func (s *HostSet) Add(item *Host) {
	s.Items = append(s.Items, item)
}

// Host模型的定义
type Host struct {
	// 资源公共属性部分
	*Resource
	// 资源独有属性部分
	*Describe
}


func NewHost() *Host {
	return &Host{
		Resource: &Resource{},
		Describe: &Describe{},
	}
}

func (h *Host) Validate() error {
	return validate.Struct(h)
}

func (h *Host) Put(obj *Host) error {
	if obj.Id != h.Id {
		return fmt.Errorf("Id invalid")
	}

	*h.Resource = *obj.Resource
	*h.Describe = *obj.Describe
	return nil
}

func (h *Host) Patch(obj *Host) error {
	return mergo.Merge(h, obj)
}

func (h *Host) InjectDefault() {
	if h.CreateAt == 0 {
		h.CreateAt = time.Now().UnixMilli()
	}
}

type Vendor int

const (
	// 枚举的默认值
	PRIVATE_IDC Vendor = iota
	// 阿里云
	ALIYUN
	// 腾讯云
	TXYUN
)

type Resource struct {
	Id          string            `json:"id"  validate:"required"`     // 全局唯一Id
	Vendor      Vendor            `json:"vendor"`                      // 厂商
	Region      string            `json:"region"  validate:"required"` // 地域
	CreateAt    int64             `json:"create_at"`                   // 创建时间
	ExpireAt    int64             `json:"expire_at"`                   // 过期时间
	Type        string            `json:"type"  validate:"required"`   // 规格
	Name        string            `json:"name"  validate:"required"`   // 名称
	Description string            `json:"description"`                 // 描述
	Status      string            `json:"status"`                      // 服务商中的状态
	Tags        map[string]string `json:"tags"`                        // 标签
	UpdateAt    int64             `json:"update_at"`                   // 更新时间
	SyncAt      int64             `json:"sync_at"`                     // 同步时间
	Account     string            `json:"accout"`                      // 资源的所属账号
	PublicIP    string            `json:"public_ip"`                   // 公网IP
	PrivateIP   string            `json:"private_ip"`                  // 内网IP
}

type Describe struct {
	CPU          int    `json:"cpu" validate:"required"`    // 核数
	Memory       int    `json:"memory" validate:"required"` // 内存
	GPUAmount    int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec      string `json:"gpu_spec"`                   // GPU类型
	OSType       string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName       string `json:"os_name"`                    // 操作系统名称
	SerialNumber string `json:"serial_number"`              // 序列号
}

type QueryHostRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	Keywords   string `json:"kws"`
}

func NewQueryHostRequest() *QueryHostRequest {
	return &QueryHostRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

func NewQueryHostFromRequest(req *http.Request) *QueryHostRequest {
	request := NewQueryHostRequest()
	qs := req.URL.Query()
	page_size := qs.Get("page_size")
	if page_size != "" {
		//Go语言的 strconv 包提供了一个Atoi()函数，该函数等效于ParseInt(str string，base int，bitSize int)用于将字符串类型转换为int类型。
		request.PageSize, _ = strconv.Atoi(page_size)
	}

	page_number := qs.Get("page_number")
	if page_number != "" {
		request.PageNumber, _ = strconv.Atoi(page_number)
	}

	request.Keywords = qs.Get("kws")

	return request
}

func (req *QueryHostRequest) GetPageSize() uint{
	return uint(req.PageSize)
}

func (req *QueryHostRequest) OffSet() int64{
	return int64((req.PageNumber - 1) * req.PageSize)
}


type DescribeHostRequest struct {
	Id string
}


func NewDescribeHostRequesttWithId(id string) *DescribeHostRequest{
	return &DescribeHostRequest{
		Id: id,
	}
}

type UpdateHostRequest struct {
	UpdateMode UPDATE_MODE `json:"update_mode"`
	*Host
}

type UPDATE_MODE string

const (
	// 全量更新
	UPDATE_MODE_PUT UPDATE_MODE = "put"
	// 局部更新
	UPDATE_MODE_PATCH UPDATE_MODE = "patch"
)

func NewPutUpdateHostRequest(id string) *UpdateHostRequest{
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		UpdateMode: UPDATE_MODE_PUT,
		Host: h,
	}
}

func NewPatchUpdateHostRequest(id string) *UpdateHostRequest{
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		UpdateMode: UPDATE_MODE_PATCH,
		Host: h,
	}
}

type DeleteHostRequest struct {
	Id string
}

func NewDeleteHostRequest(id string) *DeleteHostRequest{
	return &DeleteHostRequest{
		Id: id,
	}
}