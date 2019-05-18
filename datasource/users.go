package datasource

import (
	"errors"
	"github.com/jigang-duan/gopractices/datamodels"
)

type Engine uint32

const (
	Memory Engine = iota
	Bolt
	MySQL
)

func LoadUsers(engine Engine) (map[int64]datamodels.User, error) {
	if engine != Memory {
		return nil, errors.New("为了简单起见，我们使用一个简单的映射作为数据源")
	}

	return make(map[int64]datamodels.User), nil
}
