package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
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
	// usiamo mysql.Config per raccogliere le proprietà di connessione
	// e lo formattiamo in una stringa DNS per creare una stringa di connessione
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "@localhost:3306",
		DBName: "recordings",
	}
	// utilizziamo sql.Open per inizializzare la variabile DB passando il valore
	// restituito da FormatDSN.
	// verifichiamo la presenza di un errore, come per esempio la non riuscita della connessione
	// al DB in quanto nla stringa di connessione non sia ben formata.
	// per semplicità utilizziamo log.Fatal ma in produzione gestire meglio gli errori
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// chiamiamo DB.Ping per confermare che la connessione al Db funziona, controlliamo l'errore
	// fornito da Ping in caso di connessione fallita altrimenti stampiamo un messaggio di avvenuta connessione
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)
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