package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ContainerDeploy struct {
	Machine int
	Name    string
	Network string
	Ip      string
	Image   string
	Hport   string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Cport   string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Hvolume string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Cvolume string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
}

type DockerPort struct {
	ID          int    `gorm:"AUTO_INCREMENT"`
	Name        string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	MachineName string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Port        string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	ToIp        string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	ToPort      string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Status      int    `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Cron struct {
	ID        int    `gorm:"AUTO_INCREMENT"`
	Unique    string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Name      string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Schedule  string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Cmd       string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Machine   string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Status    int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Task struct {
	ID          int    `gorm:"AUTO_INCREMENT"`
	Name        string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	ProjectID   int
	ProjectName string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	UserName    string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Branch      string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Version     string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Step        int
	Audit       int
	Status      int
	TicketAt    time.Time `gorm:"default:null"`
	AuditAt     time.Time `gorm:"default:null"`
	DeployAt    time.Time `gorm:"default:null"`
	DoneAt      time.Time `gorm:"default:null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Project struct {
	ID         int    `gorm:"AUTO_INCREMENT"`
	Name       string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Repo       string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Deploy     string `gorm:"size:255"`
	DeployName string `gorm:"size:255"`
	Host       string `gorm:"size:255"`
	Audit      int
	RepoReal   string `gorm:"size:255"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DomainAccess struct {
	ID           int    `gorm:"AUTO_INCREMENT"`
	Domain       string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	AccessKeyId  string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	AccessSecret string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type DeployStep struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	ProjectID int
	Title     string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Action    string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID           int    `gorm:"AUTO_INCREMENT"`
	Name         string `gorm:"size:255"`        // string默认长度为255, 使用这种tag重设。
	Email        string `gorm:"unique;size:255"` // string默认长度为255, 使用这种tag重设。
	Password     string `gorm:"size:255"`        // string默认长度为255, 使用这种tag重设。
	PasswordHash string `gorm:"size:255"`
	Status       int
	Role         int
	LastLogin    time.Time `gorm:"default:null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Machine struct {
	ID        int    `gorm:"AUTO_INCREMENT"`
	Name      string `gorm:"size:255"`     // string默认长度为255, 使用这种tag重设。
	User      string `gorm:"size:255"`     // string默认长度为255, 使用这种tag重设。
	Ip        string `gorm:"size:255"`     // string默认长度为255, 使用这种tag重设。
	Port      string `gorm:"size:255"`     // string默认长度为255, 使用这种tag重设。
	Rsa       string `gorm:"type:blob(0)"` // string默认长度为255, 使用这种tag重设。
	Status    int    `gorm:"DEFAULT:1"`
	Type      int    `gorm:"DEFAULT:1"` //1.线上主机 2.容器宿主
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Machine) FindHost(project Project) (machine []Machine) {
	ms := strings.Split(project.Host, ",")
	machines := ""
	for _, m := range ms {
		machines += "'" + m + "',"
	}
	sql := fmt.Sprintf("select * from machines where name in (%s)", machines[0:len(machines)-1])
	DB.Raw(sql).Scan(&machine)
	return
}

type Service struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	Name      string
	Ip        string
	Port      string
	Auto      int `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Result struct {
	Sta  int
	Code int
	Msg  string
	Data interface{}
}

func NewResult(sta int, code int, msg string, data interface{}) []byte {
	result := Result{}
	result.Sta = sta
	result.Code = code
	result.Msg = msg
	result.Data = data
	rs, err := json.Marshal(result)
	if err != nil {
		panic(err.Error())
	}
	return rs
}
