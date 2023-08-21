import csv
from typing import Any, Iterable


def write_to_fashion_bags_csv(file_name:str, rows: Iterable[Iterable[Any]]):
    with open("csv/{file_name}.csv".format(file_name=file_name), "w", encoding="utf-8", newline="") as file:
        writer = csv.writer(file)

        writer.writerow(
            ["Category", "Title", "Pre Sale Price", "Sale Price", "Product Thumbnail","Product Link"]
        )
        writer.writerows(rows)
