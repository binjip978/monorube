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

class FuzzSeq:
    def __init__(self):
        self._cache = {}
    
    def mutate(self, pattern, d):
        if (pattern, d) in self._cache:
            return self._cache[(pattern, d)]
        
        res = self._mutate1({pattern}, d)
        self._cache[(pattern, d)] = res

        return res
        
    def _mutate1(self, xs, d):
        if d == 0:
            return xs
        s = set()
        for w in xs:
            for i in range(len(w)):
                if w[i] == 'A':
                    s.add(w[0:i] + 'C' + w[i+1:])
                    s.add(w[0:i] + 'T' + w[i+1:])
                    s.add(w[0:i] + 'G' + w[i+1:])
                if w[i] == 'C':
                    s.add(w[0:i] + 'A' + w[i+1:])
                    s.add(w[0:i] + 'T' + w[i+1:])
                    s.add(w[0:i] + 'G' + w[i+1:])
                if w[i] == 'T':
                    s.add(w[0:i] + 'A' + w[i+1:])
                    s.add(w[0:i] + 'C' + w[i+1:])
                    s.add(w[0:i] + 'G' + w[i+1:])
                if w[i] == 'G':
                    s.add(w[0:i] + 'T' + w[i+1:])
                    s.add(w[0:i] + 'C' + w[i+1:])
                    s.add(w[0:i] + 'A' + w[i+1:])

        xs.update(s)
        return self._mutate1(xs, d - 1)

def fuzzy_freq_words(text, k_size, distance):
    f = FuzzSeq()
    words = Counter()
    for i in range(len(text) - k_size + 1):
        s = text[i:i+k_size]
        mut = f.mutate(s, distance)
        for m in mut:
            words[m] += 1
    
    m = max(words.values())
    res = []
    for k, v in words.items():
        if v == m:
            res.append(k)
    
    return res

def motif_enumeration(xs, k, d):
    f = FuzzSeq()
    res = set()
    for seq in xs:
        for i in range(len(seq) - k + 1):
            cands = f.mutate(seq[i:i+k], d)
            for cand in cands:
                miss = False
                for seq2 in xs:
                    r = fuzzy_pattern_count(cand, seq2, d)
                    if len(r) == 0:
                        miss = True
                        break
                if not miss:
                    res.add(cand)
    
    return res
                
                
if __name__ == '__main__':
    lines = None
    with open('input/rosalind_ba2a.txt', 'r') as f:
        lines = f.readlines()
    res = motif_enumeration(lines[1:], 5, 1)
    print(" ".join(res))