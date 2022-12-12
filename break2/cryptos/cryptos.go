package cryptos

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"log"
)

// 默认16位密钥
const DEFKEY = "aB3?5678!@#$%^&*"

type AesCrypt struct {
	key []byte
}

func NewAesCrypt() AesCrypt {
	return AesCrypt{
		key: []byte(DEFKEY),
	}
}

func (c *AesCrypt) SetKey(key string) error {
	keys := []byte(key)
	keySize := len(keys)
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return errors.New("key size must be 16 24 32 bit")
	}
	c.key = keys
	return nil
}

func (c *AesCrypt) EnCode(text []byte) ([]byte, error) {
	cText, err := aesEncrypt(text, c.key)
	if err != nil {
		return nil, err
	}
	return cText, nil
}

func (c *AesCrypt) DeCode(text []byte) ([]byte, error) {
	oText, err := aesDecytpt(text, c.key)
	if err != nil {
		return nil, err
	}
	return oText, nil
}

// 处理aes加解密过程
func aesEncrypt(originalBytes, key []byte) ([]byte, error) {
	// 实例化密码器block 参数为密钥
	var block cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 判断加密快的大小
	blockSize := block.BlockSize()
	// 填充
	paddingBytes := pkcs7Padding(originalBytes, blockSize)
	// 初始化加密数据接收切片
	cipherBytes := make([]byte, len(paddingBytes))
	// 使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 执行加密
	blockMode.CryptBlocks(cipherBytes, paddingBytes)

	return cipherBytes, nil
}

func aesDecytpt(cipherBytes, key []byte) ([]byte, error) {
	//创建实例
	var block cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	paddingBytes := make([]byte, len(cipherBytes))
	//执行解密
	blockMode.CryptBlocks(paddingBytes, cipherBytes)
	//去除填充
	originalBytes := pkcs7UnPadding(paddingBytes)

	return originalBytes, nil
}

// padding
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// unpadding
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	if unPadding > length {
		log.Println("pkcs7UnPadding not ok")
		return data
	}
	return data[:(length - unPadding)]
}
