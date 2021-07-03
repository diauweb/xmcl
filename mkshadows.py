#!/usr/bin/python3

import os
import hashlib
import json
import shutil
from dotenv import load_dotenv

load_dotenv()

def hash(path):
    sha1 = hashlib.sha1()
    with open(path, "rb") as f:
        while True:
            chunk = f.read(sha1.block_size)
            if not chunk:
                break
            sha1.update(chunk)
    return sha1.hexdigest()

def files():

    ret = []

    for folder, _, files in os.walk("./shadows"):
        for i in files:
            key = f"{folder.replace('./shadows/', '')}/{i}"
            key = i if folder == './shadows' else key

            path = f"{folder}/{i}"

            ret.append({
                "path": key,
                "url": f"{os.environ['ENDPOINT']}files/{key}",
                "hash": hash(path)
            })
    return ret

def bundle():
    shutil.make_archive("bundle", "zip", "./shadows")
    h = hash("./bundle.zip")
    n = f"bundle.{h}.zip"
    os.rename("./bundle.zip", f"./{n}")
    return n

def main():
    print("mkshadows 0.2")
    
    ret = {
        "type": "files",
        "bundle": None,
        # "bundle": f"{os.environ['ENDPOINT']}{bundle()}",
        "sanity": [
            {
                "path": "mods/",
                "rule": "provisioned",
            }
        ],
        "files": files()
    }

    with open("./shadow_manifest.json", "w") as f:
        json.dump(ret, f, indent=2)
    


if __name__ == "__main__":
    main()
