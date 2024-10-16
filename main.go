package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/golang/freetype"
	"gopkg.in/yaml.v3"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd()
	os.Exit(0)
}

type ConfigTypes struct {
	Imports           string   `yaml:"imports"`
	Exports           string   `yaml:"exports"`
	ImageType         string   `yaml:"imagetype"`
	UnicodeStart      int32    `yaml:"unicodestart"`
	UnicodeEnd        int32    `yaml:"unicodeend"`
	Captions          bool     `yaml:"captions"`
	AdditionalPrompts []string `yaml:"additionalprompts"`
}

func cmd() {
	config := ConfigTypes{}
	configBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	if err = yaml.Unmarshal(configBytes, &config); err != nil {
		log.Fatalln(err)
	}

	fontfiles, err := filepath.Glob(config.Imports + "/*.???")
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range fontfiles {
		fontname := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		log.Println("Generating: " + fontname)
		fontBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		font, err := freetype.ParseFont(fontBytes)
		if err != nil {
			log.Fatalln(err)
		}

		additionalpromptstr := strings.Join(config.AdditionalPrompts, ", ")

		start := rune(config.UnicodeStart)
		end := rune(config.UnicodeEnd)

		for r := start; r <= end; r++ {
			if !unicode.IsGraphic(r) {
				continue
			}

			img := image.NewRGBA(image.Rect(0, 0, 512, 512))
			draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.Point{}, draw.Src)

			c := freetype.NewContext()
			c.SetDPI(72)
			c.SetFont(font)
			c.SetFontSize(448)
			c.SetClip(img.Bounds())
			c.SetDst(img)
			c.SetSrc(&image.Uniform{C: color.RGBA{R: 0, G: 0, B: 0, A: 255}})
			pt := freetype.Pt(32, 417)
			_, err := c.DrawString(string(r), pt)
			if err != nil {
				log.Fatalln(err)
			}

			filename := fmt.Sprintf("%s/%s/u%06x.%s", config.Exports, fontname, r, config.ImageType)
			dir := filepath.Dir(filename)
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalln(err)
			}
			f, err := os.Create(filename)
			if err != nil {
				log.Fatalln(err)
			}

			if config.ImageType == "png" {
				if err := png.Encode(f, img); err != nil {
					log.Fatalln(err)
				}
			} else if config.ImageType == "jpg" || config.ImageType == "jpeg" {
				if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 100}); err != nil {
					log.Fatalln(err)
				}
			}

			f.Close()

			if config.Captions {
				caption := fmt.Sprintf("%s/%s/u%06x.txt", config.Exports, fontname, r)
				f, err = os.Create(caption)
				if err != nil {
					log.Fatalln(err)
				}
				if _, err := f.WriteString(fmt.Sprintf("u%06x, %s, %s, %s\n", r, string(r), fontname, additionalpromptstr)); err != nil {
					log.Fatalln(err)
				}
				f.Close()
			}
		}
	}
	log.Println("Done")
}
