package API

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

// vérifier que la methode est POST et que toutes les données nécessaires sont bien implémenté et envoyé à la base de donnée
func createTopicHandler(w http.ResponseWriter, r *http.Request) {

	//parse data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	//check method
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
//check if all the data are here
	title := r.FormValue("title")
	body := r.FormValue("body")
	status := r.FormValue("status")
	user_id := r.FormValue("user_id") //check maj

	if title == "" || body == "" || status == "" || user_id == "" {
		http.Error(w, "Missing required form fields", http.StatusBadRequest)
		return
	}

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO Topics_Table (title, body, status) VALUES (?, ?, ?)")//check user_id



	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt)
	_, err = stmt.Exec(title, body, status)
	if err != nil {
		return err
	}
	return nil
}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
