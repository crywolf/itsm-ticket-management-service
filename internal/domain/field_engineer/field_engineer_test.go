package fieldengineer_test

import (
	"testing"

	. "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	tsession "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/time_session"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestInit initializes test suite
func TestInit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Field Engineer tests")
}

var _ = Describe("Field Engineer behavior", func() {
	var basicUser user.BasicUser
	var fieldEngineer FieldEngineer
	var actorUser actor.Actor

	BeforeEach(func() {
		basicUser = user.BasicUser{
			ExternalUserUUID: "3d334abe-f289-42a5-9742-72c3133768c2",
			Name:             "Test",
			Surname:          "User",
			OrgDisplayName:   "Some Company",
			OrgName:          "897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
		}

		fieldEngineer = FieldEngineer{
			BasicUser: basicUser,
		}
		err := fieldEngineer.SetUUID("c546d4bb-2f45-411a-8583-9d0e6fe4807a")
		Expect(err).To(BeNil())

		actorUser = actor.Actor{
			BasicUser: basicUser,
		}
	})

	Describe("StartWorking()", func() {
		incUUID := ref.UUID("3c032e34-b1a2-43a9-b1d2-eb3b241b4a78")
		inc := incident.Incident{Number: "INC123"}
		err := inc.SetUUID(incUUID)
		Expect(err).To(BeNil())

		When("called by actor that is not field engineer", func() {
			It("should return error", func() {
				err := fieldEngineer.StartWorking(actorUser, inc)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("actor is not field engineer"))
			})
		})

		When("called by some other field engineer actor", func() {
			BeforeEach(func() {
				feUUID := ref.UUID("b8d49f19-5e54-44cf-b547-f16bacb69294")
				actorUser.SetFieldEngineerID(&feUUID)
			})

			It("should return error", func() {
				err := fieldEngineer.StartWorking(actorUser, inc)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("actor is not this field engineer"))
			})
		})

		When("called by this field engineer actor", func() {
			BeforeEach(func() {
				feUUID := fieldEngineer.UUID()
				actorUser.SetFieldEngineerID(&feUUID)
			})

			Context("and the field engineer has not an open time session", func() {
				JustBeforeEach(func() {
					Expect(fieldEngineer.HasOpenTimeSession()).To(BeFalse())
					err := fieldEngineer.StartWorking(actorUser, inc)
					Expect(err).To(BeNil())
				})

				It("should open new time session", func() {
					Expect(fieldEngineer.HasOpenTimeSession()).To(BeTrue())
				})

				Describe("open time session", func() {
					var ts *tsession.TimeSession

					JustBeforeEach(func() {
						ts = fieldEngineer.OpenTimeSession()
						Expect(ts).NotTo(BeNil())
					})

					It("should contain added incident", func() {
						Expect(ts.Incidents).To(HaveLen(1))
					})

					It("should be in Work state", func() {
						Expect(ts.State()).To(Equal(tsession.StateWork))
					})
				})
			})

			Context("and the field engineer has an open time session", func() {
				alreadyAddedIncUUID := ref.UUID("d67c7799-cab5-4dbd-8a5c-2e4e19070f77")
				var ts *tsession.TimeSession

				BeforeEach(func() {
					// set open TS with one incident to the FE
					ts = &tsession.TimeSession{
						Incidents: []tsession.IncidentInfo{{
							IncidentID:         alreadyAddedIncUUID,
							HasSupplierProduct: true,
						}},
					}
					err := ts.SetState(tsession.StateTravel)
					Expect(err).To(BeNil())

					Expect(ts.Incidents).To(HaveLen(1))
					Expect(ts.Incidents[0].IncidentID).To(Equal(alreadyAddedIncUUID))
					Expect(ts.Incidents[0].HasSupplierProduct).To(BeTrue())

					fieldEngineer.SetOpenTimeSession(ts)
				})

				When("time session is in Travel or Work state", func() {
					JustBeforeEach(func() {
						Expect(fieldEngineer.HasOpenTimeSession()).To(BeTrue())
						err := fieldEngineer.StartWorking(actorUser, inc)
						Expect(err).To(BeNil())
					})

					It("should add incident to already opened time session", func() {
						ts := fieldEngineer.OpenTimeSession()
						Expect(ts).NotTo(BeNil())

						Expect(ts.State()).To(Equal(tsession.StateWork))

						Expect(ts.Incidents).To(HaveLen(2))
						Expect(ts.Incidents[0].IncidentID).To(Equal(alreadyAddedIncUUID))
						Expect(ts.Incidents[0].HasSupplierProduct).To(BeTrue())

						Expect(ts.Incidents[1].IncidentID).To(Equal(incUUID))
						//Expect(ts.Incidents[1].HasSupplierProduct).To(BeFalse()) // not implemented yet
					})
				})

				When("time session is in TravelBack state", func() {
					JustBeforeEach(func() {
						err := ts.SetState(tsession.StateTravelBack)
						Expect(err).To(BeNil())
					})

					It("should return error", func() {
						Expect(fieldEngineer.HasOpenTimeSession()).To(BeTrue())
						err := fieldEngineer.StartWorking(actorUser, inc)
						Expect(err).NotTo(BeNil())

						Expect(err.Error()).To(Equal("time session is not in Travel nor Work state"))
					})

					When("time session is in Brake state", func() {
						JustBeforeEach(func() {
							err := ts.SetState(tsession.StateBreak)
							Expect(err).To(BeNil())
						})

						It("should return error", func() {
							Expect(fieldEngineer.HasOpenTimeSession()).To(BeTrue())
							err := fieldEngineer.StartWorking(actorUser, inc)
							Expect(err).NotTo(BeNil())

							Expect(err.Error()).To(Equal("time session is not in Travel nor Work state"))
						})
					})

					When("time session is in Closed state", func() {
						JustBeforeEach(func() {
							err := ts.SetState(tsession.StateClosed)
							Expect(err).To(BeNil())
						})

						It("should return error", func() {
							Expect(fieldEngineer.HasOpenTimeSession()).To(BeTrue())
							err := fieldEngineer.StartWorking(actorUser, inc)
							Expect(err).NotTo(BeNil())

							Expect(err.Error()).To(Equal("time session is not in Travel nor Work state"))
						})
					})
				})
			})
		})
	})
})
