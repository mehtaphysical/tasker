#!/usr/bin/python

import os

print_file_path = os.environ.get('PRINT_FILE_PATH', 'DEFAULT_TEXT')

print(print_file_path)
