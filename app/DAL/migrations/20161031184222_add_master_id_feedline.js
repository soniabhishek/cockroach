
exports.up = function(knex, Promise) {
  return knex.schema
		.table('feed_line', t => {
			t.uuid('master_id')
				.index()
		})
};

exports.down = function(knex, Promise) {
	return knex.schema
		.table('feed_line', t => {
			t.dropColumn('master_id')
		})
};
