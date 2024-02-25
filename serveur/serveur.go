package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
	st "tme4-squelette/client/structures"
	"tme4-squelette/serveur/travaux"
)

var ADRESSE = "localhost"

var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "F"}
var table_association = make([]personne_serv, 0) //tableau d'association d'un identifiant de type entier(ici c'est l'indice) avec une personne_serv
var conn net.Conn

// type d'un paquet de personne stocke sur le serveur, n'implemente pas forcement personne_int (qui n'existe pas ici)
type personne_serv struct {
	// A FAIRE
	identifiant int
	statut   string
	personne st.Personne
	afaire   []func(personne st.Personne) st.Personne
}

// cree une nouvelle personne_serv, est appelé depuis le client, par le proxy, au moment ou un producteur distant
// produit une personne_dist
func creer(id int) *personne_serv {
	p := pers_vide
	pers_serv :=personne_serv{
		identifiant:id,
		statut:"V"
		personne:p,
		afaire:[]func(personne st.Personne) st.Personne{}
	}
	while(len(table_association)<id){
		table_association = append(table_association,pers_vide)	
	}	
	table_association[id] = pers_serv
	return &pers_serv
}

// Méthodes sur les personne_serv, on peut recopier des méthodes des personne_emp du client
// l'initialisation peut être fait de maniere plus simple que sur le client
// (par exemple en initialisant toujours à la meme personne plutôt qu'en lisant un fichier)
func (p *personne_serv) initialise() {
	// A FAIRE
	p.personne = pers_vide
	rand.Seed(time.Now().Unix())
	nb_alea_funs := rand.Intn(5) + 1
	for i := 0; i < nb_alea_funs; i++ {
		p.afaire = append(p.afaire, travaux.UnTravail())
	}
	p.statut = "R"
}

func (p *personne_serv) travaille() {
	// A FAIRE
	if p.statut == "C" || p.statut == "V" || len(p.afaire) == 0 {
		panic("Probleme, aucun travail ne devrait être effectué")
	}
	p.personne = p.afaire[0](p.personne)
	if len(p.afaire) > 0 {
		p.afaire = p.afaire[1:]
	}
	if len(p.afaire) == 0 {
		p.statut = "C"
	}
}

func (p *personne_serv) vers_string() string {
	// A FAIRE
	res := "Nom : " + p.personne.Nom + "\n Prenom : " + p.personne.Prenom + " \n Age : " + fmt.Sprint(p.personne.Age) + "\n Sexe : " + p.personne.Sexe
	return res
}

func (p *personne_serv) donne_statut() string {
	// A FAIRE
	return p.statut
}

// Goroutine qui maintient une table d'association entre identifiant et personne_serv
// il est contacté par les goroutine de gestion avec un nom de methode et un identifiant
// et il appelle la méthode correspondante de la personne_serv correspondante
func mainteneur(id int, methode string) string {
	// A FAIRE
	pers_serv = table_association[id]
	
	if(methode == "initialise"){
		pers_serv.initialise()
		return "ok"
	}

	if(methode == "creer"){
		creer(id)
		return "ok"
	}

	if(methode == "travaille"){
		pers_serv.travaille(&pers_serv)
		return "ok"
	}

	if(methode == "vers_string"){
		return pers_serv.vers_string(&pers_serv)
	}
	
	if(methode == "donne_statut"){
		return pers_serv.donne_statut(&pers_serv)
	}
	
}

// Goroutine de gestion des connections
// elle attend sur la socketi un message content un nom de methode et un identifiant et appelle le mainteneur avec ces arguments
// elle recupere le resultat du mainteneur et l'envoie sur la socket, puis ferme la socket
func gere_connection() {
	// A FAIRE
	defer conn.Close()
	
	for{
		buffer := make([]byte,1024)
		n,err :=conn.Read(buffer)
		if err != nil {
			fmt.Println("Erreur de lecture:", err)
			return
		}
		message := string(buffer[:n]) 
		
		// Extraction de l'identifiant et de la méthode à partir du message
		parts := strings.Split(message, ",")
		if len(parts) != 2 {
			fmt.Println("Message invalide:", message)
			return
		}
		id, err := strconv.Atoi(parts[0]) //on recupere l'identifiant
		if err != nil {
			fmt.Println("Identifiant invalide:", parts[0])
			return
		}
		
		methode := parts[1] //et on recupere la methode a donner au mainteneur

		// Vérification si l'identifiant existe dans la table d'association
		personne_serv, existe := table_association[id]
		if !existe {
			fmt.Println("Identifiant non trouvé:", id)
			return
		}

		res := mainteneur(id,methode)
		_,err = conn.Write([]byte(res))
		if(err != nil){
			fmt.Println("Erreur:",err)
			return
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Format: client <port>")
		return
	}
	port, _ := strconv.Atoi(os.Args[1]) // doit être le meme port que le client
	addr := ADRESSE + ":" + fmt.Sprint(port)
	// A FAIRE: creer les canaux necessaires, lancer un mainteneur
	ln, _ := net.Listen("tcp", addr) // ecoute sur l'internet electronique
	fmt.Println("Ecoute sur", addr)
	for {
		conn, _ := ln.Accept() // recoit une connection, cree une socket
		fmt.Println("Accepte une connection.")
		go gere_connection() // passe la connection a une routine de gestion des connections
	}
}
