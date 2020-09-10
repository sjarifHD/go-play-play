package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/vmihailenco/msgpack"

	"go.playplay.example/helper"
)

func main() {
	testGzip()

	app := fiber.New()

	app.Get("/", getFeedLogsJSON)
	app.Get("/gzip", middleware.Compress(2), getFeedLogsJSONGzip)
	app.Get("/msgpack", getFeedLogsMsgpck)

	app.Listen(2516)
}

func getFeedLogsJSON(c *fiber.Ctx) {
	feedLogsMap, err := helper.ReadJSONAsMap()

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

	feedLogsMap, err := helper.ReadJSONAsMap()

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	if err := c.JSON(feedLogsMap); err != nil {
		c.Status(500).Send(err)
		return
	}
}

func getFeedLogsMsgpck(c *fiber.Ctx) {
	feedLogsMap, err := helper.ReadJSONAsMap()

	fmt.Println("original size:\t", len(feedLogsMap))

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	feedLogsPacked, err := msgpack.Marshal(feedLogsMap)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	fmt.Println("packed size", len(feedLogsPacked))

	var feedLogsUnpacked map[string]interface{}
	err = msgpack.Unmarshal(feedLogsPacked, &feedLogsUnpacked)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	fmt.Println("unpacked size", len(feedLogsUnpacked))
	fmt.Println(reflect.DeepEqual(feedLogsMap, feedLogsUnpacked))
}

func testGzip() {
	data, err := helper.ReadJSONAsString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("original size:\t", len(data))

	// compress data
	compressedData, compressedDataErr := helper.GzipData([]byte(data))
	if compressedDataErr != nil {
		log.Fatal(compressedDataErr)
	}
	fmt.Println("compressed data len:", len(compressedData))

	// uncompress data
	uncompressedData, uncompressedDataErr := helper.GunzipData(compressedData)
	if uncompressedDataErr != nil {
		log.Fatal(uncompressedDataErr)
	}
	fmt.Println("uncompressed data len:", len(uncompressedData))
}
