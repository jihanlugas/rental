package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type targetServer struct {
	Address string
	Port    string
}

type dbServerInfo struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
}

var (
	Debug                  bool
	ListenTo               targetServer
	DBInfo                 dbServerInfo
	CryptoKey              string
	JwtSecretKey           string
	HeaderAuthName         string
	AuthTokenExpiredHour   int64
	SecureServer           bool
	CertificateFilePath    string
	CertificateKeyFilePath string
	MaxSizeUploadPhotoByte int64
	DataPerPage            int
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Failed load env Err: " + err.Error())
		panic(err)
	}

	Debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		fmt.Println("Failed parse DEBUG Err: " + err.Error())
		panic(err)
	}
	ListenTo = targetServer{
		Address: os.Getenv("LISTEN_ADDRESS"),
		Port:    os.Getenv("LISTEN_PORT"),
	}
	DBInfo = dbServerInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DbName:   os.Getenv("DB_NAME"),
	}

	hasher := md5.New()
	hasher.Write([]byte(os.Getenv("CRYPTO_KEY")))
	CryptoKey = hex.EncodeToString(hasher.Sum(nil))

	hasher.Write([]byte(os.Getenv("JWT_SECRET_KEY")))
	JwtSecretKey = hex.EncodeToString(hasher.Sum(nil))

	HeaderAuthName = os.Getenv("HEADER_AUTH_NAME")
	AuthTokenExpiredHour, err = strconv.ParseInt(os.Getenv("AUTH_TOKEN_EXPIRED_HOUR"), 10, 64)
	if err != nil {
		panic(err)
	}

	SecureServer, err = strconv.ParseBool(os.Getenv("SECURE_SERVER"))
	if err != nil {
		fmt.Println("Failed parse SECURE_SERVER Err: " + err.Error())
		panic(err)
	}
	CertificateFilePath = os.Getenv("CERTIFICATE_FILE_PATH")
	CertificateKeyFilePath = os.Getenv("CERTIFICATE_KEY_FILE_PATH")

	MaxSizeUploadPhotoByte, err = strconv.ParseInt(os.Getenv("MAX_SIZE_UPLOAD_PHOTO_BYTE"), 10, 64)
	if err != nil {
		panic(err)
	}

	DataPerPage, err = strconv.Atoi(os.Getenv("MIN_DATA_PER_PAGE"))
	if err != nil {
		panic(err)
	}
}
