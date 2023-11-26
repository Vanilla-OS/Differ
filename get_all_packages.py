#!/usr/bin/env python3

import json
import re
import subprocess

out = subprocess.check_output(["dpkg", "-l", "--no-pager"])
regex = re.compile(r"\s+")
packages = []
for line in out.decode("utf-8").split("\n"):
    if line.startswith("ii"):
        items = regex.sub(" ", line).split(" ")
        packages.append({
            "name": items[1],
            "version": items[2],
        })

print(json.dumps({"packages": packages}))
