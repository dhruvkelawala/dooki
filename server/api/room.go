package api

import "C"
import (
	"encoding/base64"
	"encoding/json"
	"github.com/daemon1024/dokidoki/server/api/utils"
	"github.com/daemon1024/dokidoki/server/entities"
	er "github.com/daemon1024/dokidoki/server/errors"
	sider "github.com/daemon1024/dokidoki/server/redis"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func CreateRoom(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name", "doki's room")

	var creator entities.Member
	if err := ctx.BodyParser(&creator); err != nil {
		return fiber.NewError(fiber.StatusNotAcceptable, er.ErrorParsingBody)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, er.GenUUIDFail)
	}

	r := entities.Room{
		ID:        id.String(),
		Name:      name,
		CreatedBy: creator,
		Status:    entities.PlayerStatus{},
		Members:   []entities.Member{},
	}

	if err := r.ToDb(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, er.DatabaseAddFailed)
	}

	if err := sider.PublishToChannel(id.String(), "Room Created"); err != nil {
		return err
	}

	return ctx.JSON(r)
}


func MyRoom(conn *websocket.Conn) {
	id := conn.Params("id")

	//Expecting a memberJson base64 encoded
	mCookies := conn.Cookies("member")

	mB, err := base64.StdEncoding.DecodeString(mCookies)
	if err != nil {
		closeSocketWithError(conn, er.MemberCookieNotFoundInvalid)
		return
	}

	var m entities.Member
	if err := json.Unmarshal(mB, &m); err != nil {
		closeSocketWithError(conn, er.MemberCookieNotFoundInvalid)
		return
	}


	if !sider.IsRoomInDb(id) {
		closeSocketWithError(conn, er.RoomNotExist)
		return
	}

	r, err := utils.AddMemberToRoom(id, m)
	if err != nil {
		closeSocketWithError(conn, err.Error())
		return
	}

	rb, _ := json.Marshal(r)
	if err := sider.PublishToChannel(id, rb); err != nil {
		closeSocketWithError(conn, err.Error())
		return
	}

	sub := sider.SubscribeToChannel(id)

	go broadcastToClient(sub.Channel(), conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			r, err = utils.RemoveMemberFromRoom(id, m.ID)
			if err != nil {
				_ = conn.Close()
				return
			}
			break
		}
	}

	_ = sider.UnsubscribeChannel(sub, id)
	rb, _ = json.Marshal(r)
	_ = sider.PublishToChannel(id, rb)
}

func ModifyRoomStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var rsm entities.RoomStatusModifier

	if err := ctx.BodyParser(&rsm); err != nil {
		return fiber.NewError(fiber.StatusNotAcceptable, er.ErrorParsingBody)
	}

	if !sider.IsRoomInDb(id) {
		return fiber.NewError(fiber.StatusBadRequest, er.RoomNotExist)
	}

	if ok, _ := utils.IsRoomCreator(id, rsm.Member); !ok  {
		return fiber.NewError(fiber.StatusBadRequest, er.NotRoomCreator)
	}

	r, err := utils.GetRoom(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	r.Status = rsm.Status

	if err := utils.SetRoom(r); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	rb, _ := json.Marshal(r)
	if err := sider.PublishToChannel(id, rb); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func closeSocketWithError(conn *websocket.Conn, err string, errType ...int) {
	var eType int
	if len(errType) == 0 {
		eType = websocket.CloseInternalServerErr
	} else {
		eType = errType[0]
	}



	_ = conn.WriteMessage(eType, []byte(err))
	_ = conn.Close()
}

func broadcastToClient(channel <-chan *redis.Message, conn *websocket.Conn) {
	for {
		select {
		case m, ok := <- channel:
			if !ok {
				break
			}

			var r entities.Room
			_ = json.Unmarshal([]byte(m.Payload), &r)
			_ = conn.WriteJSON(r)
		}
	}
}