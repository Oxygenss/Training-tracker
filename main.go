package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Plate struct {
	Id int
	Weight float64
	Amount int
}

type Train struct {
	Id int
	Name string
	Info string
	Img string
}

type history struct {
	Id int
	Name string
	Number_of_repetitions int
	Number_of_approaches int
	Img string
	Weight float64
	Handle float64
}


func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home_page/home_page.html")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var trains = []Train{}

	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `training_list`")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var train Train 
		err = rows.Scan(&train.Id, &train.Name, &train.Info, &train.Img)
		if err!= nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		trains = append(trains, train)
	}
	tmpl.ExecuteTemplate(w, "home_page", trains)
}

func Add(w http.ResponseWriter, r *http.Request) {
	weightStr := r.FormValue("weight")
	amountStr := r.FormValue("amount")

	weight, err := strconv.ParseFloat(weightStr, 64)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	

	var plates = []Plate{}
	
	rows, err := db.Query("SELECT * FROM `plates`")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var plate Plate 
		err = rows.Scan(&plate.Id, &plate.Weight, &plate.Amount)
		if err!= nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		plates = append(plates, plate)
	}

	flag := false

	for i := 0; i < len(plates); i++ {
		if weight == plates[i].Weight {
			flag = true
			plates[i].Amount += amount
			_, err = db.Exec(fmt.Sprintf("UPDATE `Plates` SET `amount` = %d WHERE id = %d", plates[i].Amount, plates[i].Id))
			if err!= nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	if !flag {
		_, err = db.Exec(fmt.Sprintf("INSERT INTO `Plates` (`weight`, `amount`) VALUES ('%f', '%d')", weight, amount))
		if err!= nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/Plates", http.StatusSeeOther)

}

func Add_training(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	var train history 
	err = db.QueryRow("SELECT id, name, photo FROM `training_list` WHERE id = ?", id).Scan(&train.Id, &train.Name, &train.Img)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(train.Name)
	
	repetitionsStr :=r.FormValue("number_of_repetitions")
	approachesStr := r.FormValue("number_of_approaches")
	weightStr := r.FormValue("weight")
	handleStr := r.FormValue("handle")

	train.Weight, err = strconv.ParseFloat(weightStr, 64)
	if err!= nil {
		http.Error(w, "Неверный формат веса", http.StatusBadRequest)
		return
	}

	train.Handle, err = strconv.ParseFloat(handleStr, 64)
	if err!= nil {
		http.Error(w, "Неверный формат веса", http.StatusBadRequest)
		return
	}

	train.Number_of_repetitions, err = strconv.Atoi(repetitionsStr)
	if err!= nil {
		http.Error(w, "Неверный формат количества", http.StatusBadRequest)
		return
	}

	train.Number_of_approaches, err = strconv.Atoi(approachesStr)
	if err!= nil {
		http.Error(w, "Неверный формат количества", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO `history` (`name`, `number_of_repetitions`, `number_of_approaches`, `img`, `weight`, `handle`) VALUES ('%s', '%d', '%d', '%s', '%f', '%f')", train.Name, train.Number_of_repetitions, train.Number_of_approaches, train.Img, train.Weight, train.Handle))
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/history", http.StatusSeeOther)

}

func Plates(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/plates/plates.html")

	var plates = []Plate{}

	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `plates`")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var plate Plate 
		err = rows.Scan(&plate.Id, &plate.Weight, &plate.Amount)
		if err!= nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		plates = append(plates, plate)
	}
	
	tmpl.ExecuteTemplate(w, "plates", plates)
}

func History(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/history/history.html")

	var trains = []history{}

	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `history`")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var his history 
		err = rows.Scan(&his.Id, &his.Name, &his.Number_of_repetitions, &his.Number_of_approaches, &his.Img, &his.Weight, &his.Handle)
		if err!= nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		trains = append(trains, his)
	}
	
	tmpl.ExecuteTemplate(w, "history", trains)
}



func Calculate(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/сalculate/сalculate.html")
	weightStr := r.FormValue("work_weight")
	handleStr := r.FormValue("handle")
	fmt.Println(weightStr)
	weight, err := strconv.ParseFloat(weightStr, 64)

	if err!= nil {
		http.Error(w, "Неверный формат веса", http.StatusBadRequest)
		return
	}

	handle, err := strconv.ParseFloat(handleStr, 64)
	if err!= nil {
		http.Error(w, "Неверный формат количества", http.StatusBadRequest)
		return
	}

	if handle < 0 || weight <= 0 {
		http.Error(w, "Отрицательные значения", http.StatusBadRequest)
		return
	}

	if handle + 1 > weight {
		http.Error(w, "Слишком большая рукоять", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		panic(err)
	}
	defer db.Close()

	

	var plates = []Plate{}
	
	rows, err := db.Query("SELECT * FROM `plates`")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var plate Plate 
		err = rows.Scan(&plate.Id, &plate.Weight, &plate.Amount)
		if err != nil {
			panic(err)
		}
		plates = append(plates, plate)
	}

	sort.Slice(plates, func(i, j int) bool {
		return plates[i].Weight > plates[j].Weight
	})

	result := make(map[float64]int, len(plates))

	for i := 0; i < len(plates); i++ {
		result[plates[i].Weight] = 0

	}

	var sum float64

	value := (weight - handle - 1) / 2

	for i := 0; i < len(plates); i++ {
		for plates[i].Weight <= value && plates[i].Amount > 1 {
			result[plates[i].Weight] += 2
			value -= plates[i].Weight
			plates[i].Amount -= 2
			sum += plates[i].Weight

		}
	}

	var minWeight float64
	flag := false

	for i := 0; i < len(plates); i++ {
		if plates[i].Amount > 1 {
			flag = true
			minWeight = plates[i].Weight
			break
		}
	}

	for i := 0; i < len(plates); i++ {
		if plates[i].Weight < minWeight &&  plates[i].Amount > 1 {
			minWeight = plates[i].Weight
		}
	}

	if !flag {
		minWeight = -1
	}


	over := sum * 2 + minWeight * 2 + 1

	fmt.Println(over, weight)

	tmpl.ExecuteTemplate(w, "result", map[string]interface{}{
        "result": result,
        "value": value * 2,
		"minWeight": minWeight,
		"sum": sum * 2 + 1,
		"over": over - (weight - handle),
    })
}

func Delete_repeat(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET number_of_repetitions = CASE WHEN number_of_repetitions > 0 THEN number_of_repetitions - 1 ELSE number_of_repetitions END WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Add_repeat(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET number_of_repetitions = number_of_repetitions + 1 WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Delete_approach(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET number_of_approaches = CASE WHEN number_of_approaches > 0 THEN number_of_approaches - 1 ELSE number_of_approaches END WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Add_approach(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET number_of_approaches = number_of_approaches + 1 WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Delete_weight(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET weight = CASE WHEN weight > 0 THEN weight - 1 ELSE weight END WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Add_weight(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET weight = weight + 1 WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Delete_handle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET handle = CASE WHEN handle > 0 THEN handle - 1 ELSE handle END WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Add_handle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
    //fmt.Println(id)
	db, err := sql.Open("sqlite3", "./data")
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE history SET handle = handle + 1 WHERE id =?", id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]


    db, err := sql.Open("sqlite3", "./data")
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("DELETE FROM history WHERE id =?", id)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func Copy(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    db, err := sql.Open("sqlite3", "./data")
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var train history 
    err = db.QueryRow("SELECT name, number_of_repetitions, number_of_approaches, img, weight, handle FROM `history` WHERE id =?", id).Scan(&train.Name, &train.Number_of_repetitions, &train.Number_of_approaches, &train.Img, &train.Weight, &train.Handle)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// var temp history 

	// repetitionsStr :=r.FormValue("number_of_repetitions")
	// approachesStr := r.FormValue("number_of_approaches")
	// weightStr := r.FormValue("weight")
	// handleStr := r.FormValue("handle")

	// temp.Weight, err = strconv.ParseFloat(weightStr, 64)
	// if err!= nil {
	// 	http.Error(w, "Неверный формат веса", http.StatusBadRequest)
	// 	return
	// }

	// temp.Handle, err = strconv.ParseFloat(handleStr, 64)
	// if err!= nil {
	// 	http.Error(w, "Неверный формат веса", http.StatusBadRequest)
	// 	return
	// }

	// temp.Number_of_repetitions, err = strconv.Atoi(repetitionsStr)
	// if err!= nil {
	// 	http.Error(w, "Неверный формат количества", http.StatusBadRequest)
	// 	return
	// }

	// temp.Number_of_approaches, err = strconv.Atoi(approachesStr)
	// if err!= nil {
	// 	http.Error(w, "Неверный формат количества", http.StatusBadRequest)
	// 	return
	// }

	// fmt.Println(temp.Number_of_approaches, temp.Number_of_repetitions, temp.Handle, temp.Weight)


    _, err = db.Exec("INSERT INTO history (name, number_of_repetitions, number_of_approaches, img, weight, handle) VALUES (?,?,?,?,?,?)", train.Name, train.Number_of_repetitions, train.Number_of_approaches, train.Img, train.Weight, train.Handle)
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	http.Redirect(w, r, "/history", http.StatusSeeOther)
}



func handleFunc() {
    r := mux.NewRouter()

	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	r.HandleFunc("/", HomePage)
	r.HandleFunc("/Add", Add)
	r.HandleFunc("/Plates", Plates)
	r.HandleFunc("/history", History)
	r.HandleFunc("/Calculate", Calculate)
	r.HandleFunc("/Copy/{id:[0-9]+}", Copy)
	r.HandleFunc("/Delete/{id:[0-9]+}", Delete)
	r.HandleFunc("/Add_training/{id:[0-9]+}", Add_training)
	r.HandleFunc("/Delete_repeat/{id:[0-9]+}", Delete_repeat)
	r.HandleFunc("/Add_repeat/{id:[0-9]+}", Add_repeat)

	r.HandleFunc("/Delete_approach/{id:[0-9]+}", Delete_approach)
	r.HandleFunc("/Add_approach/{id:[0-9]+}", Add_approach)

	r.HandleFunc("/Delete_weight/{id:[0-9]+}", Delete_weight)
	r.HandleFunc("/Add_weight/{id:[0-9]+}", Add_weight)

	r.HandleFunc("/Delete_handle/{id:[0-9]+}", Delete_handle)
	r.HandleFunc("/Add_handle/{id:[0-9]+}", Add_handle)
	http.Handle("/", r)

	

	http.ListenAndServe("localhost:8080", r)
}


func main() {
	handleFunc()
}

	