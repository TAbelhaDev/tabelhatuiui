package tuiui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// IPCArgs is the parsed form of a `bin ipc <método> [key=value...] --json`
// invocation — the scriptable-data-source convention dcal/djobs/tabelaradar
// all share, so a shell script or an LLM can ask a TUI's data without going
// through the interface itself.
type IPCArgs struct {
	Method  string
	Filters map[string]string
	JSON    bool
}

// ParseIPCArgs parses the args after `ipc`. It returns an error for an
// unknown arg (anything that's neither --json nor key=value) or a missing
// --json flag; the caller is expected to print its own usage line with it.
func ParseIPCArgs(args []string) (*IPCArgs, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("método ausente")
	}

	out := &IPCArgs{
		Method:  args[0],
		Filters: map[string]string{},
	}
	for _, arg := range args[1:] {
		if arg == "--json" {
			out.JSON = true
			continue
		}
		if k, v, ok := strings.Cut(arg, "="); ok {
			out.Filters[k] = v
			continue
		}
		return nil, fmt.Errorf("argumento inválido: %q (esperado key=value ou --json)", arg)
	}
	if !out.JSON {
		return nil, fmt.Errorf("apenas saída --json é suportada por enquanto")
	}
	return out, nil
}

// WriteJSON pretty-prints v as JSON to stdout, returning the process exit
// code (0 on success, 1 on a serialization error).
func WriteJSON(v any) int {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, "erro ao serializar json:", err)
		return 1
	}
	return 0
}
