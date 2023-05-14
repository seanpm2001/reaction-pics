import csv

import requests


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

    def internal_image_location(self) -> str:
        if not self.media_urls:
            return ''
        extension = self.media_urls.split('.')[-1].split('?')[0]
        filename = './media/%s.%s' % (self.tweet_id, extension)
        return filename

    def download_tweet(self) -> None:
        tweet_location = self.internal_image_location()
        if not tweet_location:
            return
        response = requests.get(self.media_urls)
        with open(tweet_location, 'wb') as handle:
            handle.write(response.content)


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
        # row.download_tweet()
        print(row.internal_image_location())


if __name__ == '__main__':
    main()
