package API

import "net/http"

// Handler for user logout
// verifier que la methode est POST et que l'utilisateur est connecté et que toute les donées necessaire sont bien implente
func createTopicHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
