import json
import sys
from pathlib import Path


def main():
    idir = ""
    try:
        idir = sys.argv[1]
    except ValueError:
        pass
    if not idir:
        sys.exit("error: must specify input directory name")
    process(idir)


def process(idir):
    print(f"input directory: {idir!r}")

    root = Path(idir)
    for p in root.glob('*.md'):
        if p.name.startswith("_"):
            print(f"skipping {p!r}...")
            continue
        transform(p)

def transform(p):
    print(f"transforming {p.name}")
    with open(p) as f:
        data = json.load(f)

    mutate(p, data)

    with open(p, 'w') as f:
        json.dump(data, f, indent=2)
        f.write('\n')

def mutate(p, data):
    data['location'] = [data['location']]

if __name__ == "__main__":
    main()
