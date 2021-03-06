package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"whos-that-pokemon/controllers"
	_ "whos-that-pokemon/models"

	application "whos-that-pokemon/app"

	"github.com/gorilla/mux"
)

//App instance to build the app
type App struct {
	Router *mux.Router
}

//UseRouter Method to initialize all routers of the app
func (app *App) UseRouter() {
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/", mainFunc).Methods("GET")
	app.Router.HandleFunc("/api/users/signin", controllers.SignIn).Methods("POST")

	app.Router.HandleFunc("/api/users/{id}/friendship", controllers.SearchAllFriends).Methods("GET")
	app.Router.HandleFunc("/api/users/{id}/friendship/{friend_id}", controllers.CreateFriendship).Methods("POST")
	app.Router.HandleFunc("/api/users/{friend_id}/friendship/{id}", controllers.AcceptRequest).Methods("PUT")
	app.Router.HandleFunc("/api/users/{id}/friendship/{friend_id}", controllers.DeleteFriendship).Methods("DELETE")

	app.Router.HandleFunc("/api/users/{id}/friendship/{friend_id}/games", controllers.GetAllGamesFromFriends).Methods("GET")
	app.Router.HandleFunc("/api/users/{id}/games/{friend_id}", controllers.StartGameWithFriend).Methods("POST")
	app.Router.HandleFunc("/api/users/{id}/games", controllers.GetAllUserGames).Methods("GET")
	app.Router.HandleFunc("/api/games/{id}", controllers.GetSpecificGame).Methods("GET")
	app.Router.HandleFunc("/api/games/{id}", controllers.UpdateGame).Methods("PUT")

	app.Router.Use(application.Authentication)

}

func main() {
	app := App{}

	app.UseRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("listing on port:", port)

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
