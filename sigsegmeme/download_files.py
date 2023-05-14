import csv


class CSVData:
    def __init__(
        self,
        tweet_id: str,
        text: str,
        media_urls: str,
        favorites: int,
    ) -> None:
        self.tweet_id = tweet_id
        self.text = text
        self.media_urls = media_urls
        self.favorites = favorites

    def url(self) -> str:
        return 'https://twitter.com/sigsegmeme/status/%s' % (self.tweet_id)


def parse_csv() -> list[CSVData]:
    data: list[CSVData] = []
    with open('./tweets.csv') as handle:
        csv_reader = csv.DictReader(handle)
        for row in csv_reader:
            data.append(CSVData(
                row['Tweet Id'],
                row['Text'],
                row['Media URLs'],
                int(row['Favorites']),
            ))
    return data


def main() -> None:
    csv_data = parse_csv()
    for row in csv_data:
        print(row.__dict__)


if __name__ == '__main__':
    main()
