#!/usr/bin/env python
# coding=utf-8
__author__ = 'chenfengyuan'
import csv
import sys


def main(f1, f2, out):
    d = csv.excel()
    d.lineterminator = '\n'
    writer = csv.writer(out, d)
    while True:
        try:
            l1 = next(f1)[:-2].decode('gb18030', 'ignore')
            l2 = next(f2)[:-2].decode('gb18030', 'ignore')
        except StopIteration:
            break
        assert len(l1) == len(l2)
        base = 0
        row = []
        for i in range(len(l1)):
            if l1[i] != l2[i]:
                row.append(l1[base:i])
                base = i + 1
        row.append(l1[base:])
        writer.writerow(row)


if __name__ == '__main__':
    main(open(sys.argv[1], 'rb'), open(sys.argv[2], 'rb'), sys.stdout)