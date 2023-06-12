package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "recordings"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// creiamo una variabile db di tipo *sql.DB che rappresenta il punto di accesso ad uno specifico DB
// in produzione eviteremo una variabile globale ma passeremo le variabili alle funzioni che lo
// utilizzano oppure avvolgendola in una struttura
var db *sql.DB

func main() {
	// usiamo Sprintf per raccogliere le info salvate nelle costanti e creare una stringa di connessione
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// utilizziamo sql.Open per inizializzare la variabile DB passando la stringa di connessione
	// verifichiamo la presenza di un errore, come per esempio la non riuscita della connessione
	// al DB in quanto la stringa di connessione non sia ben formata.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// chiamiamo db.Ping per confermare che la connessione al db funziona, controlliamo l'errore
	// fornito da Ping in caso di connessione fallita altrimenti stampiamo un messaggio di avvenuta connessione
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist crea una query in cui recupera gli album che hanno uno specifico artista
func albumsByArtist(name string) ([]Album, error) {
	// una slice di album per prendere i dati ritornati dalla select
	var albums []Album
	// utilizziamo db.Quesry per eseguire una query sul database,
	// il primo parametro è la query che vogliamo eseguire, dopo possiamo passare 0 o più parametri
	// di qualsiasi tipo
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// scansiona le righe restituite tramite rows.Scan per assegnare i valori di ogni riga ai campi della struttura alb
	for rows.Next() {
		var alb Album
		// Scan accetta un elenco di puntatori ai campi in cui verranno scritti i valori di ogni colonna
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		// ad ogni passaggio, dopo aver scansionato ogni riga, la aggiungo al nostro insieme delle righe restituite valide
		albums = append(albums, alb)
	}
	// dopo il ciclo, controllare la presenza di un errore della query complessiva utilizzando rows.Err, questo è l'unico modo
	// per controllare se i risultati sono incompleti
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID crea una query in cui recupera gli album che hanno uno specifico ID
func albumByID(id int64) (Album, error) {
	// utilizziamo una variabile alb di tipo album per salvare le info dell'album che recuperiamo
	var alb Album
	// utilizziamo QueryRow per eseguire una query con SELECT
	// a differenza di Query, QueryRow non restituisce un errore
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum aggiunge un album al db e restituisce il proprio ID
func addAlbum(alb Album) (int64, error) {
	// utilizziamo Exec per eseguire una query con INSERT
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	// recuperiamo l'id della riga dal databese tramite LastInsertId
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
