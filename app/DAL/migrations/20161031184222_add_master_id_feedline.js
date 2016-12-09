import co from 'co'

exports.up = function(knex, Promise) {
  return co(function * () {
	  yield knex.schema
		  .table('feed_line', t => {
			  t.uuid('master_id')
				  .index()
		  })
	  yield knex.raw(`UPDATE feed_line set master_id = id;`)
	  yield knex.raw(`ALTER TABLE feed_line ALTER COLUMN master_id SET NOT NULL;`)

	  yield knex.schema
		  .table('feed_line_log', t => {
			  t.uuid('master_flu_id')
		  })
	  yield knex.raw(`UPDATE feed_line_log set master_flu_id = flu_id;`)
  })
};

exports.down = function(knex, Promise) {
	return co(function * () {
		yield knex.schema
			.table('feed_line', t => {
				t.dropColumn('master_id')
			})
		yield knex.schema
			.table('feed_line_log', t => {
				t.dropColumn('master_flu_id')
			})
	})
};
