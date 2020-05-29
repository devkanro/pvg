package cmd

import (
	"fmt"
	"github.com/devkanro/pvg/lib"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func handleDir(input string, output string, transparent color.Color) error {
	files, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}

	stat, _ := os.Stat(output)
	if stat == nil {
		err = os.MkdirAll(output, os.ModeDir)
	} else if !stat.IsDir() {
		return fmt.Errorf("output '%s' must be a folder in batch mode", output)
	}

	if parallel {
		var waitGroup sync.WaitGroup
		cpus := runtime.NumCPU()
		waitGroup.Add(cpus)
		for i := 0; i < cpus; i++ {
			go func(group int) {
				handleFolderGroup(input, output, transparent, files, cpus, group)
				waitGroup.Done()
			}(i)
		}
		waitGroup.Wait()
	}else {
		handleFolderGroup(input, output, transparent, files, -1, -1)
	}

	return nil
}

func handleFolderGroup(input string, output string, transparent color.Color, files []os.FileInfo, groupCount int, group int) {
	for i, file := range files {
		if file.IsDir() {
			continue
		}
		if group < 0 || groupCount <= 0 || i % groupCount == group {
			basename := file.Name()
			basename = basename[:len(basename)-len(filepath.Ext(basename))]
			err := handleFile(filepath.Join(input, file.Name()), filepath.Join(output, basename+".svg"), transparent)
			if err != nil {
				fmt.Printf("Skip convert file '%s' due to error: %s\n", file.Name(), err)
			}
		}
	}
}

func handleFile(input string, output string, transparent color.Color) error {
	file, err := os.Open(input)
	if err != nil {
		return err
	}

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("Open file failed: %s ", err)
	}
	fmt.Printf("%s>>svg: '%s'.\n", format, input)

	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}

	_, err = outputFile.WriteString(lib.Convert(img, transparent))
	if err != nil {
		return err
	}

	return nil
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 9:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x%02x", &c.A, &c.R, &c.G, &c.B)
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 4, 7 or 9")
	}
	return
}