import co from 'co'

exports.up = function(knex, Promise) {
	return co(function*() {
		yield knex
			.schema
			.table('feed_line', t => {
				t.index('project_id')
				t.index('step_id')
			})
		yield knex
			.schema
			.table('feed_line_log', t => {
				t.index('flu_id')
				t.index('step_id')
				t.index('master_flu_id')
			})
	})
};

exports.down = function(knex, Promise) {
	return co(function*() {
		yield knex
			.schema
			.table('feed_line', t => {
				t.dropIndex('project_id')
				t.dropIndex('step_id')
			})
		yield knex
			.schema
			.table('feed_line_log', t => {
				t.dropIndex('flu_id')
				t.dropIndex('step_id')
				t.dropIndex('master_flu_id')
			})
	})
};
