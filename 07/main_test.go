package main

import (
	"testing"
)

func TestDetectHandType(t *testing.T) {
	tests := []struct {
		name string
		hand *Hand
		want HandType
	}{
		{hand: NewHand("AAAAJ", 0), want: fiveOfAKind},
		{hand: NewHand("AAAJJ", 0), want: fiveOfAKind},
		{hand: NewHand("AJJJJ", 0), want: fiveOfAKind},
		{hand: NewHand("AJJJJ", 0), want: fiveOfAKind},

		{hand: NewHand("AAAA1", 0), want: fourOfAKind},
		{hand: NewHand("AAAJ1", 0), want: fourOfAKind},
		{hand: NewHand("AAJJ1", 0), want: fourOfAKind},
		{hand: NewHand("AJJJ1", 0), want: fourOfAKind},

		{hand: NewHand("AAA22", 0), want: fullHouse},
		{hand: NewHand("AAJ22", 0), want: fullHouse},
		{hand: NewHand("AJA2J", 0), want: fourOfAKind},
		{hand: NewHand("AAA22", 0), want: fullHouse},

		{hand: NewHand("AAA12", 0), want: threeOfAKind},
		{hand: NewHand("AAJ12", 0), want: threeOfAKind},

		{hand: NewHand("AA113", 0), want: twoPair},
		{hand: NewHand("AA1JJ", 0), want: fourOfAKind},
		{hand: NewHand("AA1J3", 0), want: threeOfAKind},

		{hand: NewHand("12344", 0), want: onePair},
		{hand: NewHand("1234J", 0), want: onePair},

		{hand: NewHand("AT938", 0), want: highCard},
		{hand: NewHand("47QT6", 0), want: highCard},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectHandType(tt.hand); got != tt.want {
				t.Errorf("detectHandType(%s) = %v, want %v", tt.hand.cards, got, tt.want)
			}
		})
	}
}

func TestIsFourOfAKind(t *testing.T) {
	tests := []struct {
		name string
		hand *Hand
		want bool
	}{
		{hand: NewHand("AAAA1", 0), want: true},
		{hand: NewHand("AAAJ1", 0), want: true},
		{hand: NewHand("AAJJ1", 0), want: true},
		{hand: NewHand("AJJJ1", 0), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFourOfAKind(tt.hand); got != tt.want {
				t.Errorf("isFourOfAKind(%s) = %v, want %v", tt.hand.cards, got, tt.want)
			}
		})
	}
}

func TestIsFullHouse(t *testing.T) {
	tests := []struct {
		name string
		hand *Hand
		want bool
	}{
		{hand: NewHand("AA111", 0), want: true},
		{hand: NewHand("AAJ11", 0), want: true},
		{hand: NewHand("JA111", 0), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFullHouse(tt.hand); got != tt.want {
				t.Errorf("isFullHouse(%s) = %v, want %v", tt.hand.cards, got, tt.want)
			}
		})
	}
}

func TestIsThreeOfAKind(t *testing.T) {
	tests := []struct {
		name string
		hand *Hand
		want bool
	}{
		{hand: NewHand("AA111", 0), want: true},
		{hand: NewHand("AAJ11", 0), want: true},
		{hand: NewHand("1224J", 0), want: true},
		{hand: NewHand("123JJ", 0), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isThreeOfAKind(tt.hand); got != tt.want {
				t.Errorf("isThreeOfAKind(%s) = %v, want %v", tt.hand.cards, got, tt.want)
			}
		})
	}
}
