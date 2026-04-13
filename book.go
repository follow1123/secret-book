package main

type Book struct {
	Secrets         []Secret         `json:"secrets"`
	ModifiedSecrets []ModifiedSecret `json:"modified_secrets"`
}

type Secret struct {
	Platform   string `json:"platform"`
	Account    string `json:"account"`
	Password   string `json:"password"`
	Remark     string `json:"remark"`
	CreateTime string `json:"create_time"`
}

type ModifiedSecret struct {
	Secret
	ModifiedTime string `json:"modified_time"`
}
