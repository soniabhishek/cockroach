
exports.up = function(knex, Promise) {
  	return knex.raw(`INSERT INTO step_type values (10,'START_STEP')`)
};

exports.down = function(knex, Promise) {
	return knex.raw('DELETE FROM step_type where id = 10')
};
