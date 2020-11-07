from base64 import urlsafe_b64encode as btoa
from csv import DictReader
import json
import sys
from pathlib import Path


def main():
    ifname, odname = "", ""
    try:
        ifname, odname = sys.argv[1:3]
    except ValueError:
        pass
    if not ifname or not odname:
        sys.exit("error: must specify input file name and output directory name")
    process(ifname, odname)


def process(ifname, odname):
    print(f"input: {ifname!r}\noutput: {odname!r}")
    with open(ifname) as f:
        rows = list(DictReader(f))

    canonical_rows = [
        "title",
        "linktitle",
        "first",
        "middle",
        "last",
        "honorific",
        "pronoun",
        "role",
        "expertise",
        "keywords",
        "email",
        "images",
        "website",
        "facebook",
        "twitter",
        "instagram",
        "linkedin",
        "location",
        "phone",
        "bio",
        "layout",
    ]

    newrows = []
    for row in rows:
        for field in row:
            row[field] = row[field].strip()
        row["title"] = soft_join(row, "first", "middle", "last")
        row["linktitle"] = soft_join(row, "first", "last")
        row["expertise"] = trim_all(row["expertise"].title().split(","))
        row["keywords"] = trim_all(row["keywords"].split(","))
        row["email"] = btoa(row["email"].encode("utf8")).decode("ascii")
        row["layout"] = "person"
        row["images"] = []
        newrows.append({key: row.get(key, "") for key in canonical_rows})

    for row in newrows:
        fname = Path(odname, row["linktitle"] + ".md")
        with open(fname, "w") as f:
            json.dump(row, f, indent=2)


def soft_join(row, *fields):
    s = " ".join(row[field] for field in fields)
    return s.replace("  ", " ").strip()


def trim_all(ss):
    r = []
    for s in ss:
        s = s.strip()
        if s:
            r.append(s)
    return r


if __name__ == "__main__":
    main()
