package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type PokemonProduct struct {
	url, image, name, price string
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
		"url",
		"image",
		"name",
		"price",
	}
	writer.Write(headers)
	defer writer.Flush()

	// initializing the slice of structs to store the data to scrape
	var pokemonProducts []PokemonProduct

	// creating a new Colly instance
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	// visiting the target page

	// scraping logic
	c.OnHTML("li.product", func(e *colly.HTMLElement) {

		pokemonProduct := PokemonProduct{}

		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		// fmt.Printf("url: %s, image: %s, name: %s, price: %s", pokemonProduct.url, pokemonProduct.image, pokemonProduct.name, pokemonProduct.price)

		pokemonProducts = append(pokemonProducts, pokemonProduct)
		record := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}
		

		// adding a CSV record to the output file
		writer.Write(record)
	})

	c.OnScraped(func(r *colly.Response) {
		
		fmt.Println("Finished", len(pokemonProducts))

	})


	visitingErr := c.Visit("https://scrapeme.live/shop/")
	if visitingErr != nil {
		log.Fatalln("Failed to open url", visitingErr)
	}
	

}
