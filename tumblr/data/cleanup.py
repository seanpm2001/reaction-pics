import os


images = os.listdir('corrupt')

with open('posts.csv', 'r') as h:
    data = h.readlines()

data = [l for l in data if not any([i in l for i in images])]

with open('posts.csv', 'w') as h:
    h.write("".join(data))
