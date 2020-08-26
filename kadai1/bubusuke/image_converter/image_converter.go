package image_converter

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

type imageConverter struct {
	decode      func(io.Reader) (image.Image, error)
	encode      func(io.Writer, image.Image) error
	befExts     map[string]bool
	aftExt      string
	befFileName string
}

var befExtsDef = map[string]map[string]bool{
	"jpeg": {".jpg": true, ".jpeg": true, ".JPG": true, ".JPEG": true},
	"png":  {".png": true, ".PNG": true},
	"gif":  {".gif": true, ".GIF": true},
}

func New() *imageConverter {
	ic := &imageConverter{}
	ic.FromJPEG()
	ic.ToPNG()
	return ic
}

func (ic *imageConverter) FromJPEG() {
	ic.decode = jpeg.Decode
	ic.befExts = befExtsDef["jpeg"]
}
func (ic *imageConverter) ToJPEG() {
	ic.encode = func(i io.Writer, im image.Image) error { return jpeg.Encode(i, im, &jpeg.Options{Quality: 100}) }
	ic.aftExt = ".jpg"
}
func (ic *imageConverter) FromPNG() {
	ic.decode = png.Decode
	ic.befExts = befExtsDef["png"]
}
func (ic *imageConverter) ToPNG() {
	ic.encode = png.Encode
	ic.aftExt = ".png"
}
func (ic *imageConverter) FromGIF() {
	ic.decode = gif.Decode
	ic.befExts = befExtsDef["gif"]
}
func (ic *imageConverter) ToGIF() {
	ic.encode = func(i io.Writer, im image.Image) error { return gif.Encode(i, im, nil) }
	ic.aftExt = ".gif"
}

func (ic *imageConverter) convert(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	img, err := ic.decode(file)
	if err != nil {
		return "", err
	}
	newPath := filepath.Join(filepath.Dir(path), ic.newFileName())
	if _, err := os.Stat(newPath); err == nil {
		return "", errors.New("A file with the same name as the converted file already exists. FILE: " + newPath)
	}
	outBuf, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer outBuf.Close()
	if err := ic.encode(outBuf, img); err != nil {
		return "", err
	}
	return newPath, nil
}

func (ic *imageConverter) ConvertFiles(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if _, hit := ic.befExts[filepath.Ext(path)]; hit {
				ic.befFileName = filepath.Base(path)
				_, err := ic.convert(path)
				if err != nil {
					fmt.Println("FAILURE TO CONVERT. file: ", path)
					fmt.Println(err)
				} else {
					fmt.Printf("SUCCESS TO CONVERT. %v -> %v \n", ic.befFileName, ic.newFileName())
				}
			}
		}
		return nil
	})
}

func (ic *imageConverter) newFileName() string {
	return fmt.Sprintf("CONVERT_%v%v",
		ic.befFileName[:len(ic.befFileName)-len(filepath.Ext(ic.befFileName))],
		ic.aftExt)
}
