#!/usr/bin/env python2
# coding=utf-8
from __future__ import absolute_import
__author__ = 'chenfengyuan'
import requests
import re
import os
import sys
import subprocess
import time
import urllib

MAX_REDIRECT = 3


def get_filename_and_size(curl_cmd):
    url = re.findall(u"'([^']+)", curl_cmd)[0]
    headers = {k: v for k, v in re.findall(u"-H '([^:]+)\s*:\s*([^']+)'", curl_cmd)}
    for i in range(MAX_REDIRECT):
        resp = requests.get(url, headers=headers, allow_redirects=False, stream=True)
        resp_headers = resp.headers
        status_code = resp.status_code
        if status_code in {301, 302}:
            url = resp_headers[u'location']
            continue
        filename = re.findall(u'filename="(.+)"', resp_headers[u'Content-Disposition'])[0]
        filename = urllib.unquote(filename)
        filename = os.path.basename(filename).decode(u'utf-8')
        filename = re.sub(u"/|\x00", u'', filename)
        return filename, int(resp_headers[u'Content-Length'])


def encode_filename(filename):
    return u"'%s'" % re.sub(u"'", u"""'"'"'""", filename)


def main():
    with open(sys.argv[1]) as curl_cmds:
        for line in curl_cmds:
            line = line.strip()
            if not line:
                continue
            curl_prefix = line
            tmp = get_filename_and_size(line)
            raw_filename = tmp[0]
            filename = encode_filename(tmp[0])
            size = tmp[1]
            if os.path.exists(raw_filename):
                downloaded_size = os.path.getsize(raw_filename)
            else:
                downloaded_size = 0
            if size <= downloaded_size:
                print u'%s %s %s is already downloaded' % (size, downloaded_size, raw_filename)
                continue
            subprocess.call(u"echo downloading '%s'" % filename, shell=True)
            cmd = u"%s -C - -L -o '%s'" % (curl_prefix, filename)
            for i in xrange(10):
                rv = subprocess.call(cmd, shell=True)
                if rv != 0:
                    print u'retrying.....'
                    time.sleep(10)
                else:
                    break

if __name__ == u'__main__':
    main()
