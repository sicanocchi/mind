package main

import "fmt"

type Number interface {
	int64 | float64
}

func main() {
	// inizializziamo una mappa per i valori int
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// inizializziamo una mappa per i valori float64
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	// stamperemo i risultati delle funzioni non generiche.
	// in questo caso chiameremo due funzioni diverse
	fmt.Printf("Non-Generic Sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	// in questo caso stamperemo i risultati della funzione generica
	// chiameremo soltanto una funzione invece che due
	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats[string, int64](ints),
		SumIntsOrFloats[string, float64](floats))

	fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
		SumIntsOrFloats(ints),
		SumIntsOrFloats(floats))

	fmt.Printf("Generic Sums with Constraint: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))
}

//------------creiamo funzioni non generiche-------------
// creiamo due funzioni che ricevuto in input una mappa di valori, una int e una float64, ne ricava la somma

func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

//------------creiamo funzioni generiche-------------
/* in questo caso avremo una singola funzione che gestisce sia i valori interi che con la virgola invece che due funzioni distinte.
per supportare questo, quindi avremo bisogno di un modo per dichiarare quali tipi supporta e quali no.
Dall'altra parte, il codiche chiamante avrà bisogno di un modo per dichiarare quale tipo di dati andrà ad inserire.
Per fare ciò quindi, andremo ad inserire oltre i classici parametri della funzione anche un parametro tipo.
E' proprio grazie a questi parametri tipo che la funzione è generica, e consentendola di lavorare con tipi diversi.
Questi parametri tipo inoltre vincolano anche quali tipi di dato sono consentiti per il codice chiamante.
In fase di compilazione il set dei tipi di dati consentiti rappresentya un singolo tipo, ovvero quello fornito come argomento dal
codice chiamante.
*/

// all'interno delle parentesi quadre ci saranno i parametri tipo
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// inquesto caso il vincolo della generalità lo sposteremo all'interfaccia e non più nella funzione, in modo da utilizzarlo in più punti
func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
