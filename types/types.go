package types

import "time"

type Account struct {
	Email           string      `json:"email"`
	ID              string      `json:"id"`
	UUID            string      `json:"uuid"`
	Role            int         `json:"role"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
	EmailVerifiedAt time.Time   `json:"email_verified_at"`
	PaymentID       interface{} `json:"payment_id"`
	Plan            Plan        `json:"plan"`
	Card            interface{} `json:"card"`
}

type Plan struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Amount   string   `json:"amount"`
	Friendly string   `json:"friendly"`
	Dummy    bool     `json:"dummy"`
	Current  bool     `json:"current"`
	Metadata Metadata `json:"metadata"`
	Ext      string   `json:"ext"`
	Perks    []string `json:"perks"`
	Comped   bool     `json:"comped"`
}

type OnUploadProgress struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	Written int    `json:"written"`
	Total   int    `json:"total"`
	File    string `json:"file"`
	End     bool   `json:"end"`
}

type OnUploadInfo struct {
	Type      string      `json:"type"`
	Config    Config      `json:"config"`
	Certs     []Certs     `json:"certs"`
	Metadata  Metadata    `json:"metadata"`
	Urls      []Urls      `json:"urls"`
	Instances []Instances `json:"instances"`
}
type Config struct {
	Force    interface{} `json:"force"`
	Redirect interface{} `json:"redirect"`
	Cors     interface{} `json:"cors"`
	Hsts     interface{} `json:"hsts"`
	TTL      interface{} `json:"ttl"`
}

type Certs struct {
	Subject         string   `json:"subject"`
	Issuer          string   `json:"issuer"`
	NotBefore       string   `json:"notBefore"`
	NotAfter        string   `json:"notAfter"`
	ExpInDays       int      `json:"expInDays"`
	SubjectAltNames []string `json:"subjectAltNames"`
	CertName        string   `json:"certName"`
	AutoRenew       bool     `json:"autoRenew"`
}

type Metadata struct {
	Rev              int64         `json:"rev"`
	Cmd              string        `json:"cmd"`
	Email            string        `json:"email"`
	Platform         string        `json:"platform"`
	CliVersion       string        `json:"cliVersion"`
	Message          interface{}   `json:"message"`
	BuildTime        interface{}   `json:"buildTime"`
	IP               string        `json:"ip"`
	PrivateFileList  []interface{} `json:"privateFileList"`
	PublicFileCount  int           `json:"publicFileCount"`
	PublicTotalSize  int           `json:"publicTotalSize"`
	PrivateFileCount int           `json:"privateFileCount"`
	PrivateTotalSize int           `json:"privateTotalSize"`
	UploadStartTime  int64         `json:"uploadStartTime"`
	UploadEndTime    int64         `json:"uploadEndTime"`
	UploadDuration   float64       `json:"uploadDuration"`
	Current          bool          `json:"current"`
	Preview          string        `json:"preview"`
}

type Urls struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type Instances struct {
	Type              string `json:"type"`
	Provider          string `json:"provider"`
	Domain            string `json:"domain"`
	Location          string `json:"location"`
	Status            string `json:"status"`
	StatusColor       string `json:"statusColor"`
	Confirmation      string `json:"confirmation"`
	ConfirmationColor string `json:"confirmationColor"`
	IP                string `json:"ip"`
	Info              string `json:"info"`
}

type Teardown struct {
	Msg       string      `json:"msg"`
	NsDomain  string      `json:"nsDomain"`
	Instances []Instances `json:"instances"`
}

type List []struct {
	Domain     string `json:"domain"`
	PlanName   string `json:"planName"`
	Rev        int64  `json:"rev"`
	Cmd        string `json:"cmd"`
	Email      string `json:"email"`
	Platform   string `json:"platform"`
	CliVersion string `json:"cliVersion"`
	Output     struct {
	} `json:"output"`
	Config struct {
	} `json:"config"`
	Message          interface{}   `json:"message"`
	BuildTime        interface{}   `json:"buildTime"`
	IP               string        `json:"ip"`
	PrivateFileList  []interface{} `json:"privateFileList"`
	PublicFileCount  int           `json:"publicFileCount"`
	PublicTotalSize  int           `json:"publicTotalSize"`
	PrivateFileCount int           `json:"privateFileCount"`
	PrivateTotalSize int           `json:"privateTotalSize"`
	UploadStartTime  int64         `json:"uploadStartTime"`
	UploadEndTime    int64         `json:"uploadEndTime"`
	UploadDuration   float64       `json:"uploadDuration"`
	Preview          string        `json:"preview"`
	TimeAgoInWords   string        `json:"timeAgoInWords"`
}

type Token struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type TokenError struct {
	Messages []string          `json:"messages"`
	Details  map[string]string `json:"details"`
}
