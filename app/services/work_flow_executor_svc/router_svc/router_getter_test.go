package router_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouter_Config(t *testing.T) {

	flu := feed_line.FLU{FeedLineUnit: models.FeedLineUnit{ID: uuid.FromStringOrNil("eaf3b794-40c1-4132-a1b0-3e6a81c939ae"), StepId: uuid.FromStringOrNil("662cda89-ce3f-40f0-ae54-0c8dff082502"), Data: models.JsonF{"rc_check_res_success": true}, Build: models.JsonF{"rc_check_res_success": true}}}
	routerGetter := newRouteGetter()
	s, err := routerGetter.GetNextStep(flu)
	assert.NoError(t, err)
	assert.ObjectsAreEqual(s, uuid.FromStringOrNil("a3960044-d575-4b73-b208-c2764cee2062"))
}
