package bookmanager

type OperationType string

const (
	Modified OperationType = "modified"
	Deleted  OperationType = "deleted"
)

type Book struct {
	Secrets        []Secret        `json:"secrets"`
	HistorySecrets []HistorySecret `json:"history_secrets,omitempty"`
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
	OperationTime string        `json:"operation_time"`
	OperationType OperationType `json:"operation_type"`
}
