import os

images = os.listdir('static')
with open('posts.csv', 'r') as handle:
    lines = handle.readlines()
lines = [l.split(",") for l in lines]
lines = [l for l in lines if l[-2] in images]
lines = [",".join(l) for l in lines]

with open('posts.csv', 'w') as handle:
    handle.write("".join(lines))
