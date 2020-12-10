package api

import "C"
import (
	"encoding/base64"
	"encoding/json"
	"github.com/daemon1024/dokidoki/server/api/utils"
	"github.com/daemon1024/dokidoki/server/entities"
	sider "github.com/daemon1024/dokidoki/server/redis"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func CreateRoom(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name", "doki's room")

	//Expecting a base64 encode of Member JSON
	creator := ctx.FormValue("creator", "")

	if creator == "" {
		return fiber.NewError(fiber.StatusNotAcceptable, "FormKey creator not found")
	}

	cB, err := base64.StdEncoding.DecodeString(creator)
	if err != nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "FormKey creator invalid")
	}

	var m entities.Member
	if err := json.Unmarshal(cB, &m); err != nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "FormKey creator invalid")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "UUID fail")
	}

	r := entities.Room{
		ID:        id.String(),
		Name:      name,
		CreatedBy: m,
		Status:    entities.PlayerStatus{},
		Members:   []entities.Member{m},
	}

	if err := r.ToDb(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to add to database")
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
		_ = conn.WriteMessage(websocket.CloseInternalServerErr, []byte("Cookie : member invalid"))
		_ = conn.Close()
		return
	}

	var m entities.Member
	if err := json.Unmarshal(mB, &m); err != nil {
		_ = conn.WriteMessage(websocket.CloseInternalServerErr, []byte("Cookie : member invalid"))
		_ = conn.Close()
		return
	}


	if !sider.IsRoomInDb(id) {
		_ = conn.WriteMessage(websocket.CloseMessage, []byte("ROOM DOESN'T EXIST"))
		_ = conn.Close()
		return
	}

	r, err := utils.AddMemberToRoom(id, m)
	if err != nil {
		_ = conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		_ = conn.Close()
		return
	}

	rb, _ := json.Marshal(r)
	if err := sider.PublishToChannel(id, rb); err != nil {
		_ = conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		_ = conn.Close()
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