#!/bin/sh
# Renders the PNG icons in static/ from static/favicon.svg and icons/maskable.svg (needs rsvg-convert).
# Rerun after changing either SVG.
set -e
cd "$(dirname "$0")/.."
mkdir -p static/icons
rsvg-convert -w 48 -h 48 static/favicon.svg -o static/favicon.png
rsvg-convert -w 192 -h 192 static/favicon.svg -o static/icons/icon-192.png
rsvg-convert -w 512 -h 512 static/favicon.svg -o static/icons/icon-512.png
rsvg-convert -w 512 -h 512 icons/maskable.svg -o static/icons/icon-maskable-512.png
rsvg-convert -w 180 -h 180 icons/maskable.svg -o static/apple-touch-icon.png # iOS rounds the corners itself
