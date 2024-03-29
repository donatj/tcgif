# tcgif

[![Go Reference](https://pkg.go.dev/badge/github.com/donatj/tcgif.svg)](https://pkg.go.dev/github.com/donatj/tcgif)
[![Go Report Card](https://goreportcard.com/badge/github.com/donatj/tcgif)](https://goreportcard.com/report/github.com/donatj/tcgif)


Trueclor Gif Generator

## Why?

Just for kicks. I saw a similar demo many years back and was curious how hard it would be to reproduce.

## How?

Gifs have a limited pallet of 255 colors, but that limit is per frame. Gifs also use transparency as a naive compression method such that the previous frame can show through. By utilizing this, we can in fact get 16 million colors in a single Gif!

## todo

- Less naive backfill. Track if the current backfilled color would be closer than a new one.


## Samples

| Description                                                        | Image                                                                                                              | 
|--------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------| 
| No Backfill, Unsorted <br> `tcgif -backfill=false -sort=false orig.jpg` | ![Lenna, No Backfill, No Sort](https://raw.githubusercontent.com/donatj/tcgif/images/sample_nobackfill_nosort.gif) | 
| No Backfill, Sorted By Popularity <br> `tcgif -backfill=false orig.jpg` | ![Lenna, No Backfill](https://raw.githubusercontent.com/donatj/tcgif/images/sample_nobackfill.gif)                 | 
| Backfilled <br> `tcgif orig.jpg`                                        | ![Lenna, Backfilled](https://raw.githubusercontent.com/donatj/tcgif/images/sample_backfill.gif)                    | 
| Original                                                           | ![Original](https://raw.githubusercontent.com/donatj/tcgif/images/orig.jpg)                                        | 
