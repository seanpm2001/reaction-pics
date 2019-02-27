import json
import os

with open('duplicates.json', 'r') as handle:
    duplicates = json.loads(handle.read())

with open('posts.csv', 'r') as handle:
    posts = handle.readlines()

for duplicate in duplicates:
    if duplicate[0] != '03a0e98e-0f1a-46d8-99e9-67d01094bafa.gif':
        continue
    for d in duplicate:
        print(d)
        posts = [l for l in posts if d not in l]
        try:
            os.remove('static/' + d)
        except FileNotFoundError:
            pass

with open('posts.csv', 'w') as handle:
    handle.write(''.join(posts))
