package utils

import (
	"encoding/json"
	"github.com/daemon1024/dokidoki/server/entities"
	sider "github.com/daemon1024/dokidoki/server/redis"
)

func AddMemberToRoom(roomID string, m entities.Member) (entities.Room, error) {
	res, err := sider.GetDataJson(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	var r entities.Room
	_ = json.Unmarshal(res, &r)

	r.Members = append(r.Members, m)
	rB, _  := json.Marshal(r)
	if err := sider.SetData(roomID, rB, 0); err != nil {
		return entities.Room{}, err
	}

	return r, nil
}

func RemoveMemberFromRoom(roomID string, memberID string) (entities.Room, error) {
	res, err := sider.GetDataJson(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	var r entities.Room
	_ = json.Unmarshal(res, &r)

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