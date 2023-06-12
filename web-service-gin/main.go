package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	// inizializziamo un router Gin utilizando Default
	router := gin.Default()
	// utilizziamo GET per associare il metodo HTTP GET e il percorso di albums alla funzione handler getAlbums
	router.GET("/albums", getAlbums)
	// utilizziamo POST per associare il metodo HTTP POST e il percorso di albums allafunzione handler postAlbums
	router.POST("/albums", postAlbums)
	// utilizziamo GET per associare il metodo HTTP GET e il percorso di albumsy/:id alla funzione handler getAlbumBYID
	// in Gin i due punti che precedono un elemento nel percorso indicano che esso è un parametro
	router.GET("/albums/:id", getAlbumByID)
	// utilizziamo la funzione Run per collegare il router a http.server e avviare il server
	router.Run("localhost:8080")
}

// getAlbums risponde con la lista di tutti gli album come JSON.
// gin.Context trasposta i dettagli della richiesta, covalida e serializza JSON
func getAlbums(c *gin.Context) {
	// utilizziamo Context.IndentedJSON per serializzare la struttura in JSON e aggiungerla alla risposta
	// il primo argomento è il codice di stato HTTP che vogliamo inviare al client
	// si potrebbe utilizzare anche Context.JSON ma l'altra è più facile da usare durante il debug ed è più piccola
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums aggiunge un album dal JSON che riceve nel corpo della richiesta.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// utilizziamo Context.BindJSON per associare il corpo della richiesta all'indirizzo di newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// aggiungi alla struct albums già esistente la struttura iniziaallizata dalla richiesta JSON
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID individua l'album il cui ID corrisponde all'ID inviaton dal client e restituisce
// l'album come risposta
func getAlbumByID(c *gin.Context) {
	// utilizziamo Context.Param per recuperare l'ID dall'URL
	// quando mappiamo questo handler ad un percorso dovremmo includere un placeholder per questo parametro
	id := c.Param("id")

	// in questo modo scansioniamo tutte le strutture Album all'interno di Albums per vedere se ne esiste una con l'id che corrisponde
	// all'id passato nella richiesta. Questo servizio nel classico utilizzo verrebbe implementato con una query di database.
	// se viene trovata una corrispondenza invieremo l'intera struttura alla risposta altrimenti invieremo il codice $=$ che non è stata
	// trovata nessuna corrispondenza
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
