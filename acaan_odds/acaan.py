#!/usr/bin/env python3
import random

cards = ["A", "K", "Q", "J", "10", "9", "8", "7", "6", "5", "4", "3", "2"]
suits = ["♠", "♥", "♦", "♣"]
trials = 1000000
hits = 0

def create_deck():
    deck = []
    for card in cards:
        for suit in suits:
            deck.append(card + suit)
    return deck

for _ in range(trials):
    deck = create_deck()
    random.shuffle(deck)
    
    # Randomly choose a card and a position
    target_card = random.choice(deck)
    target_position = random.randint(0, len(deck) - 1)
    
    # Check if the card at the chosen position matches the target card
    if deck[target_position] == target_card:
        hits += 1

# Calculate and print the probability
probability = trials / hits if hits > 0 else 0
print(f"Probability: 1 in {probability:.2f}")
