package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	asciiChars = "@B%8WM#*oahkbdpwmZO0QCJYXzcvnxrjft/|()1{}[]-_+~<>i!lI;:,^`'. "
	scale_x      = 1.0 // 缩放因子，调整输出的 ASCII 图像大小
    scale_y      = 1.0
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请指定图片文件路径")
		return
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("无法打开图片文件:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("无法解码图片:", err)
		return
	}
	//

    w , h  :=  getWindowSize();
    if img.Bounds().Max.Y > h {
        scale_y = float64(h) * 1  / float64(img.Bounds().Max.Y)
    }
    if img.Bounds().Max.X > w {
        scale_x =float64(w) * 0.8 / float64(img.Bounds().Max.X)
    }

	width := int(float64(img.Bounds().Dx()) * scale_x)
	height := int(float64(img.Bounds().Dy()) * scale_y)

	img = resize(img, width, height)

	asciiImage := convertToASCII(img)

	fmt.Println(asciiImage)
}

func resize(img image.Image, width, height int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	dx := float64(img.Bounds().Dx()) / float64(width)
	dy := float64(img.Bounds().Dy()) / float64(height)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * dx)
			srcY := int(float64(y) * dy)
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return resized
}

func convertToASCII(img image.Image) string {
	asciiImage := ""

	bounds := img.Bounds()


	for y := bounds.Min.Y; y < bounds.Max.Y ; y++  {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := getPixelGray(img.At(x, y))
			charIndex := int(float64(gray) / 255 * float64(len(asciiChars)-1))
			asciiImage += string(asciiChars[charIndex])
		}
		asciiImage += "\n"
	}

	return asciiImage
}

func getPixelGray(pixel color.Color) uint8 {
	r, g, b, _ := pixel.RGBA()
	
	gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	
	return uint8(gray / 256)
}

func getWindowSize()(int, int){
    fd := int(os.Stdin.Fd())

    w , h , _ := terminal.GetSize(fd)
    return w , h
}
