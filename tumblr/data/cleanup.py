import os

def remove_image_from_csvs(image):
    for csv in os.listdir("."):
        if csv[-4:] != ".csv":
            continue
        with open(csv, "r") as h:
            lines = h.readlines()
        modified = False
        for x in range(len(lines)):
            if image in lines[x]:
                del lines[x]
                modified = True
                print(lines[x])
        if not modified:
            continue
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
