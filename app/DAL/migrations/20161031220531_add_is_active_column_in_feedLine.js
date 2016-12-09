
exports.up = function(knex, Promise) {
  return knex.schema
		.table('feed_line',t=>{

			t.boolean('is_active')
				.notNullable()
				.defaultTo(true)
			t.boolean('is_master')
				.notNullable()
				.defaultTo(true)
		})
};

exports.down = function(knex, Promise) {
  return knex.schema.table('feed_line',t=> {
	  t.dropColumn('is_active')
	  t.dropColumn('is_master')
  })
};
