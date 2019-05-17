package reporthelper_test

import (
	"github.com/jpatel-pivotal/si-usage-report/reporthelper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reportHelper", func() {
	var subject *reporthelper.Report

	BeforeEach(func() {
		subject = new(reporthelper.Report)
	})

	Context("new report is created", func() {
		It("is empty", func() {
			Expect(subject.Orgs).To(BeEmpty())

		})
	})
	Context("adding an org", func() {
		var org, existingOrg reporthelper.Org
		org = reporthelper.Org{
			OrgName: "some-org",
			Spaces:  nil,
		}
		It("adds the org if the list of orgs is empty", func() {
			Expect(subject.Orgs).To(BeEmpty())
			subject.AddOrg(org)
			Expect(subject.Orgs[0]).To(Equal(org))
		})
		It("adds the org if the org is not in the list", func() {
			existingOrg = reporthelper.Org{OrgName: "some-other-org"}
			subject.Orgs = append(subject.Orgs, existingOrg)
			Expect(len(subject.Orgs)).To(Equal(1))
			subject.AddOrg(org)
			Expect(len(subject.Orgs)).To(Equal(2))
			Expect(subject.Orgs[0].OrgName).To(Equal("some-other-org"))
			Expect(subject.Orgs[1].OrgName).To(Equal("some-org"))

		})
		It("does not add the org if the org is already in the list", func() {
			existingOrg = reporthelper.Org{OrgName: "some-org"}
			subject.Orgs = append(subject.Orgs, existingOrg)
			Expect(len(subject.Orgs)).To(Equal(1))
			subject.AddOrg(org)
			Expect(len(subject.Orgs)).To(Equal(1))
			Expect(subject.Orgs[0].OrgName).To(Equal("some-org"))

		})
		Context("adding spaces to orgs", func() {
			var newSpace, existingSpace reporthelper.Space
			BeforeEach(func() {
				subject.Orgs = append(subject.Orgs, existingOrg)
			})
			It("adds the space if the space list is empty", func() {
				Expect(subject.Orgs[0].Spaces).To(BeEmpty())
				newSpace = reporthelper.Space{Name: "some-space"}
				subject.AddSpace(newSpace)
				Expect(len(subject.Orgs[0].Spaces)).To(Equal(1))

			})
			It("adds the space if the space is not in the list", func() {
				existingSpace = reporthelper.Space{
					Name:     "some-other-space",
					Products: nil,
				}
				subject.Orgs[0].Spaces = append(subject.Orgs[0].Spaces, existingSpace)
				Expect(subject.Orgs[0].Spaces[0].Name).To(Equal("some-other-space"))
				newSpace = reporthelper.Space{Name: "some-space"}
				subject.AddSpace(newSpace)
				Expect(len(subject.Orgs[0].Spaces)).To(Equal(2))
			})
			It("does not add the space if the space is already in the list", func() {
				existingSpace = reporthelper.Space{
					Name:     "some-space",
					Products: nil,
				}
				subject.Orgs[0].Spaces = append(subject.Orgs[0].Spaces, existingSpace)
				Expect(subject.Orgs[0].Spaces[0].Name).To(Equal("some-space"))
				newSpace = reporthelper.Space{Name: "some-space"}
				subject.AddSpace(newSpace)
				Expect(len(subject.Orgs[0].Spaces)).To(Equal(1))
			})
		})
		Context("adding products to orgs", func() {
			var newProduct, existingProduct reporthelper.Product
			BeforeEach(func() {
				subject.Orgs = append(subject.Orgs, existingOrg)
				subject.Orgs[0].Spaces = append(subject.Orgs[0].Spaces, reporthelper.Space{Name: "some-space"})
			})
			It("adds the product if the product list is empty", func() {
				Expect(subject.Orgs[0].Spaces[0].Products).To(BeEmpty())
				newProduct = reporthelper.Product{Name: "some-product"}
				subject.AddProduct(newProduct)
				Expect(len(subject.Orgs[0].Spaces)).To(Equal(1))

			})
			It("adds the product if the product is not in the list", func() {
				existingProduct = reporthelper.Product{
					Name:  "some-other-product",
					Plans: nil,
				}
				subject.Orgs[0].Spaces[0].Products = append(subject.Orgs[0].Spaces[0].Products, existingProduct)
				Expect(subject.Orgs[0].Spaces[0].Products[0].Name).To(Equal("some-other-product"))
				newProduct = reporthelper.Product{Name: "some-product"}
				subject.AddProduct(newProduct)
				Expect(len(subject.Orgs[0].Spaces[0].Products)).To(Equal(2))
			})
			It("does not add the product if the product is already in the list", func() {
				existingProduct = reporthelper.Product{
					Name:  "some-product",
					Plans: nil,
				}
				subject.Orgs[0].Spaces[0].Products = append(subject.Orgs[0].Spaces[0].Products, existingProduct)
				Expect(subject.Orgs[0].Spaces[0].Products[0].Name).To(Equal("some-product"))
				newProduct = reporthelper.Product{Name: "some-product"}
				subject.AddProduct(newProduct)
				Expect(len(subject.Orgs[0].Spaces)).To(Equal(1))
			})
		})
		Context("adding plans to orgs", func() {
			var newPlan ,anotherPlan, existingPlan reporthelper.Plan
			BeforeEach(func() {
				subject.Orgs = append(subject.Orgs, existingOrg)
				subject.Orgs[0].Spaces = append(subject.Orgs[0].Spaces, reporthelper.Space{Name: "some-space"})
				subject.Orgs[0].Spaces[0].Products = append(subject.Orgs[0].Spaces[0].Products, reporthelper.Product{Name: "some-product"})
			})
			It("adds the plan if the product list is empty", func() {
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans).To(BeEmpty())
				newPlan = reporthelper.Plan{ProductName: "some-product", PlanName: "some-plan"}
				subject.AddPlan(newPlan)
				Expect(len(subject.Orgs[0].Spaces)).To(Equal(1))

			})
			It("adds the plan if the plan is not in the list", func() {
				existingPlan = reporthelper.Plan{
					ProductName:  "some-product",
					PlanName:     "some-other-plan",
					InstanceName: "test-si-1",
				}
				subject.Orgs[0].Spaces[0].Products[0].Plans = append(subject.Orgs[0].Spaces[0].Products[0].Plans, existingPlan)
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[0].PlanName).To(Equal("some-other-plan"))
				newPlan = reporthelper.Plan{
					ProductName: "some-product",
					PlanName:          "some-plan",
					InstanceName: "test-si-2",
				}
				anotherPlan = reporthelper.Plan{
					ProductName: "some-product",
					PlanName:          "some-plan",
					InstanceName: "test-si-3",
				}
				subject.AddPlan(newPlan)
				subject.AddPlan(anotherPlan)
				Expect(len(subject.Orgs[0].Spaces[0].Products[0].Plans)).To(Equal(3))
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[0]).To(Equal(existingPlan))
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[1]).To(Equal(newPlan))
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[2]).To(Equal(anotherPlan))
			})
			It("does not add the plan if the plan is already in the list and increments instanceocunt", func() {
				existingPlan = reporthelper.Plan{
					ProductName:  "some-product",
					PlanName:     "some-plan",
					InstanceName: "test-si-1",
				}
				subject.Orgs[0].Spaces[0].Products[0].Plans = append(subject.Orgs[0].Spaces[0].Products[0].Plans, existingPlan)
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[0].PlanName).To(Equal("some-plan"))
				newPlan = reporthelper.Plan{
					ProductName: "some-product",
					PlanName:          "some-plan",
					InstanceName: "test-si-2",
				}
				anotherPlan = reporthelper.Plan{
					ProductName: "some-product",
					PlanName:          "some-plan",
					InstanceName: "test-si-2",
				}
				subject.AddPlan(newPlan)
				Expect(len(subject.Orgs[0].Spaces[0].Products[0].Plans)).To(Equal(2))
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[0]).To(Equal(existingPlan))
				Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[1]).To(Equal(newPlan))
				//Expect(subject.Orgs[0].Spaces[0].Products[0].Plans[2]).To(Equal(anotherPlan))
			})
		})
	})
})
