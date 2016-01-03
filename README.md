# tcgif

True Color Gif Generator

![Sample](https://raw.githubusercontent.com/donatj/tcgif/master/sample.gif)

## Why?

For kicks. I saw a similar demo many years back and was curious how hard it would be to reproduce.

## How?

Gifs have a limited pallet of 255 colors, but that limit is per frame. Gifs also use transparency as a naive compression method such that the previous frame can show through. By utilizing this, we can in fact get 16 million colors in a single Gif!

## todo
- Much less naive frame generation
