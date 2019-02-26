import json

with open("duplicates.csv", "r") as handle:
    lines = handle.readlines()

groups = []
group = []
for line in lines:
    line = line.strip()
    if line == '':
        groups.append(group)
        group = []
        continue
    image = line.split(" ")[-1]
    image = image[2:]
    group.append(image)
groups.append(group)

print(json.dumps(groups, indent=2))
