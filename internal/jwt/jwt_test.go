package jwt_test

import (
	"testing"

	"github.com/ishmulyan/slotssvc/internal/jwt"
)

func TestServiceEncodeDecodeSameSecret(t *testing.T) {
	svc := jwt.NewService([]byte("boo"))
	expected := jwt.Claims{}
	tokenString, err := svc.Encode(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := svc.Decode(tokenString)
	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Fatalf("expect %+v, got actual %+v", expected, actual)
	}
}

func TestServiceEncodeDecodeAnotherSecret(t *testing.T) {
	svc := jwt.NewService([]byte("boo"))
	expected := jwt.Claims{}
	tokenString, err := svc.Encode(expected)
	if err != nil {
		t.Fatal(err)
	}

	svc = jwt.NewService([]byte("boop"))
	if _, err := svc.Decode(tokenString); err == nil {
		t.Fatal("decoding should fail because signature isn't the same")
	}
}
