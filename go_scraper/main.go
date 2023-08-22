package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	Categories, SKU, Title, Pre_Sale_Price, Sale_Price, Product_Description, Product_Post_Image, Product_Images string
}

func main() {

	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// writing the CSV headers
	headers := []string{
		"Categories", "SKU", "Title", "Pre Sale Price", "Sale Price", "Product Description", "Product Post Image", "Product Images",
	}
	writer.Write(headers)
	defer writer.Flush()

	// initializing the slice of structs to store the data to scrape
	// var products []Product

	// creating a new Colly instance
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	c.OnHTML("div.product-container", func(e *colly.HTMLElement) {
		product := Product{}

		//Get Product Categories
		var productCategories = ""
		categories := e.ChildTexts("div.product_meta>span.posted_in>a")
		for _, category := range categories {
			if productCategories == "" {
				productCategories += category
			} else {
				productCategories += "|" + category
			}
		}
		product.Categories = productCategories

		//Get Product SKU
		productSKU := e.ChildText("div.product_meta>span.sku_wrapper>span.sku")
		product.SKU = productSKU

		//Get Product Title
		title := e.ChildText("div>h1.product-title.product_title.entry-title")
		product.Title = title

		//Get Product Pre Sale Price
		var preSalePrice string = ""
		var salePrice string = ""

		prices := e.ChildTexts("span.woocommerce-Price-amount.amount")
		if len(prices) >= 2 {
			preSalePrice = prices[0]
			salePrice = prices[1]
		}
		if len(prices) == 1 {
			salePrice = prices[0]
		}
		product.Pre_Sale_Price = preSalePrice
		product.Sale_Price = salePrice

		//Get Product Description
		description := e.ChildText("div#tab-description")
		product.Product_Description = description

		//Get Product Images
		var productImage string = ""
		var productImages string = ""
		imageUrls := e.ChildAttrs("div.woocommerce-product-gallery__image.slide>a", "href")
		for index, v := range imageUrls {
			if index == 0 {
				productImage = v
			} else {
				if productImages == "" {
					productImages += v
				} else {
					productImages += "|" + v
				}
			}
		}
		product.Product_Post_Image = productImage
		product.Product_Images = productImages

		record := []string{
			product.Categories,
			product.SKU,
			product.Title,
			product.Pre_Sale_Price,
			product.Sale_Price,
			product.Product_Description,
			product.Product_Post_Image,
			product.Product_Images}

		// adding a CSV record to the output file
		err = writer.Write(record)
		if err != nil {
			log.Println(record)
			log.Fatal(err)
		}

	})
	c.OnScraped(func(r *colly.Response) {

	})
	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Failure", err)
	})

	categories := []string{"hermes","dior", "louis-vuitton", "gucci", "chanel"}

	for _, cat := range categories {
		category_file, err := os.Open(fmt.Sprintf("csv/%s.csv", cat))
		log.Println(cat)
		if err != nil {
			log.Fatal(err)
		}

		fileReader := csv.NewReader(category_file)

		count := 0
		for {

			record, err := fileReader.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Panicln(record[5])
				log.Fatal(err)
			}

			if count > 0 {
				if len(record) == 6 {
					err = c.Visit(record[5])
					if err != nil {
						log.Panicln(record[5])
						log.Fatal(err)
					}
				}

			}
			count++
		}
	}

	c.Wait()

}
