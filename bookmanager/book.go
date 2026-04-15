package bookmanager

type Book struct {
	Secrets        []Secret        `json:"secrets"`
	HistorySecrets []HistorySecret `json:"modified_secrets,omitempty"`
}

type Secret struct {
	Id         string `json:"id"`
	Platform   string `json:"platform"`
	Account    string `json:"account"`
	Password   string `json:"password"`
	Remark     string `json:"remark"`
	CreateTime string `json:"create_time"`
}

type HistorySecret struct {
	Secret
	ModifiedTime string `json:"modified_time"`
	DeletedTime  string `json:"deleted_time"`
}
