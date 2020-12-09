package entities

import (
	"encoding/json"
	"errors"
	er "github.com/daemon1024/dokidoki/server/errors"
	sider "github.com/daemon1024/dokidoki/server/redis"
	uuid2 "github.com/google/uuid"
	"time"
)

type Room struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	CreatedBy Member       `json:"created_by"`
	Status    PlayerStatus `json:"status"`
	Members   []Member     `json:"members"`
}

type PlayerStatus struct {
	Type        string    `json:"type"`
	Status      int       `json:"status"`
	CurrentTime time.Time `json:"time"`
	Duration    time.Time `json:"duration"`
	Data        string    `json:"data"`
	By          string    `json:"by"`
}

type Member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r Room) ToDb() error {
	if r.Name == "" || r.ID == "" {
		return errors.New(er.NameIdUnspecified)
	}

	rb, _ := json.Marshal(r)
	return sider.SetData(r.ID, rb, 0)
}

func GenerateMember(name string) (Member, error) {
	uuid, err := uuid2.NewRandom()
	if err != nil {
		return Member{}, err
	}

	return Member{
		ID: uuid.String(),
		Name: name,
	}, nil
}
