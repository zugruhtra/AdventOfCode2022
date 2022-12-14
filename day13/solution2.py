import sys
import json

from functools import cmp_to_key


def compare_pair(left, right):
    if isinstance(left, int) and isinstance(right, int):
        result = compare_integers(left, right)
    elif isinstance(left, list) and isinstance(right, list):
        result = compare_lists(left, right)
    elif isinstance(left, int):
        result = compare_lists([left], right)
    elif isinstance(right, int):
        result = compare_lists(left, [right])
    else:
        assert False, 'unreachable line'
    if result is not None:
        return result


def compare_integers(left, right):
    if left < right:
        return True
    elif left > right:
        return False
    else:
        return None


def compare_lists(left, right):
    pair = left, right
    for left, right in zip(*pair):
        result = compare_pair(left, right)
        if result is not None:
            return result
    d = len(pair[0]) - len(pair[1])
    if d == 0:
        return None
    else:
        return d < 0


def compare_packets(left, right):
    result = compare_pair(left, right)
    if result is None:
        return 0
    else:
        return -1 if result else 1

def main():
    raw = [line.strip() for line in sys.stdin]
    data = [json.loads(line) for line in raw if line]
    data.append([[2]])
    data.append([[6]])
    data.sort(key=cmp_to_key(compare_packets))
    idx_of_divider_packets = [
        idx+1
        for idx, packet in enumerate(data)
        if isinstance(packet, list) 
        and len(packet) == 1 
        and isinstance(packet[0], list) 
        and len(packet[0]) == 1 
        and packet[0][0] in (2,6)
    ]
    decoder_key = 1
    for idx in idx_of_divider_packets:
        decoder_key *= idx
    print('Decoder key:', decoder_key)


if __name__ == '__main__':
    main()
