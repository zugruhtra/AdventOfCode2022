import sys
import json
from pprint import pprint


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


def get_pairs_in_right_order(pairs):
    pairs_in_right_order = []
    for idx, pair in enumerate(pairs):
        result = compare_pair(*pair)
        if result:
            pairs_in_right_order.append(idx)
    return pairs_in_right_order


def main():
    raw = [line.strip() for line in sys.stdin]
    data = [json.loads(line) for line in raw if line]
    pairs = [(data[i],data[i+1]) for i in range(0, len(data), 2)]
    pairs_in_right_order = get_pairs_in_right_order(pairs)
    print(
        'Sum of indices of pairs in right order:', 
        sum(pairs_in_right_order)+len(pairs_in_right_order)
    )

if __name__ == '__main__':
    main()
