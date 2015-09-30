package magick

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkNewFromBlob(b *testing.B) {
	for n := 0; n < b.N; n++ {
		filename := "big.jpg"
		source, _ := ioutil.ReadFile(filename)
		NewFromBlob(source, "jpg")
	}
}

func loadImage(t *testing.T, filename string) (image *MagickImage) {
	image, error := NewFromFile(filename)
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)
	return
}

func setupImage(t *testing.T) (image *MagickImage) {
	image = loadImage(t, "test/heart_original.png")
	return
}

func setupBackgroundImage(t *testing.T) (image *MagickImage) {
	image = loadImage(t, "test/blue_background.png")
	return
}

func TestImageFromFile(t *testing.T) {
	filename := "test/heart_original.png"
	image, error := NewFromFile(filename)
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)

	bad_filename := "test/heart_whatwhat.png"
	image, error = NewFromFile(bad_filename)
	assert.T(t, error != nil)
}

func TestImageFromBlob(t *testing.T) {
	filename := "test/heart_original.png"
	source, _ := ioutil.ReadFile(filename)
	image, error := NewFromBlob(source, "png")
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)
	image, error = NewFromBlob([]byte{}, "png")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
}

func TestBadDataFromBlob(t *testing.T) {
	filename := "test/heart_original.png"
	source, _ := ioutil.ReadFile(filename)
	image, error := NewFromBlob(source, "")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte(""), "png")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte("blah"), "jpg")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte("   "), "image/jpg")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil

	image, error = NewFromBlob([]byte(""), " ")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil
	image, error = NewFromBlob([]byte(""), ":")
	assert.T(t, error != nil)
	assert.T(t, image == nil)
	image = nil
	error = nil
}

func TestPDFFromBlob(t *testing.T) {
	filename := "test/heart_original.pdf"
	source, _ := ioutil.ReadFile(filename)
	image, error := NewFromBlob(source, "pdf")
	assert.T(t, error == nil)
	assert.T(t, image != nil)
	assert.T(t, image.Image != nil)
	assert.T(t, image.ImageInfo != nil)
}

func TestParseGeometry(t *testing.T) {
	image := setupImage(t)
	geometry, err := image.ParseGeometry("100x100>")
	assert.T(t, err == nil)
	assert.T(t, geometry != nil)
	assert.Equal(t, 100, geometry.Width)
}

func TestResizeRatio(t *testing.T) {
	image := setupImage(t)
	ratio := image.ResizeRatio(300, 300)
	assert.T(t, ratio > 0.27)
	assert.T(t, ratio < 0.28)
}

func TestStrip(t *testing.T) {
	image := setupImage(t)
	err := image.Strip()
	assert.T(t, err == nil)
}

func TestProgressive(t *testing.T) {
	image := setupImage(t)
	image.Progressive()
	_, err := image.ToBlob("jpg")
	assert.T(t, err == nil)
}

func TestDestroy(t *testing.T) {
	image := setupImage(t)
	assert.T(t, image.Destroy() == nil)
	assert.T(t, image.Image == nil)
	assert.T(t, image.ImageInfo == nil)
}

func TestResize(t *testing.T) {
	image := setupImage(t)
	err := image.Resize("100x100!")
	assert.T(t, err == nil)
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 100, image.Height())

	image = setupImage(t)
	err = image.Resize("blurgh")
	assert.T(t, err != nil)
}

func TestPDFResize(t *testing.T) {
	filename := "test/heart_original.pdf"
	source, _ := ioutil.ReadFile(filename)
	image, err := NewFromBlob(source, "pdf")
	assert.T(t, err == nil)
	err = image.Resize("100x100!")
	assert.T(t, err == nil)
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 100, image.Height())
	if err == nil {
		filename = "test/test_from_pdf.jpg"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
}

func TestCrop(t *testing.T) {
	image := setupImage(t)
	err := image.Crop("100x100!+10+10")
	assert.T(t, err == nil)
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 100, image.Height())

	image = setupImage(t)
	err = image.Crop("blurgh")
	assert.T(t, err != nil)
}

func TestShadow(t *testing.T) {
	image := setupImage(t)
	err := image.Shadow("#000", 75, 2, 0, 0)
	assert.T(t, err == nil)
}

func TestCompose(t *testing.T) {
	image := setupBackgroundImage(t)
	overlay := setupImage(t)

	assert.T(t, image.Width() == overlay.Width())
	assert.T(t, image.Height() == overlay.Height())

	err := image.Compose(ScreenCompositeOp, overlay, 150, 150)
	assert.T(t, err == nil)
	assert.T(t, image != nil)

	filename := "test/test_compose.png"
	os.Remove(filename)
	err = image.ToFile(filename)
	assert.T(t, err == nil)
}

func TestFillBackgroundColor(t *testing.T) {
	image := setupImage(t)
	err := image.FillBackgroundColor("#CCC")
	assert.T(t, err == nil)
}

func TestSeparateAlphaChannel(t *testing.T) {
	image := setupImage(t)
	err := image.SeparateAlphaChannel()
	assert.T(t, err == nil)
	assert.T(t, image != nil)
}

func TestNegateImage(t *testing.T) {
	image := setupImage(t)
	err := image.Negate()
	assert.T(t, err == nil)
	assert.T(t, image != nil)
}

func TestToBlob(t *testing.T) {
	image := setupImage(t)
	bytes, err := image.ToBlob("png")
	assert.T(t, err == nil)
	assert.T(t, bytes != nil)
	assert.T(t, len(bytes) > 0)
	assert.Equal(t, 437185, len(bytes))
}

func TestToFile(t *testing.T) {
	image := setupImage(t)
	filename := "test/test_out.png"
	os.Remove(filename)
	err := image.ToFile(filename)
	assert.T(t, err == nil)
	file, err := os.Open(filename)
	assert.T(t, err == nil)
	defer file.Close()
	stat, err := file.Stat()
	assert.T(t, stat != nil)
	assert.Equal(t, (int64)(437185), stat.Size())
}

func TestType(t *testing.T) {
	image := setupImage(t)
	assert.Equal(t, "PNG", image.Type())
}

func TestWidth(t *testing.T) {
	image := setupImage(t)
	assert.Equal(t, 600, image.Width())
}

func TestHeight(t *testing.T) {
	image := setupImage(t)
	assert.Equal(t, 552, image.Height())
}

func TestSetProperty(t *testing.T) {
	image := setupImage(t)
	err := image.SetProperty("jpeg:sampling-factor", "4:4:4")
	assert.T(t, err == nil)
	factor := image.GetProperty("jpeg:sampling-factor")
	assert.Equal(t, "4:4:4", factor)
}

func TestCropToTiles(t *testing.T) {
	image := setupImage(t)
	defer image.Destroy()

	err := image.CropToTiles("25%x100%")

	assert.Equal(t, err, nil)
	assert.Equal(t, image.ListLength(), 4)
}

func TestSplitList(t *testing.T) {
	image := setupImage(t)
	defer image.Destroy()

	err := image.CropToTiles("25%x100%")

	assert.Equal(t, err, nil)

	images, err := image.SplitList()
	assert.Equal(t, err, nil)
	for _, im := range images {
		defer im.Destroy()
	}

	assert.Equal(t, len(images), 4)
	for _, im := range images {
		assert.Equal(t, im.ListLength(), 1)
	}
}

func TestFullStack(t *testing.T) {
	var err error
	var filename string
	var image *MagickImage
	// resize
	image = setupImage(t)
	err = image.Resize("100x100")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_resize.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// crop
	image = setupImage(t)
	err = image.Crop("100x100+10+10")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_crop.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// shadow
	image = setupImage(t)
	err = image.Shadow("#000", 90, 10, 0, 0)
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_shadow.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// fill
	image = setupImage(t)
	err = image.FillBackgroundColor("red")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_fill.png"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// combination
	image = setupImage(t)
	err = image.Resize("100x100>")
	assert.T(t, err == nil)
	err = image.Shadow("#000", 90, 10, 0, 0)
	assert.T(t, err == nil)
	err = image.FillBackgroundColor("#CCC")
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_combo.jpg"
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// alpha mask
	image = setupImage(t)
	err = image.SeparateAlphaChannel()
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_alpha.jpg"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
	// alpha mask + negate
	image = setupImage(t)
	err = image.SeparateAlphaChannel()
	assert.T(t, err == nil)

	err = image.Negate()
	assert.T(t, err == nil)
	if err == nil {
		filename = "test/test_alpha_negative.jpg"
		log.Print(filename)
		os.Remove(filename)
		err = image.ToFile(filename)
		assert.T(t, err == nil)
	}
}
