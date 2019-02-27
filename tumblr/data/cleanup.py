import json
import os

with open('duplicates.json', 'r') as handle:
    duplicates = json.loads(handle.read())

with open('posts.csv', 'r') as handle:
    posts = handle.read()

for duplicate in duplicates:
    target = duplicate[0]
    for d in duplicate[1:]:
        posts = posts.replace(d, target)
        os.remove('static/' + d)

with open('posts.csv', 'w') as handle:
    handle.write(posts)
