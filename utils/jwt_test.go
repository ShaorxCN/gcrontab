package utils

import "testing"

func TestJwt(t *testing.T) {

	c := new(Claims)
	c.DeadLine = "2024-06-19 09:47:02"
	c.Exp = "2024-06-19 09:38:02"
	c.NickName = "evan"
	c.Role = "admin"
	c.UID = "4a636db4-e4aa-466e-8bfe-a8fba8c4c72a"
	c.Secret = "SSmHSj80ir"
	token, err := GenToken(c)
	if err != nil {
		t.Log(err)
	}
	tokenwant := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZWFkTGluZSI6IjIwMjQtMDYtMTkgMDk6NDc6MDIiLCJleHAiOiIyMDI0LTA2LTE5IDA5OjM4OjAyIiwibmlja05hbWUiOiJldmFuIiwicm9sZSI6ImFkbWluIiwidWlkIjoiNGE2MzZkYjQtZTRhYS00NjZlLThiZmUtYThmYmE4YzRjNzJhIn0.bh6RSqEOy4aTLz4SCIDiMBMXEwEhBn6c2HXcLjASit4"
	if !(token == tokenwant) {
		t.Log("gen error")
	}
	_, err = ValideToken(token, "SSmHSj80ir")

	if err != nil {
		t.Log(err)
	}
}
