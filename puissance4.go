
package main

//functional-go
//import "github.com/thoas/go-funk"
import (
	"fmt"
	"errors"
//	"math"
	"strings"
	"bufio"
	"os"
	"strconv"
)


/***************************************\
|*              JETONS                 *|
\***************************************/
type Jeton int

const (
	Vide  Jeton = 0
	Jaune Jeton = 1
	Rouge Jeton = 2
)

func (j Jeton) estVide() bool{
	return j == Vide
}

func (j Jeton) getRune() rune{
		if j== Rouge {
			return 'O'
		}else if j == Jaune {
			return '#'
		}else{
			return ' '
		}
}

func (j Jeton) Equals(k Jeton) bool{
	if j == k {
		return true
	}
	return false
}

func (j Jeton) getStr() string{
	return fmt.Sprintf("%c", j.getRune())
}

//Retourne la couleur adverse du jeton
func (j Jeton) adversaire() Jeton {
	if j == Jaune {
		return Rouge
	}
	if j == Rouge {
		return Jaune
	}
	return j
}



/***************************************\
|*              GRILLE                 *|
\***************************************/
type Grille struct{
	Col int
	Row int
	Grid [][]Jeton
}

func (g Grille) JetonAt(col int, row int) Jeton{
	return g.Grid[col][row]
}

func (g Grille) RemplaceJetonAt(jeton Jeton, col int, row int) Grille {
	g.Grid[col][row] = jeton
	return g
}

func (g Grille) InsertJetonAt(jeton Jeton, col int, row int) Grille{
	grd := g.Copy()
	return grd.RemplaceJetonAt(jeton, col, row)
}

func (g Grille) Copy() Grille{
	grd := make([][]Jeton, g.Col)
	for i := 0; i < g.Col; i++ {
		grd[i] = make([]Jeton, g.Row)
		for j := 0; j < g.Row; j++{
			grd[i][j] = g.JetonAt(i, j)
		}
	}
	return Grille{
		Col: g.Col,
		Row: g.Row,
		Grid: grd,
	}
}

//Renvoie les colones, les lignes et les diagonales
func (g Grille) lines() [][]Jeton {
	//On copie les colones
	lines := g.Copy().Grid

	//On copie les lignes
	for j:= 0; j < g.Row; j++ {
		l := []Jeton{}
		for i:=0; i < g.Col; i++ {
			l = append(l,  g.JetonAt(i, j))
		}
		lines = append(lines, l)
	}

	//Il nous reste a copier les diagonales

	//diagonales en parcourant les colones
	for i:= 0; i < g.Col; i++ {
		l := []Jeton{}
		for idx:= 0; i + idx < g.Col && idx < g.Row; idx++ {
			l = append(l, g.JetonAt(i + idx, idx))
		}
		lines = append(lines, l)
	}

	//diagonales en parcourant les lignes
	for i:= 1; i< g.Row; i++ {
		l := []Jeton{}
		for idx := 0;   i + idx < g.Row && idx < g.Col; idx++ {
			l = append(l, g.JetonAt(idx, i + idx))
		}
		lines = append(lines, l)
	}

	//diagonales en descendant en parcourant les colones
	for i := 0; i < g.Col; i++ {
		l := []Jeton{}
		for idx := 0; g.Row - idx -1 >= 0 && i + idx <  g.Col; idx++ {
			//fmt.Println("Col:", i+idx, "Row:", g.Row - idx -1)
			l = append(l, g.JetonAt(i+idx, g.Row - idx -1))
		}
		lines = append(lines, l)
	}

	//diagonales en descendant en parcourant les lignes
	for i := 0; i <= g.Row - 2; i++ {
		l := []Jeton{}
		for idx := 0; g.Row - 2 - i - idx >= 0 && idx < g.Col; idx++ {
			l = append(l, g.JetonAt( idx, g.Row - 2 -i - idx))
		}
		lines = append(lines, l)
	}

	return lines
}

func (g Grille) ColIsFull(col int) bool {
	if col >= g.Col || col < 0{
		return true
	}
	return !g.JetonAt(col, g.Row - 1).estVide()
}

func (g Grille) NextVide(col int) (int, error){
	if col > g.Col || col < 0{
		strError := fmt.Sprintf("error colone %d out of bounds (0,%d)", col, g.Col)
		return 0, errors.New(strError)
	}
	if g.ColIsFull(col){
		strError := fmt.Sprintf("error colone %d is full", col)
		return 0, errors.New(strError)
	}
	for j := 0; j < g.Row; j++{
		if g.JetonAt(col, j).estVide(){
			return j, nil
		}
	}
	strError := fmt.Sprintf("error critical panic panic col = %d", col)
	return 0, errors.New(strError)
}

func (g Grille) PoseJetonCol(jeton Jeton, col int) (Grille, error){

	row, err := g.NextVide(col)
	if(err != nil){
		return g, err
	}

	return g.InsertJetonAt(jeton, col,row ), nil
}

func (g Grille) Print(){
	for j:= g.Row - 1; j >= 0; j--{
		fmt.Print("|")
		for i:= 0 ; i < g.Col; i++{
			fmt.Print(g.JetonAt(i,j).getStr())
		}
		fmt.Printf("|%d\n",j)
	}
	fmt.Println("+" + strings.Repeat("-",g.Col) + "+")

	if g.Col < 10 {
		fmt.Print(" ")
		for k := 0; k < g.Col; k++{
			fmt.Printf("%d",k)
		}
		fmt.Println(" ")
	}

	//fmt.Printf("g.Col = %v : g.Row = %v\n", g.Col, g.Row)
}

func NouvelleGrille(col int, row int) Grille{
	//Attention, on utilise le fait que 0 soit la valeur
	//par default des int et que Jeton s'appuie sur int
	//de plus la valeur vide de Jeton correspond a 0
	grd := make([][]Jeton, col)
	for i := 0; i < col ; i++{
		grd[i] = make([]Jeton, row)
	}
	return Grille{
		Col: col,
		Row: row,
		Grid: grd,
	}
}


/***************************************\
|*            REGJETON                 *|
\***************************************/

//On a une liste de jetons et l'index de depart
//et on renvois l'index pour la prochaine recherche
//et le boolean determine si la recherche fut fructueuse
type regjeton func ([]Jeton, int) (int, bool)

//Recherche un type particulier de jeton
func reg_jeton (j Jeton) regjeton {
	return func (l []Jeton, index int) (int, bool){
		if index >= len(l) {
			return index, false
		}
		if j.Equals(l[index]) {
			return index+1, true
		}
		return index, false
	}
}

func (reg1 regjeton) ET(reg2 regjeton) regjeton{
	return reg_ET(reg1, reg2)
}

func (reg1 regjeton) OU(reg2 regjeton) regjeton{
	return reg_OU(reg1, reg2)
}

func (r regjeton) FOIS(n int) regjeton{
	return reg_FOIS(r, n)
}

//regjeton1 ET regjeton2
func reg_ET (reg1 regjeton, reg2 regjeton) regjeton{
	return func (l []Jeton, index int) (int, bool){
		if i1, b1 := reg1(l, index); b1 {
			if i2, b2 := reg2(l, i1); b2 {
				return i2, b2
			}
		}
		return index, false
	}
}

//regjeton1 OU regjeton2
func reg_OU (reg1 regjeton, reg2 regjeton) regjeton {
	return func (l []Jeton, index int) (int, bool){
		if i1, b1 := reg1(l, index); b1 {
			return i1, b1
		}
		if i2, b2 := reg2(l, index); b2 {
			return i2, b2
		}
		return index, false
	}
}

//n fois regjeton
func reg_FOIS (reg regjeton, n int) regjeton {
	return func (l []Jeton, index int) (int, bool) {
		var trouve bool
		m := n
		i := index
		for ;m != 0; m-- {
			if i, trouve = reg(l, i); !trouve {
				return index, false
			}
		}
		return i, true
	}
}

/***************************************\
|*               RULES                 *|
\***************************************/
const VALEUR_GAGNANT = 100000

func evalueGagnant(joueur Jeton, ligne []Jeton) int {
	//On gagne si l'on a 4 fois la meme couleur
	numGagnant := 4
	regle_gagnant := reg_jeton(joueur).FOIS(numGagnant)

	var gagnant bool
	for i, _ := range ligne {
		_ , gagnant = regle_gagnant(ligne, i)
		if gagnant {
			return VALEUR_GAGNANT
		}
	}
	return 0
}

func evalue3Jeton(joueur Jeton, ligne []Jeton) int {

	reward := 10 //Recompense pour avoir 3 jetons a la suite
	var recompense = 0

	//Test 1 vide
	t1v := reg_jeton(Vide)

	//Test 1 jeton
	t1j := reg_jeton(joueur)

	//Test 2 jetons
	t2j := reg_jeton(joueur).FOIS(2)

	//Test 3 jetons
	t3j := reg_jeton(joueur).FOIS(3)


	test := (       //Test 3 jetons et 1 vide OU 1 vide et 3 jetons
			t3j .ET (t1v)) .OU (t1v. ET (t3j)).
		OU ( // OU
			//Test 2 jetons et 1 vide et 1 jeton
			//   OU 1 jeton et 1 vide et 2 jetons
			(t2j .ET (t1v) .ET (t1j)) .OU ( t1j .ET (t1v) .ET (t2j) ))


	for i, _ := range ligne {
		if _, trouve := test(ligne, i); trouve {
			recompense += reward
		}
	}

	return recompense
}

func evalueLigne(joueur Jeton, ligne []Jeton) int {
	//On evalue si la ligne est gagnante
	val := evalueGagnant(joueur, ligne)
	if val >= VALEUR_GAGNANT {
		return val //La valeur est gagnante, on a pas besoin d'aller plus loin
	}

	val = evalue3Jeton(joueur, ligne)
	return val
}

func evalueGrille(joueur Jeton, grille Grille) int {
	evaluation := 0
	for _, ligne := range grille.lines() {
		evaluation += evalueLigne(joueur, ligne)
	}
	return evaluation
}

//Renvoie le gain estime et la colone que joueur doit jouer
//ATTENTION: il faut renvoyer une erreur au cas ou la grille est remplie
func advice(joueur Jeton, grille Grille) (int, int) {
	r   := -500000
	pos := 0
	var err error
	var g Grille
	for i := 0 ; i < grille.Col ; i++ {
		fmt.Println("ADVICE col:", i)
		if !grille.ColIsFull(i) {
			g , err = grille.PoseJetonCol(joueur, i)
			if err != nil {
				fmt.Println("Dans advice:", err)
			}
			gains := evalueGrille(joueur, g) - evalueGrille(joueur.adversaire(), g)
			fmt.Println("GAINS:", gains)
			if gains > r {
				r   = gains
				pos = i
			}
		}
	}
	return r, pos
}

func grilleEstGagnant(joueur Jeton, grille Grille) bool {
	for _, line := range grille.lines() {
		if evalueGagnant(joueur, line) >= VALEUR_GAGNANT {
			return true
		}
	}
	return false
}

//Renvoie le gain estime et la colone que joueur doit jouer
//Et un erreur si la grille est pleine
func deep_advice(joueur Jeton, grille Grille, niveau int) (int, int, error) {
	var pos_optimale int
	max_gain  := -500000
	colsfull  := 0       //On presume qu'aucune colone est remplie
	numErrors := 0

	//On commence par regarder si le joueur est gagnant
	for pos := 0; pos < grille.Col; pos++ {
		g, err := grille.PoseJetonCol(joueur, pos)
		if colNotFull := err == nil; colNotFull {
			if grilleEstGagnant(joueur, g) {
				return VALEUR_GAGNANT, pos, nil
			}
		}else{
			colsfull++
		}
	}

	if colsfull == grille.Col {
		return 0,0, errors.New("grille pleine")
	}

	//On a pas besoin de verifier si l'adversaire est gagnant

	//testons donc chaque coup possible
	for pos := 0; pos < grille.Col; pos++ {
		g, err := grille.PoseJetonCol(joueur, pos)
		if colNotFull := err == nil; colNotFull {
			//Si le niveau est zero, on evalue le cout actuel du joueur
			if niveau == 0 {
				gainJoueur  := evalueGrille(joueur, g)
				gainAdverse := evalueGrille(joueur.adversaire(), g)
				gain        := gainJoueur - gainAdverse

				if gain > max_gain {
					max_gain = gain
					pos_optimale = pos
				}
			}else{
				//On evalue au moins un coup plus loin
				gain_adverse, _, err := deep_advice(joueur.adversaire(), g, niveau - 1)
				gain := -gain_adverse
				if err == nil {
					if gain > max_gain {
						//Si on est gagnant pas besoin de tester les autres possibilite
						if gain >= VALEUR_GAGNANT {
							return VALEUR_GAGNANT,  pos, nil
						}
						max_gain = gain
						pos_optimale = pos
					}
				}else{
					fmt.Println("Panic ERROR", err)
				}
			}
		}else{
			//Colone full
			numErrors++
		}
	}

	if numErrors == grille.Col {
		// Toutes les positions ont renvoyees une erreur!
		return 0,0, errors.New("grille pleine")
	}
	return max_gain, pos_optimale, nil
}

func main() {
//	v := Vide
	j := Jaune
	r := Rouge
	//jetons := []Jeton{v, j, r}
	//fmt.Println(v)
//	for _,v = range jetons{
//		fmt.Printf("%c\n",v.getRune())
//		fmt.Printf("%v\n", v.getStr())
//	}

	grille1 := NouvelleGrille(7,6)
	grille2 := grille1.Copy()
//	grille2.RemplaceJetonAt(j,2,3)
//	grille3 := grille2.InsertJetonAt(r,1,4)
	grille3 := grille2.RemplaceJetonAt(j,1,1).InsertJetonAt(r,0,0)
//	fmt.Println(grille1)
//	fmt.Println(grille2)
//	fmt.Println(grille3)
//	fmt.Println()
//	grille1.Print()
//	fmt.Println()
//	grille2.Print()
//	fmt.Println()
//	grille3.Print()

	grille4, e := grille3.PoseJetonCol(r, 2)
	if e != nil {
		panic(e)
	}
//	grille4.Print()
//	fmt.Println("Col 2 is full:",grille4.ColIsFull(2))
	grille5, e5 := grille4.PoseJetonCol(r,2)
	if e5 != nil {
		panic(e5)
	}
	grille5.RemplaceJetonAt(r,0,1).RemplaceJetonAt(j,1,2).RemplaceJetonAt(r,2,2).RemplaceJetonAt(r,0,2).RemplaceJetonAt(j,1,3).RemplaceJetonAt(r,2,3).RemplaceJetonAt(j,2,4).RemplaceJetonAt(r,3,0)
	grille5.Print()

	reader := bufio.NewReader(os.Stdin)
	currentGrid := grille5

	fmt.Println(grille5.lines())

	//test_rule()

	//fmt.Println("Evaluation de la grille5:")
	grille5.Print()
	//fmt.Printf("Eval pour le joueur Jaune: %v\n", evalueGrille(Jaune, grille5))
	//fmt.Printf("Eval pour le joueur Rouge: %v\n", evalueGrille(Rouge, grille5))

	bcl_run()

	return

	for {
		currentGrid.Print()
		fmt.Print("(col)-> ")
		text, errr := reader.ReadString('\n')
		
		if errr != nil{
			panic(errr)
		}

		txt := strings.Replace(text, "\n", "", -1)

		//if strings.Compare("bye", txt) == 0{
		if txt == "bye" {
			return
		}

		entier, errrr := strconv.Atoi(txt)
		if errrr != nil {
			continue
		}

		if entier < 0 || entier >= currentGrid.Col {
			fmt.Printf("%d out of bounds!\n", entier)
		}
		fmt.Printf("%d++ = %d\n", entier, entier+1)
		currentGrid, e  = currentGrid.PoseJetonCol(r, entier)
		if e != nil {
			fmt.Println(e)
			continue
		}
	}
}

func bcl_run(){
	grille := NouvelleGrille(7,6)
	joueur := select_jeton()
	difficulty := select_difficulty()
	var err error
 
	fmt.Println("Difficulty:", difficulty)

	for{
		grille.Print()
		if grilleEstGagnant(joueur.adversaire(), grille) {
			fmt.Println("GROS NAZE! J'ai gagne!!!")
			return
		}
		fmt.Printf("Valeur pour %v: %v\n", Jaune.getStr(), evalueGrille(Jaune, grille))
		fmt.Printf("Valeur pour %v: %v\n", Rouge.getStr(), evalueGrille(Rouge, grille))
		//fmt.Println("Choix:",joueur.getStr())
		position := select_position(grille)
		grille, _ = grille.PoseJetonCol(joueur, position)
		if grilleEstGagnant(joueur, grille) {
			grille.Print()
			fmt.Println("BRAVO! Vous avez gagne")
			return
		}
		//_ , posAdversaire :=  advice(joueur.adversaire(), grille)
		gain , posAdversaire, errT :=  deep_advice(joueur.adversaire(), grille, difficulty)
		if errT != nil {
			fmt.Println("ERREUR TROUBLANTE:", errT)
		}
		fmt.Printf("ESTIMATION GAIN: %v pour position %v\n", gain, posAdversaire)
		grille, err = grille.PoseJetonCol(joueur.adversaire(), posAdversaire)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}
}

func select_difficulty() int {
	var val int
	var err error
	for {
		fmt.Print("Select difficulty (0 a beaucoup)? ")
		reader  := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text,"\n", "", -1)
		val, err = strconv.Atoi(text)
		if err == nil && val >= 0 {
			break
		}
	}
	return val
}

func select_jeton() Jeton {
	var i , j Jeton
	i , j = 1 , 2
	var val int
	var err error

	for {
		fmt.Printf("Choix de la couleur: %v -> %v ou %v -> %v ? ", i, i.getStr(), j, j.getStr())
		reader := bufio.NewReader(os.Stdin)
		text, _:= reader.ReadString('\n')
		text = strings.Replace(text,"\n", "", -1)
		val, err = strconv.Atoi(text)
		if err == nil && (val == 1 || val == 2) {
				break
		}
	}
	return Jeton(val)
}

func appartient_int(tab []int, valeur int) bool {
	for i:= 0; i < len(tab); i++ {
		if tab[i] == valeur {
			return true
		}
	}
	return false
}

func select_position(grille Grille) int {
	valides := []int{}
	var val int
	var err error

	for i := 0; i < grille.Col; i++ {
		if !grille.ColIsFull(i) {
			valides = append(valides, i)
			//fmt.Printf("%d",i)
		}
	}

	for {
		fmt.Println("Positions valides:")
		fmt.Printf("|")
		for i := 0; i < len(valides); i++ {
			fmt.Printf("%d|", valides[i])
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\nChoix ? ")
		text, _:= reader.ReadString('\n')
		text = strings.Replace(text,"\n", "", -1)
		val, err = strconv.Atoi(text)
		if err == nil && appartient_int(valides, val) {
			break
		}
//		return val
	}
	return val
}

func test_rule(){

	ligne1 := []Jeton{Jaune, Jaune, Jaune, Vide, Vide, Rouge}

	ligne2 := []Jeton{Jaune, Jaune, Jaune, Jaune, Vide, Rouge}

	fmt.Printf("GAGNANT: %v pour %v\n", evalueGagnant(Jaune, ligne1), ligne1)
	fmt.Printf("GAGNANT: %v pour %v\n", evalueGagnant(Jaune, ligne2), ligne2)

	//TEST rouge ou jaune
	ligne3 := []Jeton{Jaune,Vide,Rouge}
	ligne4 := []Jeton{Rouge,Vide,Jaune}
	ligne5 := []Jeton{Vide,Jaune,Rouge}
	ligne6 := []Jeton{Vide,Jaune,Jaune,Jaune,Rouge}
	ligne7 := []Jeton{Rouge,Jaune,Jaune,Jaune,Vide,Rouge}
	ligne8 := []Jeton{Rouge,Jaune,Jaune,Jaune,Rouge,Rouge}
	ligne9 := []Jeton{Vide,Jaune,Vide,Jaune,Jaune,Rouge,Rouge}
	ligne10 := []Jeton{Rouge,Jaune,Jaune,Vide,Jaune,Rouge,Rouge}
	ligne11 := []Jeton{Rouge,Jaune,Jaune,Jaune,Jaune,Rouge,Rouge}

	test_Jaune := reg_jeton(Jaune)
	test_Rouge := reg_jeton(Rouge)
	test_Jaune_OU_Rouge := reg_OU(test_Jaune, test_Rouge)

	test_test(test_Jaune_OU_Rouge, ligne3,0, "TEST JAUNE OU ROUGE")
	test_test(test_Jaune_OU_Rouge, ligne4,0, "TEST JAUNE OU ROUGE")
	test_test(test_Jaune_OU_Rouge, ligne5,0, "TEST JAUNE OU ROUGE")

	fmt.Printf("EVALUE3JETONS: %v pour %v\n", evalue3Jeton(Jaune, ligne6), ligne6)
	fmt.Printf("EVALUE3JETONS: %v pour %v\n", evalue3Jeton(Jaune, ligne7), ligne7)
	fmt.Printf("EVALUE3JETONS: %v pour %v\n", evalue3Jeton(Jaune, ligne8), ligne8)
	fmt.Printf("EVALUE3JETONS: %v pour %v\n", evalue3Jeton(Jaune, ligne9), ligne9)
	fmt.Printf("EVALUE3JETONS: %v pour %v\n", evalue3Jeton(Jaune, ligne10), ligne10)

	fmt.Println()
	test_evalue_lignes(Jaune,[][]Jeton{ligne1,ligne2,ligne3,ligne4,ligne5,ligne6,ligne7,ligne8,ligne9,ligne10,ligne11})
}

func test_test(r regjeton, ligne []Jeton, index int, label string) {
	i, trouve := r(ligne, index)
	fmt.Printf("%v: index = %v , trouve = %v -> %v\n", label, i, trouve, ligne)
}

func test_evalue_lignes(joueur Jeton, lignes [][]Jeton) {
	for _ , val := range lignes {
		test_evalue_ligne(joueur, val)
	}
}

func test_evalue_ligne(joueur Jeton, ligne []Jeton) {
	fmt.Printf("Evaluation: %v pour ligne -> %v\n", evalueLigne(joueur, ligne), ligne)
}

