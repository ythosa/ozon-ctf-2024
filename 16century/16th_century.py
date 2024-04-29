from collections import deque
from itertools import product
import sys
import string

ALPHABET = string.ascii_lowercase


def decrypt(key, ciphertext):
    key = key.lower()
    assert all(c in ALPHABET for c in key)

    plaintext = ''
    key_index = 0

    for c in ciphertext:
        if c.lower() in ALPHABET:
            enc_c = ALPHABET[(ALPHABET.find(c.lower()) - ALPHABET.find(key[key_index])) % len(ALPHABET)]
            key_index = (key_index + 1) % len(key)
            if c.isupper():
                enc_c = enc_c.upper()

            plaintext += enc_c
        else:
            plaintext += c
    return plaintext


def encrypt(key, plaintext):
    key = key.lower()
    assert all(c in ALPHABET for c in key)

    ciphertext = ''
    key_index = 0

    for c in plaintext:
        if c.lower() in ALPHABET:
            enc_c = ALPHABET[(ALPHABET.find(key[key_index]) + ALPHABET.find(c.lower())) % len(ALPHABET)]
            key_index = (key_index + 1) % len(key)
            if c.isupper():
                enc_c = enc_c.upper()

            ciphertext += enc_c
        else:
            ciphertext += c
    return ciphertext


def build_matrix():
    deq = deque(ALPHABET)
    result = []
    for _ in ALPHABET:
        result.append(deq)
        deq = deq.copy()
        deq.rotate(-1)
    return result


almarix = build_matrix()


def hack(input_text, output_text):
    key = ""
    length = len(input_text)

    for i in range(length):
        inch = input_text[i].lower()
        ouch = output_text[i].lower()

        if inch not in ALPHABET:
            continue

        idx = list(almarix[ch2idx(inch)]).index(ouch)
        symbol = ALPHABET[idx]
        key += symbol

    return key


def ch2idx(ch: str) -> int:
    return ord(ch.lower()) - ord('a')


def main():
    with open("./gpl-3.0.txt.enc") as out_file:
        out_text = out_file.read()
    with open("./gpl-3.0.txt") as in_file:
        in_text = in_file.read()

    # in_text = 'YatqXGP'
    # out_text = 'OzonCTF'
    key = hack(in_text, out_text)
    print(key)
    print(decrypt(key, out_text))

    print(decrypt('kbfdvnkfvybpxyeqiefenooxalzrsqzb', "YatqXGP{adefcbpi_sqtmie_sbzrjokGh_qowmpbqnmimlxbuaweqcxbnusfqmrpup}"))

    # for key_length in range(7, 40):
    #     print(f'loading with length = {key_length}')
    #     products = product(ALPHABET, repeat=key_length)
    #     for key in products:
    #         if encrypt(''.join(key), in_text) == out_text:
    #             print(key)
    #             print(key)
    #             print(key)
    #             return


if __name__ == "__main__":
    main()
