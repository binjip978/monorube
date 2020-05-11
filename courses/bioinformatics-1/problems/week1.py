from collections import defaultdict

def pattern_count(text, pattern):
    cnt = 0
    for i in range(len(text) - len(pattern) + 1):
        if text[i:i+len(pattern)] == pattern:
            cnt += 1
    
    return cnt

def freq_words(text, kmer_size):
    """Input: A string Text and an integer k.
       Output: All most frequent k-mers in Text."""
    kmers = defaultdict(int)    
    for i in range(len(text) - kmer_size + 1):
        kmers[text[i:i+kmer_size]] += 1
    
    max_value = 0
    for _, v in kmers.items():
        if v > max_value:
            max_value = v
    
    res = []
    for k, v in kmers.items():
        if v == max_value:
            res.append(k)
    
    return res

def read_lines(file):
    with open(file) as f:
        raw = f.readlines()
    lines = []
    for r in raw:
        lines.append(r.strip())

    return lines

if __name__ == '__main__':
    print(pattern_count('GACCATCAAAACTGATAAACTACTTAAAAATCAGT', 'AAA'))
    print(freq_words('CGCCTAAATAGCCTCGCGGAGCCTTATGTCATACTCGTCCT', 3))