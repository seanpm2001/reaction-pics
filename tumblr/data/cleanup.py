import os
import shutil

from PIL import Image

for image in os.listdir('static'):
    if image == '.gitkeep':
        continue
    path = 'static/' + image
    try:
        im = Image.open(path)
        im.verify()
    except Exception as e:
        shutil.move(path, 'corrupt/' + image)
        # print(path, e)
        print(image)
