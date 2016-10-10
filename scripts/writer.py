#!/usr/bin/python

import os

base_data_path = os.environ.get("DATA_PATH")

text = os.environ.get("TEXT", "HI")
output = os.environ.get("OUTPUT", "writer")

with open(base_data_path + "/" + output, "w") as f:
    f.write(text)
    f.close()
