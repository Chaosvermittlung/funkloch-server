package global

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/crypto/scrypt"
)

var Execdir string
var Conf Config

const saltSize = 32

const storeItemPrefix = 200
const boxPrefix = 202

type DBConnection struct {
	Driver     string
	Connection string
}

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	Execdir = dir + "/"
	Conf.load()
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type Config struct {
	Port       int
	Connection DBConnection
}

func (c *Config) load() error {
	filename := Execdir + "config.json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePasswordHash(p, s string) (string, error) {
	dk, err := scrypt.Key([]byte(p), []byte(s), 16384, 12, 1, 32)
	if err != nil {
		return "", err
	}
	h := hex.EncodeToString(dk)
	return h, err
}

func GenerateSalt() (string, error) {
	buf := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, buf)
	return string(buf), err
}

func CreateItemEAN(id int) string {
	return CreateEAN13(storeItemPrefix, id)
}

func CreateBoxEAN(id int) string {
	return CreateEAN13(boxPrefix, id)
}

func CreateEAN13(prefix, id int) string {
	res := strconv.Itoa(prefix)
	ids := strconv.Itoa(id)
	for i := 1; i < ((12 - len(res)) - len(ids)); i++ {
		res = res + "0"
	}
	res = res + ids
	c := calculateCheckDigit(res)
	return res + c
}

func calculateCheckDigit(num string) string {
	sum := 0
	multiplier := 1
	for _, d := range num {
		di, err := strconv.Atoi(string(d))
		if err != nil {
			return ""
		}
		sum += di * multiplier
		if multiplier == 3 {
			multiplier = 1
		} else {
			multiplier = 3
		}
	}
	return mod(-sum, 10)
}

func mod(x int, y int) string {
	result := x % y
	if result < 0 {
		result += y
	}
	return strconv.Itoa(result)
}
