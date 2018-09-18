package util

import "testing"

func TestGetPostURL(t *testing.T) {
	postUrl := GetPostURL("6dykW3VqLje", "Hello World")
	if postUrl != "Hello-World-6dykW3VqLje" {
		t.Errorf("Invalid Post URL generated! %q\n", postUrl)
	}
}

func TestGetPostURL2(t *testing.T) {
	postUrl := GetPostURL("6dykW3VqLje", "Hello World This Is Too Long For Life And This Function")
	if postUrl != "Hello-World-This-Is-Too-Long-For-Life-And-6dykW3VqLje" {
		t.Errorf("Invalid Post URL generated! %q\n", postUrl)
	}
}

func TestGetPostIDFromURL(t *testing.T) {
	if resp := GetPostIDFromURL("Hello-World-6dykW3VqLje"); resp != "6dykW3VqLje" {
		t.Errorf("Invalid PostID received from URL! %q\n", resp)
	}

	if resp := GetPostIDFromURL("6dykW3VqLje"); resp != "6dykW3VqLje" {
		t.Errorf("Invalid PostID received from URL! %q\n", resp)
	}
}

func TestEncodeDecodeID(t *testing.T) {

	var id uint64
	var idStr string
	var err error

	t.Run("TestGeneratePostID", func(t *testing.T) {
		id, idStr, err = GeneratePostID()
		if err != nil {
			t.Errorf("An error occurred generating an PostID! %q\n", err)
		} else if id == 0 {
			t.Error("Generated PostID is 0!")
		} else if idStr == "" {
			t.Error("There is a nil idStr generated!")
		}

		t.Logf("%d => %q", id, idStr)
	})

	if err != nil {
		t.Error(err)
	}

	t.Run("TestDecodeID", func(t *testing.T) {
		decID, err := DecodeID(idStr)
		if err != nil {
			t.Errorf("An error occurred decoding the PostID! %q\n", err)
		} else if decID == 0 {
			t.Error("Decoded PostID is 0!")
		} else if decID != id {
			t.Errorf("PostID was not decoded properly! (enc)%d != (dec)%d", id, decID)
		}

		t.Logf("%q => %d", idStr, id)
	})
}
