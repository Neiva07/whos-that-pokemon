package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"whos-that-pokemon/controllers"
	_ "whos-that-pokemon/models"

	"github.com/gorilla/mux"
)

//App instance to build the app
type App struct {
	Router *mux.Router
}

//UseRouter Method to initilize all routers of the app
func (app *App) UseRouter() {
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/", mainFunc).Methods("GET")
	app.Router.HandleFunc("/api/users/signup", controllers.SignUp).Methods("POST")
	app.Router.HandleFunc("/api/users/{id}/friendship/{friend_id}", controllers.CreateFriendship).Methods("POST")
	app.Router.HandleFunc("/api/users/{friend_id}/friendship/{id}", controllers.AcceptRequest).Methods("POST")
	app.Router.HandleFunc("/api/users/{id}/game/{friend_id}", controllers.StartGameWithFriend).Methods("POST")
	app.Router.HandleFunc("/api/users/{id}/gameLogs", controllers.RetrieveAllGameLogsFromUser).Methods("GET")
	app.Router.HandleFunc("/api/gameLogs", controllers.Register).Methods("POST")

}

func main() {
	app := App{}

	app.UseRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, app.Router))

}

func mainFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	x := struct {
		Data string
	}{
		Data: "Matheus Otário",
	}
	err := json.NewEncoder(w).Encode(x)
	log.Println(x)

	if err != nil {
		log.Fatal(err)
	}
	return
}
