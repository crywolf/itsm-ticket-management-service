package incident

import (
	"testing"

	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestInit initializes test suite
func TestInit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Incident tests")
}

var _ = Describe("Incident behavior", func() {
	var basicUser user.BasicUser
	var fieldEngineer fieldengineer.FieldEngineer
	var actorUser actor.Actor

	BeforeEach(func() {
		basicUser = user.BasicUser{
			ExternalUserUUID: "3d334abe-f289-42a5-9742-72c3133768c2",
			Name:             "Test",
			Surname:          "User",
			OrgDisplayName:   "Some Company",
			OrgName:          "897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
		}

		fieldEngineer = fieldengineer.FieldEngineer{
			BasicUser: basicUser,
		}
		err := fieldEngineer.SetUUID("c546d4bb-2f45-411a-8583-9d0e6fe4807a")
		Expect(err).To(BeNil())

		actorUser = actor.Actor{
			BasicUser: basicUser,
		}
	})

	Describe("StartWorking()", func() {
		When("called by actor that is not field engineer", func() {
			var inc Incident

			BeforeEach(func() {
				inc = Incident{}
			})

			It("should return error", func() {
				err := inc.StartWorking(actorUser, false)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("user is not field engineer, only assigned field engineer can start working"))
			})
		})

		When("called by field engineer", func() {
			BeforeEach(func() {
				actorUser.SetFieldEngineer(&fieldEngineer)
			})

			Context("but incident has no field engineer assigned", func() {
				var inc Incident

				BeforeEach(func() {
					inc = Incident{}
				})

				It("should return error", func() {
					err := inc.StartWorking(actorUser, false)
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(Equal("ticket does not have any field engineer assigned"))
				})
			})

			Context("and incident has field engineer assigned", func() {
				Context("but the assigned field engineer is different then the actor", func() {
					var inc Incident

					BeforeEach(func() {
						feUUID := ref.UUID("63fcafcb-e0ac-490b-b67c-b6f60afeccfd")
						inc = Incident{
							FieldEngineerID: &feUUID,
						}
					})

					It("should return error", func() {
						err := inc.StartWorking(actorUser, false)
						Expect(err).NotTo(BeNil())
						Expect(err.Error()).To(Equal("user is not assigned as field engineer, only assigned field engineer can start working"))
					})
				})

				Context("and the actor is assigned as field engineer", func() {
					Context("but the incident has an open timelog", func() {
						var inc Incident

						BeforeEach(func() {
							feUUID := fieldEngineer.UUID()
							inc = Incident{
								FieldEngineerID: &feUUID,
								openTimelog:     &timelog.Timelog{},
							}
						})

						It("should return error", func() {
							err := inc.StartWorking(actorUser, false)
							Expect(err).NotTo(BeNil())
							Expect(err.Error()).To(Equal("ticket already has an open timelog"))
						})
					})

					Context("and the incident has no open timelog", func() {
						var inc Incident

						BeforeEach(func() {
							feUUID := fieldEngineer.UUID()
							inc = Incident{
								FieldEngineerID: &feUUID,
							}
						})

						JustBeforeEach(func() {
							err := inc.StartWorking(actorUser, true)
							Expect(err).To(BeNil())
						})

						It("should open new timelog", func() {
							Expect(inc.HasOpenTimelog()).To(BeTrue())
						})

						It("should set the 'remote' param in the timelog", func() {
							Expect(inc.OpenTimelog().Remote).To(Equal(true))
						})

						It("should set state to 'in progress'", func() {
							Expect(inc.State()).To(Equal(StateInProgress))
						})
					})
				})
			})
		})
	})

	Describe("Cancel()", func() {
		When("incident is in 'new' state", func() {
			var inc Incident

			BeforeEach(func() {
				inc = Incident{
					state: StateNew,
				}
			})

			JustBeforeEach(func() {
				err := inc.Cancel(actorUser)
				Expect(err).To(BeNil())
			})

			It("should set state to 'cancelled'", func() {
				Expect(inc.State()).To(Equal(StateCancelled))
			})
		})

		When("incident is NOT in 'new' state", func() {
			var inc Incident

			BeforeEach(func() {
				inc = Incident{
					state: StateInProgress,
				}
			})

			It("should return error", func() {
				err := inc.Cancel(actorUser)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("ticket can be cancelled only in NEW state"))
			})
		})
	})
})
