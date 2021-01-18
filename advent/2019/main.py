"""advent of code 2019"""

def readInts(filename):
    with open(filename) as f:
        return [int(x) for x in f.readlines()]

def fuel_req(xs):
    return sum([x // 3 - 2 for x in xs])

if __name__ == '__main__':
    xs = readInts("input/1.txt")
    print(fuel_req(xs))

