package api

import "C"
import (
	"encoding/json"
	"github.com/daemon1024/dokidoki/server/api/utils"
	"github.com/daemon1024/dokidoki/server/entities"
	sider "github.com/daemon1024/dokidoki/server/redis"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
)

func CreateRoom(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name", "doki's room")
	creator := ctx.FormValue("creator", "owner")

	member, err := entities.GenerateMember(creator)
	if err != nil {
		log.Println(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "UUID fail")
	}

	r := entities.Room{
		ID:        id.String(),
		Name:      name,
		CreatedBy: member,
		Status:    entities.PlayerStatus{},
		Members:   []entities.Member{member},
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
	m := conn.Cookies("member")

	if !sider.IsRoomInDb(id) {
		_ = conn.WriteMessage(websocket.CloseMessage, []byte("ROOM DOESN'T EXIST"))
		_ = conn.Close()
		return
	}

	member, err  := entities.GenerateMember(m)
	if err != nil {
		//I have no idea how to properly send a websocket close result. I'm Pardoning myself here - k
		_ = conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		_ = conn.Close()
		return
	}

	r, err := utils.AddMemberToRoom(id, member)
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
			r, err = utils.RemoveMemberFromRoom(id, member.ID)
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