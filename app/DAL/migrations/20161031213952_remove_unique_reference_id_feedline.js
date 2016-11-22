import co from 'co'

exports.up = function(knex, Promise) {
	return co(function * () {
		yield knex.raw(`ALTER TABLE feed_line DROP CONSTRAINT IF EXISTS feed_line_reference_id_project_id_unique;`)
		yield knex.raw(`ALTER TABLE feed_line DROP CONSTRAINT IF EXISTS feed_line_project_id_reference_id_unique;`)
	})
};

exports.down = function(knex, Promise) {
	// putting a down will cause it to fail later. why? think about it...
	//return knex.schema
	//	.table('feed_line',t.dropUnique(['project_id','reference_id']))
};
