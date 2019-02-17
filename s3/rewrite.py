import csv

with open("../tumblr/data/posts.csv", "r") as handle:
    lines = handle.readlines()

lines = [l.replace("/static/data/", "") for l in lines]

with open("../tumblr/data/posts.csv", "w") as handle:
    handle.write("".join(lines))
