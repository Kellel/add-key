package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/xerrors"
)

type CLI struct {
	Name       string   `arg:"" help:"name of the new repo to add"`
	Type       []string `required:"" default:"deb" help:"types to include (deb, deb-src)"`
	GPG        string   `help:"gpg key url" required:""`
	URI        string   `help:"debian repo url" required:""`
	Suite      string   `help:"what suite to install" required:""`
	Components []string `required:"" help:"component part"`
}

type TemplateContext struct {
	Types      []string
	URI        string
	Suite      string
	Components []string
	SignedBy   string
}

var deb822Template = `
Types: {{ range $v := .Types }}{{ $v }}{{ end }}
URIs: {{ .URI }}
Suites: {{ .Suite }}
Components: {{ range $v := .Components }}{{ $v }}{{ end }}
Signed-By: {{ .SignedBy }}
`

func AddKey(cli *CLI) error {
	resp, err := http.Get(cli.GPG)
	if err != nil {
		return xerrors.Errorf("failed to fetch GPG payload: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return xerrors.Errorf("Error fetching GPG public key: [%v] %v", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()

	entities, err := openpgp.ReadArmoredKeyRing(resp.Body)
	if err != nil {
		return xerrors.Errorf("unable to read armored keyring: %w", err)
	}

	if len(entities) != 1 {
		return xerrors.Errorf("found %v entites in gpg body, unsure what to do", len(entities))
	}

	keyPath := path.Join("/usr/share/keyrings/", cli.Name+".gpg")
	file, err := os.Create(keyPath)
	if err != nil {
		return xerrors.Errorf("unable to create keyfile: %w", err)
	}
	defer file.Close()

	entities[0].Serialize(file)
	fmt.Printf("Wrote %v\n", keyPath)

	tmpl, err := template.New("deb822").Parse(deb822Template)
	if err != nil {
		return xerrors.Errorf("unable to parse debian source template: %w", err)
	}

	sourcePath := path.Join("/etc/apt/sources.list.d/", cli.Name+".sources")
	sourceFile, err := os.Create(sourcePath)
	if err != nil {
		return xerrors.Errorf("unable to create new source file: %w", err)
	}
	defer sourceFile.Close()

	tmplContext := &TemplateContext{
		Types:      cli.Type,
		URI:        cli.URI,
		Suite:      cli.Suite,
		Components: cli.Components,
		SignedBy:   keyPath,
	}

	err = tmpl.Execute(sourceFile, tmplContext)
	if err != nil {
		return xerrors.Errorf("failed to write %v: %w", sourcePath, err)
	}

	fmt.Printf("Wrote %v\n", sourcePath)
	return nil
}

func main() {
	cli := &CLI{}
	ctx := kong.Parse(cli)
	err := AddKey(cli)
	ctx.FatalIfErrorf(err)
}
