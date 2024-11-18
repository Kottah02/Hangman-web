package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

// Structures et variables globales
type Game struct {
	Word           string
	GuessedLetters []string
	Lives          int
	Difficulty     string
	Status         string // "playing", "won", "lost"
	Player         string
}

type Score struct {
	PlayerName string
	Word       string
	Attempts   int
	Result     string
}

var (
	templates = template.Must(template.ParseGlob("templates/*.html"))
	game      *Game
	scores    []Score
	mu        sync.Mutex
)

// Fonction principale
func main() {
	// Routes principales
	http.HandleFunc("/", handleHome)         // Page d'accueil
	http.HandleFunc("/signup", handleSignup) // Page d'inscription
	http.HandleFunc("/start", handleStart)   // Démarrer le jeu
	http.HandleFunc("/game", handleGame)     // Page du jeu
	http.HandleFunc("/guess", handleGuess)   // Deviner une lettre
	http.HandleFunc("/end", handleEnd)       // Fin de partie
	http.HandleFunc("/scores", handleScores) // Tableau des scores

	// Routes additionnelles
	http.HandleFunc("/rules", handleRules)       // Règles du jeu
	http.HandleFunc("/about", handleAbout)       // À propos
	http.HandleFunc("/language", handleLanguage) // Langue du jeu

	// Gestion des fichiers statiques (CSS, JS, images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Route : Page d'accueil (redirige vers la page d'inscription)
func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html") // Charge le fichier index.html
}

// Route : Page d'inscription
// Route : Page d'inscription
func handleSignup(w http.ResponseWriter, r *http.Request) {
	// Serve the signup.html file when the user accesses /signup
	http.ServeFile(w, r, "templates/signup.html")
}

// Route : Démarrer une nouvelle partie
func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	difficulty := r.FormValue("difficulty")

	var word string
	switch difficulty {
	case "easy":
		word = "banane"
	case "medium":
		word = "ordinateur"
	case "hard":
		word = "encyclopedie"
	}

	game = &Game{
		Player:         username,
		Word:           strings.ToUpper(word),
		Lives:          len(word),
		Difficulty:     difficulty,
		Status:         "playing",
		GuessedLetters: []string{},
	}

	http.Redirect(w, r, "/game", http.StatusFound)
}

// Route : Page de jeu
func handleGame(w http.ResponseWriter, r *http.Request) {
	if game.Status != "playing" {
		http.Redirect(w, r, "/end", http.StatusFound)
		return
	}

	hiddenWord := ""
	allGuessed := true
	for _, char := range game.Word {
		if contains(game.GuessedLetters, string(char)) {
			hiddenWord += string(char) + " "
		} else {
			hiddenWord += "_ "
			allGuessed = false
		}
	}

	// Vérifiez si le joueur a deviné toutes les lettres
	if allGuessed {
		game.Status = "won"
		http.Redirect(w, r, "/end", http.StatusFound)
		return
	}

	data := struct {
		Player  string
		Word    string
		Guessed string
		Lives   int
		Guesses []string
	}{
		Player:  game.Player,
		Word:    hiddenWord,
		Guessed: strings.Join(game.GuessedLetters, ", "),
		Lives:   game.Lives,
		Guesses: game.GuessedLetters,
	}

	templates.ExecuteTemplate(w, "game.html", data)
}

// Route : Deviner une lettre ou un mot
func handleGuess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	guess := strings.ToUpper(r.FormValue("guess"))

	if game.Status != "playing" {
		http.Redirect(w, r, "/end", http.StatusFound)
		return
	}

	if len(guess) == 1 && !contains(game.GuessedLetters, guess) {
		game.GuessedLetters = append(game.GuessedLetters, guess)
		if !strings.Contains(game.Word, guess) {
			game.Lives--
		}
	} else if guess == game.Word {
		game.Status = "won"
	} else {
		game.Lives--
	}

	if game.Lives <= 0 {
		game.Status = "lost"
	}

	http.Redirect(w, r, "/game", http.StatusFound)
}

// Route : Fin de partie
// Route : Fin de partie
func handleEnd(w http.ResponseWriter, r *http.Request) {
	if game.Status == "playing" {
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	result := "Vous avez perdu !"
	if game.Status == "won" {
		result = "Félicitations ! Vous avez gagné !"
	}

	// Ajouter le score dans le tableau des scores
	mu.Lock() // Verrouiller l'accès pour éviter des problèmes de concurrence
	scores = append(scores, Score{
		PlayerName: game.Player,
		Word:       game.Word,
		Attempts:   len(game.GuessedLetters),
		Result:     result,
	})
	mu.Unlock()

	data := struct {
		Player  string
		Result  string
		Message string
	}{
		Player:  game.Player,
		Result:  result,
		Message: "Merci d'avoir joué !",
	}

	templates.ExecuteTemplate(w, "result.html", data)
}

// Route : Tableau des scores
func handleScores(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "scores.html", scores)
}

// Route : Règles du jeu
func handleRules(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "rules.html", nil)
}

// Route : À propos
func handleAbout(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "about.html", nil)
}

// Route : Choix de la langue
func handleLanguage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "language.html", nil)
}

// Helper : Vérifie si un élément est présent dans une slice
func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
