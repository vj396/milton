package types

type DatabaseMetadata struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Type     string `yaml:"type"`
	Port     uint32 `yaml:"port,omitempty"`
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
	CreateOncallRecord(o *Oncall) error
	GetOncallRecords(o *Oncall) ([]Oncall, error)
	DeleteOncallRecord(o *Oncall) error

	CreateInterruptRecord(i *Interrupt) error
	GetInterruptRecords(i *Interrupt) ([]Interrupt, error)
	GetInterruptRecordsForChannel(i *Interrupt) ([]Interrupt, error)
	DeleteInterruptRecord(i *Interrupt) error

	Close()
}
