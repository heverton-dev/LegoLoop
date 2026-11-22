package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type wixDocumento struct {
	XMLName xml.Name   `xml:"Wix"`
	Package wixPackage `xml:"Package"`
}

type wixPackage struct {
	Name              string           `xml:"Name,attr"`
	Manufacturer      string           `xml:"Manufacturer,attr"`
	Scope             string           `xml:"Scope,attr"`
	InstallPrivileges string           `xml:"InstallPrivileges,attr"`
	StandardDirs      []wixStandardDir `xml:"StandardDirectory"`
	Features          []wixFeature     `xml:"Feature"`
}

type wixStandardDir struct {
	ID   string   `xml:"Id,attr"`
	Dirs []wixDir `xml:"Directory"`
}

type wixDir struct {
	ID         string         `xml:"Id,attr"`
	Name       string         `xml:"Name,attr"`
	Components []wixComponent `xml:"Component"`
}

type wixComponent struct {
	ID    string    `xml:"Id,attr"`
	Files []wixFile `xml:"File"`
}

type wixFile struct {
	ID     string `xml:"Id,attr"`
	Name   string `xml:"Name,attr"`
	Source string `xml:"Source,attr"`
}

type wixFeature struct {
	ID string `xml:"Id,attr"`
}

func TestMSIPorUsuario(t *testing.T) {
	caminho := filepath.Join("wix", "legoloop.wxs")
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler %s: %v", caminho, err)
	}
	bruto := strings.ToLower(string(dados))
	if strings.Contains(bruto, "permachine") || strings.Contains(bruto, "elevated") {
		t.Fatalf("definicao WiX nao pode declarar perMachine ou elevacao")
	}

	var doc wixDocumento
	if err := xml.Unmarshal(dados, &doc); err != nil {
		t.Fatalf("XML WiX invalido: %v", err)
	}
	if doc.Package.Scope != "perUser" {
		t.Errorf("Package.Scope = %q, esperado %q", doc.Package.Scope, "perUser")
	}
	if doc.Package.InstallPrivileges != "" && doc.Package.InstallPrivileges != "limited" {
		t.Errorf("InstallPrivileges = %q, esperado vazio ou limited", doc.Package.InstallPrivileges)
	}
	if !contemBinarioNoPerfil(doc.Package.StandardDirs) {
		t.Errorf("definicao WiX deve instalar legoloop.exe em LocalAppDataFolder")
	}
	if len(doc.Package.Features) == 0 {
		t.Errorf("definicao WiX deve declarar ao menos uma Feature")
	}
}

func contemBinarioNoPerfil(dirs []wixStandardDir) bool {
	for _, sd := range dirs {
		if sd.ID != "LocalAppDataFolder" {
			continue
		}
		for _, d := range sd.Dirs {
			for _, c := range d.Components {
				for _, f := range c.Files {
					if f.Name == "legoloop.exe" {
						return true
					}
				}
			}
		}
	}
	return false
}
