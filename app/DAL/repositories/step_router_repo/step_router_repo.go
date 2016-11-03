package step_router_repo

import (
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type stepRouteRepo struct {
	Db repositories.IDatabase
}

var _ IStepRoutesRepo = (*stepRouteRepo)(nil)

const routesTable = "routes"

func (s *stepRouteRepo) GetRoutesByStepId(stepId uuid.UUID) (routes []models.Route, err error) {

	_, err = s.Db.Select(&routes, "select * from "+routesTable+" where step_id = $1", stepId)
	return
}

func (s *stepRouteRepo) GetRoutesWithLogicByStepId(stepId uuid.UUID) (routesWithLogic []models.RouteWithLogicGate, err error) {

	var routes []models.Route

	_, err = s.Db.Select(&routes, "select * from routes where step_id = $1", stepId)
	if err != nil {
		return
	}

	var logicGates []models.LogicGate

	_, err = s.Db.Select(&logicGates, `select l.* from logic_gate l
		inner join routes r on r.logic_gate_id = l.id
		where r.step_id = $1`, stepId)
	if err != nil {
		return
	}

	for _, route := range routes {
		for _, logicGate := range logicGates {

			if logicGate.ID == route.LogicGateId {
				routesWithLogic = append(routesWithLogic, models.RouteWithLogicGate{
					Route:     route,
					LogicGate: logicGate,
				})
				break
			}
		}
	}

	return
}

func (s *stepRouteRepo) Save(r models.Route) error {
	panic("not implemented")
	return nil
}

func (s *stepRouteRepo) GetRoutesByWorkFlowId(workFlowId uuid.UUID) (routes []models.Route, err error) {

	_, err = s.Db.Select(&routes, `select r.* from routes r inner join
	step s on s.id = r.step_id where s.work_flow_id = $1
	`, workFlowId.String())
	return
}

func (s *stepRouteRepo) AddMany(routes []models.Route) (err error) {
	var routesInterface []interface{} = make([]interface{}, len(routes))
	for i, _ := range routes {
		routesInterface[i] = &routes[i]
	}

	err = s.Db.Insert(routesInterface...)
	return
}
func (s *stepRouteRepo) UpdateMany(routes []models.Route) (response int64, err error) {
	var routesInterface []interface{} = make([]interface{}, len(routes))
	for i, _ := range routes {
		routesInterface[i] = &routes[i]
	}

	response, err = s.Db.Update(routesInterface...)
	return
}
func (s *stepRouteRepo) DeleteMany(routes []models.Route) (response int64, err error) {
	var routesInterface []interface{} = make([]interface{}, len(routes))
	for i, _ := range routes {
		routesInterface[i] = &routes[i]
	}

	response, err = s.Db.Delete(routesInterface...)
	return
}
