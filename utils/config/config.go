package config

type Config interface {
	LoadFullPath(filePath string, data interface{}) error
	SaveFullPath(filePath string, data interface{}) error
	Load(filePath, fileName string, data interface{}, save ...bool) error
	Save(filePath, fileName string, data interface{}) error
}
