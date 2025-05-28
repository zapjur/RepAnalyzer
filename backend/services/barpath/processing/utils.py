import os
from urllib.parse import urlparse

def extract_filename(url):
    path = urlparse(url).path
    base = os.path.basename(path)
    filename, _ = os.path.splitext(base)
    return filename
