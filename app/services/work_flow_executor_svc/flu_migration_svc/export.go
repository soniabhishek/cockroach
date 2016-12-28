package flu_migration_svc

import "github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"

func New() IFluMigrationService {
	return &fluMigrationService{
		fluRepo: feed_line_repo.New(),
	}
}
