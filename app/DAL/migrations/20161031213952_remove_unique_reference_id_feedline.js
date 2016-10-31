
exports.up = function(knex, Promise) {
  return knex.schema
		.table('feed_line',t=>t.dropUnique(['reference_id','project_id']))
};

exports.down = function(knex, Promise) {
	// putting a down will cause it to fail later. why? think about it...
	//return knex.schema
	//	.table('feed_line',t.dropUnique(['project_id','reference_id']))
};
