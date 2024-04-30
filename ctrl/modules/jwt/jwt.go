package jwt

import (
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/sirupsen/logrus"
)

type MyG[T any] struct {
}

func NewMyG[T any]() *MyG[T] {
	return &MyG[T]{}
}

func (m *MyG[T]) GetJwtToken(username string, data T) string {
	s := sjwt.New()
	s.Set("data", data)
	s.SetExpiresAt(time.Now().Add(72 * time.Hour))
	return s.Generate([]byte(username))
}

func (m *MyG[T]) ParseToken(token string) (map[string]interface{}, error) {
	s, err := sjwt.Parse(token)
	if err != nil {
		return nil, err
	}
	a, err := s.Get("data")
	if err != nil {
		return nil, err
	}
	logrus.Print(a)

	return a.(map[string]interface{}), nil
}
