package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
const (

	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "recordings"

)
*/
type Album struct {
	// se usiamo gorm.Model abbiamo già incorporato di default:
	// ID, createdAt, UpdatedAt, DeletedAt
	//gorm.Model
	ID     int64 `gorm:"primaryKey"`
	Title  string
	Artist string
	Price  float32
}

func (album Album) TableName() string {
	return "album"
}

var db *gorm.DB

func main() {

	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
	*/
	//------------come fa il tutorial per connettersi------------da cambiare i vari campi
	//dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	dsn := "host= localhost user= postgres password= root  dbname= recordings  port= 5432  sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	//------------------------------------------------------------

	albums, err := albumsByArtist(`John Coltrane`)
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
	if err := db.Where("artist = ?", name).Find(&albums).Error; err != nil {
		return nil, err
	}
	return albums, nil
}

// albumByID crea una query in cui recupera gli album che hanno uno specifico ID
func albumByID(id int64) (Album, error) {
	// utilizziamo una variabile alb di tipo album per salvare le info dell'album che recuperiamo
	var alb Album
	// utilizziamo QueryRow per eseguire una query con SELECT
	// a differenza di Query, QueryRow non restituisce un errore
	db.Where("id = ?", id).Find(&alb)
	return alb, nil
}

// addAlbum aggiunge un album al db e restituisce il proprio ID
func addAlbum(alb Album) (int64, error) {
	// utilizziamo Exec per eseguire una query con INSERT
	if err := db.Create(&alb).Error; err != nil {
		return 0, err
	}
	// recuperiamo l'id della riga dal databese tramite LastInsertId
	if err := db.Last(&alb).Error; err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return alb.ID, nil
}
