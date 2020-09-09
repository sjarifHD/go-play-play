package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func main() {
	testGzip()

	app := fiber.New()

	app.Get("/", getFeedLogsJSON)
	app.Get("/compress", middleware.Compress(2), getFeedLogsJSONGzip)

	app.Listen(2516)
}

func getFeedLogsJSON(c *fiber.Ctx) {
	feedLogsMap, err := readJSONAsMap()

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	if err := c.JSON(feedLogsMap); err != nil {
		c.Status(500).Send(err)
		return
	}
}

func getFeedLogsJSONGzip(c *fiber.Ctx) {
	c.Type("gzip")

	feedLogsMap, err := readJSONAsMap()

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	if err := c.JSON(feedLogsMap); err != nil {
		c.Status(500).Send(err)
		return
	}
}

func readJSONAsString() (string, error) {
	jsonFile, err := os.Open("resources/feedlogs.json")

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(byteValue), nil
}

func readJSONAsMap() (map[string]interface{}, error) {
	jsonFile, err := os.Open("resources/feedlogs.json")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func testGzip() {
	data, err := readJSONAsString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("original size:\t", len(data))

	// compress data
	compressedData, compressedDataErr := gZipData([]byte(data))
	if compressedDataErr != nil {
		log.Fatal(compressedDataErr)
	}
	fmt.Println("compressed data len:", len(compressedData))

	// uncompress data
	uncompressedData, uncompressedDataErr := gUnzipData(compressedData)
	if uncompressedDataErr != nil {
		log.Fatal(uncompressedDataErr)
	}
	fmt.Println("uncompressed data len:", len(uncompressedData))
}

func gZipData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	if err = gz.Flush(); err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}

	compressedData = b.Bytes()

	return
}

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}
