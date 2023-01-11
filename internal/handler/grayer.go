package handler

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
)

func handleGrayer(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if len(m.Attachments) != 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "you have to attach one photo")
		return err
	}

	attach := m.Attachments[0]
	resp, err := http.Get(attach.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	img, err := png.Decode(resp.Body)
	if err != nil {
		return err
	}

	bounds := img.Bounds().Bounds()
	outputImage := image.NewRGBA(bounds)
	for i := 0; i < bounds.Dx(); i++ {
		for j := 0; j < bounds.Dy(); j++ {
			pixel := img.At(i, j)
			pc := color.RGBAModel.Convert(pixel).(color.RGBA)

			grey := uint8(math.Round((float64(pc.R) + float64(pc.G) + float64(pc.B)) / 3))

			outputImage.Set(i, j, color.RGBA{
				R: grey,
				G: grey,
				B: grey,
				A: pc.A,
			})
		}
	}

	bb := new(bytes.Buffer)
	err = png.Encode(bb, outputImage)
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "grayed.png", bb)

	return err
}
