package videomotion

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
)

var (
	// ErrNoClipsForFCPXML indica ausência de clipes para a timeline.
	ErrNoClipsForFCPXML = errors.New("videomotion: nenhum clipe fornecido para a timeline FCPXML")
)

// TimelineClip representa um segmento de mídia na linha do tempo.
type TimelineClip struct {
	Name        string
	FilePath    string
	DurationSec float64
}

// FCPXMLSpec reúne as configurações para exportação da timeline.
type FCPXMLSpec struct {
	ProjectName string
	FrameRate   int
	Width       int
	Height      int
	Clips       []TimelineClip
}

// GenerateFCPXML constrói a árvore XML compatível com DaVinci Resolve Studio.
func GenerateFCPXML(spec FCPXMLSpec) ([]byte, error) {
	if len(spec.Clips) == 0 {
		return nil, ErrNoClipsForFCPXML
	}
	if spec.FrameRate <= 0 {
		spec.FrameRate = 25
	}
	var buf bytes.Buffer
	buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buf.WriteString("<!DOCTYPE fcpxml>\n")
	fcpxml := buildFCPXMLData(spec)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(fcpxml); err != nil {
		return nil, err
	}
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

type fcpxmlNode struct {
	XMLName   xml.Name         `xml:"fcpxml"`
	Version   string           `xml:"version,attr"`
	Resources resourcesNode    `xml:"resources"`
	Library   libraryEventNode `xml:"library"`
}

type resourcesNode struct {
	Format formatNode `xml:"format"`
}

type formatNode struct {
	ID        string `xml:"id,attr"`
	FrameRate string `xml:"frameDuration,attr"`
	Width     int    `xml:"width,attr"`
	Height    int    `xml:"height,attr"`
}

type libraryEventNode struct {
	Event eventNode `xml:"event"`
}

type eventNode struct {
	Name    string      `xml:"name,attr"`
	Project projectNode `xml:"project"`
}

type projectNode struct {
	Name     string       `xml:"name,attr"`
	Sequence sequenceNode `xml:"sequence"`
}

type sequenceNode struct {
	Format   string    `xml:"format,attr"`
	Duration string    `xml:"duration,attr"`
	Spine    spineNode `xml:"spine"`
}

type spineNode struct {
	Clips []videoClipNode `xml:"video"`
}

type videoClipNode struct {
	Name     string `xml:"name,attr"`
	Offset   string `xml:"offset,attr"`
	Duration string `xml:"duration,attr"`
}

func buildFCPXMLData(spec FCPXMLSpec) fcpxmlNode {
	totalSec := 0.0
	var spineClips []videoClipNode
	for _, c := range spec.Clips {
		spineClips = append(spineClips, videoClipNode{
			Name:     c.Name,
			Offset:   fmt.Sprintf("%0.2fs", totalSec),
			Duration: fmt.Sprintf("%0.2fs", c.DurationSec),
		})
		totalSec += c.DurationSec
	}
	return fcpxmlNode{
		Version: "1.9",
		Resources: resourcesNode{
			Format: formatNode{
				ID:        "r1",
				FrameRate: fmt.Sprintf("1/%ds", spec.FrameRate),
				Width:     spec.Width,
				Height:    spec.Height,
			},
		},
		Library: libraryEventNode{
			Event: eventNode{
				Name: spec.ProjectName,
				Project: projectNode{
					Name: spec.ProjectName,
					Sequence: sequenceNode{
						Format:   "r1",
						Duration: fmt.Sprintf("%0.2fs", totalSec),
						Spine:    spineNode{Clips: spineClips},
					},
				},
			},
		},
	}
}
