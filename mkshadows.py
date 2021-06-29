#!/usr/bin/python3

import os
import hashlib
import json
from dotenv import load_dotenv

load_dotenv()

def main():
    print("mkshadows 0.1")
    
    ret = {}

    for folder, _, files in os.walk("./shadows"):
        for i in files:
            key = f"{folder.replace('./shadows/', '')}/{i}"
            path = f"{folder}/{i}"
            sha1 = hashlib.sha1()
            with open(path, "rb") as f:
                while True:
                    chunk = f.read(sha1.block_size)
                    if not chunk:
                        break
                    sha1.update(chunk)
            ret[key] = {
                "url": f"{os.environ['FILE_ENDPOINT']}{key}",
                "hash": sha1.hexdigest()
            }
    print(json.dumps(ret))

if __name__ == "__main__":
    main()
