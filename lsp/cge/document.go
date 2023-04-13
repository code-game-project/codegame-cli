package cge

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/code-game-project/cge-parser/adapter"
	"github.com/code-game-project/cli-utils/components"
	"github.com/code-game-project/cli-utils/versions"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type document struct {
	uri         protocol.DocumentUri
	metadata    adapter.Metadata
	content     string
	changed     bool
	diagnostics []protocol.Diagnostic
	objects     []adapter.Object
}

func cgeDiagnosticToLSP(diagnostic adapter.Diagnostic) protocol.Diagnostic {
	var severity protocol.DiagnosticSeverity
	switch diagnostic.Type {
	case adapter.DiagInfo:
		severity = protocol.DiagnosticSeverityInformation
	case adapter.DiagWarning:
		severity = protocol.DiagnosticSeverityWarning
	case adapter.DiagError:
		severity = protocol.DiagnosticSeverityError
	}
	source := "cge-ls"
	return protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(diagnostic.StartLine),
				Character: uint32(diagnostic.StartColumn),
			},
			End: protocol.Position{
				Line:      uint32(diagnostic.EndLine),
				Character: uint32(diagnostic.EndColumn),
			},
		},
		Severity: &severity,
		Source:   &source,
		Message:  diagnostic.Message,
	}
}

func (d *document) validate(notify glsp.NotifyFunc) {
	if !d.changed {
		return
	}
	d.changed = false

	defer d.sendDiagnostics(notify)

	d.diagnostics = d.diagnostics[:0]

	file := io.Reader(bytes.NewBufferString(d.content))
	var err error
	var diagnostics []adapter.Diagnostic
	d.metadata, file, diagnostics, err = adapter.ParseMetadata(file)
	if err != nil {
		for _, diag := range diagnostics {
			d.diagnostics = append(d.diagnostics, cgeDiagnosticToLSP(diag))
		}
		return
	}

	cgeParser, err := components.CGEParser(versions.MustParse(d.metadata.CGEVersion))
	if err != nil {
		d.diagnostics = append(d.diagnostics, cgeDiagnosticToLSP(adapter.Diagnostic{
			Type:    adapter.DiagWarning,
			Message: "unsupported CGE version; LSP disabled",
		}))
		fmt.Fprintf(os.Stderr, "ERROR: failed to load cge-parser: %s", err)
		return
	}

	result, errs := adapter.ParseCGE(file, cgeParser, adapter.Config{
		IncludeComments: true,
	})
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "ERROR: failed to parse '%s':\n", d.uri)
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "- %s\n", err)
		}
		return
	}
	hadError := false
	for _, diag := range result.Diagnostics {
		if diag.Type == adapter.DiagError {
			hadError = true
		}
		d.diagnostics = append(d.diagnostics, cgeDiagnosticToLSP(diag))
	}
	if !hadError || len(d.objects) == 0 {
		d.objects = make([]adapter.Object, 0, 1+len(result.Events)+len(result.Commands)+len(result.Types)+len(result.Enums))
		d.objects = append(d.objects, result.Config)
		d.objects = append(d.objects, result.Events...)
		d.objects = append(d.objects, result.Commands...)
		d.objects = append(d.objects, result.Types...)
		d.objects = append(d.objects, result.Enums...)
	}
}

func (d *document) sendDiagnostics(notify glsp.NotifyFunc) {
	notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         d.uri,
		Diagnostics: d.diagnostics,
	})
}

func (l *lsp) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	document := &document{
		uri:         params.TextDocument.URI,
		content:     params.TextDocument.Text,
		changed:     true,
		diagnostics: make([]protocol.Diagnostic, 0),
		objects:     make([]adapter.Object, 0),
	}
	l.documents.Store(params.TextDocument.URI, document)
	go document.validate(context.Notify)
	return nil
}

func (l *lsp) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if document, ok := l.getDocument(params.TextDocument.URI); ok {
		content := document.content
		for _, change := range params.ContentChanges {
			if c, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
				start, end := c.Range.IndexesIn(content)
				content = content[:start] + c.Text + content[end:]
			} else if c, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
				content = c.Text
			}
		}
		document.content = content
		document.changed = len(params.ContentChanges) > 0
		go document.validate(context.Notify)
	}
	return nil
}

func (l *lsp) textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	_, ok := l.documents.LoadAndDelete(params.TextDocument.URI)
	if ok {
		go context.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: make([]protocol.Diagnostic, 0),
		})
	}
	return nil
}

func (l *lsp) getDocument(uri protocol.DocumentUri) (*document, bool) {
	doc, ok := l.documents.Load(uri)
	if !ok {
		return nil, false
	}
	return doc.(*document), true
}
