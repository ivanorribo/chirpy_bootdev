package main

import (
	"strings"
)

func cleanWords(text string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	text_split := strings.Split(text, " ")
	for i := 0; i < len(text_split); i++ {
		for _, item := range badWords {
			if strings.ToLower(text_split[i]) == item {
				text_split[i] = "****"
			}
		}
	}
	text = strings.Join(text_split, " ")
	return text
}
