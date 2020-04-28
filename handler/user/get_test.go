package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/testutil"

	"github.com/gin-gonic/gin"
)

// see: https://rshipp.com/go-api-integration-testing/
func TestGet(t *testing.T) {
	app := testutil.Setup()

	// Test body will be here!
	// Set up a test table.
	userTests := []model.UserModel{
		{
			ID:       12,
			Username: "user001",
			Password: "123456",
			Phone:    13810002000,
			Avatar:   "",
			Sex:      0,
		},
		{
			ID:       13,
			Username: "user002",
			Password: "123456",
			Phone:    13810002001,
			Avatar:   "",
			Sex:      0,
		},
	}

	for _, user := range userTests {
		// Create a user for us to view.
		app.DB.Create(user)

		// Set up a new request.
		req, err := http.NewRequest("GET", fmt.Sprintf("/v1/users/%d", user.ID), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// We need a mux router in order to pass in the `name` variable.
		r := gin.New()

		r.GET("/v1/users/:id", Get)
		r.ServeHTTP(rr, req)

		// Test that the status code is correct.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Status code is invalid. Expected %d. Got %d instead", http.StatusOK, status)
		}

		// Read the response body.
		data, err := ioutil.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatal(err)
		}

		// Test that the updated star is correct.
		returnedUser := model.UserModel{}
		if err := json.Unmarshal(data, &returnedUser); err != nil {
			t.Errorf("Returned user is invalid JSON. Got: %s", data)
		}
		//if returnedUser != user {
		//	t.Errorf("Returned user is invalid. Expected %+v. Got %+v instead", user, returnedUser)
		//}
	}

	testutil.Teardown(app)
}
