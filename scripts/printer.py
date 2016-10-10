#!/usr/bin/python

import os

base_data_path = os.environ.get("DATA_PATH")

print_file_path = os.environ.get('PRINT_FILE_PATH', '/writer')

with open(base_data_path + "/" + print_file_path) as f:
    print(f.read())
    f.close()
