Short summary (one sentence)
- Because each ciphertext byte is produced from one key byte in a deterministic way (XOR + a small rotation whose amount depends only on the key byte and the byte index), every known plaintext byte gives you a quick test you can run for all 256 possible key byte values; doing that for many positions and many ciphertexts narrows the possibilities until the actual key bytes are found.

Why this is easy (intuition)
- Your transform per byte looks like: rotate-left( plaintext XOR key_byte, shift ) -> ciphertext, where shift = (key_byte + index) % 8.
- For a given ciphertext byte and a guessed key byte, you can compute exactly what the plaintext would have to be. So you can test each possible key byte (0..255) and immediately rule out those that don’t match the known plaintext byte.
- That’s only 256 checks per byte — trivial for a computer. If one known plaintext byte gives several candidate key bytes, more known bytes (or other ciphertexts) typically rule down to a single candidate.

Concrete step-by-step method you can follow
1. Get a pair of matching data:
   - a ciphertext byte sequence (binary) and
   - a known plaintext fragment (a “crib”) that you know appears at a known offset inside the plaintext.
   - In the repo we produced these with `make_cipher` and fed them to `recover_key`.
2. For the ciphertext message of length L, compute the per-message start index:
   - start = L % key_len
   - This matches the way your `encrypt`/`decrypt` choose which key byte maps to message index 0.
3. For each byte of the crib at message index i (0-based from message start):
   - Compute the message-level index: msg_index = offset + i.
   - Compute the corresponding key index: kidx = (start + msg_index) % key_len.
   - For every possible k in 0..255:
     - Compute shift = (k + msg_index) % 8.
     - Undo the rotate on the ciphertext byte with that shift (rotate-right).
     - XOR with k. If the result equals the known plaintext byte, then k is a candidate for key[kidx].
4. Record all candidates for each key index. If multiple crib bytes map to the same key index, intersect candidate sets (i.e., only keep values that satisfy every mapped crib byte).
5. Repeat with additional cribs or other ciphertexts (same key but different messages/offsets) — each additional known byte further reduces the candidate sets.
6. When many key indices have single candidates, assemble the best‑effort key and decrypt the ciphertext to verify.
7. For remaining ambiguous positions:
   - Use heuristics (prefer printable ASCII, prefer alphanumeric characters, etc.)
   - Or run a small combinatorial search (try combinations of few ambiguous positions) and score the resulting plaintext for English word-likeness or printable content.
   - Or ask interactively to pick candidates while previewing the decrypted text.

Why this isn’t “magic” — what makes it work
- Deterministic mapping: one ciphertext byte points to one key byte (given the start offset formula).
- Low per-byte cost: guessing 256 possibilities per byte is cheap.
- Intersection: multiple constraints collapse many possibilities quickly.

What can fail or slow you down
- Wrong key length: the mapping kidx = (start + i) % key_len depends on the right key length. If you try the wrong key length the candidates won’t intersect usefully. That’s why the tool can try many key lengths.
- Short crib(s): a tiny crib gives each key index few constraints, leaving many candidates.
- Key bytes that are non-printable or happen to produce multiple valid candidates: heuristics (printable preference) can pick the wrong one.
- If the attacker has no known plaintext or no ability to choose plaintexts, the attack is harder but still often possible with many ciphertexts (crib-dragging, frequency/statistical analysis).

Why this proves the scheme is weak
- Key reuse: the same key bytes are used across messages (only offset shifts deterministically). That means ciphertexts leak relationships (C1 xor C2 = P1 xor P2) and known plaintexts expose key bytes.
- Deterministic rotation: the rotation depends on k and index, but that still leaves a small finite set per key byte you can brute force quickly.
- In short: it’s not a one-time random keystream; it reuses key material and so is vulnerable.

Tools we used so you can replicate
- `recover_key` — brute-forces per-byte candidates from crib(s) and prints per-position candidate lists and a best-effort key.
- `auto_resolve.py` — automatically prefers printable candidates and produces an auto-resolved key then decrypts.
- `tools.py` — helpers to show a recovered key (`show-key`) and decrypt with a recovered key.

If you want an even simpler analogy
- Imagine each key byte is a lock. Each known plaintext/ciphertext pair gives you a set of possible keys that could open that lock (often a small set). Collect many locks (positions) and many clues (cribs) and you narrow down the secret key for each lock until you’ve opened them all.
