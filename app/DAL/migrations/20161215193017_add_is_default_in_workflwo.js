
exports.up = function(knex, Promise) {
  return knex.schema
		.table('work_flow',t=> {
			t.boolean('is_default')
				.notNullable()
				.defaultTo(false)
		})
};

exports.down = function(knex, Promise) {
	return knex.schema
		.table('work_flow',t=> {
			t.dropColumn('is_default')
		})
};
