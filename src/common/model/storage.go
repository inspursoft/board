package model

//"time"

type PersistentVolume struct {
	ID         int64       `json:"pv_id" orm:"column(id)"`
	Name       string      `json:"pv_name" orm:"column(name)"`
	Type       int         `json:"pv_type" orm:"column(type)"`
	State      int         `json:"pv_state" orm:"column(state)"`
	Capacity   string      `json:"pv_capacity" orm:"column(capacity)"`
	Accessmode string      `json:"pv_accessmode" orm:"column(accessmode)"`
	Class      string      `json:"pv_class" orm:"column(class)"`
	Readonly   bool        `json:"pv_readonly" orm:"column(readonly)"`
	Reclaim    string      `json:"pv_reclaim" orm:"column(reclaim)"`
	Option     interface{} `json:"pv_options"`
}

type PersistentVolumeOptionNfs struct {
	ID     int64  `json:"pv_id" orm:"column(id)"`
	Path   string `json:"path" orm:"column(path)"`
	Server string `json:"server" orm:"column(server)"`
}

type PaginatedPersistentVolumes struct {
	Pagination           *Pagination         `json:"pagination"`
	PersistentVolumeList []*PersistentVolume `json:"pv_list"`
}

type PersistentVolumeOptionCephrbd struct {
	ID              int64  `json:"pv_id" orm:"column(id)"`
	User            string `json:"user" orm:"column(user)"`
	Keyring         string `json:"keyring" orm:"column(keyring)"`
	Pool            string `json:"pool" orm:"column(pool)"`
	Image           string `json:"image" orm:"column(image)"`
	Fstype          string `json:"fstype" orm:"column(fstype)"`
	Secretname      string `json:"secretname" orm:"column(secretname)"`
	Secretnamespace string `json:"secretnamespace" orm:"column(secretnamespace)"`
	Monitors        string `json:"monitors" orm:"column(monitors)"`
}

const (
	PVUnknown = iota
	PVNFS
	PVCephRBD
)
