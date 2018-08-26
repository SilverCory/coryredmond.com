package handlers

import "testing"

func TestGetPostURL(t *testing.T) {
	postUrl := GetPostURL("6dykW3VqLje", "Hello World")
	if postUrl != "Hello-World-6dykW3VqLje" {
		t.Errorf("Invalid Post URL generated! %q\n", postUrl)
	}
}

func TestGetPostIDFromURL(t *testing.T) {
	if resp := GetPostIDFromURL("Hello-World-6dykW3VqLje"); resp != "6dykW3VqLje" {
		t.Errorf("Invalid ID received from URL! %q\n", resp)
	}
}

func TestGeneratePostID(t *testing.T) {
	id, idStr, err := GeneratePostID()
	if err != nil {
		t.Errorf("An error occurred generating an ID! %q\n", err)
	} else if id == 0 {
		t.Error("Generated ID is 0!")
	} else if idStr == "" {
		t.Error("There is a nil idStr generated!")
	}

	t.Logf("%d => %q", id, idStr)
}

func TestDecodeID(t *testing.T) {
	id, err := DecodeID("6dykW3VqLje")
	if err != nil {
		t.Errorf("An error occurred decoding the ID! %q\n", err)
	} else if id == 0 {
		t.Error("Decoded ID is 0!")
	}

	t.Logf("\"6dykW3VqLje\" => %d", id)
}
