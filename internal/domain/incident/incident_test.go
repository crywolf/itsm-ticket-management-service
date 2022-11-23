package incident_test

import (
	"testing"
	"time"

	fieldengineer "github.com/crywolf/itsm-ticket-management-service/internal/domain/field_engineer"
	. "github.com/crywolf/itsm-ticket-management-service/internal/domain/incident"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/crywolf/itsm-ticket-management-service/internal/mocks"
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
	var clock *mocks.FixedClock

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

		clock = mocks.NewFixedClock()
	})

	Describe("StartWorking()", func() {
		When("called by actor that is not field engineer", func() {
			var inc Incident

			BeforeEach(func() {
				inc = Incident{}
			})

			It("should return error", func() {
				err := inc.StartWorking(actorUser, clock, false)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("user is not field engineer, only assigned field engineer can start working"))
			})
		})

		When("called by field engineer actor", func() {
			BeforeEach(func() {
				feUUID := fieldEngineer.UUID()
				actorUser.SetFieldEngineerID(&feUUID)
			})

			Context("but incident has no field engineer assigned", func() {
				var inc Incident

				BeforeEach(func() {
					inc = Incident{}
				})

				It("should return error", func() {
					err := inc.StartWorking(actorUser, clock, false)
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
						err := inc.StartWorking(actorUser, clock, false)
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
							}
							inc.SetOpenTimelog(&timelog.Timelog{})
						})

						It("should return error", func() {
							err := inc.StartWorking(actorUser, clock, false)
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
							err := inc.SetState(StateNew)
							Expect(err).To(BeNil())
						})

						JustBeforeEach(func() {
							Expect(inc.State()).To(Equal(StateNew))
							err := inc.StartWorking(actorUser, clock, true)
							Expect(err).To(BeNil())
						})

						It("should open new timelog", func() {
							Expect(inc.HasOpenTimelog()).To(BeTrue())
						})

						It("should set the Remote param in the open timelog", func() {
							Expect(inc.OpenTimelog().Remote).To(Equal(true))
						})

						It("should set state to InProgress", func() {
							Expect(inc.State()).To(Equal(StateInProgress))
						})
					})

					Context("but the incident is not in New, InProgress or OnHold state", func() {
						var inc Incident

						BeforeEach(func() {
							feUUID := fieldEngineer.UUID()
							inc = Incident{
								FieldEngineerID: &feUUID,
							}
							err := inc.SetState(StatePreOnHold)
							Expect(err).To(BeNil())
						})

						It("should return error", func() {
							err := inc.StartWorking(actorUser, clock, false)
							Expect(err).NotTo(BeNil())
							Expect(err.Error()).To(Equal("ticket is not in New, InProgress nor OnHold state"))
						})
					})
				})
			})
		})
	})

	Describe("StopWorking()", func() {
		When("called by actor that is not field engineer", func() {
			var inc Incident

			BeforeEach(func() {
				inc = Incident{}
			})

			It("should return error", func() {
				err := inc.StopWorking(actorUser, clock, "summary")
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("user is not field engineer, only assigned field engineer can stop working"))
			})
		})

		When("called by field engineer actor", func() {
			BeforeEach(func() {
				feUUID := fieldEngineer.UUID()
				actorUser.SetFieldEngineerID(&feUUID)
			})

			Context("but incident has no field engineer assigned", func() {
				var inc Incident

				BeforeEach(func() {
					inc = Incident{}
				})

				It("should return error", func() {
					err := inc.StopWorking(actorUser, clock, "summary")
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
						err := inc.StopWorking(actorUser, clock, "summary")
						Expect(err).NotTo(BeNil())
						Expect(err.Error()).To(Equal("user is not assigned as field engineer, only assigned field engineer can stop working"))
					})
				})

				Context("and the actor is assigned as field engineer", func() {
					Context("but the incident has no open timelog", func() {
						var inc Incident

						BeforeEach(func() {
							feUUID := fieldEngineer.UUID()
							inc = Incident{
								FieldEngineerID: &feUUID,
							}
							err := inc.SetState(StateInProgress)
							Expect(err).To(BeNil())
						})

						It("should return error", func() {
							err := inc.StopWorking(actorUser, clock, "summary")
							Expect(err).NotTo(BeNil())
							Expect(err.Error()).To(Equal("ticket does not have an open timelog"))
						})
					})

					Context("and the incident has an open timelog", func() {
						var inc Incident

						BeforeEach(func() {
							feUUID := fieldEngineer.UUID()
							inc = Incident{
								FieldEngineerID: &feUUID,
							}
							err := inc.SetState(StateInProgress)
							Expect(err).To(BeNil())

							inc.SetOpenTimelog(&timelog.Timelog{
								Start: clock.NowFormatted(),
							})
						})

						JustBeforeEach(func() {
							clock.AddTime(time.Hour)
							Expect(inc.HasOpenTimelog()).To(BeTrue())
							err := inc.StopWorking(actorUser, clock, "summary")
							Expect(err).To(BeNil())
						})

						It("should calculate correct Work in the open timelog", func() {
							Expect(inc.OpenTimelog().End).To(Equal(clock.NowFormatted()))
							Expect(inc.OpenTimelog().Work).To(Equal(uint(3600)))
						})

						It("should set the VisitSummary param in the open timelog", func() {
							Expect(inc.OpenTimelog().VisitSummary).To(Equal("summary"))
						})

						It("should not change state", func() {
							Expect(inc.State()).To(Equal(StateInProgress))
						})
					})
				})
			})
		})
	})

	Describe("Cancel()", func() {
		When("incident is in New' state", func() {
			var inc Incident

			BeforeEach(func() {
				inc = Incident{}
				err := inc.SetState(StateNew)
				Expect(err).To(BeNil())
			})

			JustBeforeEach(func() {
				err := inc.Cancel(actorUser)
				Expect(err).To(BeNil())
			})

			It("should set state to Cancelled", func() {
				Expect(inc.State()).To(Equal(StateCancelled))
			})
		})

		When("incident is NOT in New state", func() {
			var inc Incident

			BeforeEach(func() {
				inc = Incident{}
				err := inc.SetState(StateInProgress)
				Expect(err).To(BeNil())
			})

			It("should return error", func() {
				err := inc.Cancel(actorUser)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("ticket can be cancelled only in New state"))
			})
		})
	})
})
