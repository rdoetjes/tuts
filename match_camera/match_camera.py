#!/usr/bin/env python3

import cv2
import numpy as np
import glob
import sys

def extract_noise_residual(img):
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY).astype(np.float32)

    # Better than simple blur: use median + Gaussian combo
    denoised = cv2.GaussianBlur(gray, (3, 3), 0)
    denoised = cv2.medianBlur(denoised.astype(np.uint8), 3).astype(np.float32)

    noise = gray - denoised

    # Normalize
    noise = (noise - np.mean(noise)) / (np.std(noise) + 1e-8)
    return noise

def load_images(folder):
    paths = glob.glob(folder + "/*.jpg")
    images = [cv2.imread(p) for p in paths if cv2.imread(p) is not None]
    return images

def build_fingerprint(images):
    noises = []

    for img in images:
        noise = extract_noise_residual(img)
        noises.append(noise)

    # Resize all to smallest shape
    h = min(n.shape[0] for n in noises)
    w = min(n.shape[1] for n in noises)

    noises = [n[:h, :w] for n in noises]

    # Stack and average
    fingerprint = np.mean(noises, axis=0)

    # Normalize fingerprint
    fingerprint = (fingerprint - np.mean(fingerprint)) / (np.std(fingerprint) + 1e-8)

    return fingerprint

def correlation(a, b):
    return np.corrcoef(a.flatten(), b.flatten())[0, 1]

def evaluate(fingerprint, test_images):
    scores = []

    for img in test_images:
        noise = extract_noise_residual(img)

        h = min(fingerprint.shape[0], noise.shape[0])
        w = min(fingerprint.shape[1], noise.shape[1])

        f = fingerprint[:h, :w]
        n = noise[:h, :w]

        score = correlation(f, n)
        scores.append(score)

    return scores

def normalize_scores(scores, min_ref=-0.01, max_ref=0.05):
    # Map correlation range to 0–100%
    norm = []
    for s in scores:
        pct = (s - min_ref) / (max_ref - min_ref)
        pct = max(0, min(1, pct))
        norm.append(pct * 100)
    return norm

# -------------------------
# Example usage
# -------------------------

train_images = load_images("cameraA_train")   # ~30–50 images
test_same   = load_images("cameraA_test")     # same camera
test_other  = load_images("cameraB_test")     # different camera

if len(train_images) == 0 or len(test_same) == 0 or len(test_other) == 0:
    print("no training images found")
    sys.exit(1)
fingerprint = build_fingerprint(train_images)

scores_same = evaluate(fingerprint, test_same)
scores_other = evaluate(fingerprint, test_other)

pct_same = normalize_scores(scores_same)
pct_other = normalize_scores(scores_other)

print("Same camera scores:", scores_same)
print("Same camera match %:", pct_same)

print("Other camera scores:", scores_other)
print("Other camera match %:", pct_other)

print("\nAverage same-camera match: {:.2f}%".format(np.mean(pct_same)))
print("Average different-camera match: {:.2f}%".format(np.mean(pct_other)))
