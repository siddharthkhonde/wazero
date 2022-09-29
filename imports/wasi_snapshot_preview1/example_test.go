package wasi_snapshot_preview1

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

// exitOnStartWasm was generated by the following:
//
//	cd testdata; wat2wasm --debug-names exit_on_start.wat
//
//go:embed testdata/exit_on_start.wasm
var exitOnStartWasm []byte

// This is an example of how to use WebAssembly System Interface (WASI) with its simplest function: "proc_exit".
//
// See https://github.com/tetratelabs/wazero/tree/main/examples/wasi for another example.
func Example() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)

	// Instantiate WASI, which implements system I/O such as console output.
	wm, err := Instantiate(ctx, r)
	if err != nil {
		log.Panicln(err)
	}
	defer wm.Close(testCtx)

	// Compile the WebAssembly module using the default configuration.
	code, err := r.CompileModule(ctx, exitOnStartWasm)
	if err != nil {
		log.Panicln(err)
	}
	defer code.Close(ctx)

	// InstantiateModule runs the "_start" function which is like a "main" function.
	// Override default configuration (which discards stdout).
	mod, err := r.InstantiateModule(ctx, code, wazero.NewModuleConfig().WithStdout(os.Stdout).WithName("wasi-demo"))
	if mod != nil {
		defer r.Close(ctx)
	}

	// Note: Most compilers do not exit the module after running "_start", unless
	// there was an error. This allows you to call exported functions.
	if exitErr, ok := err.(*sys.ExitError); ok {
		fmt.Printf("exit_code: %d\n", exitErr.ExitCode())
	}

	// Output:
	// exit_code: 2
}
