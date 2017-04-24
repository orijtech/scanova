package scanova_test

import (
	"io"
	"log"
	"os"

	"github.com/orijtech/scanova"
)

func Example_client_NewQRCode() {
	client, err := scanova.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	req := &scanova.Request{
		URL: "https://github.com/orijtech/scanova",

		Size:            scanova.LargeSize,
		ErrorCorrection: scanova.LevelQ,

		EyePattern: scanova.CircularCircle,
	}

	res, err := client.NewQR(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	f, err := os.Create("./testdata/new-qr.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	io.Copy(f, res)
}
