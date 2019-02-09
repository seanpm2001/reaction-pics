import os

data = ""
for csv in os.listdir("."):
    if csv[-4:] != ".csv":
        continue
    with open(csv, "r") as h:
        data += h.read()


images = os.listdir("static")
for image in images:
    if image == ".gitkeep":
        continue
    path = os.path.join("static", image)
    size = os.path.getsize(path)
    if image not in data:
        print(image)
        os.remove(path)
