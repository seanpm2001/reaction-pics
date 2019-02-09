import os

def remove_image_from_csvs(image):
    for csv in os.listdir("."):
        if csv[-4:] != ".csv":
            continue
        with open(csv, "r") as h:
            lines = h.readlines()
        for line in lines:
            if image in line:
                print(line)
        lines = [x for x in lines if image not in x]
        with open(csv, "w") as h:
            h.write("".join(lines))


images = os.listdir("static")
for image in images:
    if image == ".gitkeep":
        continue
    path = os.path.join("static", image)
    size = os.path.getsize(path)
    if size != 0:
        continue
    print(image)
    remove_image_from_csvs(image)
    os.remove(path)
