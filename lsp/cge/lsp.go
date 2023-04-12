package cge

import (
	"fmt"
	"sync"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"

	"github.com/code-game-project/codegame-cli/version"
)

type lsp struct {
	handler   protocol.Handler
	documents sync.Map
}

func RunLSP() error {
	logging.Configure(0, nil)

	lsp := &lsp{}
	lsp.handler = protocol.Handler{
		Initialize:             lsp.initialize,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentDidOpen:    lsp.textDocumentDidOpen,
		TextDocumentDidChange:  lsp.textDocumentDidChange,
		TextDocumentDidClose:   lsp.textDocumentDidClose,
		TextDocumentCompletion: lsp.textDocumentCompletion,
	}

	srv := server.NewServer(&lsp.handler, "cge-ls", false)

	err := srv.RunStdio()
	if err != nil {
		return fmt.Errorf("failed to run CGE language server: %w", err)
	}
	return nil
}

func (l *lsp) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := l.handler.CreateServerCapabilities()
	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindIncremental
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{},
	}
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "cge-ls",
			Version: &version.Version,
		},
	}, nil
}

func shutdown(_ *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(_ *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
