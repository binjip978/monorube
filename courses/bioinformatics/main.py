from collections import Counter

def pattern_pos(text, pattern):
    pos = []
    for i in range(len(text) - len(pattern) + 1):
        if pattern == text[i:i + len(pattern)]:
            pos.append(i)
    
    return pos

def complement(text):
    rev = ''
    for v in reversed(text):
        if v == 'A':
            rev += 'T'
        elif v == 'T':
            rev += 'A'
        elif v == 'C':
            rev += 'G'
        else:
            rev += 'C'
    return rev

def pattern_count(text, pattern):
    cnt = 0
    for i in range(len(text) - len(pattern) + 1):
        if text[i:i+len(pattern)] == pattern:
            cnt += 1
    
    return cnt

def freq_word(text, k):
    cnt = Counter()
    for i in range(len(text) - k + 1):
        cnt[text[i:i+k]] += 1
    
    max_l = 0
    for _, v in cnt.items():
        if v > max_l:
            max_l = v
    
    res = []
    for k, v in cnt.items():
        if v == max_l:
            res.append(k)

    return res

def clump(text, k, clump, times):
    res = set()
    freq = Counter()
    for i in range(len(text) - k + 1):
        if i + 1 >= clump:
            s = text[i + 1 - clump : i + 1 - clump + k]
            freq[s] -= 1
        s = text[i : i + k]
        freq[s] += 1
        if freq[s] >= times:
            res.add(s)

    return res

def skew(text):
    skew_pos = [0]
    for i, e in enumerate(text, start=1):
        if e == 'C':
            skew_pos.append(skew_pos[i - 1] - 1)
        elif e == 'G':
            skew_pos.append(skew_pos[i - 1] + 1)
        else:
            skew_pos.append(skew_pos[i - 1])

    res = []
    m = min(skew_pos)
    for i, e in enumerate(skew_pos):
        if e == m:
            res.append(i)

    return res

def hamming_distance(s1, s2):
    d = 0
    for e1, e2 in zip(s1, s2):
        if e1 != e2:
            d += 1
    
    return d

def fuzzy_pattern_count(pattern, text, distance):
    res = []
    for i in range(len(text) - len(pattern) + 1):
        if hamming_distance(text[i:i+len(pattern)], pattern) <= distance:
            res.append(i)

    return res

if __name__ == '__main__':
    pass
