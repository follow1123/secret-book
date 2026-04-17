package bookmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BookManager struct {
	book     *Book
	bookPath string
}

func New(bookPath string) (*BookManager, error) {
	_, err := os.Stat(bookPath)
	notExists := false
	if err != nil {
		if os.IsNotExist(err) {
			notExists = true
		} else {
			return nil, fmt.Errorf("check book path %s error:\n\t%w", bookPath, err)
		}
	}

	book := &Book{}
	if !notExists {
		data, err := os.ReadFile(bookPath)
		if err != nil {
			return nil, fmt.Errorf("read book path: %s error:\n\t%w", bookPath, err)
		}
		if err := json.Unmarshal(data, book); err != nil {
			return nil, fmt.Errorf("unmarshal book path: %s error:\n\t%w", bookPath, err)
		}
	}

	return &BookManager{book: book, bookPath: bookPath}, nil
}

func (m *BookManager) Save() error {
	data, err := json.MarshalIndent(m.book, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal data to json error:\n\t%w", err)
	}
	if err := os.WriteFile(m.bookPath, data, 0664); err != nil {
		return fmt.Errorf("save to %s error:\n\t%w", m.bookPath, err)
	}
	return nil
}

func (m *BookManager) ListPlatforms() []string {
	platformMap := make(map[string]struct{})

	for _, secret := range m.book.Secrets {
		if _, exists := platformMap[secret.Platform]; !exists {
			platformMap[secret.Platform] = struct{}{}
		}
	}

	var platforms []string

	for key := range platformMap {
		platforms = append(platforms, key)
	}
	return platforms
}

func (m *BookManager) ListByPlatform(platform string) []Secret {
	secrets := make([]Secret, 0)
	for _, secret := range m.book.Secrets {
		if secret.Platform == platform {
			secrets = append(secrets, secret)
		}
	}

	if len(secrets) > 1 {
		slices.SortFunc(secrets, func(a, b Secret) int {
			aTime, err := time.Parse("2006-01-02 15:04:05", a.CreateTime)
			if err != nil {
				panic(fmt.Errorf("parse time %s error:\n\t%w", a.CreateTime, err))
			}

			bTime, err := time.Parse("2006-01-02 15:04:05", b.CreateTime)
			if err != nil {
				panic(fmt.Errorf("parse time %s error:\n\t%w", b.CreateTime, err))
			}
			if aTime.After(bTime) {
				return -1
			} else if aTime.Before(bTime) {
				return 1
			}
			return 0
		})
	}
	return secrets
}

func (m *BookManager) ListHistory(secret Secret) []HistorySecret {
	historySecrets := make([]HistorySecret, 0)
	var (
		hasPlatformCond = secret.Platform != ""
		hasAccountCond  = secret.Account != ""
		hasPasswordCond = secret.Password != ""
		hasRemarkCond   = secret.Remark != ""
	)

	for _, hs := range m.book.HistorySecrets {
		if hasPlatformCond {
			if hs.Platform == secret.Platform {
				historySecrets = append(historySecrets, hs)
			}
		}
		if hasAccountCond {
			if strings.Contains(hs.Account, secret.Account) {
				historySecrets = append(historySecrets, hs)
			}
		}
		if hasPasswordCond {
			if strings.Contains(hs.Password, secret.Password) {
				historySecrets = append(historySecrets, hs)
			}
		}
		if hasRemarkCond {
			if strings.Contains(hs.Remark, secret.Remark) {
				historySecrets = append(historySecrets, hs)
			}
		}
		if !(hasPlatformCond || hasAccountCond || hasPasswordCond || hasRemarkCond) {
			historySecrets = append(historySecrets, hs)
		}
	}

	slices.SortFunc(historySecrets, func(a, b HistorySecret) int {
		aTime, err := time.Parse("2006-01-02 15:04:05", a.OperationTime)
		if err != nil {
			panic(fmt.Errorf("parse time %s error:\n\t%w", a.OperationTime, err))
		}

		bTime, err := time.Parse("2006-01-02 15:04:05", b.OperationTime)
		if err != nil {
			panic(fmt.Errorf("parse time %s error:\n\t%w", b.OperationTime, err))
		}
		if aTime.After(bTime) {
			return -1
		} else if aTime.Before(bTime) {
			return 1
		}
		return 0
	})
	return historySecrets
}

func (m *BookManager) GetSecretByIdPerfix(idPrefix string) (*Secret, error) {
	idx, err := m.GetSecretIndexByIdPerfix(idPrefix)
	if err != nil {
		return nil, err
	}
	return &m.book.Secrets[idx], nil
}

func (m *BookManager) GetSecretIndexByIdPerfix(idPrefix string) (int, error) {
	index := -1
	for i, s := range m.book.Secrets {
		if strings.HasPrefix(s.Id, idPrefix) {
			if index >= 0 {
				return index, fmt.Errorf("duplicated id prefix %s", idPrefix)
			}
			index = i
		}
	}
	if index < 0 {
		return index, fmt.Errorf("data for id prefix %s is not exists", idPrefix)
	}
	return index, nil
}

func (m *BookManager) GetByPlatformAccount(platform string, account string) *Secret {
	for _, secret := range m.book.Secrets {
		if secret.Platform == platform && secret.Account == account {
			return &secret
		}
	}
	return nil
}

func (m *BookManager) Add(secret Secret) error {
	if strings.TrimSpace(secret.Platform) == "" {
		return fmt.Errorf("platform cannot be empty")
	}
	if strings.TrimSpace(secret.Account) == "" {
		return fmt.Errorf("account cannot be empty")
	}

	if m.GetByPlatformAccount(secret.Platform, secret.Account) != nil {
		return fmt.Errorf("account: %s is duplicated on the platform: %s", secret.Account, secret.Platform)
	}

	secret.Id = nextId()
	secret.CreateTime = currentTime()

	m.book.Secrets = append(m.book.Secrets, secret)
	return nil
}

func (m *BookManager) deleteByIndex(index int, saveToHistory bool) {
	deleteSecret := m.book.Secrets[index]
	m.book.Secrets = slices.Delete(m.book.Secrets, index, index+1)

	if saveToHistory {
		// 添加到历史列表内
		historySecret := HistorySecret{
			Secret:        deleteSecret,
			OperationTime: currentTime(),
			OperationType: Deleted,
		}
		m.book.HistorySecrets = append(m.book.HistorySecrets, historySecret)

		// 重置id
		historySecret.Id = nextId()
	}
}

func (m *BookManager) DeleteByIdPrefix(idPrefix string, saveToHistory bool) error {
	idx, err := m.GetSecretIndexByIdPerfix(idPrefix)
	if err != nil {
		return err
	}
	m.deleteByIndex(idx, saveToHistory)
	return nil
}

func (m *BookManager) DeleteHistoryByIdPrefix(idPrefix string) error {
	index := -1
	for i, hs := range m.book.HistorySecrets {
		if strings.HasPrefix(hs.Id, idPrefix) {
			if index >= 0 {
				return fmt.Errorf("duplicated id prefix %s", idPrefix)
			}
			index = i
		}
	}
	if index < 0 {
		return fmt.Errorf("data for id prefix %s is not exists", idPrefix)
	}
	m.book.HistorySecrets = slices.Delete(m.book.HistorySecrets, index, index+1)
	return nil
}

func (m *BookManager) updateByIndex(index int, secret Secret) {
	platform := strings.TrimSpace(secret.Platform)
	account := strings.TrimSpace(secret.Account)
	password := strings.TrimSpace(secret.Password)
	remark := strings.TrimSpace(secret.Remark)

	historySecret := HistorySecret{
		Secret:        m.book.Secrets[index],
		OperationTime: currentTime(),
		OperationType: Modified,
	}
	// 重置id
	historySecret.Id = nextId()

	updateFields := 0
	if platform != "" {
		m.book.Secrets[index].Platform = platform
		updateFields += 1
	}
	if account != "" {
		m.book.Secrets[index].Account = account
		updateFields += 1
	}
	if password != "" {
		m.book.Secrets[index].Password = password
		updateFields += 1
	}
	if remark != "" {
		m.book.Secrets[index].Remark = remark
		updateFields += 1
	}

	// 有属性被修改才添加到历史列表内
	if updateFields > 0 {
		// 添加到历史列表内
		m.book.HistorySecrets = append(m.book.HistorySecrets, historySecret)
	}
}

func (m *BookManager) UpdateById(id string, secret Secret) error {
	var updateIdx int = -1
	for i, s := range m.book.Secrets {
		if s.Id == id {
			updateIdx = i
		}
	}
	if updateIdx < 0 {
		return fmt.Errorf("data for id %s is not exists", id)
	}

	m.updateByIndex(updateIdx, secret)
	return nil
}

func (m *BookManager) UpdateByIdPrefix(idPrefix string, secret Secret) error {
	idx, err := m.GetSecretIndexByIdPerfix(idPrefix)
	if err != nil {
		return err
	}

	m.updateByIndex(idx, secret)
	return nil
}

func currentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func nextId() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func DefaultSecretsFile() string {
	path, err := os.Getwd()
	if err != nil {
		panic("read current working directory error")
	}
	return filepath.Join(path, "secrets.json")
}
