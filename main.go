package main

import (
	"fmt"
	"os"
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

type Palavra struct {
	string  string
	tamanho int
}

type BufferPalavra struct {
	bufferArquivo   [LIMITE_MAXIMO_DE_PALAVRAS]Palavra
	quantidadeAtual int
}

func (palavra *Palavra) New(string string) Palavra {

	// Atribui a palavra o somente o limite de caractéres
	// para isso que o :LIMITE_TAMANHO_PALAVRA serve

	if len(string) > LIMITE_TAMANHO_PALAVRA {
		palavra.string = string[:LIMITE_TAMANHO_PALAVRA]
	} else {
		palavra.string = string
	}

	palavra.tamanho = len(palavra.string)
	return *palavra
}

func (palavra Palavra) String() string {
	return palavra.string
}

var buffer BufferPalavra = *new(BufferPalavra)

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

	if escolhaMenu < 0 || escolhaMenu > 6 {
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
	case MENU_INICIAL:
		return true

	default:
		fmt.Println("Saindo... Até mais")
		os.Exit(0)
	}

	return false
}

func main() {

	palavra := new(Palavra)
	palavra.New("Jailson")
	fmt.Println(palavra)

	for {
		escolha := escolhaMenuPrincipal()

		switch escolha {

		case NOVO_JOGO:
		case VER_PALAVRAS:
		case ADICIONAR_PALAVRAS:
		case ATUALIZAR_PALAVRAS:
		case DELETAR_PALAVRAS:
			break

		case SAIR:
			fmt.Println("Saindo... Até mais!")
			return
		}
	}
}
