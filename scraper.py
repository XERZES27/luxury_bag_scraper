from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.remote.webdriver import WebDriver
from csv_writer import *

chrome_options = webdriver.ChromeOptions()
prefs = {"profile.managed_default_content_settings.images": 2}
chrome_options.add_experimental_option("prefs", prefs)
browser = webdriver.Chrome(chrome_options)


def pageDoesNotExist(driver: WebDriver):
    return "Page Not Found" in driver.title.title()


def scrape_products_from_category(category_name: str):
    products_metadata = []
    page_count = 0
    while True:
        page_count += 1
        url = "https://luxy-bags.ru/product-category/{category_name}/page/{page_count}/".format(
            category_name=category_name, page_count=page_count
        )
        print(url)
        browser.get(url)
        if pageDoesNotExist(browser):
            break
        box_images = browser.find_elements(
            by=By.CSS_SELECTOR, value="div[class='product-small box ']"
        )
        
        
        for i in range(len(box_images)):
            try:
                product_metadata = []
                box = box_images[i]
                
                # Get product category
                category = box.find_element(by=By.CLASS_NAME, value="product-cat")
                product_metadata.append(category.text)
                # Get product title
                name = box.find_element(by=By.CLASS_NAME, value="product-title")
                title_link = name.find_element(by=By.TAG_NAME, value="a")
                product_metadata.append(title_link.text)

                # Get product presale price andn sale price
                prices = box.find_elements(
                    by=By.CLASS_NAME, value="woocommerce-Price-amount"
                )
                for price in prices:
                    product_metadata.append(price.text)

                # Get product url
                img_box = box.find_element(by=By.CLASS_NAME,value="image-fade_in_back")
                img = img_box.find_element(
                    by=By.CLASS_NAME, value="size-woocommerce_thumbnail"
                )
                image_src = img.get_attribute("src")
                if image_src is not None:
                    product_metadata.append(image_src)
                print(img.get_attribute('class'))
                products_metadata.append(product_metadata)
            except NoSuchElementException:
                print(url,str(i))
                continue

    return products_metadata


categories = ["hermes","dior", "louis-vuitton", "gucci", "chanel"]

for category in categories:
    products_of_category = scrape_products_from_category(category)
    write_to_fashion_bags_csv(category, products_of_category)
