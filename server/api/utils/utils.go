package utils

import (
	"encoding/json"
	"github.com/tokidooki/dooki/server/entities"
	sider "github.com/tokidooki/dooki/server/redis"
)

func AddMemberToRoom(roomID string, m entities.Member) (entities.Room, error) {
	r, err := GetRoom(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	r.Members = append(r.Members, m)
	rB, _  := json.Marshal(r)
	if err := sider.SetData(roomID, rB, 0); err != nil {
		return entities.Room{}, err
	}

	return r, nil
}

func RemoveMemberFromRoom(roomID string, memberID string) (entities.Room, error) {
	r, err := GetRoom(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	//Empty a slice
	var mList []entities.Member
	for _, mem := range r.Members {
		//Might just get the ID at this point, will see later.
		if !(mem.ID == memberID) {
			mList = append(r.Members, mem)
		}
	}
	r.Members = mList

	rB, _  := json.Marshal(r)
	if err := sider.SetData(roomID, rB, 0); err != nil {
		return entities.Room{}, err
	}

	return r, nil
}

func IsRoomCreator(roomID string, m entities.Member) (bool, error) {
	r, err := GetRoom(roomID)
	if err != nil {
		return false, err
	}

	if r.CreatedBy.ID == m.ID {
		return true, nil
	}

	return false, nil
}

func GetRoom(roomID string) (entities.Room, error) {
	res, err := sider.GetDataJson(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	var r entities.Room
	_ = json.Unmarshal(res, &r)

	return r, nil
}

func SetRoom(room entities.Room)  error {
	rB, _  := json.Marshal(room)
	if err := sider.SetData(room.ID, rB, 0); err != nil {
		return  err
	}

	return nil
}