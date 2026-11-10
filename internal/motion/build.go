package motion

import (
	"fmt"
	"strings"
)

// BuildHTML generates an HTML/WebGL page for the selected repos.
func BuildHTML(repos []string, prompt string, resolution string) (string, error) {
	if len(repos) == 0 {
		return "", ErrNoRepoSelected
	}
	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><html><head><meta charset="utf-8">`)
	b.WriteString(fmt.Sprintf(`<title>%s</title>`, prompt))
	b.WriteString(`<style>body{margin:0;overflow:hidden;background:#000}</style>`)
	b.WriteString(`</head><body>`)
	b.WriteString(`<canvas id="c"></canvas>`)
	for _, repo := range repos {
		switch repo {
		case "threejs":
			b.WriteString(`<script type="importmap">{"imports":{"three":"https://unpkg.com/three@0.160.0/build/three.module.js"}}</script>`)
			b.WriteString(`<script type="module">import * as THREE from "three";`)
			b.WriteString(`const s=new THREE.Scene();const cam=new THREE.PerspectiveCamera(75,innerWidth/innerHeight,0.1,1000);`)
			b.WriteString(`const r=new THREE.WebGLRenderer({canvas:document.getElementById("c")});`)
			b.WriteString(`r.setSize(innerWidth,innerHeight);cam.position.z=5;`)
			b.WriteString(`function a(){requestAnimationFrame(a);r.render(s,cam)}a();</script>`)
		case "pixelfx":
			b.WriteString(`<script>const c=document.getElementById("c");c.style.imageRendering="pixelated";</script>`)
		}
	}
	b.WriteString(`<script>document.title=`)
	b.WriteString(fmt.Sprintf(`"%s";`, prompt))
	b.WriteString(`</script></body></html>`)
	return b.String(), nil
}
