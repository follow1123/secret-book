package encryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	saltSize int    = 16
	header   uint16 = 0x5bfb
	attempts byte   = 3
)

type Encryptor struct {
	dataPath string
	key      []byte
	pack     *Package
}

func New(dataPath string, password string) (*Encryptor, error) {
	var pack *Package

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		salt, err := generateSalt()
		if err != nil {
			return nil, fmt.Errorf("generate salt error:\n\t%w", err)
		}
		pack = &Package{
			Salt:     salt,
			Attempts: attempts,
		}
	} else {
		data, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, fmt.Errorf("read data path: %s error:\n\t%w", dataPath, err)
		}
		pack, err = unpack(data)
		if err != nil {
			return nil, fmt.Errorf("unpack data error:\n\t%w", err)
		}
	}

	key, err := generateKey(password, pack.Salt)
	if err != nil {
		return nil, fmt.Errorf("genearte key error:\n\t%w", err)
	}

	return &Encryptor{dataPath: dataPath, pack: pack, key: key}, nil
}

func (e *Encryptor) Update(data []byte) error {
	encryptData, err := encrypt(data, e.key)
	if err != nil {
		return fmt.Errorf("encrypt data error:\n\t%w", err)
	}

	e.pack.Data = encryptData
	return nil
}

func (e *Encryptor) GetData() ([]byte, error) {
	// 文件不存在时，默认没有文件
	if e.pack.Data == nil {
		return nil, nil
	}
	decryptedData, err := decrypt(e.pack.Data, e.key)
	if err != nil {
		return nil, fmt.Errorf("decrypt data error:\n\t%w", err)
	}
	return decryptedData, nil
}

func (e *Encryptor) Save() error {
	packedData, err := pack(e.pack)
	if err != nil {
		return fmt.Errorf("pack data error:\n\t%w", err)
	}
	if err := os.WriteFile(e.dataPath, packedData, 0600); err != nil {
		return fmt.Errorf("save to %s error:\n\t%w", e.dataPath, err)
	}
	return nil
}

func (e *Encryptor) AttemptFailed() {
	if e.pack.Data != nil {
		e.pack.Attempts = e.pack.Attempts - 1
		_ = e.Save()
	}
}

func (e *Encryptor) AttemptSucceed() {
	if e.pack.Data != nil {
		if e.pack.Attempts != attempts {
			e.pack.Attempts = attempts
			_ = e.Save()
		}
	}
}

type Package struct {
	Salt     []byte
	Data     []byte
	Attempts byte
}

// 生成盐
func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func unpack(data []byte) (*Package, error) {
	datBuf := bytes.NewBuffer(data)
	headerBytes := make([]byte, 2)

	if _, err := datBuf.Read(headerBytes); err != nil {
		return nil, fmt.Errorf("read header error:\n\t%w", err)
	}

	if header != uint16(headerBytes[0])<<8|uint16(headerBytes[1]) {
		return nil, fmt.Errorf("invalid data")
	}
	saltBytes := make([]byte, saltSize)
	if _, err := datBuf.Read(saltBytes); err != nil {
		return nil, fmt.Errorf("read salt error:\n\t%w", err)
	}

	dataBytes := make([]byte, datBuf.Len()-1)
	if _, err := datBuf.Read(dataBytes); err != nil {
		return nil, fmt.Errorf("read data error:\n\t%w", err)
	}
	atps, err := datBuf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read attempts error:\n\t%w", err)
	}
	if atps == 0 {
		_, _ = rand.Read(dataBytes)
		atps = attempts
	}

	return &Package{
		Salt:     saltBytes,
		Data:     dataBytes,
		Attempts: atps,
	}, nil
}

func pack(pack *Package) ([]byte, error) {
	if len(pack.Salt) > saltSize {
		return nil, fmt.Errorf("invalid salt size")
	}
	var datBuf bytes.Buffer
	headerBytes := []byte{byte(header >> 8), byte(header & 0xFF)}
	if _, err := datBuf.Write(headerBytes); err != nil {
		return nil, fmt.Errorf("write header error:\n\t%w", err)
	}
	if _, err := datBuf.Write(pack.Salt); err != nil {
		return nil, fmt.Errorf("write salt error:\n\t%w", err)
	}
	if _, err := datBuf.Write(pack.Data); err != nil {
		return nil, fmt.Errorf("write data error:\n\t%w", err)
	}
	if err := datBuf.WriteByte(pack.Attempts); err != nil {
		return nil, fmt.Errorf("write attempts error:\n\t%w", err)
	}

	return datBuf.Bytes(), nil
}

// 使用 PBKDF2 密钥派生函数生成密钥
func generateKey(password string, salt []byte) ([]byte, error) {
	k, err := pbkdf2.Key(sha256.New, password, salt, 1000000, 32)
	if err != nil {
		return nil, fmt.Errorf("generate key by password error:\n\t%w", err)
	}
	return k, nil
}

// 使用 CTR 模式进行解密
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	// 提取 IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 使用 AES CTR 模式进行解密
	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// 使用 CTR 模式进行加密
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个新的密文数组，前 16 字节是 IV
	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)

	// 随机生成 IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 使用 AES CTR 模式进行加密
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	// 将 IV 和密文连接一起返回
	return append(iv, ciphertext...), nil
}
