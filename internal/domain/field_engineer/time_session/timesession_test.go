package tsession_test

import (
	"testing"

	. "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/time_session"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestInit initializes test suite
func TestInit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Time Session tests")
}

var _ = Describe("Time Session behavior", func() {
	Describe("AddIncident()", func() {
		var inc incident.Incident
		var ts TimeSession

		BeforeEach(func() {
			inc = incident.Incident{}
			err := inc.SetUUID("40018f49-e7dd-4afa-86f5-021b44ad33ad")
			Expect(err).To(BeNil())

			ts = TimeSession{}
			err = ts.SetState(StateWork)
			Expect(err).To(BeNil())
		})

		It("should add incident to the time session", func() {
			Expect(ts.Incidents).To(HaveLen(0))

			err := ts.AddIncident(inc)
			Expect(err).To(BeNil())

			Expect(ts.Incidents).To(HaveLen(1))
		})

		It("should allow adding more incidents when called multiple times", func() {
			inc2 := incident.Incident{}
			err := inc2.SetUUID("c1ddbaf5-4d10-4181-b6b0-c2a7ff714989")
			Expect(err).To(BeNil())

			Expect(ts.Incidents).To(HaveLen(0))

			err = ts.AddIncident(inc)
			Expect(err).To(BeNil())
			err = ts.AddIncident(inc2)
			Expect(err).To(BeNil())

			Expect(ts.Incidents).To(HaveLen(2))
		})

		It("should not add the same incident twice", func() {
			Expect(ts.Incidents).To(HaveLen(0))

			err := ts.AddIncident(inc)
			Expect(err).To(BeNil())
			err = ts.AddIncident(inc)
			Expect(err).To(BeNil())

			Expect(ts.Incidents).To(HaveLen(1))
		})

		When("time session is not in Work state", func() {
			BeforeEach(func() {
				err := ts.SetState(StateTravel)
				Expect(err).To(BeNil())
			})

			It("should return error", func() {
				err := ts.AddIncident(inc)
				Expect(err).NotTo(BeNil())

				Expect(err.Error()).To(Equal("time session is not in Work state"))
			})
		})
	})
})
