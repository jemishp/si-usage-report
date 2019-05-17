package reporthelper

type Report struct {
	Orgs []Org
}

type Org struct {
	OrgName string
	Spaces  []Space
}

type Space struct {
	Name     string
	Products []Product
}

type Product struct {
	Name  string
	Plans []Plan
}

type Plan struct {
	ProductName  string
	PlanName     string
	InstanceName string
}

func (r *Report) AddOrg(orgToAdd Org) {
	if len(r.Orgs) == 0 {
		r.Orgs = append(r.Orgs, orgToAdd)
	}
	for _, org := range r.Orgs {
		if org.OrgName != orgToAdd.OrgName {
			//org does not exist so add it
			r.Orgs = append(r.Orgs, orgToAdd)
		}
	}
}

func (r *Report) AddSpace(spaceToAdd Space) {
	for k, org := range r.Orgs {
		if len(org.Spaces) == 0 {
			r.Orgs[k].Spaces = append(r.Orgs[k].Spaces, spaceToAdd)
		}
		for _, space := range org.Spaces {
			if space.Name != spaceToAdd.Name {
				//space does not exist so add it
				r.Orgs[k].Spaces = append(r.Orgs[k].Spaces, spaceToAdd)
			}
		}
	}

}

func (r *Report) AddProduct(productToAdd Product) {
	for k, org := range r.Orgs {
		for a, space := range org.Spaces {
			if len(space.Products) == 0 {
				r.Orgs[k].Spaces[a].Products = append(r.Orgs[k].Spaces[a].Products, productToAdd)
			}
			for _, product := range space.Products {
				if product.Name != productToAdd.Name {
					//product does not exist so add it
					r.Orgs[k].Spaces[a].Products = append(r.Orgs[k].Spaces[a].Products, productToAdd)
				}
			}

		}
	}

}

func (r *Report) AddPlan(planToAdd Plan) {
	for k, org := range r.Orgs {
		for a, space := range org.Spaces {
			for b, product := range space.Products {
				if len(product.Plans) == 0 {
					r.Orgs[k].Spaces[a].Products[b].Plans = append(r.Orgs[k].Spaces[a].Products[b].Plans, planToAdd)
				}
				for _, plan := range product.Plans {

					if plan.ProductName != planToAdd.ProductName {
						//plan does not exist so add it
						r.Orgs[k].Spaces[a].Products[b].Plans = append(r.Orgs[k].Spaces[a].Products[b].Plans, planToAdd)
						continue
					}
					if plan.ProductName == planToAdd.ProductName &&
						plan.PlanName != planToAdd.PlanName {
						//plan names match so increment instances
						r.Orgs[k].Spaces[a].Products[b].Plans = append(r.Orgs[k].Spaces[a].Products[b].Plans, planToAdd)
						continue
					}
					//if plan.ProductName == planToAdd.ProductName &&
					//	plan.PlanName == planToAdd.PlanName &&
					//	plan.InstanceName != planToAdd.InstanceName {
					//	//plan names match so increment instances
					//	r.Orgs[k].Spaces[a].Products[b].Plans = append(r.Orgs[k].Spaces[a].Products[b].Plans, planToAdd)
					//
					//}
				}

			}

		}
	}

}
