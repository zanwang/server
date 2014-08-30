package tests

import (
	"fmt"
	"net/http"

	. "github.com/onsi/gomega"

	"github.com/tommy351/maji.moe/models"
)

func (s *TestSuite) Activation() {
	s.Describe("Activation", func() {
		s.Before(func() {
			s.createUser1()
		})

		s.It("Success", func() {
			s.setUserActivated("user", false)
			user := s.Get("user").(*models.User)
			models.DB.SelectOne(user, "SELECT * FROM users WHERE id=?", user.ID)

			url := fmt.Sprintf("/activation?id=%d&token=%s", user.ID, user.ActivationToken)
			r := s.Request("GET", url, nil)

			Expect(r.Code).To(Equal(http.StatusFound))
			Expect(r.Header().Get("Location")).To(Equal("/app"))

			user.Activated = true
		})

		s.It("User does not exist", func() {
			url := fmt.Sprintf("/activation?id=%d&token=%s", 99999999, "abcdefg")
			r := s.Request("GET", url, nil)

			Expect(r.Code).To(Equal(http.StatusNotFound))
		})

		s.It("Wrong token", func() {
			s.setUserActivated("user", false)
			user := s.Get("user").(*models.User)
			url := fmt.Sprintf("/activation?id=%d&token=%s", user.ID, "abcdefg")
			r := s.Request("GET", url, nil)

			Expect(r.Code).To(Equal(http.StatusBadRequest))
		})

		s.It("User has been activated", func() {
			s.setUserActivated("user", true)
			user := s.Get("user").(*models.User)
			models.DB.SelectOne(user, "SELECT * FROM users WHERE id=?", user.ID)

			url := fmt.Sprintf("/activation?id=%d&token=%s", user.ID, user.ActivationToken)
			r := s.Request("GET", url, nil)

			Expect(r.Code).To(Equal(http.StatusFound))
			Expect(r.Header().Get("Location")).To(Equal("/app"))
		})

		s.After(func() {
			s.deleteUser1()
		})
	})
}
