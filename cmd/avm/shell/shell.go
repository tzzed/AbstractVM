package shell

import (
	"avm/evaluator"
	"avm/parser"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
)

const PROMPT = "avm>"

const (
	historyFilename = ".avm_history"
)

type Shell struct {
	prompt  string
	history []string
	st *evaluator.Stack
}

func (sh *Shell) dumpHistory() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fname := filepath.Join(homeDir, historyFilename)

	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, h := range sh.history {
		_, err = w.WriteString(h + "\n")
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

func (sh *Shell) loadHistory() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fname := filepath.Join(homeDir, historyFilename)

	_, err = os.Stat(fname)
	if err != nil {
		return nil, nil
	}

	f, err := os.Open(fname)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	var history []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		history = append(history, s.Text())
	}

	return history, s.Err()
}

func (sh *Shell) runInstruction(in string) error {
	p := parser.NewParser(in)
	pg, err := p.ParseInstruction()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if pg == nil {
		return errors.New("empty program")
	}

	_, err = sh.st.Eval(pg)
	if err != nil {
		fmt.Println(err.Error() + in)
	} else {
		_, _ = sh.st.Dump()
	}

	return nil
}

func (sh *Shell) executeInput(in string) error {
	in = strings.TrimSpace(in)
	if in == "help" {
		return displayHelpCommand()
	}
	err := sh.runInstruction(in)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (sh *Shell) createStack() *evaluator.Stack {
	return evaluator.NewStack()
}

func (sh *Shell) execute(in string, ) {
	sh.history = append(sh.history, in)

	err := sh.executeInput(in)
	if err != nil {
		fmt.Println(err)
	}
}

func (sh *Shell) completer(in prompt.Document) []prompt.Suggest {
	_, err := parser.NewParser(in.Text).ParseInstruction()
	if err != nil {
		return []prompt.Suggest{}
	}
	expected := getCommands()
	var suggestions []prompt.Suggest
	for _, e := range expected {
		suggestions = append(suggestions, prompt.Suggest{
			Text: e.name,
		})
	}

	w := in.GetWordBeforeCursor()
	if w == "" {
		return suggestions
	}

	return prompt.FilterHasPrefix(suggestions, w, true)
}

// Run start shell
func Run(in io.Reader, out io.Writer) error {

	var sh Shell

	history, err := sh.loadHistory()
	if err != nil {
		return err
	}
	fmt.Println("Abstract VM")
	fmt.Println("Enter \".help\" for usage hints.")
	registerCommands()
	sh.st = evaluator.NewStack()
	e := prompt.New(sh.execute,
		sh.completer,
		prompt.OptionTitle("AVM"),
		prompt.OptionPrefix(PROMPT),
		prompt.OptionHistory(history),
	)

	e.Run()

	return nil

}
