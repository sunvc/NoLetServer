package config

import (
	"embed"
	"encoding/base64"
	"strings"
)

const LOGORAW = `
<svg id="noletLogo" data-name="logo" xmlns="http://www.w3.org/2000/svg" version="1.1" viewBox="0 0 500 500" fill="currentColor">
	<path class="cls-1" d="M133.1,441.7c-.1,0-.4,0-.4-.4v-132.9l235.1-240.4c8.9-9.2,20.8-14.2,33.5-14.2s24.4,5.1,33,14.2c8.4,8.9,13,20.7,12.9,33.2s-4.9,24.2-13.5,33.1L133.3,441.6c-.1.1-.2.1-.2.1h0Z"/>
	<polygon class="cls-1" points="85.9 411.7 39.2 400.6 6.5 436.4 20.4 483.2 67.1 494.3 99.8 458.5 85.9 411.7"/>
	<path class="cls-1" d="M447.3,7.6H114L49.4,73.7l58.1,28.5h339.8c25.5,0,46.3-21.2,46.3-47.3s-20.7-47.3-46.2-47.3h0Z"/>
	<path class="cls-1" d="M445.4,447.5l48.1-49.2V54.9c0-26.1-20.7-47.3-46.3-47.3s-46.3,21.2-46.3,47.3v322.6l44.4,69.9h0Z"/>
</svg>
`

var StaticFS *embed.FS

func LogoSvgImage(color string, svg bool) string {
	color1 := "#ff0000"
	if color != "" {
		color1 = color
	}
	logosvg := strings.Replace(LOGORAW, "currentColor", color1, 1)
	if svg {
		return logosvg
	}
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(logosvg))
}
