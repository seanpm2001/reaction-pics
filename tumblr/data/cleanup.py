import os

images = os.listdir("static")
with open("posts.csv", "r") as h:
    lines = h.readlines()

x = 0
while x < len(lines):
    line = lines[x]
    line_image = line[line.rfind("/")+1:line.rfind(",")]
    if line_image not in images:
        print(line)
        del lines[x]
        continue
    x += 1
