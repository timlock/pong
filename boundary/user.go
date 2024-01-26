package boundary

import (
	"encoding/json"
	"log"
	"net/http"
	"pong/entity"
	"pong/repository"

	"github.com/google/uuid"
)

type UserHandler struct {
	userStore *repository.UserStore
}

func NewUserHandler(u *repository.UserStore) UserHandler {
	return UserHandler{userStore: u}
}

func (u UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name :=	r.URL.Query().Get("name")
	id := uuid.New()
	player := entity.Player{Id: id, Name: name}
	u.userStore.Add(&player)
	content, err := json.Marshal(player.Id)
	if err != nil {
		log.Println("Could not create json object of player: ", player, " err: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		log.Println("Registered new player: ", player.Name)
		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
	}

}
