package router_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouter_Config(t *testing.T) {

	flu := feed_line.FLU{FeedLineUnit: models.FeedLineUnit{ID: uuid.NewV4(), StepId: uuid.FromStringOrNil("662cda89-ce3f-40f0-ae54-0c8dff082502"), Data: models.JsonF{"rc_check_res_success": true}, Build: models.JsonF{"rc_check_res_success": true}}}
	flu2 := feed_line.FLU{FeedLineUnit: models.FeedLineUnit{ID: uuid.NewV4(), StepId: uuid.FromStringOrNil("662cda89-ce3f-40f0-ae54-0c8dff082502"), Data: models.JsonF{"rc_check_res_success": true}, Build: models.JsonF{"rc_check_res_success": false}}}

	routerGetter := newRouteGetter()
	s, err := routerGetter.GetNextStep(flu)
	s2, err2 := routerGetter.GetNextStep(flu2)
	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.ObjectsAreEqual(s, uuid.FromStringOrNil("a3960044-d575-4b73-b208-c2764cee2062"))
	assert.ObjectsAreEqual(s2, uuid.FromStringOrNil("93367a4b-bdb4-4a00-803b-ca52b5fb6c7b"))
}
