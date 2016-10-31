import co from 'co'

exports.up = function(knex, Promise) {
  return co(function * () {
	  yield knex.schema
		  .table('feed_line', t => {
			  t.uuid('master_id')
				  .index()
		  })
	  yield knex.schema
		  .table('feed_line_log', t => {
			  t.uuid('master_id')
		  })
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
				t.dropColumn('master_id')
			})
	})
};
