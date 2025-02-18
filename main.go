package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const LIMITE_MINIMO_DE_PALAVRAS = 10
const LIMITE_MAXIMO_DE_PALAVRAS = 100

const TAMANHO_ALFABETO = 26
const LIMITE_TAMANHO_PALAVRA = 15

const ARQUIVO_PALAVRAS = "palavras.bin"

// iota -> inicia no zero e soma mais um p/ cada const
// iota + 1 faz ele iniciar no 1
const (
	NOVO_JOGO = iota + 1
	VER_PALAVRAS
	ADICIONAR_PALAVRAS
	ATUALIZAR_PALAVRAS
	DELETAR_PALAVRAS
	SAIR
)

const (
	CONTINUAR_JOGO = iota + 1
	MENU_INICIAL
	SAIR_JOGO
)

// Palavra ---
type Palavra struct {
	Conteudo string
	Tamanho  int16
}

func (palavra Palavra) String() string {
	return palavra.Conteudo + " tamanho: " + strconv.FormatInt(int64(palavra.Tamanho), 10)
}

func calcularTamanhoString(string string) int16 {
	return int16(utf8.RuneCountInString(string))
}

func NewPalavra(string string) (Palavra, error) {

	// Atribui a palavra o somente o limite de caractéres
	// para isso que o :LIMITE_TAMANHO_PALAVRA serve

	palavra := new(Palavra)
	if validarDuplicidade(string) {
		return *palavra, fmt.Errorf("Essa palavra já foi cadastrada.")
	}

	tamanhoPalavra := calcularTamanhoString(string)

	if tamanhoPalavra < 5 {
		return *palavra, fmt.Errorf("A palavra deve ter no mínimo 5 letras")
	}

	if tamanhoPalavra > LIMITE_TAMANHO_PALAVRA {
		palavra.Conteudo = string[:LIMITE_TAMANHO_PALAVRA]
	} else {
		palavra.Conteudo = string
	}

	palavra.Tamanho = tamanhoPalavra
	return *palavra, nil
}

// Palavra ---

// BufferPalavra ---
type BufferPalavra struct {
	BufferArquivo   []Palavra
	QuantidadeAtual int16
}

func (buffer *BufferPalavra) AdicionarPalavraNoFinal(palavra Palavra) error {
	if buffer.QuantidadeAtual >= 100 {
		return fmt.Errorf("Não é possível adicionar mais do que %d palavras", LIMITE_MAXIMO_DE_PALAVRAS)
	}

	buffer.BufferArquivo = append(buffer.BufferArquivo, palavra)
	buffer.QuantidadeAtual++

	salvarBufferPalavrasArquivo()
	return nil
}

func (buffer *BufferPalavra) AtualizarPalavra(palavra Palavra, index int) {

	buffer.BufferArquivo[index] = palavra
	salvarBufferPalavrasArquivo()
}

func (buffer *BufferPalavra) DeletarPalavra(index int) {

	buffer.BufferArquivo = append(buffer.BufferArquivo[:index], buffer.BufferArquivo[index+1:]...)
	buffer.QuantidadeAtual--
	salvarBufferPalavrasArquivo()
}

// BufferPalavra ---

var buffer BufferPalavra = *new(BufferPalavra)

// Arquivo ---
func lerPalavrasArquivo() {

	arquivo, erro := os.Open(ARQUIVO_PALAVRAS)

	defer arquivo.Close()
	if erro != nil {
		arquivo, erro = os.Create(ARQUIVO_PALAVRAS)

		if erro != nil {
			log.Fatal("Falha ao gerar arquivo:" + erro.Error())
		}

		lerPalavrasArquivo()
		return
	}

	decoder := gob.NewDecoder(arquivo)
	erro = decoder.Decode(&buffer)

	if erro != nil {
		fmt.Printf("Falha ao ler palavras\n")
		fmt.Printf("Ou o arquivo estava vazio  ¯\\_(ツ)_/¯\n\n")
	}
}

func salvarBufferPalavrasArquivo() {

	arquivo, erro := os.Create(ARQUIVO_PALAVRAS)

	if erro != nil {
		log.Fatal("Falha ao abrir arquivo: " + erro.Error())
	}

	defer arquivo.Close()

	encoder := gob.NewEncoder(arquivo)
	erro = encoder.Encode(buffer)

	if erro != nil {
		log.Fatal("Falha de escrita: " + erro.Error())
	}
}

// Arquivo ---

func validarDuplicidade(stringUsuario string) bool {
	if buffer.QuantidadeAtual == 0 {
		return false
	}

	for i := 0; i < int(buffer.QuantidadeAtual); i++ {

		stringArquivo := buffer.BufferArquivo[i].Conteudo

		if strings.EqualFold(stringUsuario, stringArquivo) {
			return true
		}
	}

	return false
}

func lerInputUser() string {

	string := ""

	fmt.Print("Escreva uma palavra: ")
	fmt.Scan(&string)

	return string
}

func adicionarPalavra() {

	stringUsuario := lerInputUser()

	palavra, error := NewPalavra(stringUsuario)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	buffer.AdicionarPalavraNoFinal(palavra)
}

func lerInputUserIndex() int8 {
	index := 0

	fmt.Print("Selecione o código de uma palavra: ")
	fmt.Scan(&index)

	return int8(index)
}

func atualizarPalavra() {

	if buffer.QuantidadeAtual == 0 {
		fmt.Println("Não existem palavras cadastradas para atualizar, tente cadastrar uma palavra")
		return
	}

	exibirPalavras()
	index := lerInputUserIndex()

	if index <= 0 || int16(index) > buffer.QuantidadeAtual {
		fmt.Println("O código inserido está inválido, tente novamente")

		atualizarPalavra()
		return
	}

	palavra, erro := NewPalavra(lerInputUser())

	if erro != nil {
		fmt.Println(erro.Error())
		return
	}

	buffer.AtualizarPalavra(palavra, int(index)-1)
}

func deletarPalavras() {

	if buffer.QuantidadeAtual == 0 {
		fmt.Println("Não existem palavras cadastradas para deletar, tente cadastrar uma palavra")
		return
	}

	exibirPalavras()
	index := lerInputUserIndex()

	if index <= 0 || int16(index) > buffer.QuantidadeAtual {
		fmt.Println("O código inserido está inválido, tente novamente")

		deletarPalavras()
		return
	}

	buffer.DeletarPalavra(int(index - 1))
}

func exibirPalavras() {

	if buffer.QuantidadeAtual == 0 {
		fmt.Println("Em questão de palavras, não temos palavras")
		return
	}

	for i, palavra := range buffer.BufferArquivo {

		if palavra.Tamanho == 0 {
			return
		}

		fmt.Printf("%d. %s\n", i+1, palavra.Conteudo)
	}
}

func escolhaMenuPrincipal() int {

	escolhaMenu := 0

	fmt.Println("==Menu Inicial==")
	fmt.Println("1. Novo Jogo")
	fmt.Println("2. Ver Palavras")
	fmt.Println("3. Adicionar Palavras")
	fmt.Println("4. Atualizar Palavras")
	fmt.Println("5. Deletar Palavras")
	fmt.Println("6. Sair")
	fmt.Print("Escolha uma das opções: ")
	fmt.Scan(&escolhaMenu)

	erro := validaEscolhaMenu(escolhaMenu)

	if erro != nil {
		fmt.Println(erro)
		return escolhaMenuPrincipal()
	}

	return escolhaMenu
}

func validaEscolhaMenu(escolhaMenu int) error {

	if escolhaMenu < NOVO_JOGO || escolhaMenu > SAIR {
		return fmt.Errorf("Insira uma opção válida!")
	}

	return nil
}

func escolhaMenuFimJogo() bool {

	escolhaMenu := 0

	fmt.Println("==Menu de Fim de Jogo==")
	fmt.Println("1. Novo Jogo")
	fmt.Println("2. Menu Inicial")
	fmt.Println("3. Sair")

	fmt.Print("Escolha uma das opções: ")
	fmt.Scan(&escolhaMenu)

	if escolhaMenu == 0 || escolhaMenu > 3 {
		fmt.Println("Escolha inválida tente novamente")
		return escolhaMenuFimJogo()
	}

	switch escolhaMenu {

	case CONTINUAR_JOGO:
    jogo()
    break
	case MENU_INICIAL:
		return true

	case SAIR_JOGO:
		fmt.Println("Saindo... Até mais")
		os.Exit(0)
	}

	return false
}

func definirPalavraAleatoria() Palavra {

	randomNumber := rand.Int() % int(buffer.QuantidadeAtual)
	palavraSecreta := buffer.BufferArquivo[randomNumber]

	return palavraSecreta
}

func inicializarStringPalpite(tamanhoPalavraSecreta int) string {

	stringPalpite := ""
	for i := 0; i > tamanhoPalavraSecreta; i++ {
		stringPalpite += "_"

	}
	return stringPalpite
}

func letraFoiUsada(letra rune, letrasUsadas string) bool {

	if letrasUsadas == "" {
		return false
	}

	for _, letraUsada := range letrasUsadas {

		if letra == letraUsada {
			return true
		}
	}

	return false
}

func validarAcerto(palpiteLetra rune, palavraSecreta Palavra) bool {

	return strings.ContainsRune(palavraSecreta.Conteudo, palpiteLetra) || strings.ContainsRune(palavraSecreta.Conteudo, unicode.ToUpper(palpiteLetra))
}

func validarSeGanhou(palpite string, palavraSecreta Palavra) bool {
	return strings.EqualFold(palpite, palavraSecreta.Conteudo)
}

func jogo() {

	if buffer.QuantidadeAtual < LIMITE_MINIMO_DE_PALAVRAS {
		fmt.Printf("São necessárias no mínimo %d palavras para poder jogar\n", LIMITE_MINIMO_DE_PALAVRAS)
		fmt.Printf("No momento faltam: %d\n", LIMITE_MINIMO_DE_PALAVRAS-buffer.QuantidadeAtual)
		return
	}

	palavraSecreta := definirPalavraAleatoria()
	palpite := inicializarStringPalpite(int(palavraSecreta.Tamanho))
	letrasUsadas := ""

	tentativas := 6
	pontuacao := 0
	acertosConsecutivos := 0
	qtdLetrasUsadas := 0

	fmt.Println("Hora de Adivinhar. Você tem 6 tentativas...")
	fmt.Println(palpite)

	for tentativas > 0 {

		palpiteLetra := '\n'
		fmt.Print("Adivinhe uma letra da palavra: ")
		fmt.Scan(&palpiteLetra)

		if letraFoiUsada(palpiteLetra, letrasUsadas) {
			fmt.Printf("A letra %c já foi utilizada, tente novamente\n", palpiteLetra)
			continue
		}

		//Adiciona uma runa a string
		letrasUsadas += string(palpiteLetra)
		qtdLetrasUsadas++

		fmt.Printf("%d letras foram usadas\n", qtdLetrasUsadas)

		if validarAcerto(palpiteLetra, palavraSecreta) {

			acertosConsecutivos++
			pontuacao = (pontuacao + 10) * acertosConsecutivos
			fmt.Println("Muito bem!")
			if validarSeGanhou(palpite, palavraSecreta) {

        if (tentativas == 6) {

        }

       if escolhaMenuFimJogo() {
         return;
       }
			}
		} else {
			acertosConsecutivos = 0
			tentativas--
			pontuacao -= 5
			fmt.Println("Não foi dessa vez :(")
		}

	}
}

func main() {

	lerPalavrasArquivo()
	for {
		escolha := escolhaMenuPrincipal()

		switch escolha {
		case NOVO_JOGO:
			jogo()

		case VER_PALAVRAS:
			exibirPalavras()

		case ADICIONAR_PALAVRAS:
			adicionarPalavra()

		case ATUALIZAR_PALAVRAS:
			atualizarPalavra()

		case DELETAR_PALAVRAS:
			deletarPalavras()

		case SAIR:
			fmt.Println("Saindo... Até mais!")
			return
		}
	}
}
