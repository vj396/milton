package types

type DatabaseMetadata struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Type     string `yaml:"type"`
	Port     uint32 `yaml:"port,omitempty"`
}

type Membership struct {
	Id          string
	SLA         uint
	ChannelName string
}

type Documents struct {
	Id    uint
	Title string
	Link  string
}

type Oncall struct {
	Id   string
	Name string
	Type string
}

type Interrupt struct {
	Id          uint
	Item        string
	SubmittedBy string
	SubmittedAt int64
	ChannelId   string
}

type Backend interface {
	CreateMembershipRecord(m *Membership) error
	GetMemebershipRecords(m *Membership) ([]Membership, error)
	DeleteMembershipRecord(m *Membership) error

	CreateDocmentsRecord(d *Documents) error
	GetDocumentsRecords(d *Documents) ([]Documents, error)
	DeleteDocumentsRecord(d *Documents) error

	CreateOncallRecord(o *Oncall) error
	GetOncallRecords(o *Oncall) ([]Oncall, error)
	DeleteOncallRecord(o *Oncall) error

	CreateInterruptRecord(i *Interrupt) error
	GetInterruptRecords(i *Interrupt) ([]Interrupt, error)
	DeleteInterruptRecord(i *Interrupt) error
}
